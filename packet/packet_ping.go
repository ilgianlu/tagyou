package packet

func (p *Packet) pingReq() int {
	p.Event = EVENT_PING
	return 0
}

func PingResp() Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PINGRES) << 4
	p.remainingLength = 0
	return p
}
