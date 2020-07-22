package mqtt

import (
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onPublish(db *gorm.DB, e Event, outQueue chan<- OutData) {
	if conf.ACL_ON && !CheckAcl(e.topic, e.session.PublishAcl) {
		if e.packet.QoS() == 1 {
			sendAck(db, e, PUBACK_NOT_AUTHORIZED, outQueue)
		} else {
			sendPubrec(db, e, PUBREC_NOT_AUTHORIZED, outQueue)
		}
		return
	}

	if e.packet.Retain() {
		saveRetain(db, e)
	}
	sendForward(db, e.session.ProtocolVersion, e.topic, e.packet, outQueue)
	if e.packet.QoS() == 1 {
		sendAck(db, e, PUBACK_SUCCESS, outQueue)
	} else {
		sendPubrec(db, e, PUBREC_SUCCESS, outQueue)
	}
}

func sendAck(db *gorm.DB, e Event, reasonCode uint8, outQueue chan<- OutData) {
	sendSimple(e.clientId, Puback(e.packet.PacketIdentifier(), reasonCode, e.session.ProtocolVersion), outQueue)
}

func sendPubrec(db *gorm.DB, e Event, reasonCode uint8, outQueue chan<- OutData) {
	r := model.Retry{
		ClientId:           e.clientId,
		PacketIdentifier:   e.packet.PacketIdentifier(),
		Qos:                e.packet.QoS(),
		Dup:                e.packet.Dup(),
		ApplicationMessage: e.packet.ApplicationMessage(),
		AckStatus:          model.WAIT_FOR_PUB_REL,
		CreatedAt:          time.Now(),
	}
	db.Save(&r)
	sendSimple(e.clientId, Pubrec(e.packet.PacketIdentifier(), reasonCode, e.session.ProtocolVersion), outQueue)
}
