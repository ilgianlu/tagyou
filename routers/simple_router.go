package routers

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

type SimpleRouter struct {
	Conns Connections
}

func (s SimpleRouter) AddDestination(clientId string, conn model.TagyouConn) {
	s.Conns.Add(clientId, conn)
}

func (s SimpleRouter) RemoveDestination(clientId string) {
	err := s.Conns.Close(clientId)
	if err != nil {
		slog.Debug("could not close connection", "client-id", clientId, "err", err)
	}
	s.Conns.Remove(clientId)
}

func (s SimpleRouter) DestinationExists(clientId string) bool {
	_, exists := s.Conns.Exists(clientId)
	return exists
}

func (s SimpleRouter) Send(clientId string, payload []byte) {
	conn, exists := s.Conns.Exists(clientId)
	if exists {
		if conn == nil {
			slog.Error("cannot write to net.Conn, c is nil (removing)", "client-id", clientId)
			s.Conns.Remove(clientId)
			return
		}
		_, err := conn.Write(payload)
		if err != nil {
			slog.Debug("cannot write to net.Conn", "client-id", clientId, "err", err)
		}
	} else {
		slog.Debug("client is not connected", "client-id", clientId)
	}
}

func (s SimpleRouter) Forward(senderId string, topic string, p *packet.Packet) {
	destSubs := []string{topic}
	s.sendSubscribers(topic, destSubs, p)
	s.sendSharedSubscribers(topic, destSubs, p)
}

func (s SimpleRouter) SendRetain(protocolVersion uint8, subscription model.Subscription) {
	dests := []string{subscription.Topic}
	retains := persistence.RetainRepository.FindRetains(dests)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := packet.Publish(protocolVersion, subscription.Qos, true, r.Topic, packet.NewPacketIdentifier(), r.ApplicationMessage)
		s.Send(subscription.ClientId, p.ToByteSlice())
	}
}

func (s SimpleRouter) sendSubscribers(topic string, destSubs []string, p *packet.Packet) {
	subs := persistence.SubscriptionRepository.FindSubscriptions(destSubs, false)
	for _, sub := range subs {
		s.forwardSend(topic, sub, p)
	}
}

func (s SimpleRouter) sendSharedSubscribers(topic string, destSubs []string, p *packet.Packet) {
	subs := persistence.SubscriptionRepository.FindOrderedSubscriptions(destSubs, true)
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
	slog.Debug("picked client", "client-id", group[i].ClientId)
	return group[i]
}

func groupSubscribers(subs []model.Subscription) model.SubscriptionGroup {
	grouped := model.SubscriptionGroup{}
	for _, s := range subs {
		online := persistence.SessionRepository.IsOnline(s.SessionID)
		if !online {
			continue
		}
		if val, ok := grouped[s.ShareName]; ok {
			grouped[s.ShareName] = append(val, s)
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
		persistence.RetryRepository.InsertOne(r)
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
		persistence.RetryRepository.InsertOne(r)
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
