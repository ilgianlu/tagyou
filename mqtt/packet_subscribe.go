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
		pl, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			event.err = uint8(err)
			events <- event
			return
		}
		i = i + pl
	}
	// payload
	j := 0
	event.subscriptions = make([]model.Subscription, 0)
	for {
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		s := string(p.remainingBytes[i : i+sl])
		sub := model.Subscription{
			ClientId: session.ClientId,
			Topic:    s,
		}
		i = i + sl
		if p.remainingBytes[i]&0x12 != 0 {
			log.Println("ignore this subscription & stop")
			break
		}
		sub.RetainHandling = p.remainingBytes[i] & 0x30 >> 4
		sub.RetainAsPublished = p.remainingBytes[i] & 0x08 >> 3
		sub.NoLocal = p.remainingBytes[i] & 0x04 >> 2
		sub.QoS = p.remainingBytes[i] & 0x03
		sub.Enabled = true
		sub.CreatedAt = time.Now()
		i++
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > conf.MAX_TOPIC_SINGLE_SUBSCRIBE {
			break
		}
	}
	event.packet = p
	events <- event
}
