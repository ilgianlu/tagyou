package mqtt

import (
	"fmt"
	"log"
	"net"
	"time"
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
		log.Printf("error after first byte read: %s\n", err)
		return err
	}
	if bytesCount < 1 {
		return fmt.Errorf("header: read %d bytes, expected more than 1 byte.. discarding", bytesCount)
	}
	if bytesCount > PACKET_MAX_SIZE {
		log.Printf("oversize packet %d > %d, discarding...\n", bytesCount, PACKET_MAX_SIZE)
	}
	p.packetType = (buffer[0] & 0xF0) >> 4
	p.flags = buffer[0] & 0x0F
	remainingLength, k, err0 := ReadVarInt(buffer[1:])
	if err0 != nil {
		log.Printf("malformed remainingLength format: %s\n", err0)
		return err0
	}
	if remainingLength > PACKET_MAX_SIZE {
		log.Printf("oversize packet %d > %d, discarding...\n", remainingLength, PACKET_MAX_SIZE)
	}
	p.header = buffer[:1+k]
	p.remainingBytes = buffer[1+k : bytesCount]
	bytesCount = bytesCount - k
	for bytesCount < remainingLength {
		buffer = make([]byte, 1024)
		n, err := conn.Read(buffer)
		bytesCount = bytesCount + n
		if err != nil {
			log.Printf("error after %d byte read: %s\n", bytesCount, err)
			return err
		}
		p.remainingBytes = append(p.remainingBytes, buffer[:n]...)
	}
	log.Printf("read %d bytes packet\n", bytesCount)
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
	if p.flags != 0 {
		log.Println("malformed packet")
		event.err = MALFORMED_PACKET
		e <- event
		return nil
	}
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	// log.Printf("%d bytes, protocolName %s\n", pl, string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	// log.Println("protocolVersion", v)
	connection.protocolVersion = v
	i++
	if int(v) < MINIMUM_SUPPORTED_PROTOCOL {
		log.Println("unsupported protocol version err", v)
		event.err = UNSUPPORTED_PROTOCOL_VERSION
		e <- event
		return nil
	}
	connection.connectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	connection.keepAlive = Read2BytesInt(ka, 0)
	// log.Println("keepAlive", Read2BytesInt(ka, 0))
	i = i + 2
	cil := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	clientId := string(p.remainingBytes[i : i+cil])
	log.Printf("%d bytes, clientId %s\n", cil, string(p.remainingBytes[i:i+cil]))
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
		log.Printf("will topic \"%s\"\nwith message \"%s\"\n", connection.willTopic, connection.willMessage)
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
		log.Printf("username \"%s\"\nlogging with password \"%s\"\n", username, password)
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
	if p.flags != 2 {
		log.Println("malformed packet")
		event.err = MALFORMED_PACKET
		e <- event
		return nil
	}
	i := 0
	pi := Read2BytesInt(p.remainingBytes, i)
	p.packetIdentifier = pi
	i = i + 2
	if c.protocolVersion == 5 {
		pl := int(p.remainingBytes[i])
		log.Println("property length", pl)
		if pl > 0 {
			props := p.remainingBytes[i : i+pl]
			log.Println("properties", props)
			si, _, err := ReadVarInt(props)
			if err != nil {
				log.Println("err decoding sub ident", err)
				return nil
			}
			log.Println("subscription identifier", si)
			i = i + pl
		}
		i++
	}
	j := 0
	for {
		var subevent Event
		subevent.eventType = EVENT_SUBSCRIPTION
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		s := string(p.remainingBytes[i : i+sl])
		subevent.subscription.clientId = c.clientId
		subevent.subscription.topic = s
		i = i + sl
		if p.remainingBytes[i]&0x12 != 0 {
			log.Println("ignore this subscription & stop")
			break
		}
		subevent.subscription.RetainHandling = p.remainingBytes[i] & 0x30 >> 4
		subevent.subscription.RetainAsPublished = p.remainingBytes[i] & 0x08 >> 3
		subevent.subscription.NoLocal = p.remainingBytes[i] & 0x04 >> 2
		subevent.subscription.QoS = p.remainingBytes[i] & 0x03
		subevent.subscription.enabled = true
		subevent.subscription.createdAt = time.Now()
		e <- subevent
		i++
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > MAX_TOPIC_SINGLE_SUBSCRIBE {
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

func Publish(qos uint8, retain bool, topic string, payload []byte) Packet {
	var p Packet
	h := make([]byte, 1)
	h[0] = uint8(PACKET_TYPE_PUBLISH) << 4
	h[0] = h[0] | qos<<1
	if retain {
		h[0] = h[0] | 1
	}
	// write var int length 2 + len(topic) + len(payload)
	remainingLength := 2 + len(topic) + len(payload)
	h = append(h, WriteVarInt(remainingLength)...)
	p.header = h

	// write topic length
	rb := Write2BytesInt(len(topic))
	// write topic string
	rb = append(rb, []byte(topic)...)
	// write packet identifier only if qos > 0

	// write payload
	rb = append(rb, payload...)
	p.remainingBytes = rb
	return p
}
