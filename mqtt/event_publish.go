package mqtt

import (
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onPublish(db *gorm.DB, p Packet, outQueue chan<- OutData) {
	if (conf.ACL_ON || !p.session.FromLocalhost()) && !CheckAcl(p.topic, p.session.PublishAcl) {
		if p.QoS() == 1 {
			sendAck(db, p, PUBACK_NOT_AUTHORIZED, outQueue)
		} else {
			sendPubrec(db, p, PUBREC_NOT_AUTHORIZED, outQueue)
		}
		return
	}

	if p.Retain() {
		saveRetain(db, p)
	}
	sendForward(db, p.topic, p, outQueue)
	if p.QoS() == 1 {
		sendAck(db, p, PUBACK_SUCCESS, outQueue)
	} else {
		sendPubrec(db, p, PUBREC_SUCCESS, outQueue)
	}
}

func sendAck(db *gorm.DB, p Packet, reasonCode uint8, outQueue chan<- OutData) {
	sendSimple(p.session.ClientId, Puback(p.PacketIdentifier(), reasonCode, p.session.ProtocolVersion), outQueue)
}

func sendPubrec(db *gorm.DB, p Packet, reasonCode uint8, outQueue chan<- OutData) {
	r := model.Retry{
		ClientId:           p.session.ClientId,
		PacketIdentifier:   p.PacketIdentifier(),
		Qos:                p.QoS(),
		Dup:                p.Dup(),
		ApplicationMessage: p.ApplicationMessage(),
		AckStatus:          model.WAIT_FOR_PUB_REL,
		CreatedAt:          time.Now(),
	}
	db.Save(&r)
	sendSimple(p.session.ClientId, Pubrec(p.PacketIdentifier(), reasonCode, p.session.ProtocolVersion), outQueue)
}
