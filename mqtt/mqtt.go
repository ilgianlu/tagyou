package mqtt

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartMQTT(port string) {
	conf.FORBID_ANONYMOUS_LOGIN = os.Getenv("FORBID_ANONYMOUS_LOGIN") == "true"
	conf.ACL_ON = os.Getenv("ACL_ON") == "true"
	conf.CLEAN_EXPIRED_SESSIONS = os.Getenv("CLEAN_EXPIRED_SESSIONS") == "true"
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_PATH")+os.Getenv("DB_NAME")), &gorm.Config{})
	if err != nil {
		log.Fatalf("[MQTT] failed to connect database %s", err)
	}
	log.Println("[MQTT] db connected !")
	defer closeDb(db)

	model.Migrate(db)

	connections := make(model.Connections)
	events := make(chan *packet.Packet, 1)
	outQueue := make(chan *out.OutData, 1)

	go event.RangeEvents(connections, db, events, outQueue)
	go out.RangeOutQueue(connections, db, outQueue)

	if conf.CLEAN_EXPIRED_SESSIONS {
		StartSessionCleaner(db)
	}
	startTCP(events, port)
}

func closeDb(db *gorm.DB) {
	sql, err := db.DB()
	if err != nil {
		log.Println("could not close DB", err)
		return
	}
	sql.Close()
}

func startTCP(events chan<- *packet.Packet, port string) {
	// start tcp socket
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("[MQTT] tcp listen error", err)
		return
	}
	log.Println("[MQTT] mqtt listening on", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[MQTT] tcp accept error", err)
		}
		go handleConnection(conn, events)
	}
}

func handleConnection(conn net.Conn, events chan<- *packet.Packet) {
	defer conn.Close()

	session := model.Session{
		Connected:      true,
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: conf.SESSION_MAX_DURATION_SECONDS,
		Conn:           conn,
	}

	scanner := bufio.NewScanner(conn)
	packetSplit := func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(b) == 0 && atEOF {
			// socket down - closed
			if session.ClientId != "" {
				willEvent(&session, events)
				disconnectClient(&session, events)
				return
			}
			return 0, b, bufio.ErrFinalToken
		}
		pb, err := packet.ReadFromByteSlice(b)
		if err != nil {
			log.Printf("[MQTT] %s\n", err)
			if !atEOF {
				return 0, nil, nil
			}
			return 0, pb, bufio.ErrFinalToken
		}
		return len(pb), pb, nil
	}
	scanner.Split(packetSplit)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				// socket up but silent
				log.Println("[MQTT] keepalive not respected!")
				if session.ClientId != "" {
					willEvent(&session, events)
					disconnectClient(&session, events)
					return
				}
			}
		}

		derr := conn.SetReadDeadline(time.Now().Add(time.Duration(session.KeepAlive*2) * time.Second))
		if derr != nil {
			log.Println("[MQTT] cannot set read deadline", derr)
			defer conn.Close()
		}

		b := scanner.Bytes()
		p, err := packet.Start(b)
		if err != nil {
			log.Printf("[MQTT] Start err %s\n", err)
			return
		}
		p.Session = &session
		parseErr := p.Parse()
		if parseErr != 0 {
			log.Printf("[MQTT] parseErr %d\n", parseErr)
		}
		events <- &p
	}

	// log.Println("Out of Scan loop!")
}

func willEvent(session *model.Session, e chan<- *packet.Packet) {
	p := packet.Packet{Session: session, Event: packet.EVENT_WILL_SEND}
	e <- &p
}

func disconnectClient(session *model.Session, e chan<- *packet.Packet) {
	p := packet.Packet{Session: session, Event: packet.EVENT_DISCONNECT}
	e <- &p
}
