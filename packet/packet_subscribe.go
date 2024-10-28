package packet

import (
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/topic"
)

func (p *Packet) subscribeReq(session *model.RunningSession) int {
	session.Mu.RLock()
	defer session.Mu.RUnlock()
	// variable header
	i := 2 // 2 bytes for packet identifier
	if session.ProtocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			slog.Error("err reading properties", "err", err)
			return err
		}
		i = i + pl
	}
	// payload
	j := 0
	p.Subscriptions = []model.Subscription{}
	for {
		sl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		s := string(p.remainingBytes[i : i+sl])
		sub := model.Subscription{
			SessionID: session.SessionID,
			ClientId:  session.ClientId,
		}
		if topic.SharedSubscription(s) {
			sub.Shared = true
			sub.ShareName, sub.Topic = topic.SharedSubscriptionTopicParse(s)
		} else {
			sub.Topic = s
		}
		i = i + sl
		if p.remainingBytes[i]&0x12 != 0 {
			slog.Debug("ignore this subscription & stop")
			break
		}
		sub.RetainHandling = p.remainingBytes[i] & 0x30 >> 4
		sub.RetainAsPublished = p.remainingBytes[i] & 0x08 >> 3
		sub.NoLocal = p.remainingBytes[i] & 0x04 >> 2
		sub.Qos = p.remainingBytes[i] & 0x03
		sub.ProtocolVersion = session.ProtocolVersion
		sub.Enabled = true
		sub.CreatedAt = time.Now().Unix()
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
