package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func onUnsubscribe(router routers.Router, session *model.RunningSession, p *packet.Packet) {
	reasonCodes := []uint8{}
	for _, unsub := range p.Subscriptions {
		rCode := clientUnsubscription(session.GetClientId(), unsub)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientUnsubscribed(router, session, p.PacketIdentifier(), reasonCodes)
}

func clientUnsubscribed(router routers.Router, session *model.RunningSession, packetIdentifier int, reasonCodes []uint8) {
	toSend := packet.Unsuback(packetIdentifier, reasonCodes, session.GetProtocolVersion())
	router.Send(session.GetClientId(), toSend.ToByteSlice())
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
