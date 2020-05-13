package mqtt

import (
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func clientPublish(db *gorm.DB, e Event, outQueue chan<- OutData) {
	if e.published.retain {
		saveRetain(db, e)
	}
	sendForward(db, e.published.topic, e.packet, outQueue)
	if e.published.qos == 1 {
		r := model.Retry{
			ClientId:           e.clientId,
			PacketIdentifier:   e.packet.packetIdentifier,
			Qos:                e.packet.QoS(),
			Dup:                e.packet.Dup(),
			ApplicationMessage: e.packet.ApplicationMessage(),
			AckStatus:          model.WAIT_FOR_PUB_ACK,
			CreatedAt:          time.Now(),
		}
		db.Save(&r)
		sendSimple(e.clientId, Puback(r.PacketIdentifier, PUBACK_SUCCESS), outQueue)
	} else if e.published.qos == 2 {
		r := model.Retry{
			ClientId:           e.clientId,
			PacketIdentifier:   e.packet.packetIdentifier,
			Qos:                e.packet.QoS(),
			Dup:                e.packet.Dup(),
			ApplicationMessage: e.packet.ApplicationMessage(),
			AckStatus:          model.WAIT_FOR_PUB_REL,
			CreatedAt:          time.Now(),
		}
		db.Save(&r)
		sendSimple(e.clientId, Pubrec(e.packet.packetIdentifier, PUBREC_SUCCESS), outQueue)
	}
}
