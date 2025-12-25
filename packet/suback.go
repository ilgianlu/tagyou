package packet

import (
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/format"
)

func Suback(packetIdentifier int, reasonCodes []uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = header(uint8(PACKET_TYPE_SUBACK) << 4)
	p.remainingBytes = format.Write2BytesInt(packetIdentifier)
	if protocolVersion >= conf.MQTT_V5 {
		// TODO: encode properties ...
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingBytes = append(p.remainingBytes, reasonCodes...)
	p.remainingLength = len(p.remainingBytes)
	return p
}

func (p *Packet) subackReq() int {
	// this mqtt server is not able to subscribe
	return 99
}
