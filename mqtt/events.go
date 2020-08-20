package mqtt

import (
	"log"
	"strings"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func rangeEvents(connections Connections, db *gorm.DB, events <-chan Packet, outQueue chan<- OutData) {
	for p := range events {
		switch p.event {
		case EVENT_CONNECT:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "client connect")
			onConnect(db, connections, p, outQueue)
		case EVENT_SUBSCRIBED:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "client subscribed")
			onSubscribe(db, p, outQueue)
		case EVENT_UNSUBSCRIBED:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "client unsubscribed")
			onUnsubscribe(db, p, outQueue)
		case EVENT_PUBLISH:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "client published to", p.topic)
			onPublish(db, p, outQueue)
		case EVENT_PUBACKED:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "client acked message", p.PacketIdentifier())
			clientPuback(db, p)
		case EVENT_PUBRECED:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "pub received message", p.PacketIdentifier())
			clientPubrec(db, p, outQueue)
		case EVENT_PUBRELED:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "pub releases message", p.PacketIdentifier())
			clientPubrel(db, p, outQueue)
		case EVENT_PUBCOMPED:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "pub complete message", p.PacketIdentifier())
			clientPubcomp(db, p)
		case EVENT_PING:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "client ping")
			onPing(p, outQueue)
		case EVENT_DISCONNECT:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "client disconnect")
			clientDisconnect(db, connections, p.session.ClientId)
		case EVENT_WILL_SEND:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "sending will message")
			sendWill(db, p, outQueue)
		case EVENT_PACKET_ERR:
			log.Println("//!! EVENT type", p.event, p.session.ClientId, "packet error")
			clientDisconnect(db, connections, p.session.ClientId)
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

func onPing(p Packet, outQueue chan<- OutData) {
	var o OutData
	o.clientId = p.session.ClientId
	o.packet = PingResp()
	outQueue <- o
}

func clientDisconnect(db *gorm.DB, connections Connections, clientId string) {
	if _, ok := connections.Exists(clientId); ok {
		connections.Close(clientId)
		connections.Remove(clientId)
		model.DisconnectSession(db, clientId)
	}
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
		qos := getQos(packet.QoS(), s.Qos)
		if qos == conf.QOS0 {
			// prepare publish packet qos 0 no packet identifier
			p := Publish(s.ProtocolVersion, conf.QOS0, packet.Retain(), topic, 0, packet.ApplicationMessage())
			sendSimple(s.ClientId, p, outQueue)
		} else if qos == conf.QOS1 {
			// prepare publish packet qos 1 (if sub permit) new packet identifier
			p := Publish(s.ProtocolVersion, qos, packet.Retain(), topic, newPacketIdentifier(), packet.ApplicationMessage())
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
			p := Publish(s.ProtocolVersion, qos, packet.Retain(), topic, newPacketIdentifier(), packet.ApplicationMessage())
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

func saveRetain(db *gorm.DB, p Packet) {
	var r model.Retain
	r.Topic = p.topic
	r.ApplicationMessage = p.ApplicationMessage()
	r.CreatedAt = time.Now()
	db.Delete(&r)
	if len(r.ApplicationMessage) > 0 {
		db.Create(&r)
	}
}

func sendWill(db *gorm.DB, p Packet, outQueue chan<- OutData) {
	var s model.Session
	if db.First(&s, "client_id = ?", p.session.ClientId).RecordNotFound() {
		return
	}
	if s.WillTopic != "" {
		p := Publish(p.session.ProtocolVersion, s.WillQoS(), s.WillRetain(), s.WillTopic, newPacketIdentifier(), s.WillMessage)
		sendForward(db, s.WillTopic, p, outQueue)
	}
}
