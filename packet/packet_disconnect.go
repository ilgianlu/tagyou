package packet

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
)

func (p *Packet) disconnectReq(protocolVersion uint8) int {
	p.Event = EVENT_DISCONNECT
	if len(p.remainingBytes) > 0 {
		i := 0
		p.ReasonCode = p.remainingBytes[i]
		if protocolVersion >= conf.MQTT_V5 {
			_, err := p.parseProperties(i)
			if err != 0 {
				slog.Error("err reading properties", "err", err)
				return err
			}
		}
	} else {
		p.ReasonCode = 0
	}
	return 0
}
