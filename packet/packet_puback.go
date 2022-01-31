package packet

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
)

func Puback(packetIdentifier int, ReasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBACK) << 4
	// var header
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	p.remainingBytes = append(p.remainingBytes, ReasonCode)
	if protocolVersion >= conf.MQTT_V5 {
		// TODO: encode properties
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingLength = len(p.remainingBytes)
	return p
}

func (p *Packet) pubackReq() int {
	p.Event = EVENT_PUBACKED
	i := 2 // expect packet identifier in first 2 bytes
	if i < len(p.remainingBytes) {
		p.ReasonCode = p.remainingBytes[i]
	}
	if p.Session.GetProtocolVersion() >= conf.MQTT_V5 {
		_, err := p.parseProperties(i)
		if err != 0 {
			log.Error().Msgf("err reading properties %d", err)
			return err
		}
	}
	return 0
}
