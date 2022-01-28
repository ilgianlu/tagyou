package packet

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
)

const PACKET_TYPE_CONNECT = 1
const PACKET_TYPE_CONNACK = 2
const PACKET_TYPE_PUBLISH = 3
const PACKET_TYPE_PUBACK = 4
const PACKET_TYPE_PUBREC = 5
const PACKET_TYPE_PUBREL = 6
const PACKET_TYPE_PUBCOMP = 7
const PACKET_TYPE_SUBSCRIBE = 8
const PACKET_TYPE_SUBACK = 9
const PACKET_TYPE_UNSUBSCRIBE = 10
const PACKET_TYPE_UNSUBACK = 11
const PACKET_TYPE_PINGREQ = 12
const PACKET_TYPE_PINGRES = 13
const PACKET_TYPE_DISCONNECT = 14
const PACKET_MAX_SIZE = 65000

// connect response
const CONNECT_OK = 0
const UNSPECIFIED_ERROR = 0x80
const MALFORMED_PACKET = 0x81
const UNSUPPORTED_PROTOCOL_VERSION = 0x84
const SESSION_TAKEN_OVER = 0x8E

// publish ack in QoS 1
const PUBACK_SUCCESS = 0x00
const PUBACK_NO_MATCHING_SUBSCRIBERS = 0x10
const PUBACK_NOT_AUTHORIZED = 0x87

// publish in QoS 2
const PUBCOMP_SUCCESS = 0x00
const PUBREL_SUCCESS = 0x00
const PUBREC_SUCCESS = 0x00
const PUBREC_NOT_AUTHORIZED = 0x87

// event type
const EVENT_CONNECT = 0
const EVENT_PUBLISH = 2
const EVENT_PUBACKED = 3
const EVENT_PUBRECED = 4
const EVENT_PUBRELED = 5
const EVENT_PUBCOMPED = 6
const EVENT_SUBSCRIBED = 10
const EVENT_SUBSCRIPTION = 11
const EVENT_PING = 12
const EVENT_UNSUBSCRIBED = 20
const EVENT_UNSUBSCRIPTION = 21
const EVENT_DISCONNECT = 100
const EVENT_KEEPALIVE_TIMEOUT = 101
const EVENT_WILL_SEND = 102

const EVENT_PACKET_ERR = 1000

type Packet struct {
	// header
	header               byte
	remainingLength      int
	remainingLengthBytes int
	// packet remaining bytes
	remainingBytes []byte

	// variable header offset in remaining bytes
	// varHeaderOffset int always 0
	properties Properties
	// CONNACK, PUBACK, PUBREC, PUBREL, PUBCOMP, DISCONNECT
	ReasonCode uint8

	// payload
	payloadOffset  int
	willProperties Properties

	// metadata
	Session       *model.RunningSession
	Subscriptions []model.Subscription
	Topic         string
	Event         int
}

func (p *Packet) RemainingBytes() []byte {
	return p.remainingBytes
}

func (p *Packet) PacketType() byte {
	return (p.header & 0x00F0) >> 4
}

func (p *Packet) Flags() byte {
	return p.header & 0x0F
}

func (p *Packet) QoS() byte {
	if p.PacketType() == PACKET_TYPE_PUBLISH {
		return p.Flags() & 0x06 >> 1
	}
	return 0x10
}

func (p *Packet) Dup() bool {
	if p.PacketType() == PACKET_TYPE_PUBLISH {
		return (p.Flags() & 0x08 >> 3) == 1
	}
	return false
}

func (p *Packet) Retain() bool {
	if p.PacketType() == PACKET_TYPE_PUBLISH {
		return (p.Flags() & 0x01) == 1
	}
	return false
}

func (p *Packet) PublishTopic() []byte {
	if p.PacketType() == PACKET_TYPE_PUBLISH {
		tl := Read2BytesInt(p.remainingBytes, 0)
		return p.remainingBytes[:tl]
	}
	return []byte{}
}

func (p *Packet) PacketIdentifier() int {
	var offset int
	if p.PacketType() == PACKET_TYPE_PUBLISH {
		topicLength := Read2BytesInt(p.remainingBytes, 0)
		offset = 2 + topicLength
	}
	return Read2BytesInt(p.remainingBytes, offset)
}

