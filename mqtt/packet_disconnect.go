package mqtt

import (
	"log"
)

func (p *Packet) disconnectReq() int {
	p.event = EVENT_DISCONNECT
	if len(p.remainingBytes) > 0 {
		i := 0
		p.reasonCode = p.remainingBytes[i]
		if p.session.ProtocolVersion >= MQTT_V5 {
			_, err := p.parseProperties(i)
			if err != 0 {
				log.Println("err reading properties", err)
				return err
			}
		}
	} else {
		p.reasonCode = 0
	}
	return 0
}
