package event

import (
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func onSubscribe(router routers.Router, p *packet.Packet) {
	reasonCodes := []uint8{}
	for _, subscription := range p.Subscriptions {
		rCode := clientSubscription(router, p.Session, subscription)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(router, p, reasonCodes)
}

func clientSubscribed(router routers.Router, p *packet.Packet, reasonCodes []uint8) {
	p.Session.Mu.RLock()
	toSend := packet.Suback(p.PacketIdentifier(), reasonCodes, p.Session.ProtocolVersion)
	router.Send(p.Session.ClientId, toSend.ToByteSlice())
	p.Session.Mu.RUnlock()
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
		sendRetain(router, protocolVersion, subscription)
	}
	return 0
}

func sendRetain(router routers.Router, protocolVersion uint8, subscription model.Subscription) {
	retains := persistence.RetainRepository.FindRetains(subscription.Topic)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := packet.Publish(protocolVersion, subscription.Qos, true, r.Topic, packet.NewPacketIdentifier(), r.ApplicationMessage)
		router.Send(subscription.ClientId, p.ToByteSlice())
	}
}
