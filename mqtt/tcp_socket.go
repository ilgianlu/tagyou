package mqtt

import (
	"bufio"
	"log/slog"
	"net"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func StartMQTT(port string, connections model.Connections) {
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
		go handleTcpConnection(conn, connections)
	}
}

func handleTcpConnection(conn net.Conn, connections model.Connections) {
	defer conn.Close()

	sessionTimestamp := time.Now().Unix()
	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		Conn:           conn,
		Router:         routers.NewDefault(conf.ROUTER_MODE, connections),
		LastConnect:    sessionTimestamp,
		LastSeen:       sessionTimestamp,
	}

	events := make(chan *packet.Packet)
	go event.RangePackets(&session, events)
	defer disconnectClient(&session, events)

	scanner := bufio.NewScanner(conn)
	scanner.Split(packetSplit(&session))

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				onSocketUpButSilent(&session)
			}
		}

		b := scanner.Bytes()
		p, err := packet.PacketParse(&session, b)
		if err != nil {
			return
		}

		session.Mu.RLock()
		clientId := session.ClientId
		keepAlive := session.KeepAlive
		session.Mu.RUnlock()

		slog.Debug("[MQTT] session setting read deadline of seconds", "client-id", clientId, "keep-alive", keepAlive*2)
		derr := conn.SetReadDeadline(time.Now().Add(time.Duration(keepAlive*2) * time.Second))
		if derr != nil {
			slog.Error("[MQTT] cannot set read deadline", "err", derr)
		}

		events <- &p
	}
}

func packetSplit(session *model.RunningSession) func(b []byte, atEOF bool) (int, []byte, error) {
	return func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(b) == 0 && atEOF {
			wasConnected := onSocketDownClosed(session)
			if wasConnected {
				return 0, nil, nil
			}
			return 0, b, bufio.ErrFinalToken
		}
		pb, err := packet.ReadFromByteSlice(b)
		if err != nil {
			if !atEOF {
				return 0, nil, nil
			}
			slog.Debug("[MQTT] error reading bytes - session", "client-id", session.GetClientId(), "err", err.Error())
			return 0, pb, bufio.ErrFinalToken
		}
		return len(pb), pb, nil
	}
}

func onSocketUpButSilent(session *model.RunningSession) bool {
	slog.Debug("[MQTT] keepalive not respected!", "keep-alive", session.KeepAlive*2)
	if session.GetClientId() != "" {
		slog.Debug("[MQTT] will due to keepalive not respected!", "client-id", session.GetClientId(), "last-connect", session.LastConnect)
		event.SendWill(session)
		return true
	}
	return false
}

func onSocketDownClosed(session *model.RunningSession) bool {
	slog.Debug("[MQTT] socket closed!", "client-id", session.GetClientId())
	if session.GetClientId() != "" {
		slog.Debug("[MQTT] client was connected!", "client-id", session.GetClientId(), "last-connect", session.LastConnect)
		event.SendWill(session)
		return true
	}
	return false
}

func disconnectClient(session *model.RunningSession, e chan<- *packet.Packet) {
	session.Mu.RLock()
	clientId := session.ClientId
	session.Mu.RUnlock()
	if clientId != "" {
		session.Router.RemoveDestination(clientId)
		persistence.SessionRepository.DisconnectSession(clientId)
	}
	close(e)
}
