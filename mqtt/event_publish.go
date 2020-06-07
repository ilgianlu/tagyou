package mqtt

import (
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func clientPublish(db *gorm.DB, e Event, outQueue chan<- OutData) {
	if e.packet.Retain() {
		saveRetain(db, e)
	}
	sendForward(db, e.topic, e.packet, outQueue)
	if e.packet.QoS() > 0 {
		sendAck(db, e, outQueue)
	}
}

func sendAck(db *gorm.DB, e Event, outQueue chan<- OutData) {
	if e.packet.QoS() == 1 {
		sendSimple(e.clientId, Puback(e.packet.PacketIdentifier(), PUBACK_SUCCESS, e.session.ProtocolVersion), outQueue)
	} else if e.packet.QoS() == 2 {
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
		sendSimple(e.clientId, Pubrec(e.packet.PacketIdentifier(), PUBREC_SUCCESS, e.session.ProtocolVersion), outQueue)
	}
}
