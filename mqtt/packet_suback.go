package mqtt

func Suback(packetIdentifier int, subscribed int, qosAccepted uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_SUBACK) << 4
	p.remainingLength = 2 + subscribed
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	if protocolVersion >= MQTT_V5 {
		// properties
		p.remainingLength = p.remainingLength + 1
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	subRes := make([]byte, subscribed)
	for i := 0; i < subscribed; i++ {
		subRes[i] = qosAccepted
	}
	p.remainingBytes = append(p.remainingBytes, subRes...)
	return p
}
