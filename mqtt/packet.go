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

func (p Packet) ProtocolVersion() int {
	l := p.ProtocolNameLength()
	return int(p.Payload()[2+l])
}

func (p Packet) ConnectFlags() int {
	l := p.ProtocolNameLength()
	return int(p.Payload()[2+l+1])
}

func (p Packet) KeepAlive() int {
	pay := p.Payload()
	l := p.ProtocolNameLength()
	return int(pay[2+l+1+1])<<8 + int(pay[2+l+1+2])
}

func (p Packet) ClientId() string {
	pay := p.Payload()
	l := p.ProtocolNameLength()
	lcid := int(pay[2+l+1+2+1])<<8 + int(pay[2+l+1+2+2])
	return string(pay[2+l+1+2+2+1 : 2+l+1+2+2+1+lcid])
}

func Connack() Packet {
	p := make(Packet, 5)
	p[0] = uint8(2) << 4
	p[1] = uint8(3)
	return p
}
