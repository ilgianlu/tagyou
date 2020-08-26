package packet

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
)

func (p *Packet) disconnectReq() int {
	p.Event = EVENT_DISCONNECT
	if len(p.remainingBytes) > 0 {
		i := 0
		p.ReasonCode = p.remainingBytes[i]
		if p.Session.ProtocolVersion >= conf.MQTT_V5 {
			_, err := p.parseProperties(i)
			if err != 0 {
				log.Println("err reading properties", err)
				return err
			}
		}
	} else {
		p.ReasonCode = 0
	}
	return 0
}
