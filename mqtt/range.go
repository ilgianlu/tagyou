package mqtt

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

func rangePackets(session *model.RunningSession, packets <-chan *packet.Packet) {
	for p := range packets {
		managePacket(session, p)
	}
}

func managePacket(session *model.RunningSession, p model.Packet) {
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
			event.OnClientDisconnect(session, session.GetClientId())
			return
		}
		event.OnConnect(session)
		return
	case packet.PACKET_TYPE_DISCONNECT:
		slog.Debug("//!! EVENT client disconnect", "client-id", session.GetClientId())
		event.OnClientDisconnect(session, session.GetClientId())
		return
	case packet.PACKET_TYPE_PINGREQ:
		slog.Debug("//!! EVENT client ping", "client-id", session.GetClientId())
		event.OnPing(session)
	case packet.PACKET_TYPE_SUBSCRIBE:
		slog.Debug("//!! EVENT client subscribed", "client-id", session.GetClientId())
		event.OnSubscribe(session, p)
	case packet.PACKET_TYPE_UNSUBSCRIBE:
		slog.Debug("//!! EVENT client unsubscribed", "client-id", session.GetClientId())
		event.OnUnsubscribe(session, p)
	case packet.PACKET_TYPE_PUBLISH:
		slog.Debug("//!! EVENT client published", "topic", p.GetPublishTopic(), "client-id", session.GetClientId(), "qos", p.QoS())
		event.OnPublish(session, p)
	case packet.PACKET_TYPE_PUBACK:
		slog.Debug("//!! EVENT client acked message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		event.OnClientPuback(session, p)
	case packet.PACKET_TYPE_PUBREC:
		slog.Debug("//!! EVENT pub received message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		event.OnClientPubrec(session, p)
	case packet.PACKET_TYPE_PUBREL:
		slog.Debug("//!! EVENT pub releases message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		event.OnClientPubrel(session, p)
	case packet.PACKET_TYPE_PUBCOMP:
		slog.Debug("//!! EVENT pub complete message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		event.OnClientPubcomp(session.GetClientId(), p.PacketIdentifier(), p.GetReasonCode())
	}
}
