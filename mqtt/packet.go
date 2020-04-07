package mqtt

import (
	"fmt"
	"net"
)

type Packet []byte

func readHeader(conn net.Conn, event *Event) error {
	event.header = make([]byte, 2)
	n, err := conn.Read(event.header)
	if err != nil {
		return err
	}
	if n < 2 {
		return fmt.Errorf("header: read %d bytes, expected 2 bytes.. discarding", n)
	}
	return decodeHeader(event)
}

func decodeHeader(event *Event) error {
	i := 0
	event.packetType = (event.header[i] & 0xF0) >> 4
	event.flags = event.header[i] & 0x0F
	i++
	event.remainingLength = uint8(event.header[i])
	return nil
}

func readRemainingBytes(conn net.Conn, event *Event) error {
	event.remainingBytes = make([]byte, event.remainingLength)
	n, err := conn.Read(event.remainingBytes)
	if err != nil {
		return err
	}
	if n < int(event.remainingLength) {
		return fmt.Errorf("remainingBytes: read %d bytes, expected %d bytes.. discarding", n, event.remainingLength)
	}
	return nil
}

func manageEvent(e chan<- Event, connStatus *ConnStatus, event *Event) error {
	switch event.packetType {
	case PACKET_TYPE_CONNECT:
		return connectReq(e, connStatus, event)
	case PACKET_TYPE_PUBLISH:
		return publishReq(e, connStatus, event)
	case PACKET_TYPE_SUBSCRIBE:
		return subscribeReq(e, connStatus, event)
	case PACKET_TYPE_DISCONNECT:
		event.eventType = 100
		event.clientId = connStatus.clientId
		e <- *event
		return nil
	default:
		return fmt.Errorf("Unknown message type %d", event.packetType)
	}
}

func connectReq(e chan<- Event, connStatus *ConnStatus, event *Event) error {
	event.eventType = EVENT_CONNECT
	i := 0
	pl := Read2BytesInt(event.remainingBytes, i)
	i = i + 2
	// fmt.Println("protocolName", string(event.remainingBytes[i:i+pl]))
	i = i + pl
	v := event.remainingBytes[i]
	// fmt.Println("protocolVersion", v)
	connStatus.protocolVersion = v
	i++
	if int(v) < MINIMUM_SUPPORTED_PROTOCOL {
		fmt.Println("unsupported protocol version err", v)
		event.err = CONNECT_UNSUPPORTED_PROTOCOL_VERSION
		e <- *event
		return nil
	}
	connStatus.connectFlags = event.remainingBytes[i]
	i++
	ka := event.remainingBytes[i : i+2]
	fmt.Println("clean session", connStatus.cleanSession())

	// fmt.Println("keepAlive", Read2BytesInt(ka, 0))
	connStatus.keepAlive = ka
	i = i + 2
	cil := Read2BytesInt(event.remainingBytes, i)
	i = i + 2
	clientId := string(event.remainingBytes[i : i+cil])
	// fmt.Println("clientId", clientId)
	connStatus.clientId = clientId
	event.clientId = clientId
	e <- *event
	return nil
}

func publishReq(e chan<- Event, connStatus *ConnStatus, event *Event) error {
	event.eventType = EVENT_PUBLISH
	i := 0
	tl := Read2BytesInt(event.remainingBytes, i)
	i = i + 2
	topic := string(event.remainingBytes[i : i+tl])
	event.topic = topic
	e <- *event
	return nil
}

func subscribeReq(e chan<- Event, connStatus *ConnStatus, event *Event) error {
	event.eventType = EVENT_SUBSCRIBED
	i := 0
	pi := Read2BytesInt(event.remainingBytes, i)
	event.packetIdentifier = pi
	i = i + 2
	if connStatus.protocolVersion == 5 {
		pl := int(event.remainingBytes[i])
		fmt.Println("property length", pl)
		if pl > 0 {
			props := event.remainingBytes[i : i+pl]
			fmt.Println("properties", props)
			si, _, err := ReadVarInt(props)
			if err != nil {
				fmt.Println("err decoding sub ident", err)
				return nil
			}
			fmt.Println("subscription identifier", si)
			i = i + pl
		}
		i++
	}
	subs := make([]string, 10)
	j := 0
	for {
		var subevent Event
		subevent.eventType = EVENT_SUBSCRIPTION
		sl := Read2BytesInt(event.remainingBytes, i)
		i = i + 2
		subs[j] = string(event.remainingBytes[i : i+sl])
		subevent.clientId = connStatus.clientId
		subevent.topic = subs[j]
		e <- subevent
		i = i + sl
		if i >= len(event.remainingBytes)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	event.subscribedCount = j
	e <- *event
	return nil
}

func Suback(packetIdentifier int, subscribed int) Packet {
	p := make(Packet, 4+subscribed)
	p[0] = uint8(9) << 4
	p[1] = uint8(2 + subscribed)
	p[2] = byte(packetIdentifier & 0xFF00 >> 8)
	p[3] = byte(packetIdentifier & 0x00FF)
	return p
}

func Connack(reasonCode uint8) Packet {
	p := make(Packet, 5)
	p[0] = uint8(2) << 4
	p[1] = uint8(3)
	p[2] = uint8(0)
	p[3] = reasonCode
	return p
}
