package event

import (
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sender"
	"github.com/rs/zerolog/log"
)

func OnPublish(sender sender.Sender, p *packet.Packet) {
	if conf.ACL_ON && !p.Session.FromLocalhost() && !CheckAcl(p.Topic, p.Session.PublishAcl) {
		if p.QoS() == 1 {
			sendAck(sender, p, packet.PUBACK_NOT_AUTHORIZED)
		} else if p.QoS() == 2 {
			sendPubrec(sender, p, packet.PUBREC_NOT_AUTHORIZED)
		}
		return
	}

	if p.Retain() {
		log.Debug().Msgf("[PUBLISH] to retain")
		saveRetain(p)
	}
	sender.Forward(p.Topic, p)
	if p.QoS() == 1 {
		log.Debug().Msgf("[PUBLISH] QoS 1 return ACK %d", p.PacketIdentifier())
		sendAck(sender, p, packet.PUBACK_SUCCESS)
	} else if p.QoS() == 2 {
		log.Debug().Msgf("[PUBLISH] QoS 2 return PUBREC")
		sendPubrec(sender, p, packet.PUBREC_SUCCESS)
	}
}

func sendAck(sender sender.Sender, p *packet.Packet, reasonCode uint8) {
	puback := packet.Puback(p.PacketIdentifier(), reasonCode, p.Session.ProtocolVersion)
	sender.Send(p.Session.ClientId, puback.ToByteSlice())
}

func sendPubrec(sender sender.Sender, p *packet.Packet, reasonCode uint8) {
	p.Session.Mu.RLock()
	clientId := p.Session.ClientId
	protocolVersion := p.Session.ProtocolVersion
	p.Session.Mu.RUnlock()
	r := model.Retry{
		ClientId:           clientId,
		PacketIdentifier:   p.PacketIdentifier(),
		Qos:                p.QoS(),
		Dup:                p.Dup(),
		ApplicationMessage: p.ApplicationMessage(),
		AckStatus:          model.WAIT_FOR_PUB_REL,
		CreatedAt:          time.Now().Unix(),
	}
	persistence.RetryRepository.SaveOne(r)

	pubrec := packet.Pubrec(p.PacketIdentifier(), reasonCode, protocolVersion)
	sender.Send(clientId, pubrec.ToByteSlice())
}
