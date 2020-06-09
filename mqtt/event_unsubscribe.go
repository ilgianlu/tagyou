package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onUnsubscribe(db *gorm.DB, e Event, outQueue chan<- OutData) {
	reasonCodes := []uint8{}
	for _, unsub := range e.subscriptions {
		rCode := clientUnsubscription(db, e.session.ClientId, unsub.Topic)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientUnsubscribed(e, reasonCodes, outQueue)
}

func clientUnsubscribed(e Event, reasonCodes []uint8, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.session.ClientId
	o.packet = Unsuback(e.packet.PacketIdentifier(), reasonCodes, e.session.ProtocolVersion)
	outQueue <- o
}

func clientUnsubscription(db *gorm.DB, clientId string, topic string) uint8 {
	var sub model.Subscription
	if db.Where("topic = ? and client_id = ?", topic, clientId).First(&sub).RecordNotFound() {
		log.Println("no subscription to unsubscribe", topic, clientId)
	}
	db.Delete(sub)
	return 0
}
