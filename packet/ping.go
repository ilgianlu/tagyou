package packet

func (p *Packet) pingReq() int {
	return 0
}

func PingResp() Packet {
	var p Packet
	p.header = header(uint8(PACKET_TYPE_PINGRES) << 4)
	p.remainingLength = 0
	return p
}
