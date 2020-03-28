package mqtt

import (
	"errors"
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
	i := 0
	t := (req[i] & 0xF0) >> 4
	fmt.Printf("packet type %d\n", t)
	fmt.Printf("flags %d\n", req[i]&0x0F)
	i++
	l := int(req[i])
	fmt.Printf("remaining length %d\n", l)
	i++
	fmt.Println("payload", req[i:i+l])
	if 1 == t {
		return connectReq(db, req[i:i+l])
	} else if 8 == t {
		fmt.Println("Subscribe message")
		return subscribeReq(db, req[i:i+l])
	} else {
		return Connack(CONNECT_UNSPECIFIED_ERROR), nil
	}
}

func subscribeReq(db *bolt.DB, req Packet) (Packet, error) {
	i := 0
	pi := int(req[i]) << 8
	i++
	pi = pi + int(req[i])
	fmt.Println("packet identifier", pi)
	i++
	pl := int(req[i])
	fmt.Println("property length", pl)
	if pl > 0 {
		props := req[i : i+pl]
		fmt.Println("properties", props)
		si, _, err := decodeSubscriptionIdentifier(props)
		if err != nil {
			fmt.Println("err decoding sub ident", err)
			return Connack(CONNECT_UNSPECIFIED_ERROR), nil
		}
		fmt.Println("subscription identifier", si)
		i = i + pl
	}
	i++
	pay := req[i:len(req)]
	fmt.Println("subscribe payload", pay)
	return Suback(pi), nil
}

func decodeSubscriptionIdentifier(props []byte) (int, int, error) {
	multiplier := 1
	value := 0
	i := 0
	encodedByte := props[i]
	for ok := true; ok; ok = int(encodedByte&128) != 0 {
		value = value + int(encodedByte&127)*multiplier
		if multiplier > 128*128*128 {
			return 0, 0, errors.New("malformed value")
		}
		multiplier *= 128
		i++
		encodedByte = props[i]
	}
	return value, i - 1, nil
}

func connectReq(db *bolt.DB, req Packet) (Packet, error) {
	i := 0
	pl := int(req[i]) << 8
	i++
	pl = pl + int(req[i])
	i++
	fmt.Println("protocolName", string(req[i:i+pl]))
	i = i + pl
	v := int(req[i])
	fmt.Println("protocolVersion", v)
	i++
	if v < 4 {
		fmt.Println("unsupported protocol version err", v)
		return Connack(CONNECT_UNSUPPORTED_PROTOCOL_VERSION), nil
	}
	fmt.Println("connectFlags", req[i])
	i++
	ka := int(req[i]) << 8
	i++
	ka = ka + int(req[i])
	i++
	fmt.Println("keepAlive", ka)
	cil := int(req[i]) << 8
	i++
	cil = cil + int(req[i])
	i++
	clientId := string(req[i : i+cil])
	fmt.Println("clientId", clientId)
	newClient(db, clientId)
	return Connack(CONNECT_OK), nil
}

func Suback(packetIdentifier int) Packet {
	p := make(Packet, 5)
	p[0] = uint8(2) << 4
	p[1] = uint8(3)
	p[2] = byte(packetIdentifier & 0xFF >> 8)
	p[3] = byte(packetIdentifier & 0xFF)
	fmt.Println(p)
	return p
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
