package mqtt

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func rangeEvents(connections Connections, db *gorm.DB, events <-chan Event, outQueue chan<- OutData) {
	for e := range events {
		switch e.eventType {
		case EVENT_CONNECT:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client connect")
			clientConnection(db, connections, e, outQueue)
		case EVENT_SUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client subscribed")
			clientSubscribed(e, outQueue)
		case EVENT_SUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.subscription.ClientId, "client subscription", e.subscription.Topic)
			clientSubscription(db, e, outQueue)
		case EVENT_UNSUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscribed")
			clientUnsubscribed(e, outQueue)
		case EVENT_UNSUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscription", e.topic)
			clientUnsubscription(db, e)
		case EVENT_PUBLISH:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client published to", e.published.topic)
			clientPublish(db, e, outQueue)
		case EVENT_PUBACKED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client acked message", e.packet.packetIdentifier)
			clientPuback(db, e)
		case EVENT_PUBRECED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub received message", e.packet.packetIdentifier)
			clientPubrec(e, outQueue)
		case EVENT_PUBRELED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub releases message", e.packet.packetIdentifier)
			clientPubrel(e, outQueue)
		case EVENT_PUBCOMPED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub complete message", e.packet.packetIdentifier)
			clientPubcomp(db, e)
		case EVENT_PING:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client ping")
			clientPing(e, outQueue)
		case EVENT_DISCONNECT:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client disconnect")
			clientDisconnect(db, connections, e)
		case EVENT_WILL_SEND:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "sending will message")
			sendWill(db, e, outQueue)
		case EVENT_PACKET_ERR:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "packet error")
			clientDisconnect(db, connections, e)
		}
	}
}

func clientConnection(db *gorm.DB, connections Connections, e Event, outQueue chan<- OutData) {
	if DISALLOW_ANONYMOUS_LOGIN && !model.CheckAuth(db, e.clientId, e.session.Username, e.session.Password) {
		log.Println("wrong connect credentials")
		return
	}
	if c, ok := connections[e.clientId]; ok {
		log.Println("session taken over")
		p := Connack(false, SESSION_TAKEN_OVER)
		sendSimple(e.clientId, p, outQueue)
		closeClient(c)
		removeClient(e.clientId, connections)
	}
	connections[e.clientId] = e.session.Conn
	if e.session.CleanStart() {
		db.Delete("client_id = ?", e.clientId)
	} else {
		db.Model(&model.Subscription{}).Where("client_id = ?", e.clientId).UpdateColumn("enabled", true)
	}
	if e.err != 0 {
		sendSimple(e.clientId, Connack(false, e.err), outQueue)
	} else {
		sendSimple(e.clientId, Connack(false, CONNECT_OK), outQueue)
	}
}

func clientSubscribed(e Event, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.clientId
	o.packet = Suback(e.packet.packetIdentifier, e.packet.subscribedCount, e.subscription.QoS)
	outQueue <- o
}

func clientSubscription(db *gorm.DB, e Event, outQueue chan<- OutData) {
	db.Create(&e.subscription)
	sendRetain(db, e, outQueue)
}

func sendRetain(db *gorm.DB, e Event, outQueue chan<- OutData) {
	var retains []model.Retain
	db.Where("topic = ?", e.subscription.Topic).Find(&retains)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		sendSimple(e.clientId, Publish(e.subscription.QoS, true, r.Topic, r.ApplicationMessage), outQueue)
	}
}

func clientUnsubscribed(e Event, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.clientId
	o.packet = Unsuback(e.packet.packetIdentifier, e.packet.subscribedCount)
	outQueue <- o
}

func clientUnsubscription(db *gorm.DB, e Event) {
	var sub model.Subscription
	if db.Where("topic = ? and client_id = ?", e.topic, e.clientId).First(&sub).RecordNotFound() {
		log.Println("no subscription to unsubscribe", e.topic, e.clientId)
	}
	db.Delete(sub)
}

