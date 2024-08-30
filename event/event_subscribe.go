package event

import (
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func onSubscribe(router routers.Router, session *model.RunningSession, p *packet.Packet) {
	reasonCodes := []uint8{}
	for _, subscription := range p.Subscriptions {
		rCode := clientSubscription(router, session, subscription)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(router, session, p.PacketIdentifier(), reasonCodes)
}

func clientSubscribed(router routers.Router, session *model.RunningSession, packetIdentifier int, reasonCodes []uint8) {
	session.Mu.RLock()
	toSend := packet.Suback(packetIdentifier, reasonCodes, session.ProtocolVersion)
	router.Send(session.ClientId, toSend.ToByteSlice())
	session.Mu.RUnlock()
}

func clientSubscription(router routers.Router, session *model.RunningSession, subscription model.Subscription) uint8 {
	session.Mu.RLock()
	fromLocalhost := session.FromLocalhost()
	subscribeAcl := session.SubscribeAcl
	protocolVersion := session.ProtocolVersion
	session.Mu.RUnlock()
	// check subscr qos, topic valid...
	if conf.ACL_ON && !fromLocalhost && !CheckAcl(subscription.Topic, subscribeAcl) {
		return conf.SUB_TOPIC_FILTER_INVALID
	}
	// db.Create(&subscription)
	persistence.SubscriptionRepository.CreateOne(subscription)
	if !subscription.Shared {
		router.SendRetain(protocolVersion, subscription)
	}
	return 0
}
