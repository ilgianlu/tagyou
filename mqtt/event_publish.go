package mqtt

import (
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/jinzhu/gorm"
)

func onPublish(db *gorm.DB, p *packet.Packet, outQueue chan<- *OutData) {
	if (conf.ACL_ON || !p.Session.FromLocalhost()) && !CheckAcl(p.Topic, p.Session.PublishAcl) {
		if p.QoS() == 1 {
			sendAck(db, p, packet.PUBACK_NOT_AUTHORIZED, outQueue)
		} else {
			sendPubrec(db, p, packet.PUBREC_NOT_AUTHORIZED, outQueue)
		}
		return
	}

	if p.Retain() {
		saveRetain(db, p)
	}
	sendForward(db, p.Topic, p, outQueue)
	if p.QoS() == 1 {
		sendAck(db, p, packet.PUBACK_SUCCESS, outQueue)
	} else {
		sendPubrec(db, p, packet.PUBREC_SUCCESS, outQueue)
	}
}

func sendAck(db *gorm.DB, p *packet.Packet, reasonCode uint8, outQueue chan<- *OutData) {
	puback := packet.Puback(p.PacketIdentifier(), reasonCode, p.Session.ProtocolVersion)
	sendSimple(p.Session.ClientId, &puback, outQueue)
}

func sendPubrec(db *gorm.DB, p *packet.Packet, reasonCode uint8, outQueue chan<- *OutData) {
	r := model.Retry{
		ClientId:           p.Session.ClientId,
		PacketIdentifier:   p.PacketIdentifier(),
		Qos:                p.QoS(),
		Dup:                p.Dup(),
		ApplicationMessage: p.ApplicationMessage(),
		AckStatus:          model.WAIT_FOR_PUB_REL,
		CreatedAt:          time.Now(),
	}
	db.Save(&r)

	pubrec := packet.Pubrec(p.PacketIdentifier(), reasonCode, p.Session.ProtocolVersion)
	sendSimple(p.Session.ClientId, &pubrec, outQueue)
}
