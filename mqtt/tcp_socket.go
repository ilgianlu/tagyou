package mqtt

import (
	"bufio"
	"net"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

func StartMQTT(port string, connections *model.Connections) {
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
		go handleTcpConnection(connections, conn)
	}
}

func handleTcpConnection(connections *model.Connections, conn net.Conn) {
	defer conn.Close()

	events := make(chan *packet.Packet)
	go event.RangeEvents(connections, events)

	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		Conn:           conn,
		LastConnect:    time.Now().Unix(),
	}

	scanner := bufio.NewScanner(conn)
	scanner.Split(packetSplit(&session, events))

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				// socket up but silent
				log.Debug().Msgf("[MQTT] keepalive of %d seconds not respected!", session.KeepAlive*2)
				if session.GetClientId() != "" {
					log.Debug().Msgf("[MQTT] (%s:%d) will due to keepalive not respected!", session.GetClientId(), session.LastConnect)
					willEvent(&session, events)
					disconnectClient(&session, events)
					return
				}
			}
		}

		b := scanner.Bytes()
		p, err := packetParse(&session, b)
		if err != nil {
			return
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