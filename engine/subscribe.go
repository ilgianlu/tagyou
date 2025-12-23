package engine

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func (s StandardEngine) OnSubscribe(session *model.RunningSession, p model.Packet) {
	reasonCodes := []uint8{}
	for _, subscription := range p.GetSubscriptions() {
		rCode := clientSubscription(session, subscription)
		slog.Debug("[ENGINE] new subscription", "result", rCode)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(session, p.PacketIdentifier(), reasonCodes)
}

func clientSubscribed(session *model.RunningSession, packetIdentifier int, reasonCodes []uint8) {
	toSend := packet.Suback(packetIdentifier, reasonCodes, session.ProtocolVersion)
	bs, err := toSend.ToByteSlice()
	if err != nil {
		slog.Warn("[ENGINE] error encoding suback", "err", err)
		return
	}
	session.Router.Send(session.ClientId, bs)
}

func clientSubscription(session *model.RunningSession, subscription model.Subscription) uint8 {
	fromLocalhost := session.FromLocalhost()
	subscribeAcl := session.SubscribeAcl
	protocolVersion := session.ProtocolVersion
	// check subscr qos, topic valid...
	if conf.ACL_ON && !fromLocalhost && !CheckAcl(subscription.Topic, subscribeAcl) {
		return conf.SUB_TOPIC_FILTER_INVALID
	}
	// db.Create(&subscription)
	persistence.SubscriptionRepository.CreateOne(subscription)
	if !subscription.Shared {
		session.Router.SendRetain(protocolVersion, subscription)
	}
	return 0
}
