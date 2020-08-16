package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onUnsubscribe(db *gorm.DB, p Packet, outQueue chan<- OutData) {
	reasonCodes := []uint8{}
	for _, unsub := range p.subscriptions {
		rCode := clientUnsubscription(db, p.session.ClientId, unsub.Topic)
		reasonCodes = append(reasonCodes, rCode)
	}
	clientUnsubscribed(p, reasonCodes, outQueue)
}

func clientUnsubscribed(p Packet, reasonCodes []uint8, outQueue chan<- OutData) {
	var o OutData
	o.clientId = p.session.ClientId
	o.packet = Unsuback(p.PacketIdentifier(), reasonCodes, p.session.ProtocolVersion)
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
