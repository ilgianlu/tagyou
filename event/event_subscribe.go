package event

import (
	"fmt"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/rs/zerolog/log"
)

func onSubscribe(p *packet.Packet, outQueue chan<- out.OutData) {
	reasonCodes := []uint8{}
	for _, subscription := range p.Subscriptions {
		rCode := clientSubscription(p.Session, subscription, outQueue)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(p, reasonCodes, outQueue)
}

func clientSubscribed(p *packet.Packet, reasonCodes []uint8, outQueue chan<- out.OutData) {
	p.Session.Mu.RLock()
	clientId := p.Session.ClientId
	protocolVersion := p.Session.ProtocolVersion
	p.Session.Mu.RUnlock()
	var o out.OutData
	o.ClientId = clientId
	toSend := packet.Suback(p.PacketIdentifier(), reasonCodes, protocolVersion)
	o.Packet = toSend.ToByteSlice()
	outQueue <- o
}

func clientSubscription(session *model.RunningSession, subscription model.Subscription, outQueue chan<- out.OutData) uint8 {
	session.Mu.RLock()
	fromLocalhost := session.FromLocalhost()
	subscribeAcl := session.SubscribeAcl
	protocolVersion := session.ProtocolVersion
	session.Mu.RUnlock()
	// check subscr qos, topic valid...
	if conf.ACL_ON && !fromLocalhost && !CheckAcl(subscription.Topic, subscribeAcl) {
		return conf.SUB_TOPIC_FILTER_INVALID
	}
	log.Debug().Msg(fmt.Sprintf("%s client is subscribing %s", subscription.ClientId, subscription.Topic))
	persistence.SubscriptionRepository.Create(subscription)
	if !subscription.Shared {
		sendRetain(protocolVersion, subscription, outQueue)
	}
	return 0
}

func sendRetain(protocolVersion uint8, subscription model.Subscription, outQueue chan<- out.OutData) {
	retains := persistence.RetainRepository.FindRetains(subscription.Topic)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := packet.Publish(protocolVersion, subscription.Qos, true, r.Topic, packet.NewPacketIdentifier(), r.ApplicationMessage)
		sendForward(r.Topic, &p, outQueue)
	}
}
