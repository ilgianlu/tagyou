package mqtt

import (
	"bufio"
	"log/slog"
	"net"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/engine"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func StartMQTT(port string, connections model.Connections) {
	slog.Info("[MQTT] default router mode", "mode", conf.ROUTER_MODE)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("[MQTT] tcp listen error", "err", err)
		return
	}
	slog.Info("[MQTT] mqtt listening on", "tcp-port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("[MQTT] tcp accept error", "err", err)
		}
		go handleTCPConnection(conn, connections)
	}
}

func handleTCPConnection(conn net.Conn, connections model.Connections) {
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Warn("[MQTT] could not clean close connection", "err", err)
		}
	}()

	sessionTimestamp := time.Now().Unix()
	r := routers.NewDefault(conf.ROUTER_MODE, connections)
	e := engine.NewEngine()
	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		Conn:           conn,
		Router:         r,
		Engine:         e,
		LastConnect:    sessionTimestamp,
		LastSeen:       sessionTimestamp,
	}

	// the size of the channel is equivalent to:
	// 0 -> store client messages queue on the socket cache of the Operating System
	// 1 - Infinite -> store "some" client messages queue on channel (application) memory, rest on OS
	// packets := make(chan *packet.Packet, 5)
	packets := make(chan *packet.Packet)
	go rangePackets(&session, packets)
	defer disconnectClient(&session, packets)

	for {
		p := packet.Packet{}

		reader := bufio.NewReader(conn)
		err := p.ReadHeader(reader)
		if err != nil {
			slog.Debug("[MQTT] error reading header byte", "client-id", session.GetClientId(), "err", err)
			disconnectClient(&session, packets)
			return
		}

		err = p.ReadRemainingLength(reader)
		if err != nil {
			slog.Debug("[MQTT] error reading remaining length bytes", "client-id", session.GetClientId(), "err", err)
			disconnectClient(&session, packets)
			return
		}

		err = p.ReadRemainingBytes(reader)
		if err != nil {
			slog.Debug("[MQTT] error reading remaining bytes", "client-id", session.GetClientId(), "err", err)
			disconnectClient(&session, packets)
			return
		}

		errCode := p.Parse(&session)
		if errCode != 0 {
			slog.Debug("[MQTT] error parsing remaining bytes", "client-id", session.GetClientId())
			disconnectClient(&session, packets)
			return
		}
		clientID := session.ClientId
		keepAlive := session.KeepAlive

		slog.Debug("[MQTT] session setting read deadline of seconds", "client-id", clientID, "keep-alive", keepAlive*2)
		derr := conn.SetReadDeadline(time.Now().Add(time.Duration(keepAlive*2) * time.Second))
		if derr != nil {
			slog.Error("[MQTT] cannot set read deadline", "err", derr)
		}

		packets <- &p
	}
}

func disconnectClient(session *model.RunningSession, e chan<- *packet.Packet) {
	clientID := session.ClientId
	if clientID != "" {
		session.Router.RemoveDestination(clientID)
		persistence.SessionRepository.DisconnectSession(clientID)
	}
	close(e)
}
