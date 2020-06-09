package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
)

func Pubcomp(packetIdentifier int, reasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBCOMP) << 4
	if reasonCode == 0 {
		p.remainingLength = 2
		p.remainingBytes = Write2BytesInt(packetIdentifier)
	} else {
		p.remainingLength = 3
		p.remainingBytes = Write2BytesInt(packetIdentifier)
		p.remainingBytes = append(p.remainingBytes, reasonCode)
	}
	if protocolVersion >= MQTT_V5 {
		// properties
		p.remainingLength = p.remainingLength + 1
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	return p
}
func pubcompReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PUBCOMPED
	event.clientId = session.ClientId
	event.session = session
	i := 2 // expect packet identifier in first 2 bytes
	if i < len(p.remainingBytes) {
		p.reasonCode = p.remainingBytes[i]
	}
	if session.ProtocolVersion >= MQTT_V5 {
		_, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			event.err = uint8(err)
			events <- event
			return
		}
	}
	event.packet = p
	events <- event
}
