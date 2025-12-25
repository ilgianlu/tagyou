package packet

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/format"
)

func Puback(packetIdentifier int, ReasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = header(uint8(PACKET_TYPE_PUBACK) << 4)
	// var header
	p.remainingBytes = format.Write2BytesInt(packetIdentifier)
	p.remainingBytes = append(p.remainingBytes, ReasonCode)
	if protocolVersion >= conf.MQTT_V5 {
		// TODO: encode properties
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingLength = len(p.remainingBytes)
	return p
}

func (p *Packet) pubackReq(protocolVersion uint8) int {
	packetIdentifier, _ := format.Read2BytesInt(p.remainingBytes, 0)
	p.packetIdentifier = packetIdentifier
	i := 2 // expect packet identifier in first 2 bytes
	if i < len(p.remainingBytes) {
		p.ReasonCode = p.remainingBytes[i]
	}
	if protocolVersion >= conf.MQTT_V5 {
		_, err := p.parseProperties(i)
		if err != 0 {
			slog.Error("err reading properties", "err", err)
			return err
		}
	}
	return 0
}
