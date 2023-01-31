package event

import (
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/rs/zerolog/log"
)

func onPublish(p *packet.Packet, outQueue chan<- out.OutData) {
	if conf.ACL_ON && !p.Session.FromLocalhost() && !CheckAcl(p.Topic, p.Session.PublishAcl) {
		if p.QoS() == 1 {
			sendAck(p, packet.PUBACK_NOT_AUTHORIZED, outQueue)
		} else if p.QoS() == 2 {
			sendPubrec(p, packet.PUBREC_NOT_AUTHORIZED, outQueue)
		}
		return
	}

	if p.Retain() {
		log.Debug().Msgf("[PUBLISH] to retain")
		saveRetain(p)
	}
	sendForward(p.Topic, p, outQueue)
	if p.QoS() == 1 {
		log.Debug().Msgf("[PUBLISH] QoS 1 return ACK %d", p.PacketIdentifier())
		sendAck(p, packet.PUBACK_SUCCESS, outQueue)
	} else if p.QoS() == 2 {
		log.Debug().Msgf("[PUBLISH] QoS 2 return PUBREC")
		sendPubrec(p, packet.PUBREC_SUCCESS, outQueue)
	}
}

func sendAck(p *packet.Packet, reasonCode uint8, outQueue chan<- out.OutData) {
	puback := packet.Puback(p.PacketIdentifier(), reasonCode, p.Session.ProtocolVersion)
	sendSimple(p.Session.ClientId, &puback, outQueue)
}

func sendPubrec(p *packet.Packet, reasonCode uint8, outQueue chan<- out.OutData) {
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
	sendSimple(clientId, &pubrec, outQueue)
}
