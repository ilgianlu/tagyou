package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func clientPuback(db *gorm.DB, p *packet.Packet) {
	onRetryFound := func(db *gorm.DB, retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_ACK {
			db.Delete(&retry)
		} else {
			log.Info().Msgf("puback for invalid retry status %s %s %s", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}

	retry := model.Retry{
		ClientId:         p.Session.ClientId,
		PacketIdentifier: p.PacketIdentifier(),
	}
	if err := db.Find(&retry).Error; err != nil {
		log.Info().Msgf("puback for invalid retry %s %s", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(db, retry)
	}
}
