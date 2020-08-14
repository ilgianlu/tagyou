package mqtt

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const MQTT_V5 = 5
const MQTT_V3_11 = 4

const TOPIC_SEPARATOR = "/"
const TOPIC_WILDCARD = "#"

func StartMQTT(port string) {
	conf.FORBID_ANONYMOUS_LOGIN = os.Getenv("FORBID_ANONYMOUS_LOGIN") == "true"
	conf.ACL_ON = os.Getenv("ACL_ON") == "true"
	db, err := gorm.Open("sqlite3", os.Getenv("DB_PATH")+os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalf("[MQTT] failed to connect database %s", err)
	}
	log.Println("[MQTT] db connected !")
	defer db.Close()

	model.Migrate(db)

	connections := make(Connections)
	events := make(chan Event, 1)
	outQueue := make(chan OutData, 1)

	go rangeEvents(connections, db, events, outQueue)
	go rangeOutQueue(connections, db, outQueue)

	startTCP(events, port)
}

func startTCP(events chan<- Event, port string) {
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

func handleConnection(conn net.Conn, events chan<- Event) {
	defer conn.Close()

	session := model.Session{
		Connected:      true,
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: conf.SESSION_MAX_DURATION_SECONDS,
		Conn:           conn,
	}

	derr := conn.SetReadDeadline(time.Now().Add(time.Duration(session.KeepAlive*2) * time.Second))
	if derr != nil {
		log.Println("[MQTT] cannot set read deadline", derr)
		defer conn.Close()
	}

	packets := make(chan Packet)
	go rangePackets(packets, events, &session)

	scanner := bufio.NewScanner(conn)
	packetSplit := func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		log.Println(len(b), b, atEOF)
		pb, err := ReadFromByteSlice(b)
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
		b := scanner.Bytes()
		// packet as a []bytes
		p, err := Start(b)
		if err != nil {
			log.Printf("[MQTT] %s\n", err)
			return
		}
		packets <- p
	}

	// err := scanner.Err()
	// if err != nil {
	// 	log.Println("could read packet", err)
	// 	if err, ok := err.(net.Error); ok && err.Timeout() {
	// 		log.Println("keepalive not respected!")
	// 		if session.ClientId != "" {
	// 			willEvent(session.ClientId, events)
	// 			disconnectClient(session.ClientId, events)
	// 		}
	// 	}

	// }

	// buffers := make(chan []byte, 1)

	// for {
	// 	buffer := make([]byte, 1024)
	// 	bytesCount, err := conn.Read(buffer)

	// 	if err != nil {
	// 		log.Println("could read packet", err)
	// 		if err, ok := err.(net.Error); ok && err.Timeout() {
	// 			log.Println("keepalive not respected!")
	// 			if session.ClientId != "" {
	// 				willEvent(session.ClientId, events)
	// 				disconnectClient(session.ClientId, events)
	// 			}
	// 			break
	// 		}
	// 		if err == io.EOF {
	// 			log.Println("connection closed!")
	// 			if session.ClientId != "" {
	// 				willEvent(session.ClientId, events)
	// 				disconnectClient(session.ClientId, events)
	// 			}
	// 			break
	// 		}
	// 	}
	// 	buffers <- buffer[:bytesCount]

	// 	go rangeBuffers(buffers, packets)

	// }
	// log.Println("abandon closed connection!")
}

func willEvent(clientId string, e chan<- Event) {
	var event Event
	event.eventType = EVENT_WILL_SEND
	event.clientId = clientId
	e <- event
}

func disconnectClient(clientId string, e chan<- Event) {
	var event Event
	event.eventType = EVENT_DISCONNECT
	event.clientId = clientId
	e <- event
}
