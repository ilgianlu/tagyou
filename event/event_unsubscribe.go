package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func onUnsubscribe(db *gorm.DB, p *packet.Packet, outQueue chan<- *out.OutData) {
	reasonCodes := []uint8{}
	for _, unsub := range p.Subscriptions {
		rCode := clientUnsubscription(db, p.Session.ClientId, unsub.Topic)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientUnsubscribed(p, reasonCodes, outQueue)
}

func clientUnsubscribed(p *packet.Packet, reasonCodes []uint8, outQueue chan<- *out.OutData) {
	var o out.OutData
	o.ClientId = p.Session.ClientId
	o.Packet = packet.Unsuback(p.PacketIdentifier(), reasonCodes, p.Session.ProtocolVersion)
	outQueue <- &o
}

func clientUnsubscription(db *gorm.DB, clientId string, topic string) uint8 {
	var sub model.Subscription
	if err := db.Where("topic = ? and client_id = ?", topic, clientId).First(&sub).Error; err != nil {
		log.Info().Msgf("no subscription to unsubscribe %s %s", topic, clientId)
		log.Error().Err(err).Msg("error unsubscribing")
		return conf.UNSUB_NO_SUB_EXISTED
	}
	db.Delete(&sub)
	return conf.SUCCESS
}
