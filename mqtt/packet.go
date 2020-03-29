package mqtt

import (
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type Packet []byte

func (req Packet) Respond(db *bolt.DB, connStatus *ConnStatus) (Packet, error) {
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
		fmt.Println("Connect message")
		return connectReq(db, connStatus, req[i:i+l])
	} else if 8 == t {
		fmt.Println("Subscribe message")
		return subscribeReq(db, connStatus, req[i:i+l])
	} else {
		return Connack(CONNECT_UNSPECIFIED_ERROR), nil
	}
}

func subscribeReq(db *bolt.DB, connStatus *ConnStatus, req Packet) (Packet, error) {
	i := 0
	pi := read2BytesInt(req, i)
	i = i + 2
	fmt.Println("packet identifier", pi)
	if connStatus.protocolVersion == 5 {
		pl := int(req[i])
		fmt.Println("property length", pl)
		if pl > 0 {
			props := req[i : i+pl]
			fmt.Println("properties", props)
			si, _, err := readVarInt(props)
			if err != nil {
				fmt.Println("err decoding sub ident", err)
				return Connack(CONNECT_UNSPECIFIED_ERROR), nil
			}
			fmt.Println("subscription identifier", si)
			i = i + pl
		}
		i++
	}
	pay := req[i:]
	fmt.Println("subscribe payload", pay)
	subs := make([]string, 10)
	j := 0
	for {
		sl := int(req[i]) << 8
		i++
		sl = sl + int(req[i])
		i++
		subs[j] = string(req[i : i+sl])
		fmt.Println(j, "subscribtion:", subs[j])
		i = i + sl
		if i >= len(req)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	subacks := make([]byte, j)
	return Suback(pi, subacks), nil
}

func Suback(packetIdentifier int, acks []byte) Packet {
	p := make(Packet, 4+len(acks))
	p[0] = uint8(9) << 4
	p[1] = uint8(0)
	p[2] = byte(packetIdentifier & 0xFF >> 8)
	p[3] = byte(packetIdentifier & 0xFF)
	fmt.Println(p)
	return p
}

func readVarInt(props []byte) (int, int, error) {
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

func read2BytesInt(a []byte, i int) int {
	v := int(a[i]) << 8
	i++
	return v + int(a[i])
}

func connectReq(db *bolt.DB, connStatus *ConnStatus, req Packet) (Packet, error) {
	i := 0
	pl := int(req[i]) << 8
	i++
	pl = pl + int(req[i])
	i++
	fmt.Println("protocolName", string(req[i:i+pl]))
	i = i + pl
	v := req[i]
	fmt.Println("protocolVersion", v)
	connStatus.protocolVersion = v
	i++
	if int(v) < MINIMUM_SUPPORTED_PROTOCOL {
		fmt.Println("unsupported protocol version err", v)
		return Connack(CONNECT_UNSUPPORTED_PROTOCOL_VERSION), nil
	}
	fmt.Println("connectFlags", req[i])
	connStatus.connectFlags = req[i]
	i++
	ka := req[i : i+2]
	i = i + 2
	fmt.Println("keepAlive", read2BytesInt(ka, 0))
	connStatus.keepAlive = ka
	cil := read2BytesInt(req, i)
	i = i + 2
	clientId := string(req[i : i+cil])
	fmt.Println("clientId", clientId)
	connStatus.clientId = clientId
	newClient(db, connStatus)
	return Connack(CONNECT_OK), nil
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
