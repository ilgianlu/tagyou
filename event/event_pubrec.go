package event

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func clientPubrec(db *gorm.DB, p *packet.Packet, outQueue chan<- *out.OutData) {
	sendPubrel := func(retry model.Retry) {
		var o out.OutData
		o.ClientId = retry.ClientId
		o.Packet = packet.Pubrel(retry.PacketIdentifier, retry.ReasonCode, p.Session.ProtocolVersion)
		outQueue <- &o
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
		ClientId:         p.Session.ClientId,
		PacketIdentifier: p.PacketIdentifier(),
		ReasonCode:       p.ReasonCode,
	}
	if err := db.First(&retry).Error; err != nil {
		log.Println("pubrec for invalid retry", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
