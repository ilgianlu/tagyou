package event

import (
	"strings"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func onSubscribe(db *gorm.DB, p *packet.Packet, outQueue chan<- *out.OutData) {
	reasonCodes := []uint8{}
	for _, subscription := range p.Subscriptions {
		rCode := clientSubscription(db, p.Session, &subscription, outQueue)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(p, reasonCodes, outQueue)
}

func clientSubscribed(p *packet.Packet, reasonCodes []uint8, outQueue chan<- *out.OutData) {
	var o out.OutData
	o.ClientId = p.Session.ClientId
	o.Packet = packet.Suback(p.PacketIdentifier(), reasonCodes, p.Session.ProtocolVersion)
	outQueue <- &o
}

func clientSubscription(db *gorm.DB, session *model.RunningSession, subscription *model.Subscription, outQueue chan<- *out.OutData) uint8 {
	// check subscr qos, topic valid...
	if conf.ACL_ON && !session.FromLocalhost() && !CheckAcl(subscription.Topic, session.SubscribeAcl) {
		return conf.SUB_TOPIC_FILTER_INVALID
	}
	db.Create(subscription)
	if !subscription.Shared {
		sendRetain(db, session.ProtocolVersion, subscription, outQueue)
	}
	return 0
}

func sendRetain(db *gorm.DB, protocolVersion uint8, subscription *model.Subscription, outQueue chan<- *out.OutData) {
	retains := findRetains(db, subscription.Topic)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := packet.Publish(protocolVersion, subscription.Qos, true, r.Topic, packet.NewPacketIdentifier(), r.ApplicationMessage)
		sendForward(db, r.Topic, &p, outQueue)
	}
}

func findRetains(db *gorm.DB, subscribedTopic string) []model.Retain {
	trimmedTopic := trimWildcard(subscribedTopic)
	var retains []model.Retain
	db.Where("topic LIKE ?", strings.Join([]string{trimmedTopic, "%"}, "")).Find(&retains)
	return retains
}
