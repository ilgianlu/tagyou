package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func clientPubrec(db *gorm.DB, e Event, outQueue chan<- OutData) {
	sendPubrel := func(retry model.Retry) {
		var o OutData
		o.clientId = retry.ClientId
		o.packet = Pubrel(retry.PacketIdentifier, retry.ReasonCode, e.session.ProtocolVersion)
		outQueue <- o
	}

	onExpectedPubrec := func(retry model.Retry) {
		sendPubrel(retry)
		// change retry state to wait for pubcomp
		retry.AckStatus = model.WAIT_FOR_PUB_COMP
		db.Save(&retry)
	}

	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_REC {
			onExpectedPubrec(retry)
		} else {
			log.Println("pubrec for invalid retry status", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}

	retry := model.Retry{
		ClientId:         e.clientId,
		PacketIdentifier: e.packet.packetIdentifier,
		ReasonCode:       e.packet.reasonCode,
	}
	if db.Find(&retry).RecordNotFound() {
		log.Println("pubrec for invalid retry", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
