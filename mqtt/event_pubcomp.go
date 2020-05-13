package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func clientPubcomp(db *gorm.DB, e Event) {
	onRetryFound := func(db *gorm.DB, retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_COMP {
			db.Delete(&retry)
		} else {
			log.Println("pubcomp for invalid retry status", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}

	retry := model.Retry{
		ClientId:         e.clientId,
		PacketIdentifier: e.packet.packetIdentifier,
		ReasonCode:       e.packet.reasonCode,
	}
	if db.Find(&retry).RecordNotFound() {
		log.Println("pubcomp for invalid retry", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(db, retry)
	}

}
