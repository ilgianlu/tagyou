package routers

import (
	"encoding/csv"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

type DebugRouter struct {
	Conns        model.Connections
	DebugClients string
}

func (s DebugRouter) GetConns() model.Connections {
	return s.Conns
}

func (s DebugRouter) AddDestination(clientId string, conn model.TagyouConn) {
	s.Conns.Add(clientId, conn)
}

func (s DebugRouter) RemoveDestination(clientId string) {
	err := s.Conns.Close(clientId)
	if err != nil {
		slog.Debug("[MQTT] could not close connection", "client-id", clientId, "err", err)
	}
	s.Conns.Remove(clientId)
}

func (s DebugRouter) DestinationExists(clientId string) bool {
	_, exists := s.Conns.Exists(clientId)
	return exists
}

func (s DebugRouter) Send(clientId string, payload []byte) {
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

func (s DebugRouter) Forward(senderId string, topic string, p model.Packet) {
	destSubs := explodeFull(topic)
	s.sendSubscribers(topic, destSubs, p)
	s.sendSharedSubscribers(topic, destSubs, p)
	s.sendDebug(senderId, topic, p)
}

func (s DebugRouter) SendRetain(protocolVersion uint8, subscription model.Subscription) {
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

func (s DebugRouter) sendDebug(senderId string, topic string, p model.Packet) {
	filename := conf.DebugDataFilepath(senderId)
	debugFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("error writing to debug file", "err", err, "filename", filename)
		return
	}
	defer debugFile.Close()

	writer := csv.NewWriter(debugFile)
	defer writer.Flush()

	debugLine := []string{
		strconv.FormatInt(time.Now().Unix(), 10),
		senderId,
		topic,
		string(p.ApplicationMessage()),
	}
	err = writer.Write(debugLine)
	if err != nil {
		slog.Error("error writing to debug file", "err", err, "filename", filename)
		return
	}
}

func (s DebugRouter) sendSubscribers(topic string, destSubs []string, p model.Packet) {
	subs := persistence.SubscriptionRepository.FindSubscriptions(destSubs, false)
	for _, sub := range subs {
		s.forwardSend(topic, sub, p)
	}
}

func (s DebugRouter) sendSharedSubscribers(topic string, destSubs []string, p model.Packet) {
	subs := persistence.SubscriptionRepository.FindOrderedSubscriptions(destSubs, true)
	grouped := groupSubscribers(subs)
	for _, group := range grouped {
		dest := pickDest(group, 1)
		s.forwardSend(topic, dest, p)
	}
}

func (s DebugRouter) forwardSend(topic string, sub model.Subscription, p model.Packet) {
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
