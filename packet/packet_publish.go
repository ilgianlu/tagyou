package packet

import (
	"log/slog"
	"math/rand"

	"github.com/ilgianlu/tagyou/conf"
)

func Publish(protocolVersion uint8, qos uint8, retain bool, topic string, packetIdentifier int, payload []byte) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBLISH) << 4
	p.header = p.header | qos<<1
	if retain {
		p.header = p.header | 1
	}
	// variable header
	// write topic length
	p.remainingBytes = Write2BytesInt(len(topic))
	// write topic string
	p.remainingBytes = append(p.remainingBytes, []byte(topic)...)
	// write packet identifier only if qos > 0
	if qos != 0 {
		p.remainingBytes = append(p.remainingBytes, Write2BytesInt(packetIdentifier)...)
	}
	if protocolVersion >= conf.MQTT_V5 {
		p.remainingBytes = append(p.remainingBytes, p.encodeProperties()...)
	}
	// write payload
	p.payloadOffset = len(p.remainingBytes)
	p.remainingBytes = append(p.remainingBytes, payload...)
	p.remainingLength = len(p.remainingBytes)
	return p
}

func (p *Packet) publishReq(protocolVersion uint8) int {
	p.Event = EVENT_PUBLISH
	i := 0
	tl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	// variable header
	p.Topic = string(p.remainingBytes[i : i+tl])
	i = i + tl
	if p.QoS() > 0 {
		i = i + 2 // + 2 for packet identifier
	}
	if protocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			slog.Error("err reading properties", "err", err)
			return err
		}
		i = i + pl
	}
	// payload
	p.payloadOffset = i
	return 0
}

func NewPacketIdentifier() int {
	return rand.Intn(65534)
}
