package packet

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/topic"
)

func (p *Packet) unsubscribeReq() int {
	p.Session.Mu.RLock()
	defer p.Session.Mu.RUnlock()
	p.Event = EVENT_UNSUBSCRIBED
	i := 2 // 2 bytes for packet identifier
	if p.Session.ProtocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			log.Error().Msgf("err reading properties %d", err)
			return err
		}
		i = i + pl
	}
	p.Subscriptions = make([]model.Subscription, 0)
	j := 0
	for {
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		tpc := string(p.remainingBytes[i : i+sl])
		unsub := model.Subscription{}
		if topic.SharedSubscription(tpc) {
			shareName, unsubTopic := topic.SharedSubscriptionTopicParse(tpc)
			unsub.Shared = true
			unsub.ShareName = shareName
			unsub.Topic = unsubTopic
		} else {
			unsub.Topic = tpc
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
