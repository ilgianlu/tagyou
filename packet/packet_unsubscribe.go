package packet

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func (p *Packet) unsubscribeReq() int {
	p.Event = EVENT_UNSUBSCRIBED
	i := 2 // 2 bytes for packet identifier
	if p.Session.ProtocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			return err
		}
		i = i + pl
	}
	p.Subscriptions = make([]model.Subscription, 0)
	j := 0
	for {
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		unsub := model.Subscription{
			Topic: string(p.remainingBytes[i : i+sl]),
		}
		p.Subscriptions = append(p.Subscriptions, unsub)
		i = i + sl
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > 10 {
			break
		}
	}
	return 0
}
