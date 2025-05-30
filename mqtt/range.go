package mqtt

import (
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

func rangePackets(session *model.RunningSession, packets <-chan *packet.Packet) {
	for p := range packets {
		managePacket(session, p)
	}
}

func managePacket(session *model.RunningSession, p model.Packet) {
	defer time.Sleep(100 * time.Millisecond)
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
			session.Engine.OnClientDisconnect(session, session.GetClientId())
			return
		}
		session.Engine.OnConnect(session)
		return
	case packet.PACKET_TYPE_DISCONNECT:
		slog.Debug("//!! EVENT client disconnect", "client-id", session.GetClientId())
		session.Engine.OnClientDisconnect(session, session.GetClientId())
		return
	case packet.PACKET_TYPE_PINGREQ:
		slog.Debug("//!! EVENT client ping", "client-id", session.GetClientId())
		session.Engine.OnPing(session)
	case packet.PACKET_TYPE_SUBSCRIBE:
		slog.Debug("//!! EVENT client subscribed", "client-id", session.GetClientId())
		session.Engine.OnSubscribe(session, p)
	case packet.PACKET_TYPE_UNSUBSCRIBE:
		slog.Debug("//!! EVENT client unsubscribed", "client-id", session.GetClientId())
		session.Engine.OnUnsubscribe(session, p)
	case packet.PACKET_TYPE_PUBLISH:
		slog.Debug("//!! EVENT client published", "topic", p.GetPublishTopic(), "client-id", session.GetClientId(), "qos", p.QoS())
		session.Engine.OnPublish(session, p)
	case packet.PACKET_TYPE_PUBACK:
		slog.Debug("//!! EVENT client acked message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPuback(session, p)
	case packet.PACKET_TYPE_PUBREC:
		slog.Debug("//!! EVENT pub received message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPubrec(session, p)
	case packet.PACKET_TYPE_PUBREL:
		slog.Debug("//!! EVENT pub releases message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPubrel(session, p)
	case packet.PACKET_TYPE_PUBCOMP:
		slog.Debug("//!! EVENT pub complete message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPubcomp(session.GetClientId(), p.PacketIdentifier(), p.GetReasonCode())
	}
}
