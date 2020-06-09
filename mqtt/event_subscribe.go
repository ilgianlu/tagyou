package mqtt

import (
	"strings"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onSubscribe(db *gorm.DB, e Event, outQueue chan<- OutData) {
	reasonCodes := []uint8{}
	for _, subscription := range e.subscriptions {
		rCode := clientSubscription(db, &subscription, outQueue)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientSubscribed(e, reasonCodes, outQueue)
}

func clientSubscribed(e Event, reasonCodes []uint8, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.clientId
	o.packet = Suback(e.packet.PacketIdentifier(), reasonCodes, e.session.ProtocolVersion)
	outQueue <- o
}

func clientSubscription(db *gorm.DB, subscription *model.Subscription, outQueue chan<- OutData) uint8 {
	// check subscr qos, topic valid...
	db.Create(subscription)
	sendRetain(db, subscription, outQueue)
	return 0
}

func sendRetain(db *gorm.DB, subscription *model.Subscription, outQueue chan<- OutData) {
	retains := findRetains(db, subscription.Topic)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := Publish(subscription.QoS, true, r.Topic, newPacketIdentifier(), r.ApplicationMessage)
		sendForward(db, r.Topic, p, outQueue)
	}
}

func findRetains(db *gorm.DB, subscribedTopic string) []model.Retain {
	trimmedTopic := trimWildcard(subscribedTopic)
	var retains []model.Retain
	db.Where("topic LIKE ?", strings.Join([]string{trimmedTopic, "%"}, "")).Find(&retains)
	return retains
}
