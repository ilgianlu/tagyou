package packet

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
)

func (p *Packet) disconnectReq() int {
	p.Event = EVENT_DISCONNECT
	if len(p.remainingBytes) > 0 {
		i := 0
		p.ReasonCode = p.remainingBytes[i]
		if p.Session.GetProtocolVersion() >= conf.MQTT_V5 {
			_, err := p.parseProperties(i)
			if err != 0 {
				log.Error().Msgf("err reading properties %d", err)
				return err
			}
		}
	} else {
		p.ReasonCode = 0
	}
	return 0
}
