package packet

import (
	"log"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func (p *Packet) subscribeReq() int {
	p.Event = EVENT_SUBSCRIBED
	// variable header
	i := 2 // 2 bytes for packet identifier
	if p.Session.ProtocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			return err
		}
		i = i + pl
	}
	// payload
	j := 0
	p.Subscriptions = make([]model.Subscription, 0)
	for {
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		s := string(p.remainingBytes[i : i+sl])
		sub := model.Subscription{
			ClientId: p.Session.ClientId,
			Topic:    s,
		}
		i = i + sl
		if p.remainingBytes[i]&0x12 != 0 {
			log.Println("ignore this subscription & stop")
			break
		}
		sub.RetainHandling = p.remainingBytes[i] & 0x30 >> 4
		sub.RetainAsPublished = p.remainingBytes[i] & 0x08 >> 3
		sub.NoLocal = p.remainingBytes[i] & 0x04 >> 2
		sub.Qos = p.remainingBytes[i] & 0x03
		sub.ProtocolVersion = p.Session.ProtocolVersion
		sub.Enabled = true
		sub.CreatedAt = time.Now()
		p.Subscriptions = append(p.Subscriptions, sub)
		i++
		if i >= len(p.remainingBytes)-1 {
			break
		}
		j++
		if j > conf.MAX_TOPIC_SINGLE_SUBSCRIBE {
			break
		}
	}
	return 0
}
