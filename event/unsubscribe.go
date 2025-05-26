package event

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func OnUnsubscribe(session *model.RunningSession, p *packet.Packet) {
	reasonCodes := []uint8{}
	for _, unsub := range p.Subscriptions {
		rCode := clientUnsubscription(session.GetClientId(), unsub)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientUnsubscribed(session, p.PacketIdentifier(), reasonCodes)
}

func clientUnsubscribed(session *model.RunningSession, packetIdentifier int, reasonCodes []uint8) {
	toSend := packet.Unsuback(packetIdentifier, reasonCodes, session.GetProtocolVersion())
	session.Router.Send(session.GetClientId(), toSend.ToByteSlice())
}

func clientUnsubscription(clientId string, unsub model.Subscription) uint8 {
	if sub, err := persistence.SubscriptionRepository.FindToUnsubscribe(unsub.ShareName, unsub.Topic, clientId); err != nil {
		slog.Info("no subscription to unsubscribe", "topic", unsub.Topic, "client-id", clientId)
		slog.Error("error unsubscribing", "err", err)
		return conf.UNSUB_NO_SUB_EXISTED
	} else {
		persistence.SubscriptionRepository.DeleteByClientIdTopicShareName(clientId, sub.Topic, sub.ShareName)
		return conf.SUCCESS
	}
}
