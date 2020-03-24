package mqtt

type Packet []byte

func (p Packet) PacketType() uint8 {
	return (p[0] & 0xF0) >> 4
}

func (p Packet) Flags() uint8 {
	return p[0] & 0x0F
}

func (p Packet) RemainingLength() uint8 {
	return p[1]
}

func (p Packet) Payload() []byte {
	l := p.RemainingLength()
	return p[2 : 2+l]
}

func (p Packet) ProtocolNameLength() int {
	pay := p.Payload()
	return int(pay[0])<<8 + int(pay[1])
}

func (p Packet) ProtocolName() []byte {
	l := p.ProtocolNameLength()
	return p.Payload()[2 : 2+l]
}
