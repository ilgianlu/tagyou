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
	i := 0
	pi := Read2BytesInt(p.remainingBytes, i)
	p.packetIdentifier = pi
	i = i + 2
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
	unsubs := make([]string, 10)
	j := 0
	for {
		var unsubevent Event
		unsubevent.eventType = EVENT_UNSUBSCRIPTION
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		unsubs[j] = string(p.remainingBytes[i : i+sl])
		unsubevent.clientId = session.ClientId
		unsubevent.topic = unsubs[j]
		events <- unsubevent
		i = i + sl
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	p.subscribedCount = j
	events <- event
}
