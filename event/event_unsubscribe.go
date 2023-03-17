package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func onUnsubscribe(router routers.Router, p *packet.Packet) {
	reasonCodes := []uint8{}
	for _, unsub := range p.Subscriptions {
		rCode := clientUnsubscription(p.Session.GetClientId(), unsub)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientUnsubscribed(router, p, reasonCodes)
}

func clientUnsubscribed(router routers.Router, p *packet.Packet, reasonCodes []uint8) {
	toSend := packet.Unsuback(p.PacketIdentifier(), reasonCodes, p.Session.GetProtocolVersion())
	router.Send(p.Session.GetClientId(), toSend.ToByteSlice())
}

func clientUnsubscription(clientId string, unsub model.Subscription) uint8 {
	if sub, err := persistence.SubscriptionRepository.FindToUnsubscribe(unsub.ShareName, unsub.Topic, clientId); err != nil {
		log.Info().Msgf("no subscription to unsubscribe %s %s", unsub.Topic, clientId)
		log.Error().Err(err).Msg("error unsubscribing")
		return conf.UNSUB_NO_SUB_EXISTED
	} else {
		persistence.SubscriptionRepository.DeleteByClientIdTopicShareName(clientId, sub.Topic, sub.ShareName)
		return conf.SUCCESS
	}
}
