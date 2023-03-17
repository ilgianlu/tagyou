package event

import (
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sender"
)

func onSubscribe(sender sender.Sender, p *packet.Packet) {
	reasonCodes := []uint8{}
	for _, subscription := range p.Subscriptions {
		rCode := clientSubscription(sender, p.Session, subscription)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(sender, p, reasonCodes)
}

func clientSubscribed(sender sender.Sender, p *packet.Packet, reasonCodes []uint8) {
	p.Session.Mu.RLock()
	toSend := packet.Suback(p.PacketIdentifier(), reasonCodes, p.Session.ProtocolVersion)
	sender.Send(p.Session.ClientId, toSend.ToByteSlice())
	p.Session.Mu.RUnlock()
}

func clientSubscription(sender sender.Sender, session *model.RunningSession, subscription model.Subscription) uint8 {
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
		sendRetain(sender, protocolVersion, subscription)
	}
	return 0
}

func sendRetain(sender sender.Sender, protocolVersion uint8, subscription model.Subscription) {
	retains := persistence.RetainRepository.FindRetains(subscription.Topic)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := packet.Publish(protocolVersion, subscription.Qos, true, r.Topic, packet.NewPacketIdentifier(), r.ApplicationMessage)
		sender.Send(subscription.ClientId, p.ToByteSlice())
	}
}
