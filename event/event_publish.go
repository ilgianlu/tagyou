package event

import (
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/nowherecloud"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/rs/zerolog/log"
	kgo "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func onPublish(db *gorm.DB, kwriter *kgo.Writer, p *packet.Packet, outQueue chan<- *out.OutData) {
	if conf.ACL_ON && !p.Session.FromLocalhost() && !CheckAcl(p.Topic, p.Session.PublishAcl) {
		if p.QoS() == 1 {
			sendAck(db, p, packet.PUBACK_NOT_AUTHORIZED, outQueue)
		} else if p.QoS() == 2 {
			sendPubrec(db, p, packet.PUBREC_NOT_AUTHORIZED, outQueue)
		}
		return
	}

	if p.Retain() {
		log.Debug().Msgf("[PUBLISH] to retain")
		saveRetain(db, p)
	}
	sendForward(db, p.Topic, p, outQueue)
	if nowherecloud.KAFKA_ON {
		nowherecloud.Publish(kwriter, p.Topic, p)
	}
	if p.QoS() == 1 {
		log.Debug().Msgf("[PUBLISH] QoS 1 return ACK %d", p.PacketIdentifier())
		sendAck(db, p, packet.PUBACK_SUCCESS, outQueue)
	} else if p.QoS() == 2 {
		log.Debug().Msgf("[PUBLISH] QoS 2 return PUBREC")
		sendPubrec(db, p, packet.PUBREC_SUCCESS, outQueue)
	}
}

func sendAck(db *gorm.DB, p *packet.Packet, reasonCode uint8, outQueue chan<- *out.OutData) {
	puback := packet.Puback(p.PacketIdentifier(), reasonCode, p.Session.ProtocolVersion)
	sendSimple(p.Session.ClientId, &puback, outQueue)
}

func sendPubrec(db *gorm.DB, p *packet.Packet, reasonCode uint8, outQueue chan<- *out.OutData) {
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
	db.Save(&r)

	pubrec := packet.Pubrec(p.PacketIdentifier(), reasonCode, protocolVersion)
	sendSimple(clientId, &pubrec, outQueue)
}
