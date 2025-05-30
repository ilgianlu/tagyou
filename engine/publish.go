package engine

import (
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func (s StandardEngine) OnPublish(session *model.RunningSession, p model.Packet) {
	if conf.ACL_ON && !session.FromLocalhost() && !CheckAcl(p.GetPublishTopic(), session.PublishAcl) {
		if p.QoS() == conf.QOS1 {
			sendAck(session, p.PacketIdentifier(), packet.PUBACK_NOT_AUTHORIZED)
		} else if p.QoS() == conf.QOS2 {
			sendPubrec(session, p, packet.PUBREC_NOT_AUTHORIZED)
		}
		return
	}

	if p.Retain() {
		slog.Debug("[PUBLISH] to retain")
		saveRetain(session.GetClientId(), p.GetPublishTopic(), p.ApplicationMessage())
	}
	session.Router.Forward(session.GetClientId(), p.GetPublishTopic(), p)
	if p.QoS() == conf.QOS1 {
		slog.Debug("[PUBLISH] QoS 1 return ACK", "packet-identifier", p.PacketIdentifier())
		sendAck(session, p.PacketIdentifier(), packet.PUBACK_SUCCESS)
	} else if p.QoS() == conf.QOS2 {
		slog.Debug("[PUBLISH] QoS 2 return PUBREC")
		sendPubrec(session, p, packet.PUBREC_SUCCESS)
	} else if p.QoS() == conf.QOS0 {
		slog.Debug("[PUBLISH] QoS 0 no return")
	}
}

func sendAck(session *model.RunningSession, packetIdentifier int, reasonCode uint8) {
	puback := packet.Puback(packetIdentifier, reasonCode, session.ProtocolVersion)
	session.Router.Send(session.ClientId, puback.ToByteSlice())
}

func sendPubrec(session *model.RunningSession, p model.Packet, reasonCode uint8) {
	clientId := session.ClientId
	protocolVersion := session.ProtocolVersion
	r := model.Retry{
		ClientId:           clientId,
		PacketIdentifier:   p.PacketIdentifier(),
		Qos:                p.QoS(),
		Dup:                p.Dup(),
		ApplicationMessage: p.ApplicationMessage(),
		AckStatus:          model.WAIT_FOR_PUB_REL,
		CreatedAt:          time.Now().Unix(),
	}
	persistence.RetryRepository.InsertOne(r)

	pubrec := packet.Pubrec(p.PacketIdentifier(), reasonCode, protocolVersion)
	session.Router.Send(clientId, pubrec.ToByteSlice())
}
