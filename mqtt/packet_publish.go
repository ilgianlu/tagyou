package mqtt

import (
	"log"
	"math/rand"

	"github.com/ilgianlu/tagyou/model"
)

func Publish(qos uint8, retain bool, topic string, packetIdentifier int, payload []byte) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBLISH) << 4
	p.header = p.header | qos<<1
	if retain {
		p.header = p.header | 1
	}
	// write var int length 2 + len(topic) + len(payload)
	p.remainingLength = 2 + len(topic) + len(payload)

	// write topic length
	p.remainingBytes = Write2BytesInt(len(topic))
	// write topic string
	p.remainingBytes = append(p.remainingBytes, []byte(topic)...)
	// write packet identifier only if qos > 0
	if qos != 0 {
		p.remainingBytes = append(p.remainingBytes, Write2BytesInt(packetIdentifier)...)
		p.remainingLength = p.remainingLength + 2
	}
	// write payload
	p.remainingBytes = append(p.remainingBytes, payload...)
	return p
}

func publishReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PUBLISH
	event.clientId = session.ClientId
	event.session = session
	event.published.dup = p.Dup()
	event.published.qos = p.QoS()
	event.published.retain = p.Retain()
	i := 0
	tl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	topic := string(p.remainingBytes[i : i+tl])
	event.published.topic = topic
	i = i + tl
	if event.published.qos != 0 {
		pi := Read2BytesInt(p.remainingBytes, i)
		p.packetIdentifier = pi
		i = i + 2
	}
	if session.ProtocolVersion >= MQTT_V5 {
		pl, pp, err := p.readProperties(i)
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
	p.applicationMessage = i
	event.packet = p
	events <- event
}

func newPacketIdentifier() int {
	return rand.Intn(65534)
}
