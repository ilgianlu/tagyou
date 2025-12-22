package packet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/ilgianlu/tagyou/format"
	"github.com/ilgianlu/tagyou/model"
)

const (
	PACKET_TYPE_CONNECT     = 1
	PACKET_TYPE_CONNACK     = 2
	PACKET_TYPE_PUBLISH     = 3
	PACKET_TYPE_PUBACK      = 4
	PACKET_TYPE_PUBREC      = 5
	PACKET_TYPE_PUBREL      = 6
	PACKET_TYPE_PUBCOMP     = 7
	PACKET_TYPE_SUBSCRIBE   = 8
	PACKET_TYPE_SUBACK      = 9
	PACKET_TYPE_UNSUBSCRIBE = 10
	PACKET_TYPE_UNSUBACK    = 11
	PACKET_TYPE_PINGREQ     = 12
	PACKET_TYPE_PINGRES     = 13
	PACKET_TYPE_DISCONNECT  = 14
	PACKET_MAX_SIZE         = 65000
)

// connect response
const (
	CONNECT_OK                   = 0
	UNSPECIFIED_ERROR            = 0x80
	MALFORMED_PACKET             = 0x81
	UNSUPPORTED_PROTOCOL_VERSION = 0x84
	SESSION_TAKEN_OVER           = 0x8E
)

// publish ack in QoS 1
const (
	PUBACK_SUCCESS                 = 0x00
	PUBACK_NO_MATCHING_SUBSCRIBERS = 0x10
	PUBACK_NOT_AUTHORIZED          = 0x87
)

// publish in QoS 2
const (
	PUBCOMP_SUCCESS       = 0x00
	PUBREL_SUCCESS        = 0x00
	PUBREC_SUCCESS        = 0x00
	PUBREC_NOT_AUTHORIZED = 0x87
)

type Packet struct {
	// header
	header               header
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
	Subscriptions []model.Subscription
	PublishTopic  string
}

func (p *Packet) ReadHeader(r *bufio.Reader) error {
	h, err := r.ReadByte()
	if err != nil {
		return errors.New("[PACKET] error on client side")
	}
	if CheckHeader(h) {
		p.header = header(h)
		return nil
	}
	return errors.New("invalid header")
}

func (p *Packet) ReadRemainingLength(r *bufio.Reader) error {
	v, n, err := format.ReadVarInt(r)
	if err != nil {
		return err
	}
	p.remainingLength = v
	p.remainingLengthBytes = n
	return nil
}

func (p *Packet) ReadRemainingBytes(r *bufio.Reader) error {
	buf := make([]byte, p.remainingLength)
	n, err := io.ReadFull(r, buf)
	slog.Debug("[PACKET] read %d bytes, expected remaining length %d", n, p.remainingLength)
	if err != nil || n < p.remainingLength {
		return errors.New("[PACKET] fewer bytes read")
	}
	p.remainingBytes = buf
	return nil
}

func (p *Packet) RemainingBytes() []byte {
	return p.remainingBytes
}

func (p *Packet) GetPublishTopic() string {
	return p.PublishTopic
}

func (p *Packet) GetReasonCode() uint8 {
	return p.ReasonCode
}

func (p *Packet) PacketIdentifier() (int, error) {
	var offset int
	if p.header.PacketType() == PACKET_TYPE_PUBLISH {
		topicLength, err := format.Read2BytesInt(p.remainingBytes, 0)
		if err != nil {
			return 0, err
		}
		offset = 2 + topicLength
	}
	return format.Read2BytesInt(p.remainingBytes, offset)
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
	if p.header.PacketType() == PACKET_TYPE_PUBLISH && p.PacketComplete() {
		return p.remainingBytes[p.payloadOffset:]
	}
	return []byte{}
}

func (p *Packet) GetSubscriptions() []model.Subscription {
	return p.Subscriptions
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

func (p *Packet) Parse(session *model.RunningSession) int {
	switch p.header.PacketType() {
	case PACKET_TYPE_CONNECT:
		return p.connectReq(session)
	case PACKET_TYPE_DISCONNECT:
		return p.disconnectReq(session.GetProtocolVersion())
	case PACKET_TYPE_PUBLISH:
		return p.publishReq(session.GetProtocolVersion())
	case PACKET_TYPE_PUBACK:
		return p.pubackReq(session.GetProtocolVersion())
	case PACKET_TYPE_PUBREC:
		return p.pubrecReq(session.GetProtocolVersion())
	case PACKET_TYPE_PUBREL:
		return p.pubrelReq(session.GetProtocolVersion())
	case PACKET_TYPE_PUBCOMP:
		return p.pubcompReq(session.GetProtocolVersion())
	case PACKET_TYPE_SUBSCRIBE:
		return p.subscribeReq(session)
	case PACKET_TYPE_UNSUBSCRIBE:
		return p.unsubscribeReq(session)
	case PACKET_TYPE_PINGREQ:
		return p.pingReq()
	default:
		slog.Warn("[MQTT] Unknown packet type", "packet-type", p.header.PacketType())
		return MALFORMED_PACKET
	}
}

func ReadFromByteSlice(buff []byte) ([]byte, error) {
	if len(buff) < 2 {
		return nil, fmt.Errorf("header: not enough bytes in buffer")
	}
	i := 1
	rl, k, err := format.ReadVarIntFromBytes(buff[i:])
	if err != nil {
		slog.Error("header: malformed remainingLength format", "err", err)
		return nil, err
	}
	i = i + k
	if len(buff[i:]) < rl {
		return nil, fmt.Errorf("remaining bytes: not enough bytes in buffer")
	}
	i = i + rl
	return buff[:i], nil
}

func (p *Packet) ToByteSlice() ([]byte, error) {
	res := make([]byte, 1)
	res[0] = byte(p.header)
	encodedLength, err := format.WriteVarInt(p.remainingLength)
	if err != nil {
		return []byte{}, err
	}
	res = append(res, encodedLength...)
	res = append(res, p.remainingBytes...)
	return res, nil
}