func clientPublish(db *gorm.DB, e Event, outQueue chan<- OutData) {
	if e.published.retain {
		saveRetain(db, e)
	}
	sendForward(db, e.published.topic, e.packet, outQueue)
	if e.published.qos == 1 {
		sendSimple(e.clientId, Puback(e.packet.packetIdentifier, PUBACK_SUCCESS), outQueue)
	}
	if e.published.qos == 2 {
		sendSimple(e.clientId, Pubrec(e.packet.packetIdentifier, PUBREC_SUCCESS), outQueue)
	}
}

func sendForward(db *gorm.DB, topic string, packet Packet, outQueue chan<- OutData) {
	topicSegments := strings.Split(topic, TOPIC_SEPARATOR)
	subs := findDests(db, topicSegments)
	sendSubscribers(subs, packet, outQueue)
}

func findDests(db *gorm.DB, topicSegments []string) []model.Subscription {
	subs := []model.Subscription{}
	for i := 1; i <= len(topicSegments); i++ {
		subT := append(make([]string, 0), topicSegments[:i]...)
		if len(subT) < len(topicSegments) {
			subT = append(subT, TOPIC_WILDCARD)
		}
		t := strings.Join(subT, TOPIC_SEPARATOR)
		ss := []model.Subscription{}
		db.Where("topic = ?", t).Find(&ss)
		subs = append(subs, ss...)
	}
	return subs
}

func sendSubscribers(subscribers []model.Subscription, packet Packet, outQueue chan<- OutData) {
	for _, s := range subscribers {
		sendSimple(s.ClientId, packet, outQueue)
	}
}

func sendSimple(clientId string, p Packet, outQueue chan<- OutData) {
	var o OutData
	o.clientId = clientId
	o.packet = p
	outQueue <- o
}

func clientPuback(db *gorm.DB, e Event) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok remove retry
	removeRetry(db, e.session.ClientId, e.packet.packetIdentifier)
}

func removeRetry(db *gorm.DB, clientId string, packetIdentifier int) {
	var r model.Retry
	if db.Where("client_id = ? and packet_identifier = ?").First(&r).RecordNotFound() {
		log.Println("ack for invalid retry", clientId, packetIdentifier)
	} else {
		log.Println("retry found, removing...", clientId, packetIdentifier)
		db.Delete(&r)
	}
}

func clientPubrec(e Event, outQueue chan<- OutData) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok
	// if retry in wait for pub rec -> send pub rel
	var o OutData
	o.clientId = e.clientId
	o.packet = Pubrel(e.packet.packetIdentifier, PUBREL_SUCCESS)
	outQueue <- o
	// change retry state to wait for pubcomp
}

func clientPubrel(e Event, outQueue chan<- OutData) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok
	// if retry in wait for pubrel -> send pub comp
	var o OutData
	o.clientId = e.clientId
	o.packet = Pubcomp(e.packet.packetIdentifier, PUBCOMP_SUCCESS)
	outQueue <- o
}

func clientPubcomp(db *gorm.DB, e Event) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok
	// if retry in wait for pubcomp -> remove retry
	removeRetry(db, e.session.ClientId, e.packet.packetIdentifier)
}

func saveRetain(db *gorm.DB, e Event) {
	var r model.Retain
	r.Topic = e.published.topic
	r.ApplicationMessage = e.packet.remainingBytes[e.packet.applicationMessage:]
	r.CreatedAt = time.Now()
	db.Create(&r)
}

func clientPing(e Event, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.clientId
	o.packet = PingResp()
	outQueue <- o
}

func clientDisconnect(db *gorm.DB, connections Connections, e Event) {
	if conn, ok := connections[e.clientId]; ok {
		closeClient(conn)
		removeClient(e.clientId, connections)
	}
}

func sendWill(db *gorm.DB, e Event, outQueue chan<- OutData) {
	var s model.Session
	if db.First(&s, "client_id = ?", e.clientId).RecordNotFound() {
		return
	}
	if s.WillTopic != "" {
		p := Publish(s.WillQoS(), s.WillRetain(), s.WillTopic, s.WillMessage)
		sendForward(db, s.WillTopic, p, outQueue)
	}
}

func closeClient(connection net.Conn) {
	err := connection.Close()
	if err != nil {
		log.Println("could not close conn", err)
	}
}

func removeClient(clientId string, connections Connections) {
	delete(connections, clientId)
}
