package mqtt

func Unsuback(packetIdentifier int, unsubscribed int, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_UNSUBACK) << 4
	p.remainingLength = 2 + unsubscribed
	p.remainingBytes = make([]byte, 2+unsubscribed)
	p.remainingBytes[0] = byte(packetIdentifier & 0xFF00 >> 8)
	p.remainingBytes[1] = byte(packetIdentifier & 0x00FF)
	if protocolVersion >= MQTT_V5 {
		// properties
		p.remainingLength = p.remainingLength + 1
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	return p
}
