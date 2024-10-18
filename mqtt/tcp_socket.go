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

func StartMQTT(port string, router routers.Router) {
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
		go handleTcpConnection(router, conn)
	}
}

func handleTcpConnection(router routers.Router, conn net.Conn) {
	defer conn.Close()

	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		Conn:           conn,
		LastConnect:    time.Now().Unix(),
	}

	events := make(chan *packet.Packet)
	go event.RangeEvents(router, &session, events)
	defer disconnectClient(router, &session, events)

	scanner := bufio.NewScanner(conn)
	scanner.Split(packetSplit(router,&session))

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				onSocketUpButSilent(router, &session)
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

func packetSplit(router routers.Router, session *model.RunningSession) func(b []byte, atEOF bool) (int, []byte, error) {
	return func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(b) == 0 && atEOF {
			wasConnected := onSocketDownClosed(router, session)
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

func onSocketUpButSilent(router routers.Router, session *model.RunningSession) bool {
	slog.Debug("[MQTT] keepalive not respected!", "keep-alive", session.KeepAlive*2)
	if session.GetClientId() != "" {
		slog.Debug("[MQTT] will due to keepalive not respected!", "client-id", session.GetClientId(), "last-connect", session.LastConnect)
    event.SendWill(router, session)
		return true
	}
	return false
}

func onSocketDownClosed(router routers.Router, session *model.RunningSession) bool {
	slog.Debug("[MQTT] socket down closed!")
	if session.GetClientId() != "" {
		slog.Debug("[MQTT] will due to socket down!", "client-id", session.GetClientId(), "last-connect", session.LastConnect)
		event.SendWill(router, session)
		return true
	}
	return false
}

func disconnectClient(router routers.Router, session *model.RunningSession, e chan<- *packet.Packet) {
	session.Mu.RLock()
	clientId := session.ClientId
	session.Mu.RUnlock()
	if clientId != "" {
		router.RemoveDestination(clientId)
		persistence.SessionRepository.DisconnectSession(clientId)
	}
	p := packet.Packet{Event: packet.EVENT_DISCONNECT}
	e <- &p
	close(e)
}
