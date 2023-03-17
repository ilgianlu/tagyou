package routers

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

type SimpleRouter struct {
	Connections *model.Connections
}

func (s SimpleRouter) AddDestination(clientId string, conn model.TagyouConn) {
	s.Connections.Add(clientId, conn)
}

func (s SimpleRouter) RemoveDestination(clientId string) {
	err := s.Connections.Close(clientId)
	if err != nil {
		log.Debug().Err(err).Msgf("could not clone connection %s", clientId)
	}
	s.Connections.Remove(clientId)
}

func (s SimpleRouter) DestinationExists(clientId string) bool {
	_, exists := s.Connections.Exists(clientId)
	return exists
}

func (s SimpleRouter) Send(clientId string, payload []byte) {
	conn, exists := s.Connections.Exists(clientId)
	if exists {
		if conn == nil {
			log.Error().Msgf("cannot write to %s net.Conn, c is nil (removing)", clientId)
			s.Connections.Remove(clientId)
			return
		}
		// packetBytes := p.ToByteSlice()
		_, err := conn.Write(payload)
		if err != nil {
			log.Debug().Err(err).Msgf("cannot write to %s", clientId)
		}
		// else {
		// 	log.Println("published", n, "bytes to", clientId)
		// }
	} else {
		log.Debug().Msgf("%s is not connected", clientId)
	}
}

func (s SimpleRouter) Forward(topic string, p *packet.Packet) {
	destSubs := tpc.Explode(topic)
	s.sendSubscribers(topic, destSubs, p)
	s.sendSharedSubscribers(topic, destSubs, p)
}

func (s SimpleRouter) sendSubscribers(topic string, destSubs []string, p *packet.Packet) {
	subs := persistence.SubscriptionRepository.FindSubscriptions(destSubs, false)
	for _, sub := range subs {
		s.forwardSend(topic, sub, p)
	}
}

func (s SimpleRouter) sendSharedSubscribers(topic string, destSubs []string, p *packet.Packet) {
	subs := persistence.SubscriptionRepository.FindOrderedSubscriptions(destSubs, true, "share_name")
	grouped := groupSubscribers(subs)
	for _, group := range grouped {
		dest := pickDest(group, 1)
		s.forwardSend(topic, dest, p)
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

func (s SimpleRouter) forwardSend(topic string, sub model.Subscription, p *packet.Packet) {
	qos := getQos(p.QoS(), sub.Qos)
	if qos == conf.QOS0 {
		// prepare publish packet qos 0 no packet identifier
		p := packet.Publish(sub.ProtocolVersion, conf.QOS0, p.Retain(), topic, 0, p.ApplicationMessage())
		s.Send(sub.ClientId, p.ToByteSlice())
	} else if qos == conf.QOS1 {
		// prepare publish packet qos 1 (if sub permit) new packet identifier
		p := packet.Publish(sub.ProtocolVersion, qos, p.Retain(), topic, packet.NewPacketIdentifier(), p.ApplicationMessage())
		r := model.Retry{
			ClientId:           sub.ClientId,
			PacketIdentifier:   p.PacketIdentifier(),
			Qos:                qos,
			Dup:                p.Dup(),
			ApplicationMessage: p.ApplicationMessage(),
			AckStatus:          model.WAIT_FOR_PUB_ACK,
			CreatedAt:          time.Now().Unix(),
		}
		persistence.RetryRepository.SaveOne(r)
		s.Send(r.ClientId, p.ToByteSlice())
	} else if qos == 2 {
		// prepare publish packet qos 2 (if sub permit) new packet identifier
		p := packet.Publish(sub.ProtocolVersion, qos, p.Retain(), topic, packet.NewPacketIdentifier(), p.ApplicationMessage())
		r := model.Retry{
			ClientId:           sub.ClientId,
			PacketIdentifier:   p.PacketIdentifier(),
			Qos:                qos,
			Dup:                p.Dup(),
			ApplicationMessage: p.ApplicationMessage(),
			AckStatus:          model.WAIT_FOR_PUB_REL,
			CreatedAt:          time.Now().Unix(),
		}
		persistence.RetryRepository.SaveOne(r)
		s.Send(r.ClientId, p.ToByteSlice())
	}
}

func getQos(pubQos uint8, subQos uint8) uint8 {
	if pubQos > subQos {
		return subQos
	} else {
		return pubQos
	}
}
