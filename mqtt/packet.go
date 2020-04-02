package mqtt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type Packet []byte

func (req Packet) Respond(db *bolt.DB, e chan<- Event, connStatus *ConnStatus, event *Event) (Packet, error) {
	i := 0
	t := (req[i] & 0xF0) >> 4
	fmt.Printf("packet type %d\n", t)
	fmt.Printf("flags %d\n", req[i]&0x0F)
	i++
	l := int(req[i])
	fmt.Printf("remaining length %d\n", l)
	i++
	fmt.Println("payload", req[i:i+l])
	switch t {
	case 1:
		fmt.Println("Connect message")
		return connectReq(db, e, connStatus, req[i:i+l], event)
	case 8:
		fmt.Println("Subscribe message")
		return subscribeReq(db, e, connStatus, req[i:i+l])
	case 3:
		fmt.Println("Publish message")
		return publishReq(db, e, connStatus, req[i:i+l], event)
	default:
		return Connack(CONNECT_UNSPECIFIED_ERROR), nil
	}
}

func publishReq(db *bolt.DB, e chan<- Event, connStatus *ConnStatus, req Packet, event *Event) (Packet, error) {
	i := 0
	event.eventType = 2
	tl := Read2BytesInt(req, i)
	fmt.Println("topic length", tl)
	i = i + 2
	topic := string(req[i : i+tl])
	fmt.Println("pub topic", topic)
	event.topic = topic
	i = i + tl
	// pi := Read2BytesInt(req, i)
	// fmt.Println("packet identifier", pi)
	// i = i + 2
	pay := req[i:]
	fmt.Println("payload", pay)
	event.message = string(pay)
	event.clientId = connStatus.clientId
	e <- *event
	return nil, nil
}

func subscribeReq(db *bolt.DB, e chan<- Event, connStatus *ConnStatus, req Packet) (Packet, error) {
	i := 0
	pi := Read2BytesInt(req, i)
	i = i + 2
	fmt.Println("packet identifier", pi)
	if connStatus.protocolVersion == 5 {
		pl := int(req[i])
		fmt.Println("property length", pl)
		if pl > 0 {
			props := req[i : i+pl]
			fmt.Println("properties", props)
			si, _, err := ReadVarInt(props)
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
		var subevent Event
		subevent.eventType = 1
		sl := Read2BytesInt(req, i)
		i = i + 2
		subs[j] = string(req[i : i+sl])
		fmt.Println(j, "subscribtion:", subs[j])
		subevent.clientId = connStatus.clientId
		subevent.topic = subs[j]
		e <- subevent
		i = i + sl
		if i >= len(req)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	return Suback(pi, j+1), nil
}

func Suback(packetIdentifier int, subscribed int) Packet {
	p := make(Packet, 4+subscribed)
	p[0] = uint8(9) << 4
	p[1] = uint8(2 + subscribed)
	p[2] = byte(packetIdentifier & 0xFF00 >> 8)
	p[3] = byte(packetIdentifier & 0x00FF)
	return p
}

func connectReq(db *bolt.DB, e chan<- Event, connStatus *ConnStatus, req Packet, event *Event) (Packet, error) {
	event.eventType = 0
	i := 0
	pl := Read2BytesInt(req, i)
	i = i + 2
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

	fmt.Println("keepAlive", Read2BytesInt(ka, 0))
	connStatus.keepAlive = ka
	i = i + 2
	cil := Read2BytesInt(req, i)
	i = i + 2
	clientId := string(req[i : i+cil])
	fmt.Println("clientId", clientId)
	connStatus.clientId = clientId
	event.clientId = clientId
	// connStatus.persist(db)
	e <- *event
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
