package mqtt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

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

func (req Packet) Respond(db *bolt.DB) (Packet, error) {
	if req.PacketType() == 1 {
		if req.ProtocolVersion() < 4 {
			fmt.Println("unsupported protocol version err", req.ProtocolVersion())
			return Connack(CONNECT_UNSUPPORTED_PROTOCOL_VERSION), nil
		}

		return Connack(CONNECT_OK), nil
	} else {
		return Connack(CONNECT_UNSPECIFIED_ERROR), nil
	}
}

func (p Packet) PrettyLog() {
	i := 0
	t := (p[i] & 0xF0) >> 4
	fmt.Printf("packet type %d\n", t)
	fmt.Printf("flags %d\n", p[i]&0x0F)
	i++
	l := int(p[i])
	fmt.Printf("remaining length %d\n", l)
	i++
	fmt.Println("payload", p[i:i+l])
	if t == 1 {
		pl := int(p[i]) << 8
		i++
		pl = pl + int(p[i])
		i++
		fmt.Println("protocolName", string(p[i:i+pl]))
		i = i + pl
		fmt.Println("protocolVersion", int(p[i]))
		i++
		fmt.Println("connectFlags", p[i])
		i++
		ka := int(p[i]) << 8
		i++
		ka = ka + int(p[i])
		i++
		fmt.Println("keepAlive", ka)
		cil := int(p[i]) << 8
		i++
		cil = cil + int(p[i])
		i++
		fmt.Println("clientId", string(p[i:i+cil]))
	}
}

func Connack(reasonCode uint8) Packet {
	p := make(Packet, 5)
	p[0] = uint8(2) << 4
	p[1] = uint8(3)
	p[2] = uint8(0)
	p[3] = reasonCode
	fmt.Println(p)
	return p
}
