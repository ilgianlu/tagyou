package mqtt

import (
	"log"
)

func Pubrel(packetIdentifier int, reasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBREL) << 4
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	if reasonCode != 0 {
		p.remainingBytes = append(p.remainingBytes, reasonCode)
	}
	if protocolVersion >= MQTT_V5 {
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingLength = len(p.remainingBytes)
	return p
}

func (p *Packet) pubrelReq() int {
	p.event = EVENT_PUBRELED
	i := 2 // 2 bytes for packet identifier
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
