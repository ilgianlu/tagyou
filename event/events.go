package event

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	tpc "github.com/ilgianlu/tagyou/topic"
)

func RangeEvents(connections *model.Connections, events <-chan *packet.Packet) {
	for p := range events {
		clientId := p.Session.GetClientId()
		switch p.Event {
		case packet.EVENT_CONNECT:
			log.Debug().Msgf("//!! EVENT type %d client connect %s", p.Event, clientId)
			onConnect(connections, p)
		case packet.EVENT_SUBSCRIBED:
			log.Debug().Msgf("//!! EVENT type %d client subscribed %s", p.Event, clientId)
			onSubscribe(connections, p)
		case packet.EVENT_UNSUBSCRIBED:
			log.Debug().Msgf("//!! EVENT type %d client unsubscribed %s", p.Event, clientId)
			onUnsubscribe(connections, p)
		case packet.EVENT_PUBLISH:
			log.Debug().Msgf("//!! EVENT type %d client published to %s %s QoS %d", p.Event, p.Topic, clientId, p.QoS())
			OnPublish(connections, p)
		case packet.EVENT_PUBACKED:
			log.Debug().Msgf("//!! EVENT type %d client acked message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPuback(p)
		case packet.EVENT_PUBRECED:
			log.Debug().Msgf("//!! EVENT type %d pub received message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPubrec(connections, p)
		case packet.EVENT_PUBRELED:
			log.Debug().Msgf("//!! EVENT type %d pub releases message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPubrel(connections, p)
		case packet.EVENT_PUBCOMPED:
			log.Debug().Msgf("//!! EVENT type %d pub complete message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPubcomp(p)
		case packet.EVENT_PING:
			log.Debug().Msgf("//!! EVENT type %d client ping %s", p.Event, clientId)
			onPing(connections, p)
		case packet.EVENT_DISCONNECT:
			log.Debug().Msgf("//!! EVENT type %d client disconnect %s", p.Event, clientId)
			clientDisconnect(p, connections, clientId)
		case packet.EVENT_WILL_SEND:
			log.Debug().Msgf("//!! EVENT type %d sending will message %s", p.Event, clientId)
			sendWill(connections, p)
		case packet.EVENT_PACKET_ERR:
			log.Debug().Msgf("//!! EVENT type %d packet error %s", p.Event, clientId)
			clientDisconnect(p, connections, clientId)
		}
	}
}

func onPing(connections *model.Connections, p *packet.Packet) {
	toSend := packet.PingResp()
	SimpleSend(connections, p.Session.GetClientId(), toSend.ToByteSlice())
}

func clientDisconnect(p *packet.Packet, connections *model.Connections, clientId string) {
	if _, ok := connections.Exists(clientId); ok {
		needDisconnection := needDisconnection(p)
		if !needDisconnection {
			return
		}
		connections.Close(clientId)
		connections.Remove(clientId)
		persistence.SessionRepository.DisconnectSession(clientId)
	}
}

func sendForward(connections *model.Connections, topic string, p *packet.Packet) {
	destSubs := tpc.Explode(topic)
	sendSubscribers(connections, topic, destSubs, p)
	sendSharedSubscribers(connections, topic, destSubs, p)
}

func sendSubscribers(connections *model.Connections, topic string, destSubs []string, p *packet.Packet) {
	subs := persistence.SubscriptionRepository.FindSubscriptions(destSubs, false)
	for _, s := range subs {
		send(connections, topic, s, p)
	}
}

func send(connections *model.Connections, topic string, s model.Subscription, p *packet.Packet) {
	qos := getQos(p.QoS(), s.Qos)
	if qos == conf.QOS0 {
		// prepare publish packet qos 0 no packet identifier
		p := packet.Publish(s.ProtocolVersion, conf.QOS0, p.Retain(), topic, 0, p.ApplicationMessage())
		SimpleSend(connections, s.ClientId, p.ToByteSlice())
	} else if qos == conf.QOS1 {
		// prepare publish packet qos 1 (if sub permit) new packet identifier
		p := packet.Publish(s.ProtocolVersion, qos, p.Retain(), topic, packet.NewPacketIdentifier(), p.ApplicationMessage())
		r := model.Retry{
			ClientId:           s.ClientId,
			PacketIdentifier:   p.PacketIdentifier(),
			Qos:                qos,
			Dup:                p.Dup(),
			ApplicationMessage: p.ApplicationMessage(),
			AckStatus:          model.WAIT_FOR_PUB_ACK,
			CreatedAt:          time.Now().Unix(),
		}
		persistence.RetryRepository.SaveOne(r)
		SimpleSend(connections, r.ClientId, p.ToByteSlice())
	} else if qos == 2 {
		// prepare publish packet qos 2 (if sub permit) new packet identifier
		p := packet.Publish(s.ProtocolVersion, qos, p.Retain(), topic, packet.NewPacketIdentifier(), p.ApplicationMessage())
		r := model.Retry{
			ClientId:           s.ClientId,
			PacketIdentifier:   p.PacketIdentifier(),
			Qos:                qos,
			Dup:                p.Dup(),
			ApplicationMessage: p.ApplicationMessage(),
			AckStatus:          model.WAIT_FOR_PUB_REL,
			CreatedAt:          time.Now().Unix(),
		}
		persistence.RetryRepository.SaveOne(r)
		SimpleSend(connections, r.ClientId, p.ToByteSlice())
	}
}

func sendSharedSubscribers(connections *model.Connections, topic string, destSubs []string, p *packet.Packet) {
	subs := persistence.SubscriptionRepository.FindOrderedSubscriptions(destSubs, true, "share_name")
	grouped := groupSubscribers(subs)
	for _, group := range grouped {
		dest := pickDest(group, 1)
		send(connections, topic, dest, p)
	}
}

func pickDest(group []model.Subscription, mode int8) model.Subscription {
	if mode == 0 {
		return group[0]
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(len(group))
	log.Debug().Msgf("picked %s", group[i].ClientId)
	return group[i]
}

func groupSubscribers(subs []model.Subscription) model.SubscriptionGroup {
	grouped := model.SubscriptionGroup{}
	for _, s := range subs {
		if val, ok := grouped[s.ShareName]; ok {
			if persistence.SessionRepository.IsOnline(s.SessionID) {
				grouped[s.ShareName] = append(val, s)
			}
		} else {
			grouped[s.ShareName] = []model.Subscription{s}
		}
	}
	return grouped
}

func getQos(pubQos uint8, subQos uint8) uint8 {
	if pubQos > subQos {
		return subQos
	} else {
		return pubQos
	}
}

func saveRetain(p *packet.Packet) {
	var r model.Retain
	r.Topic = p.Topic
	r.ApplicationMessage = p.ApplicationMessage()
	r.CreatedAt = time.Now().Unix()
	persistence.RetainRepository.Delete(r)
	if len(r.ApplicationMessage) > 0 {
		persistence.RetainRepository.Create(r)
	}
}

func needDisconnection(p *packet.Packet) bool {
	if session, ok := persistence.SessionRepository.SessionExists(p.Session.ClientId); ok {
		log.Debug().Msgf("[MQTT] (%s) Persisted session LastConnect %d running session %d", p.Session.ClientId, session.LastConnect, p.Session.LastConnect)
		if session.LastConnect > p.Session.LastConnect {
			// session persisted is newer then running memory session... device reconnected!
			// no need to send will
			log.Debug().Msgf("[MQTT] (%s) avoid disconnect! (device reconnected)", p.Session.ClientId)
			return false
		}
	}
	return true
}
