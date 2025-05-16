package routers

import (
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

type StandardRouter struct {
	Conns model.Connections
}

func (s StandardRouter) GetConns() model.Connections {
	return s.Conns
}

func (s StandardRouter) AddDestination(clientId string, conn model.TagyouConn) {
	s.Conns.Add(clientId, conn)
}

func (s StandardRouter) RemoveDestination(clientId string) {
	err := s.Conns.Close(clientId)
	if err != nil {
		slog.Debug("[MQTT] could not close connection", "client-id", clientId, "err", err)
	}
	s.Conns.Remove(clientId)
}

func (s StandardRouter) DestinationExists(clientId string) bool {
	_, exists := s.Conns.Exists(clientId)
	return exists
}

func (s StandardRouter) Send(clientId string, payload []byte) {
	conn, exists := s.Conns.Exists(clientId)
	if exists {
		if conn == nil {
			slog.Error("[MQTT] cannot write to net.Conn, c is nil (removing)", "client-id", clientId)
			s.Conns.Remove(clientId)
			return
		}
		_, err := conn.Write(payload)
		if err != nil {
			slog.Debug("[MQTT] cannot write to net.Conn", "client-id", clientId, "err", err)
		}
	} else {
		slog.Debug("[MQTT] client is not connected", "client-id", clientId)
	}
}

func (s StandardRouter) Forward(senderId string, topic string, p model.Packet) {
	destSubs := explodeFull(topic)
	s.sendSubscribers(topic, destSubs, p)
	s.sendSharedSubscribers(topic, destSubs, p)
}

func (s StandardRouter) SendRetain(protocolVersion uint8, subscription model.Subscription) {
	dests := explodeFull(subscription.Topic)
	retains := persistence.RetainRepository.FindRetains(dests)
	if len(retains) == 0 {
		return
	}
	for _, r := range retains {
		p := packet.Publish(protocolVersion, subscription.Qos, true, r.Topic, packet.NewPacketIdentifier(), r.ApplicationMessage)
		s.Send(subscription.ClientId, p.ToByteSlice())
	}
}

func (s StandardRouter) sendSubscribers(topic string, destSubs []string, p model.Packet) {
	subs := persistence.SubscriptionRepository.FindSubscriptions(destSubs, false)
	for _, sub := range subs {
		s.forwardSend(topic, sub, p)
	}
}

func (s StandardRouter) sendSharedSubscribers(topic string, destSubs []string, p model.Packet) {
	subs := persistence.SubscriptionRepository.FindOrderedSubscriptions(destSubs, true)
	grouped := groupSubscribers(subs)
	for _, group := range grouped {
		dest := pickDest(group, 1)
		s.forwardSend(topic, dest, p)
	}
}

func (s StandardRouter) forwardSend(topic string, sub model.Subscription, p model.Packet) {
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

func explodeFull(topic string) []string {
	road := strings.Split(topic, conf.LEVEL_SEPARATOR)
	res := []string{"#"}
	for i := 1; i <= len(road); i++ {
		subRoads := explodeSingleLevel(road[:i])
		for _, subRoad := range subRoads {
			if i != len(road) {
				subRoad = append(subRoad, conf.WILDCARD_MULTI_LEVEL)
			}
			res = append(res, strings.Join(subRoad, "/"))
		}
	}
	return res
}

func explodeSingleLevel(road []string) [][]string {
	res := [][]string{}
	l := math.Pow(2, float64(len(road)))
	for i := 0; i < int(l); i++ {
		res = append(res, singleLevel(road, i))
	}
	return res
}

func singleLevel(road []string, i int) []string {
	ss := []string{}
	for p, e := range road {
		o := i & (1 << p)
		if o > 0 {
			ss = append(ss, conf.WILDCARD_SINGLE_LEVEL)
		} else {
			ss = append(ss, e)
		}
	}
	return ss
}
