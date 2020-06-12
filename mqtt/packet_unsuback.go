package mqtt

func Unsuback(packetIdentifier int, reasonCodes []uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_UNSUBACK) << 4
	// var header
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	if protocolVersion >= MQTT_V5 {
		// TODO: encode properties ...
		// no properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingBytes = append(p.remainingBytes, reasonCodes...)
	p.remainingLength = len(p.remainingBytes)
	return p
}
