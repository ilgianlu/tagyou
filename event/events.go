package event

import (
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func RangePackets(router routers.Router, session *model.RunningSession, packets <-chan *packet.Packet) {
	for p := range packets {
		managePacket(router, session, p)
	}
}

func managePacket(router routers.Router, session *model.RunningSession, p *packet.Packet) {
	slog.Debug("//!! is session connected?", "connected?", session.GetConnected())
	slog.Debug("//!! packet arriving", "packet-type", p.PacketType())
	if !session.GetConnected() && p.PacketType() != packet.PACKET_TYPE_CONNECT {
		slog.Warn("can accept only connect, disconnecting...", "packet-type", p.PacketType())
		session.Conn.Close()
		return
	}
	switch p.PacketType() {
	case packet.PACKET_TYPE_CONNECT:
		slog.Debug("//!! EVENT client connect", "packet-type", p.PacketType(), "client-id", session.GetClientId())
		if session.GetConnected() {
			slog.Debug("//!! EVENT double connect event, disconnecting...", "client-id", session.GetClientId())
			clientDisconnect(router, session, session.GetClientId())
			return
		}
		onConnect(router, session)
		return
	case packet.PACKET_TYPE_DISCONNECT:
		slog.Debug("//!! EVENT client disconnect", "packet-type", p.PacketType(), "client-id", session.GetClientId())
		clientDisconnect(router, session, session.GetClientId())
		return
	}
	switch p.Event {
	case packet.EVENT_SUBSCRIBED:
		slog.Debug("//!! EVENT client subscribed", "event-type", p.Event, "client-id", session.GetClientId())
		onSubscribe(router, session, p)
	case packet.EVENT_UNSUBSCRIBED:
		slog.Debug("//!! EVENT client unsubscribed", "event-type", p.Event, "client-id", session.GetClientId())
		onUnsubscribe(router, session, p)
	case packet.EVENT_PUBLISH:
		slog.Debug("//!! EVENT client published", "event-type", p.Event, "topic", p.Topic, "client-id", session.GetClientId(), "qos", p.QoS())
		OnPublish(router, session, p)
	case packet.EVENT_PUBACKED:
		slog.Debug("//!! EVENT client acked message", "event-type", p.Event, "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPuback(session, p)
	case packet.EVENT_PUBRECED:
		slog.Debug("//!! EVENT pub received message", "event-type", p.Event, "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPubrec(router, session, p)
	case packet.EVENT_PUBRELED:
		slog.Debug("//!! EVENT pub releases message", "event-type", p.Event, "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPubrel(router, session, p)
	case packet.EVENT_PUBCOMPED:
		slog.Debug("//!! EVENT pub complete message", "event-type", p.Event, "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPubcomp(session.GetClientId(), p.PacketIdentifier(), p.ReasonCode)
	case packet.EVENT_PING:
		slog.Debug("//!! EVENT client ping", "event-type", p.Event, "client-id", session.GetClientId())
		onPing(router, session)
	}
}

func onPing(router routers.Router, session *model.RunningSession) {
	toSend := packet.PingResp()
	router.Send(session.GetClientId(), toSend.ToByteSlice())
}

func clientDisconnect(router routers.Router, session *model.RunningSession, clientId string) {
	session.SetConnected(false)
	if router.DestinationExists(clientId) {
		needDisconnection := needDisconnection(session)
		if !needDisconnection {
			return
		}
		router.RemoveDestination(clientId)
		persistence.SessionRepository.DisconnectSession(clientId)
	}
}

func saveRetain(session *model.RunningSession, p *packet.Packet) {
	var r model.Retain
	r.ClientID = session.ClientId
	r.Topic = p.Topic
	r.ApplicationMessage = p.ApplicationMessage()
	r.CreatedAt = time.Now().Unix()
	persistence.RetainRepository.Delete(r)
	if len(r.ApplicationMessage) > 0 {
		persistence.RetainRepository.Create(r)
	}
}

func needDisconnection(runningSession *model.RunningSession) bool {
	if session, ok := persistence.SessionRepository.SessionExists(runningSession.ClientId); ok {
		if session.GetLastConnect() > runningSession.LastConnect {
			// session persisted is newer then running memory session... device reconnected!
			// no need to send will
			slog.Debug("[MQTT] avoid disconnect! (client reconnected)", "client-id", session.GetClientId())
			return false
		}
	}
	return true
}
