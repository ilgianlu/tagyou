package mqtt

import (
	"fmt"
)

type Packet struct {
	header               byte
	remainingLength      int
	remainingLengthBytes int
	remainingBytes       []byte
	applicationMessage   int
	packetIdentifier     int
	subscribedCount      int
	err                  error
}

func (p *Packet) PacketType() byte {
	return (p.header & 0x00F0) >> 4
}

func (p *Packet) Flags() byte {
	return p.header & 0x0F
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
	i := 0
	p.header = buff[i]
	i++
	if p.checkHeader() {
		rl, k, err := ReadVarInt(buff[i:])
		p.remainingLengthBytes = k
		if err != nil {
			return p, fmt.Errorf("header: malformed remainingLength format: %s\n", err)
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

func (p *Packet) toByteSlice() []byte {
	res := make([]byte, 1)
	res[0] = p.header
	res = append(res, WriteVarInt(p.remainingLength)...)
	res = append(res, p.remainingBytes...)
	return res
}
