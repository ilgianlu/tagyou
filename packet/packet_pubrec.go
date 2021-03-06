package packet

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
)

func Pubrec(packetIdentifier int, ReasonCode uint8, protocolVersion uint8) Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PUBREC) << 4
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

func (p *Packet) pubrecReq() int {
	p.Event = EVENT_PUBRECED
	i := 2 // 2 bytes for packet identifier
	if i < len(p.remainingBytes) {
		p.ReasonCode = p.remainingBytes[i]
	}
	if p.Session.ProtocolVersion >= conf.MQTT_V5 {
		_, err := p.parseProperties(i)
		if err != 0 {
			log.Error().Msgf("err reading properties %d", err)
			return err
		}
	}
	return 0
}
