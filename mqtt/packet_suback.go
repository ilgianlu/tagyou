package mqtt

func Suback(packetIdentifier int, reasonCodes []uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_SUBACK) << 4
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	if protocolVersion >= MQTT_V5 {
		// TODO: encode properties ...
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingBytes = append(p.remainingBytes, reasonCodes...)
	p.remainingLength = len(p.remainingBytes)
	return p
}
