package event

import (
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func RangePackets(session *model.RunningSession, packets <-chan *packet.Packet) {
	for p := range packets {
		managePacket(session, p)
	}
}

func managePacket(session *model.RunningSession, p *packet.Packet) {
	slog.Debug("//!! is session connected?", "connected?", session.GetConnected())
	slog.Debug("//!! packet arriving", "packet-type", p.PacketType())
	if !session.GetConnected() && p.PacketType() != packet.PACKET_TYPE_CONNECT {
		slog.Warn("can accept only connect, disconnecting...", "packet-type", p.PacketType())
		session.Conn.Close()
		return
	}
	switch p.PacketType() {
	case packet.PACKET_TYPE_CONNECT:
		slog.Debug("//!! EVENT client connect", "client-id", session.GetClientId())
		if session.GetConnected() {
			slog.Debug("//!! EVENT double connect event, disconnecting...", "client-id", session.GetClientId())
			clientDisconnect(session, session.GetClientId())
			return
		}
		onConnect(session)
		return
	case packet.PACKET_TYPE_DISCONNECT:
		slog.Debug("//!! EVENT client disconnect", "client-id", session.GetClientId())
		clientDisconnect(session, session.GetClientId())
		return
	case packet.PACKET_TYPE_PINGREQ:
		slog.Debug("//!! EVENT client ping", "client-id", session.GetClientId())
		onPing(session)
	case packet.PACKET_TYPE_SUBSCRIBE:
		slog.Debug("//!! EVENT client subscribed", "client-id", session.GetClientId())
		onSubscribe(session, p)
	case packet.PACKET_TYPE_UNSUBSCRIBE:
		slog.Debug("//!! EVENT client unsubscribed", "client-id", session.GetClientId())
		onUnsubscribe(session, p)
	case packet.PACKET_TYPE_PUBLISH:
		slog.Debug("//!! EVENT client published", "topic", p.Topic, "client-id", session.GetClientId(), "qos", p.QoS())
		OnPublish(session, p)
	case packet.PACKET_TYPE_PUBACK:
		slog.Debug("//!! EVENT client acked message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPuback(session, p)
	case packet.PACKET_TYPE_PUBREC:
		slog.Debug("//!! EVENT pub received message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPubrec(session, p)
	case packet.PACKET_TYPE_PUBREL:
		slog.Debug("//!! EVENT pub releases message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPubrel(session, p)
	case packet.PACKET_TYPE_PUBCOMP:
		slog.Debug("//!! EVENT pub complete message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		clientPubcomp(session.GetClientId(), p.PacketIdentifier(), p.ReasonCode)
	}
}

func onPing(session *model.RunningSession) {
	toSend := packet.PingResp()
	session.Router.Send(session.GetClientId(), toSend.ToByteSlice())
}

func clientDisconnect(session *model.RunningSession, clientId string) {
	session.SetConnected(false)
	if session.Router.DestinationExists(clientId) {
		needDisconnection := needDisconnection(session)
		if !needDisconnection {
			return
		}
		session.Router.RemoveDestination(clientId)
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
