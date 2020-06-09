package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
)

func unsubscribeReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_UNSUBSCRIBED
	event.clientId = session.ClientId
	event.session = session
	event.packet = p
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
	event.subscriptions = make([]model.Subscription, 0)
	j := 0
	for {
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		unsub := model.Subscription{
			Topic: string(p.remainingBytes[i : i+sl]),
		}
		event.subscriptions = append(event.subscriptions, unsub)
		i = i + sl
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	events <- event
}
