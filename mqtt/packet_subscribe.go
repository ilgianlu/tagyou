package mqtt

import (
	"log"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func subscribeReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_SUBSCRIBED
	event.clientId = session.ClientId
	event.session = session
	// variable header
	i := 2 // 2 bytes for packet identifier
	if session.ProtocolVersion >= MQTT_V5 {
		pl, pp, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			event.err = uint8(err)
			events <- event
			return
		}
		p.propertiesLength = pl
		p.propertiesPos = pp
		i = i + pl
	}
	// payload
	j := 0
	for {
		var subevent Event
		subevent.eventType = EVENT_SUBSCRIPTION
		subevent.session = session
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		s := string(p.remainingBytes[i : i+sl])
		subevent.subscription.ClientId = session.ClientId
		subevent.subscription.Topic = s
		i = i + sl
		if p.remainingBytes[i]&0x12 != 0 {
			log.Println("ignore this subscription & stop")
			break
		}
		subevent.subscription.RetainHandling = p.remainingBytes[i] & 0x30 >> 4
		subevent.subscription.RetainAsPublished = p.remainingBytes[i] & 0x08 >> 3
		subevent.subscription.NoLocal = p.remainingBytes[i] & 0x04 >> 2
		subevent.subscription.QoS = p.remainingBytes[i] & 0x03
		subevent.subscription.Enabled = true
		subevent.subscription.CreatedAt = time.Now()
		events <- subevent
		i++
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > conf.MAX_TOPIC_SINGLE_SUBSCRIBE {
			break
		}
	}
	p.subscribedCount = j
	event.packet = p
	events <- event
}
