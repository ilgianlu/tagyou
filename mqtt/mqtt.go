package mqtt

import (
	"bufio"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func StartMQTT(port string) {
	conf.Loader()
	db, err := openDb()
	if err != nil {
		log.Fatal().Err(err).Msg("[MQTT] failed to connect database")
	}
	log.Info().Msg("[MQTT] db connected !")
	defer closeDb(db)

	model.Migrate(db)

	connections := model.Connections{}
	connections.Conns = make(map[string]net.Conn)
	events := make(chan *packet.Packet)
	outQueue := make(chan out.OutData)

	go event.RangeEvents(&connections, db, events, outQueue)
	go out.RangeOutQueue(&connections, db, outQueue)

	if conf.CLEAN_EXPIRED_SESSIONS {
		StartSessionCleaner(db)
	}
	if conf.CLEAN_EXPIRED_RETRIES {
		StartRetryCleaner(db)
	}
	startTCP(events, port)
}

func openDb() (*gorm.DB, error) {
	logLevel := logger.Silent
	if os.Getenv("DEBUG") != "" {
		logLevel = logger.Info
	}
	return gorm.Open(sqlite.Open(os.Getenv("DB_PATH")+os.Getenv("DB_NAME")), &gorm.Config{
		Logger: logger.New(
			&log.Logger,
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logLevel,
				Colorful:      true,
			},
		),
	})
}

func closeDb(db *gorm.DB) {
	sql, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("could not close DB")
		return
	}
	sql.Close()
}

func startTCP(events chan<- *packet.Packet, port string) {
	// start tcp socket
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Error().Err(err).Msg("[MQTT] tcp listen error")
		return
	}
	log.Info().Msgf("[MQTT] mqtt listening on %s", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Err(err).Msg("[MQTT] tcp accept error")
		}
		go handleConnection(conn, events)
	}
}

func handleConnection(conn net.Conn, events chan<- *packet.Packet) {
	defer conn.Close()

	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		Conn:           conn,
	}

	scanner := bufio.NewScanner(conn)
	packetSplit := func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(b) == 0 && atEOF {
			// socket down - closed
			if session.GetClientId() != "" {
				willEvent(&session, events)
				disconnectClient(&session, events)
				return
			}
			return 0, b, bufio.ErrFinalToken
		}
		pb, err := packet.ReadFromByteSlice(b)
		if err != nil {
			log.Error().Err(err).Msg("[MQTT] error reading bytes")
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
				log.Debug().Msgf("[MQTT] keepalive of %d seconds not respected!", session.KeepAlive*2)
				if session.GetClientId() != "" {
					willEvent(&session, events)
					disconnectClient(&session, events)
					return
				}
			}
		}

		b := scanner.Bytes()
		p, err := packet.Start(b)
		if err != nil {
			log.Error().Err(err).Msg("[MQTT] Start err")
			return
		}
		p.Session = &session
		parseErr := p.Parse()
		if parseErr != 0 {
			log.Error().Msgf("[MQTT] parse err %d", parseErr)
		}

		session.Mu.RLock()
		clientId := session.ClientId
		keepAlive := session.KeepAlive
		session.Mu.RUnlock()

		log.Debug().Msgf("[MQTT] session %s setting read deadline of %d seconds", clientId, keepAlive*2)
		derr := conn.SetReadDeadline(time.Now().Add(time.Duration(keepAlive*2) * time.Second))
		if derr != nil {
			log.Error().Err(derr).Msg("[MQTT] cannot set read deadline")
			defer conn.Close()
		}

		events <- &p
	}

	// log.Println("Out of Scan loop!")
}

func willEvent(session *model.RunningSession, e chan<- *packet.Packet) {
	p := packet.Packet{Session: session, Event: packet.EVENT_WILL_SEND}
	e <- &p
}

func disconnectClient(session *model.RunningSession, e chan<- *packet.Packet) {
	p := packet.Packet{Session: session, Event: packet.EVENT_DISCONNECT}
	e <- &p
}
