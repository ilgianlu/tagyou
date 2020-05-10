package mqtt

import (
	"log"
	"time"

	"github.com/ilgianlu/tagyou/model"
)

func rangePackets(packets <-chan Packet, events chan<- Event, session *model.Session) {
	for p := range packets {
		// log.Println(p)
		switch p.PacketType() {
		case PACKET_TYPE_CONNECT:
			connectReq(p, events, session)
		case PACKET_TYPE_DISCONNECT:
			disconnectReq(p, events, session)
		case PACKET_TYPE_PUBLISH:
			publishReq(p, events, session)
		case PACKET_TYPE_PUBACK:
			pubackReq(p, events, session)
		case PACKET_TYPE_PUBREC:
			pubrecReq(p, events, session)
		case PACKET_TYPE_PUBREL:
			pubrelReq(p, events, session)
		case PACKET_TYPE_PUBCOMP:
			pubcompReq(p, events, session)
		case PACKET_TYPE_SUBSCRIBE:
			subscribeReq(p, events, session)
		case PACKET_TYPE_UNSUBSCRIBE:
			unsubscribeReq(p, events, session)
		case PACKET_TYPE_PINGREQ:
			pingReq(events, session)
		default:
			var event Event
			event.eventType = EVENT_PACKET_ERR
			log.Printf("Unknown packet type %d\n", p.PacketType())
			events <- event
		}
	}
}

func connectReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_CONNECT
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	// log.Printf("%d bytes, protocolName %s\n", pl, string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	// log.Println("protocolVersion", v)
	session.ProtocolVersion = v
	i++
	if int(v) < MINIMUM_SUPPORTED_PROTOCOL {
		log.Println("unsupported protocol version err", v)
		event.err = UNSUPPORTED_PROTOCOL_VERSION
		events <- event
	}
	session.ConnectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	session.KeepAlive = Read2BytesInt(ka, 0)
	// log.Println("keepAlive", Read2BytesInt(ka, 0))
	i = i + 2
	cil := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	event.clientId = string(p.remainingBytes[i : i+cil])
	session.ClientId = event.clientId
	log.Printf("%d bytes, clientId %s\n", cil, event.clientId)
	i = i + cil
	if session.WillFlag() {
		// read will topic
		wtl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		session.WillTopic = string(p.remainingBytes[i : i+wtl])
		i = i + wtl
		// will message
		wml := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		session.WillMessage = p.remainingBytes[i : i+wml]
		log.Printf("will topic \"%s\"\nwith message \"%s\"\n", session.WillTopic, session.WillMessage)
		i = i + wml
	}
	if session.HaveUser() {
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
	event.session = session
	events <- event
}

func disconnectReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_DISCONNECT
	event.clientId = session.ClientId
	event.session = session
	events <- event
}

func publishReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PUBLISH
	event.clientId = session.ClientId
	event.session = session
	event.published.dup = p.Dup()
	event.published.qos = p.QoS()
	event.published.retain = p.Retain()
	i := 0
	tl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	topic := string(p.remainingBytes[i : i+tl])
	event.published.topic = topic
	i = i + tl
	if event.published.qos != 0 {
		pi := Read2BytesInt(p.remainingBytes, i)
		p.packetIdentifier = pi
		i = i + 2
	}
	p.applicationMessage = i
	event.packet = p
	events <- event
}

func pubackReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PUBACKED
	event.clientId = session.ClientId
	event.session = session
	i := 0
	p.packetIdentifier = Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	if i < len(p.remainingBytes) {
		p.reasonCode = p.remainingBytes[i]
	}
	if session.ProtocolVersion == MQTT_V5 {
		// i++
		// read properties
	}
	event.packet = p
	events <- event
}

func pubrelReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PUBRELED
	event.clientId = session.ClientId
	event.session = session
	i := 0
	p.packetIdentifier = Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	if i < len(p.remainingBytes) {
		p.reasonCode = p.remainingBytes[i]
	}
	if session.ProtocolVersion == MQTT_V5 {
		// i++
		// read properties
	}
	event.packet = p
	events <- event
}

func pubrecReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PUBRECED
	event.clientId = session.ClientId
	event.session = session
	i := 0
	p.packetIdentifier = Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	if i < len(p.remainingBytes) {
		p.reasonCode = p.remainingBytes[i]
	}
	if session.ProtocolVersion == MQTT_V5 {
		// i++
		// read properties
	}
	event.packet = p
	events <- event
}

func pubcompReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PUBCOMPED
	event.clientId = session.ClientId
	event.session = session
	i := 0
	p.packetIdentifier = Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	if i < len(p.remainingBytes) {
		p.reasonCode = p.remainingBytes[i]
	}
	if session.ProtocolVersion == MQTT_V5 {
		// i++
		// read properties
	}
	event.packet = p
	events <- event
}

func subscribeReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_SUBSCRIBED
	event.clientId = session.ClientId
	event.session = session
	i := 0
	pi := Read2BytesInt(p.remainingBytes, i)
	p.packetIdentifier = pi
	i = i + 2
	if session.ProtocolVersion == 5 {
		pl := int(p.remainingBytes[i])
		log.Println("property length", pl)
		if pl > 0 {
			props := p.remainingBytes[i : i+pl]
			log.Println("properties", props)
			si, _, err := ReadVarInt(props)
			if err != nil {
				log.Println("err decoding sub ident", err)
				return
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
		subevent.session = session
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		s := string(p.remainingBytes[i : i+sl])
		subevent.subscription.ClientId = session.ClientId
		subevent.subscription.Topic = s
		i = i + sl
		if p.remainingBytes[i]&0x12 != 0 {
			log.Println("ignore this subscription & stop")
			break
		}
		subevent.subscription.RetainHandling = p.remainingBytes[i] & 0x30 >> 4
		subevent.subscription.RetainAsPublished = p.remainingBytes[i] & 0x08 >> 3
		subevent.subscription.NoLocal = p.remainingBytes[i] & 0x04 >> 2
		subevent.subscription.QoS = p.remainingBytes[i] & 0x03
		subevent.subscription.Enabled = true
		subevent.subscription.CreatedAt = time.Now()
		events <- subevent
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
	event.packet = p
	events <- event
}

func unsubscribeReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_UNSUBSCRIBED
	event.clientId = session.ClientId
	event.session = session
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
		unsubevent.clientId = session.ClientId
		unsubevent.topic = unsubs[j]
		events <- unsubevent
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
	events <- event
}

func pingReq(events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PING
	event.clientId = session.ClientId
	event.session = session
	events <- event
}

func Suback(packetIdentifier int, subscribed int, qosAccepted uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_SUBACK) << 4
	p.remainingLength = 2 + subscribed
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	subRes := make([]byte, subscribed)
	for i := 0; i < subscribed; i++ {
		subRes[i] = qosAccepted
	}
	p.remainingBytes = append(p.remainingBytes, subRes...)
	return p
}

func Unsuback(packetIdentifier int, unsubscribed int) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_UNSUBACK) << 4
	p.remainingLength = 2 + unsubscribed
	p.remainingBytes = make([]byte, 2+unsubscribed)
	p.remainingBytes[0] = byte(packetIdentifier & 0xFF00 >> 8)
	p.remainingBytes[1] = byte(packetIdentifier & 0x00FF)
	return p
}

func Connack(sessionPresent bool, reasonCode uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_CONNACK) << 4
	p.remainingLength = 2
	p.remainingLengthBytes = 1
	p.remainingBytes = make([]byte, 2)
	if sessionPresent {
		p.remainingBytes[0] = 1
	} else {
		p.remainingBytes[0] = 0
	}
	p.remainingBytes[1] = reasonCode
	return p
}

func PingResp() Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PINGRES) << 4
	p.remainingLength = 0
	return p
}

func Publish(qos uint8, retain bool, topic string, payload []byte) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBLISH) << 4
	p.header = p.header | qos<<1
	if retain {
		p.header = p.header | 1
	}
	// write var int length 2 + len(topic) + len(payload)
	p.remainingLength = 2 + len(topic) + len(payload)

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

func Puback(packetIdentifier int, reasonCode uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBACK) << 4
	if reasonCode == 0 {
		p.remainingLength = 2
		p.remainingBytes = Write2BytesInt(packetIdentifier)
	} else {
		p.remainingLength = 3
		p.remainingBytes = Write2BytesInt(packetIdentifier)
		p.remainingBytes = append(p.remainingBytes, reasonCode)
	}
	return p
}

func Pubrel(packetIdentifier int, reasonCode uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBREL) << 4
	if reasonCode == 0 {
		p.remainingLength = 2
		p.remainingBytes = Write2BytesInt(packetIdentifier)
	} else {
		p.remainingLength = 3
		p.remainingBytes = Write2BytesInt(packetIdentifier)
		p.remainingBytes = append(p.remainingBytes, reasonCode)
	}
	return p
}

func Pubrec(packetIdentifier int, reasonCode uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBREC) << 4
	if reasonCode == 0 {
		p.remainingLength = 2
		p.remainingBytes = Write2BytesInt(packetIdentifier)
	} else {
		p.remainingLength = 3
		p.remainingBytes = Write2BytesInt(packetIdentifier)
		p.remainingBytes = append(p.remainingBytes, reasonCode)
	}
	return p
}

func Pubcomp(packetIdentifier int, reasonCode uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBCOMP) << 4
	if reasonCode == 0 {
		p.remainingLength = 2
		p.remainingBytes = Write2BytesInt(packetIdentifier)
	} else {
		p.remainingLength = 3
		p.remainingBytes = Write2BytesInt(packetIdentifier)
		p.remainingBytes = append(p.remainingBytes, reasonCode)
	}
	return p
}
