package mqtt

import (
	"log"
)

func Puback(packetIdentifier int, reasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBACK) << 4
	// var header
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	p.remainingBytes = append(p.remainingBytes, reasonCode)
	if protocolVersion >= MQTT_V5 {
		// TODO: encode properties
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingLength = len(p.remainingBytes)
	return p
}

func (p *Packet) pubackReq() int {
	p.event = EVENT_PUBACKED
	i := 2 // expect packet identifier in first 2 bytes
	if i < len(p.remainingBytes) {
		p.reasonCode = p.remainingBytes[i]
	}
	if p.session.ProtocolVersion >= MQTT_V5 {
		_, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			return err
		}
	}
	return 0
}
