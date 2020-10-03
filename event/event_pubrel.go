package event

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func clientPubrel(db *gorm.DB, p *packet.Packet, outQueue chan<- *out.OutData) {
	sendPubcomp := func(retry model.Retry) {
		var o out.OutData
		o.ClientId = p.Session.ClientId
		o.Packet = packet.Pubcomp(p.PacketIdentifier(), retry.ReasonCode, p.Session.ProtocolVersion)
		outQueue <- &o
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
		ClientId:         p.Session.ClientId,
		PacketIdentifier: p.PacketIdentifier(),
		ReasonCode:       p.ReasonCode,
	}
	if err := db.Find(&retry).Error; err != nil {
		log.Println("pubrel for invalid retry", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
