package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func onUnsubscribe(p *packet.Packet, outQueue chan<- out.OutData) {
	reasonCodes := []uint8{}
	for _, unsub := range p.Subscriptions {
		rCode := clientUnsubscription(p.Session.GetClientId(), unsub)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientUnsubscribed(p, reasonCodes, outQueue)
}

func clientUnsubscribed(p *packet.Packet, reasonCodes []uint8, outQueue chan<- out.OutData) {
	var o out.OutData
	o.ClientId = p.Session.GetClientId()
	toSend := packet.Unsuback(p.PacketIdentifier(), reasonCodes, p.Session.GetProtocolVersion())
	o.Packet = toSend.ToByteSlice()
	outQueue <- o
}

func clientUnsubscription(clientId string, unsub model.Subscription) uint8 {
	if sub, err := persistence.SubscriptionRepository.FindToUnsubscribe(unsub.ShareName, unsub.Topic, clientId); err != nil {
		log.Info().Msgf("no subscription to unsubscribe %s %s", unsub.Topic, clientId)
		log.Error().Err(err).Msg("error unsubscribing")
		return conf.UNSUB_NO_SUB_EXISTED
	} else {
		persistence.SubscriptionRepository.DeleteOne(sub)
		return conf.SUCCESS
	}
}
