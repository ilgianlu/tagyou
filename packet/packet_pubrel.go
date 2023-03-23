package packet

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
)

func Pubrel(packetIdentifier int, ReasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBREL) << 4
	p.remainingBytes = Write2BytesInt(packetIdentifier)
	if ReasonCode != 0 {
		p.remainingBytes = append(p.remainingBytes, ReasonCode)
	}
	if protocolVersion >= conf.MQTT_V5 {
		// properties
		p.remainingBytes = append(p.remainingBytes, 0)
	}
	p.remainingLength = len(p.remainingBytes)
	return p
}

func (p *Packet) pubrelReq(protocolVersion uint8) int {
	p.Event = EVENT_PUBRELED
	i := 2 // 2 bytes for packet identifier
	if i < len(p.remainingBytes) {
		p.ReasonCode = p.remainingBytes[i]
	}
	if protocolVersion >= conf.MQTT_V5 {
		_, err := p.parseProperties(i)
		if err != 0 {
			log.Error().Msgf("err reading properties %d", err)
			return err
		}
	}
	return 0
}
