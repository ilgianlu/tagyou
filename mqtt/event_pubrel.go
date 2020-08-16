package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func clientPubrel(db *gorm.DB, p Packet, outQueue chan<- OutData) {
	sendPubcomp := func(retry model.Retry) {
		var o OutData
		o.clientId = p.session.ClientId
		o.packet = Pubcomp(p.PacketIdentifier(), retry.ReasonCode, p.session.ProtocolVersion)
		outQueue <- o
	}

	onExpectedPubrel := func(retry model.Retry) {
		sendPubcomp(retry)
		db.Delete(&retry)
	}

	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_REL {
			onExpectedPubrel(retry)
		} else {
			log.Println("pubrel for invalid retry status", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}

	retry := model.Retry{
		ClientId:         p.session.ClientId,
		PacketIdentifier: p.PacketIdentifier(),
		ReasonCode:       p.reasonCode,
	}
	if db.Find(&retry).RecordNotFound() {
		log.Println("pubrel for invalid retry", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
