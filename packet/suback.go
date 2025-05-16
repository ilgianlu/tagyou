package packet

import "github.com/ilgianlu/tagyou/conf"

func Suback(packetIdentifier int, ReasonCodes []uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_SUBACK) << 4
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	if protocolVersion >= conf.MQTT_V5 {
		// TODO: encode properties ...
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingBytes = append(p.remainingBytes, ReasonCodes...)
	p.remainingLength = len(p.remainingBytes)
	return p
}