func (p *Packet) PacketLength() int {
	return 1 + p.remainingLengthBytes + len(p.remainingBytes)
}

func (p *Packet) PacketComplete() bool {
	return len(p.remainingBytes) == p.remainingLength
}

func (p *Packet) missingBytes() int {
	return p.remainingLength - len(p.remainingBytes)
}

func (p *Packet) Payload() []byte {
	return p.remainingBytes[p.payloadOffset:]
}

func (p *Packet) ApplicationMessage() []byte {
	if p.PacketType() == PACKET_TYPE_PUBLISH && p.PacketComplete() {
		return p.remainingBytes[p.payloadOffset:]
	}
	return []byte{}
}

func (p *Packet) CompletePacket(buff []byte) int {
	if len(buff) >= p.missingBytes() {
		p.remainingBytes = append(p.remainingBytes, buff[:p.missingBytes()]...)
		return p.missingBytes()
	} else {
		p.remainingBytes = append(p.remainingBytes, buff...)
		return len(buff)
	}
}

func Start(buff []byte) (Packet, error) {
	var p Packet
	if len(buff) < 2 {
		return p, fmt.Errorf("Start: buffer too short")
	}
	i := 0
	p.header = buff[i]
	i++
	if p.checkHeader() {
		rl, k, err := ReadVarInt(buff[i:])
		p.remainingLengthBytes = k
		if err != nil {
			return p, fmt.Errorf("header: malformed remainingLength format: %s", err)
		}
		p.remainingLength = rl
		i = i + k
		p.CompletePacket(buff[i:])
		return p, nil
	} else {
		return p, fmt.Errorf("header: invalid %b", buff[0])
	}
}

func (p *Packet) checkHeader() bool {
	switch p.PacketType() {
	case PACKET_TYPE_CONNECT:
		if p.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_PUBLISH:
		if p.QoS() > 2 {
			return false
		}
		return true
	case PACKET_TYPE_PUBACK:
		if p.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_PUBREC:
		if p.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_PUBREL:
		if p.Flags() != 2 {
			return false
		}
		return true
	case PACKET_TYPE_PUBCOMP:
		if p.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_SUBSCRIBE:
		if p.Flags() != 2 {
			return false
		}
		return true
	case PACKET_TYPE_UNSUBSCRIBE:
		return true
	case PACKET_TYPE_PINGREQ:
		return true
	case PACKET_TYPE_DISCONNECT:
		return true
	default:
		return false
	}
}

func (p *Packet) Parse() int {
	switch p.PacketType() {
	case PACKET_TYPE_CONNECT:
		return p.connectReq()
	case PACKET_TYPE_DISCONNECT:
		return p.disconnectReq()
	case PACKET_TYPE_PUBLISH:
		return p.publishReq()
	case PACKET_TYPE_PUBACK:
		return p.pubackReq()
	case PACKET_TYPE_PUBREC:
		return p.pubrecReq()
	case PACKET_TYPE_PUBREL:
		return p.pubrelReq()
	case PACKET_TYPE_PUBCOMP:
		return p.pubcompReq()
	case PACKET_TYPE_SUBSCRIBE:
		return p.subscribeReq()
	case PACKET_TYPE_UNSUBSCRIBE:
		return p.unsubscribeReq()
	case PACKET_TYPE_PINGREQ:
		return p.pingReq()
	default:
		log.Printf("Unknown packet type %d\n", p.PacketType())
		return MALFORMED_PACKET
	}
}

func ReadFromByteSlice(buff []byte) ([]byte, error) {
	if len(buff) < 2 {
		return nil, fmt.Errorf("header: not enough bytes in buffer")
	}
	i := 1
	rl, k, err := ReadVarInt(buff[i:])
	if err != nil {
		log.Printf("header: malformed remainingLength format: %s", err)
		return nil, err
	}
	i = i + k
	if len(buff[i:]) < rl {
		log.Debug().Msg("remaining bytes: not enough bytes in buffer")
		return nil, fmt.Errorf("remaining bytes: not enough bytes in buffer")
	}
	i = i + rl
	return buff[:i], nil
}

func (p *Packet) ToByteSlice() []byte {
	res := make([]byte, 1)
	res[0] = p.header
	encodedLength := WriteVarInt(p.remainingLength)
	res = append(res, encodedLength...)
	res = append(res, p.remainingBytes...)
	return res
}
