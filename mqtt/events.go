package mqtt

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/ilgianlu/tagyou/conf"
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
			onSubscribe(db, e, outQueue)
		case EVENT_UNSUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscribed")
			onUnsubscribe(db, e, outQueue)
		case EVENT_PUBLISH:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client published to", e.topic)
			clientPublish(db, e, outQueue)
		case EVENT_PUBACKED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client acked message", e.packet.PacketIdentifier())
			clientPuback(db, e)
		case EVENT_PUBRECED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub received message", e.packet.PacketIdentifier())
			clientPubrec(db, e, outQueue)
		case EVENT_PUBRELED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub releases message", e.packet.PacketIdentifier())
			clientPubrel(db, e, outQueue)
		case EVENT_PUBCOMPED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub complete message", e.packet.PacketIdentifier())
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
	if conf.DISALLOW_ANONYMOUS_LOGIN && !model.CheckAuth(db, e.clientId, e.session.Username, e.session.Password) {
		log.Println("wrong connect credentials")
		return
	}

	if c, ok := connections[e.clientId]; ok {
		log.Println("session taken over")
		p := Connack(false, SESSION_TAKEN_OVER, e.session.ProtocolVersion)
		sendSimple(e.clientId, p, outQueue)
		closeClient(c)
		removeClient(e.clientId, connections)
	}
	connections[e.clientId] = e.session.Conn
	sendSimple(e.clientId, Connack(false, CONNECT_OK, e.session.ProtocolVersion), outQueue)

	startSession(db, e.session)
}

func startSession(db *gorm.DB, session *model.Session) {
	if db.Where("client_id = ?", session.ClientId).First(&session).RecordNotFound() {
		db.Create(&session)
	} else {
		if session.CleanStart() {
			model.CleanSession(db, session.ClientId)
			db.Create(&session)
		} else {
			db.Save(&session)
		}
	}
}

func trimWildcard(topic string) string {
	lci := len(topic) - 1
	lc := topic[lci]
	if string(lc) == TOPIC_WILDCARD {
		topic = topic[:lci]
	}
	return topic
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
		model.DisconnectSession(db, e.clientId)
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

func sendForward(db *gorm.DB, topic string, packet Packet, outQueue chan<- OutData) {
	topicSegments := strings.Split(topic, TOPIC_SEPARATOR)
	subs := findDests(db, topicSegments)
	sendSubscribers(db, topic, subs, packet, outQueue)
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

func sendSubscribers(db *gorm.DB, topic string, subscribers []model.Subscription, packet Packet, outQueue chan<- OutData) {
	for _, s := range subscribers {
		qos := getQos(packet.QoS(), s.QoS)
		if qos == conf.QOS0 {
			// prepare publish packet qos 0 no packet identifier
			p := Publish(conf.QOS0, packet.Retain(), topic, 0, packet.ApplicationMessage())
			sendSimple(s.ClientId, p, outQueue)
		} else if qos == conf.QOS1 {
			// prepare publish packet qos 1 (if sub permit) new packet identifier
			p := Publish(qos, packet.Retain(), topic, newPacketIdentifier(), packet.ApplicationMessage())
			r := model.Retry{
				ClientId:           s.ClientId,
				PacketIdentifier:   packet.PacketIdentifier(),
				Qos:                qos,
				Dup:                packet.Dup(),
				ApplicationMessage: packet.ApplicationMessage(),
				AckStatus:          model.WAIT_FOR_PUB_ACK,
				CreatedAt:          time.Now(),
			}
			db.Save(&r)
			sendSimple(r.ClientId, p, outQueue)
		} else if qos == 2 {
			// prepare publish packet qos 2 (if sub permit) new packet identifier
			p := Publish(qos, packet.Retain(), topic, newPacketIdentifier(), packet.ApplicationMessage())
			r := model.Retry{
				ClientId:           s.ClientId,
				PacketIdentifier:   packet.PacketIdentifier(),
				Qos:                qos,
				Dup:                packet.Dup(),
				ApplicationMessage: packet.ApplicationMessage(),
				AckStatus:          model.WAIT_FOR_PUB_REL,
				CreatedAt:          time.Now(),
			}
			db.Save(&r)
			sendSimple(r.ClientId, p, outQueue)
		}

	}
}

func getQos(pubQos uint8, subQos uint8) uint8 {
	if pubQos > subQos {
		return subQos
	} else {
		return pubQos
	}
}

func sendSimple(clientId string, p Packet, outQueue chan<- OutData) {
	var o OutData
	o.clientId = clientId
	o.packet = p
	outQueue <- o
}

func saveRetain(db *gorm.DB, e Event) {
	var r model.Retain
	r.Topic = e.topic
	r.ApplicationMessage = e.packet.ApplicationMessage()
	r.CreatedAt = time.Now()
	db.Delete(&r)
	if len(r.ApplicationMessage) > 0 {
		db.Create(&r)
	}
}

func sendWill(db *gorm.DB, e Event, outQueue chan<- OutData) {
	var s model.Session
	if db.First(&s, "client_id = ?", e.clientId).RecordNotFound() {
		return
	}
	if s.WillTopic != "" {
		p := Publish(s.WillQoS(), s.WillRetain(), s.WillTopic, newPacketIdentifier(), s.WillMessage)
		sendForward(db, s.WillTopic, p, outQueue)
	}
}
