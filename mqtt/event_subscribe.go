package mqtt

import (
	"strings"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onSubscribe(db *gorm.DB, p Packet, outQueue chan<- OutData) {
	reasonCodes := []uint8{}
	for _, subscription := range p.subscriptions {
		rCode := clientSubscription(db, p.session, &subscription, outQueue)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(p, reasonCodes, outQueue)
}

func clientSubscribed(p Packet, reasonCodes []uint8, outQueue chan<- OutData) {
	var o OutData
	o.clientId = p.session.ClientId
	o.packet = Suback(p.PacketIdentifier(), reasonCodes, p.session.ProtocolVersion)
	outQueue <- o
}

func clientSubscription(db *gorm.DB, session *model.Session, subscription *model.Subscription, outQueue chan<- OutData) uint8 {
	// check subscr qos, topic valid...
	if (conf.ACL_ON || !session.FromLocalhost()) && !CheckAcl(subscription.Topic, session.SubscribeAcl) {
		return conf.SUB_TOPIC_FILTER_INVALID
	}
	db.Create(subscription)
	sendRetain(db, session.ProtocolVersion, subscription, outQueue)
	return 0
}

func sendRetain(db *gorm.DB, protocolVersion uint8, subscription *model.Subscription, outQueue chan<- OutData) {
	retains := findRetains(db, subscription.Topic)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := Publish(protocolVersion, subscription.Qos, true, r.Topic, newPacketIdentifier(), r.ApplicationMessage)
		sendForward(db, r.Topic, p, outQueue)
	}
}

func findRetains(db *gorm.DB, subscribedTopic string) []model.Retain {
	trimmedTopic := trimWildcard(subscribedTopic)
	var retains []model.Retain
	db.Where("topic LIKE ?", strings.Join([]string{trimmedTopic, "%"}, "")).Find(&retains)
	return retains
}
