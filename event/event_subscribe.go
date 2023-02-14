package event

import (
	"fmt"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/rs/zerolog/log"
)

func onSubscribe(connections *model.Connections, p *packet.Packet) {
	reasonCodes := []uint8{}
	for _, subscription := range p.Subscriptions {
		rCode := clientSubscription(connections, p.Session, subscription)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(connections, p, reasonCodes)
}

func clientSubscribed(connections *model.Connections, p *packet.Packet, reasonCodes []uint8) {
	p.Session.Mu.RLock()
	toSend := packet.Suback(p.PacketIdentifier(), reasonCodes, p.Session.ProtocolVersion)
	SimpleSend(connections, p.Session.ClientId, toSend.ToByteSlice())
	p.Session.Mu.RUnlock()
}

func clientSubscription(connections *model.Connections, session *model.RunningSession, subscription model.Subscription) uint8 {
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
		sendRetain(connections, protocolVersion, subscription)
	}
	return 0
}

func sendRetain(connections *model.Connections, protocolVersion uint8, subscription model.Subscription) {
	retains := persistence.RetainRepository.FindRetains(subscription.Topic)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := packet.Publish(protocolVersion, subscription.Qos, true, r.Topic, packet.NewPacketIdentifier(), r.ApplicationMessage)
		SimpleSend(connections, subscription.ClientId, p.ToByteSlice())
	}
}
