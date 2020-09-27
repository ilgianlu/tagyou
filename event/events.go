package event

import (
	"log"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	tpc "github.com/ilgianlu/tagyou/topic"
	"gorm.io/gorm"
)

func RangeEvents(connections model.Connections, db *gorm.DB, events <-chan *packet.Packet, outQueue chan<- *out.OutData) {
	for p := range events {
		switch p.Event {
		case packet.EVENT_CONNECT:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "client connect")
			onConnect(db, connections, p, outQueue)
		case packet.EVENT_SUBSCRIBED:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "client subscribed")
			onSubscribe(db, p, outQueue)
		case packet.EVENT_UNSUBSCRIBED:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "client unsubscribed")
			onUnsubscribe(db, p, outQueue)
		case packet.EVENT_PUBLISH:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "client published to", p.Topic)
			onPublish(db, p, outQueue)
		case packet.EVENT_PUBACKED:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "client acked message", p.PacketIdentifier())
			clientPuback(db, p)
		case packet.EVENT_PUBRECED:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "pub received message", p.PacketIdentifier())
			clientPubrec(db, p, outQueue)
		case packet.EVENT_PUBRELED:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "pub releases message", p.PacketIdentifier())
			clientPubrel(db, p, outQueue)
		case packet.EVENT_PUBCOMPED:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "pub complete message", p.PacketIdentifier())
			clientPubcomp(db, p)
		case packet.EVENT_PING:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "client ping")
			onPing(p, outQueue)
		case packet.EVENT_DISCONNECT:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "client disconnect")
			clientDisconnect(db, connections, p.Session.ClientId)
		case packet.EVENT_WILL_SEND:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "sending will message")
			sendWill(db, p, outQueue)
		case packet.EVENT_PACKET_ERR:
			log.Println("//!! EVENT type", p.Event, p.Session.ClientId, "packet error")
			clientDisconnect(db, connections, p.Session.ClientId)
		}
	}
}

func trimWildcard(topic string) string {
	lci := len(topic) - 1
	lc := topic[lci]
	if string(lc) == conf.TOPIC_WILDCARD {
		topic = topic[:lci]
	}
	return topic
}

func onPing(p *packet.Packet, outQueue chan<- *out.OutData) {
	var o out.OutData
	o.ClientId = p.Session.ClientId
	o.Packet = packet.PingResp()
	outQueue <- &o
}

func clientDisconnect(db *gorm.DB, connections model.Connections, clientId string) {
	if _, ok := connections.Exists(clientId); ok {
		connections.Close(clientId)
		connections.Remove(clientId)
		model.DisconnectSession(db, clientId)
	}
}

func sendForward(db *gorm.DB, topic string, p *packet.Packet, outQueue chan<- *out.OutData) {
	destSubs := tpc.Explode(topic)
	subs := []model.Subscription{}
	if err := db.Where("topic IN (?)", destSubs).Find(&subs).Error; err != nil {
		log.Println("could not query for subscriptions:", err)
		return
	}
	sendSubscribers(db, topic, subs, p, outQueue)
}

func sendSubscribers(db *gorm.DB, topic string, subscribers []model.Subscription, p *packet.Packet, outQueue chan<- *out.OutData) {
	for _, s := range subscribers {
		qos := getQos(p.QoS(), s.Qos)
		if qos == conf.QOS0 {
			// prepare publish packet qos 0 no packet identifier
			p := packet.Publish(s.ProtocolVersion, conf.QOS0, p.Retain(), topic, 0, p.ApplicationMessage())
			sendSimple(s.ClientId, &p, outQueue)
		} else if qos == conf.QOS1 {
			// prepare publish packet qos 1 (if sub permit) new packet identifier
			p := packet.Publish(s.ProtocolVersion, qos, p.Retain(), topic, packet.NewPacketIdentifier(), p.ApplicationMessage())
			r := model.Retry{
				ClientId:           s.ClientId,
				PacketIdentifier:   p.PacketIdentifier(),
				Qos:                qos,
				Dup:                p.Dup(),
				ApplicationMessage: p.ApplicationMessage(),
				AckStatus:          model.WAIT_FOR_PUB_ACK,
				CreatedAt:          time.Now(),
			}
			db.Save(&r)
			sendSimple(r.ClientId, &p, outQueue)
		} else if qos == 2 {
			// prepare publish packet qos 2 (if sub permit) new packet identifier
			p := packet.Publish(s.ProtocolVersion, qos, p.Retain(), topic, packet.NewPacketIdentifier(), p.ApplicationMessage())
			r := model.Retry{
				ClientId:           s.ClientId,
				PacketIdentifier:   p.PacketIdentifier(),
				Qos:                qos,
				Dup:                p.Dup(),
				ApplicationMessage: p.ApplicationMessage(),
				AckStatus:          model.WAIT_FOR_PUB_REL,
				CreatedAt:          time.Now(),
			}
			db.Save(&r)
			sendSimple(r.ClientId, &p, outQueue)
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

func sendSimple(clientId string, p *packet.Packet, outQueue chan<- *out.OutData) {
	var o out.OutData
	o.ClientId = clientId
	o.Packet = *p
	outQueue <- &o
}

func saveRetain(db *gorm.DB, p *packet.Packet) {
	var r model.Retain
	r.Topic = p.Topic
	r.ApplicationMessage = p.ApplicationMessage()
	r.CreatedAt = time.Now()
	db.Delete(&r)
	if len(r.ApplicationMessage) > 0 {
		db.Create(&r)
	}
}

func sendWill(db *gorm.DB, p *packet.Packet, outQueue chan<- *out.OutData) {
	var s model.Session
	if err := db.First(&s, "client_id = ?", p.Session.ClientId).Error; err != nil {
		return
	}
	if s.WillTopic != "" {
		p := packet.Publish(p.Session.ProtocolVersion, s.WillQoS(), s.WillRetain(), s.WillTopic, packet.NewPacketIdentifier(), s.WillMessage)
		sendForward(db, s.WillTopic, &p, outQueue)
	}
}
