package mqtt

import (
	"fmt"
	"net"
)

type Packet struct {
	header           []byte
	packetType       uint8
	flags            uint8
	remainingBytes   []byte
	packetIdentifier int
	subscribedCount  int
}

func ReadPacket(conn net.Conn, connection *Connection, e chan<- Event) (Packet, error) {
	var p Packet
	err0 := p.read(conn)
	if err0 != nil {
		return p, err0
	}

	err1 := p.emit(connection, e)
	if err1 != nil {
		return p, err1
	}
	return p, nil
}

func (p *Packet) read(conn net.Conn) error {
	buffer := make([]byte, 1024)
	bytesCount, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("error after first byte read: %s\n", err)
		return err
	}
	if bytesCount < 1 {
		return fmt.Errorf("header: read %d bytes, expected more than 1 byte.. discarding", bytesCount)
	}
	if bytesCount > PACKET_MAX_SIZE {
		fmt.Printf("oversize packet %d > %d, discarding...\n", bytesCount, PACKET_MAX_SIZE)
	}
	p.packetType = (buffer[0] & 0xF0) >> 4
	p.flags = buffer[0] & 0x0F
	remainingLength, k, err0 := ReadVarInt(buffer[1:])
	if err0 != nil {
		fmt.Printf("malformed remainingLength format: %s\n", err0)
		return err0
	}
	if remainingLength > PACKET_MAX_SIZE {
		fmt.Printf("oversize packet %d > %d, discarding...\n", remainingLength, PACKET_MAX_SIZE)
	}
	p.header = buffer[:1+k]
	p.remainingBytes = buffer[1+k : bytesCount]
	bytesCount = bytesCount - k
	for bytesCount < remainingLength {
		buffer = make([]byte, 1024)
		n, err := conn.Read(buffer)
		bytesCount = bytesCount + n
		if err != nil {
			fmt.Printf("error after %d byte read: %s\n", bytesCount, err)
			return err
		}
		p.remainingBytes = append(p.remainingBytes, buffer[:n]...)
	}
	fmt.Printf("read %d bytes packet\n", bytesCount)
	return nil
}

func (p *Packet) emit(connection *Connection, e chan<- Event) error {
	switch p.packetType {
	case PACKET_TYPE_CONNECT:
		return p.connectReq(e, connection)
	case PACKET_TYPE_PUBLISH:
		return p.publishReq(e)
	case PACKET_TYPE_SUBSCRIBE:
		return p.subscribeReq(e, connection)
	case PACKET_TYPE_UNSUBSCRIBE:
		return p.unsubscribeReq(e, connection)
	case PACKET_TYPE_PINGREQ:
		return p.pingReq(e, connection)
	case PACKET_TYPE_DISCONNECT:
		return p.disconnectReq(e, connection)
	default:
		return fmt.Errorf("Unknown packet type %d", p.packetType)
	}
}

func (p *Packet) connectReq(e chan<- Event, connection *Connection) error {
	var event Event
	event.eventType = EVENT_CONNECT
	event.connection = connection
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	// fmt.Printf("%d bytes, protocolName %s\n", pl, string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	// fmt.Println("protocolVersion", v)
	connection.protocolVersion = v
	i++
	if int(v) < MINIMUM_SUPPORTED_PROTOCOL {
		fmt.Println("unsupported protocol version err", v)
		event.err = CONNECT_UNSUPPORTED_PROTOCOL_VERSION
		e <- event
		return nil
	}
	connection.connectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	connection.keepAlive = Read2BytesInt(ka, 0)
	// fmt.Println("keepAlive", Read2BytesInt(ka, 0))
	i = i + 2
	cil := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	clientId := string(p.remainingBytes[i : i+cil])
	fmt.Printf("%d bytes, clientId %s\n", cil, string(p.remainingBytes[i:i+cil]))
	event.clientId = clientId
	connection.clientId = clientId
	i = i + cil
	if connection.willFlag() {
		// read will topic
		wtl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		connection.willTopic = string(p.remainingBytes[i : i+wtl])
		i = i + wtl
		// will message
		wml := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		connection.willMessage = p.remainingBytes[i : i+wml]
		fmt.Printf("will topic \"%s\"\nwith message \"%s\"\n", connection.willTopic, connection.willMessage)
		i = i + wml
	}
	if connection.haveUser() {
		// read username
		unl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		username := string(p.remainingBytes[i : i+unl])
		i = i + unl
		// read password
		pwdl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		password := string(p.remainingBytes[i : i+pwdl])
		fmt.Printf("username \"%s\"\nlogging with password \"%s\"\n", username, password)
	}
	e <- event
	return nil
}

