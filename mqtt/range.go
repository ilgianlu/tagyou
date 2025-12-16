package mqtt

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

func rangePackets(session *model.RunningSession, packets <-chan *packet.Packet) {
	for p := range packets {
		managePacket(session, p)
	}
}

func managePacket(session *model.RunningSession, p model.Packet) {
	slog.Debug("[MQTT] packet arriving", "packet-type", p.PacketType())
	if !session.GetConnected() && p.PacketType() != packet.PACKET_TYPE_CONNECT {
		slog.Warn("[MQTT] session is disconnected,I can accept only connect, closing...", "packet-type", p.PacketType())
		err := session.Conn.Close()
		if err != nil {
			slog.Warn("[MQTT] could not clean close connection", "err", err)
		}
		return
	}
	switch p.PacketType() {
	case packet.PACKET_TYPE_CONNECT:
		slog.Debug("[MQTT] client connect", "client-id", session.GetClientId())
		if session.GetConnected() {
			slog.Debug("[MQTT] double connect event, disconnecting...", "client-id", session.GetClientId())
			session.Engine.OnClientDisconnect(session, session.GetClientId())
			return
		}
		session.Engine.OnConnect(session)
		return
	case packet.PACKET_TYPE_DISCONNECT:
		slog.Debug("[MQTT] client disconnect", "client-id", session.GetClientId())
		session.Engine.OnClientDisconnect(session, session.GetClientId())
		return
	case packet.PACKET_TYPE_PINGREQ:
		slog.Debug("[MQTT] client ping", "client-id", session.GetClientId())
		session.Engine.OnPing(session)
	case packet.PACKET_TYPE_SUBSCRIBE:
		slog.Debug("[MQTT] client subscribed", "client-id", session.GetClientId())
		session.Engine.OnSubscribe(session, p)
	case packet.PACKET_TYPE_UNSUBSCRIBE:
		slog.Debug("[MQTT] client unsubscribed", "client-id", session.GetClientId())
		session.Engine.OnUnsubscribe(session, p)
	case packet.PACKET_TYPE_PUBLISH:
		slog.Debug("[MQTT] client published", "topic", p.GetPublishTopic(), "client-id", session.GetClientId(), "qos", p.QoS())
		session.Engine.OnPublish(session, p)
	case packet.PACKET_TYPE_PUBACK:
		slog.Debug("[MQTT] client acked message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPuback(session, p)
	case packet.PACKET_TYPE_PUBREC:
		slog.Debug("[MQTT] pub received message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPubrec(session, p)
	case packet.PACKET_TYPE_PUBREL:
		slog.Debug("[MQTT] pub releases message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPubrel(session, p)
	case packet.PACKET_TYPE_PUBCOMP:
		slog.Debug("[MQTT] pub complete message", "packet-identifier", p.PacketIdentifier(), "client-id", session.GetClientId())
		session.Engine.OnClientPubcomp(session.GetClientId(), p.PacketIdentifier(), p.GetReasonCode())
	}
}
