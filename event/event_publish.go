package event

import (
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func OnPublish(router routers.Router, session *model.RunningSession, p *packet.Packet) {
	if conf.ACL_ON && !session.FromLocalhost() && !CheckAcl(p.Topic, session.PublishAcl) {
		if p.QoS() == conf.QOS1 {
			sendAck(router, session, p.PacketIdentifier(), packet.PUBACK_NOT_AUTHORIZED)
		} else if p.QoS() == conf.QOS2 {
			sendPubrec(router, session, p, packet.PUBREC_NOT_AUTHORIZED)
		}
		return
	}

	if p.Retain() {
		slog.Debug("[PUBLISH] to retain")
		saveRetain(session, p)
	}
	router.Forward(session.GetClientId(), p.Topic, p)
	if p.QoS() == conf.QOS1 {
		slog.Debug("[PUBLISH] QoS 1 return ACK", "packet-identifier", p.PacketIdentifier())
		sendAck(router, session, p.PacketIdentifier(), packet.PUBACK_SUCCESS)
	} else if p.QoS() == conf.QOS2 {
		slog.Debug("[PUBLISH] QoS 2 return PUBREC")
		sendPubrec(router, session, p, packet.PUBREC_SUCCESS)
	} else if p.QoS() == conf.QOS0 {
		slog.Debug("[PUBLISH] QoS 0 no return")
	}
}

func sendAck(router routers.Router, session *model.RunningSession, packetIdentifier int, reasonCode uint8) {
	puback := packet.Puback(packetIdentifier, reasonCode, session.ProtocolVersion)
	router.Send(session.ClientId, puback.ToByteSlice())
}

func sendPubrec(router routers.Router, session *model.RunningSession, p *packet.Packet, reasonCode uint8) {
	session.Mu.RLock()
	clientId := session.ClientId
	protocolVersion := session.ProtocolVersion
	session.Mu.RUnlock()
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
	router.Send(clientId, pubrec.ToByteSlice())
}