func (p *Packet) disconnectReq(e chan<- Event, connection *Connection) error {
	var event Event
	event.eventType = EVENT_DISCONNECT
	event.clientId = connection.clientId
	event.connection = connection
	e <- event
	return nil
}

func (p *Packet) publishReq(e chan<- Event) error {
	var event Event
	event.eventType = EVENT_PUBLISH
	i := 0
	tl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	topic := string(p.remainingBytes[i : i+tl])
	event.topic = topic
	event.packet = p
	e <- event
	return nil
}

func (p *Packet) subscribeReq(e chan<- Event, c *Connection) error {
	var event Event
	event.eventType = EVENT_SUBSCRIBED
	event.clientId = c.clientId
	event.connection = c
	event.packet = p
	i := 0
	pi := Read2BytesInt(p.remainingBytes, i)
	p.packetIdentifier = pi
	i = i + 2
	if c.protocolVersion == 5 {
		pl := int(p.remainingBytes[i])
		fmt.Println("property length", pl)
		if pl > 0 {
			props := p.remainingBytes[i : i+pl]
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
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		subs[j] = string(p.remainingBytes[i : i+sl])
		subevent.clientId = c.clientId
		subevent.topic = subs[j]
		e <- subevent
		i = i + sl
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	p.subscribedCount = j
	e <- event
	return nil
}

func (p *Packet) unsubscribeReq(e chan<- Event, c *Connection) error {
	var event Event
	event.eventType = EVENT_UNSUBSCRIBED
	event.clientId = c.clientId
	event.connection = c
	event.packet = p
	i := 0
	pi := Read2BytesInt(p.remainingBytes, i)
	p.packetIdentifier = pi
	i = i + 2
	unsubs := make([]string, 10)
	j := 0
	for {
		var unsubevent Event
		unsubevent.eventType = EVENT_UNSUBSCRIPTION
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		unsubs[j] = string(p.remainingBytes[i : i+sl])
		unsubevent.clientId = c.clientId
		unsubevent.topic = unsubs[j]
		e <- unsubevent
		i = i + sl
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	p.subscribedCount = j
	e <- event
	return nil
}

func (p *Packet) pingReq(e chan<- Event, c *Connection) error {
	var event Event
	event.eventType = EVENT_PING
	event.clientId = c.clientId
	event.connection = c
	e <- event
	return nil
}

func Suback(packetIdentifier int, subscribed int) []byte {
	p := make([]byte, 4+subscribed)
	p[0] = uint8(PACKET_TYPE_SUBACK) << 4
	p[1] = uint8(2 + subscribed)
	p[2] = byte(packetIdentifier & 0xFF00 >> 8)
	p[3] = byte(packetIdentifier & 0x00FF)
	return p
}

func Unsuback(packetIdentifier int, unsubscribed int) []byte {
	p := make([]byte, 4+unsubscribed)
	p[0] = uint8(PACKET_TYPE_UNSUBACK) << 4
	p[1] = uint8(2 + unsubscribed)
	p[2] = byte(packetIdentifier & 0xFF00 >> 8)
	p[3] = byte(packetIdentifier & 0x00FF)
	return p
}

func Connack(reasonCode uint8) []byte {
	p := make([]byte, 4)
	p[0] = uint8(PACKET_TYPE_CONNACK) << 4
	p[1] = uint8(2)
	p[2] = uint8(0)
	p[3] = reasonCode
	return p
}

func PingResp() []byte {
	p := make([]byte, 2)
	p[0] = uint8(PACKET_TYPE_PINGRES) << 4
	p[1] = uint8(0)
	return p
}
