package mqtt

func Connack(sessionPresent bool, reasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_CONNACK) << 4

	p.remainingBytes = make([]byte, 2)
	if sessionPresent {
		p.remainingBytes[0] = 1
	} else {
		p.remainingBytes[0] = 0
	}
	p.remainingBytes[1] = reasonCode
	if protocolVersion >= MQTT_V5 {
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingLength = len(p.remainingBytes)
	return p
}
