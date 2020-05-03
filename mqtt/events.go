package mqtt

import (
	"log"
	"time"
)

func rangeEvents(subscriptions Subscriptions, retains Retains, connections Connections, auths Auths, events <-chan Event, outQueue chan<- OutData) {
	for e := range events {
		switch e.eventType {
		case EVENT_CONNECT:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client connect")
			clientConnection(connections, subscriptions, auths, e, outQueue)
		case EVENT_SUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client subscribed")
			clientSubscribed(e, outQueue)
		case EVENT_SUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.subscription.clientId, "client subscription", e.subscription.topic)
			clientSubscription(subscriptions, retains, e, outQueue)
		case EVENT_UNSUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscribed")
			clientUnsubscribed(e, outQueue)
		case EVENT_UNSUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscription", e.topic)
			clientUnsubscription(subscriptions, e)
		case EVENT_PUBLISH:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client published to", e.published.topic)
			clientPublish(subscriptions, retains, connections, e, outQueue)
		case EVENT_PUBACKED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client acked message", e.packet.packetIdentifier)
			clientPuback(e)
		case EVENT_PING:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client ping")
			clientPing(e)
		case EVENT_DISCONNECT:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client disconnect")
			clientDisconnect(subscriptions, connections, e)
		case EVENT_PACKET_ERR:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "packet error")
			clientDisconnect(subscriptions, connections, e)
		}
	}
}

func clientConnection(connections Connections, subscriptions Subscriptions, auths Auths, e Event, outQueue chan<- OutData) {
	if DISALLOW_ANONYMOUS_LOGIN && !auths.checkAuth(e.clientId, e.connection.username, e.connection.password) {
		log.Println("wrong connect credentials")
		return
	}
	aerr := connections.addConn(e.clientId, *e.connection)
	if aerr != nil {
		log.Println("could not add connection", aerr)
		return
	}
	if e.connection.cleanStart() {
		subscriptions.remSubscriptionsByClientId(e.clientId)
	} else {
		subscriptions.enableClientSubscriptions(e.clientId)
	}
	var o OutData
	o.clientId = e.clientId
	if e.err != 0 {
		o.packet = Connack(false, e.err)
	} else {
		o.packet = Connack(false, 0)
	}
	outQueue <- o
}

func clientSubscribed(e Event, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.clientId
	o.packet = Suback(e.packet.packetIdentifier, e.packet.subscribedCount, e.subscription.QoS)
	outQueue <- o
}

func clientSubscription(subscriptions Subscriptions, retains Retains, e Event, outQueue chan<- OutData) {
	err := subscriptions.addSubscription(e.subscription)
	if err != nil {
		log.Println("cannot persist subscription:", err)
	}
	sendRetain(retains, e, outQueue)
}

func sendRetain(retains Retains, e Event, outQueue chan<- OutData) {
	rs := retains.findRetainsByTopic(e.subscription.topic)
	if len(rs) == 0 {
		return
	}
	for _, r := range rs {
		var o OutData
		o.clientId = e.clientId
		o.packet = Publish(e.subscription.QoS, true, r.topic, r.applicationMessage)
		outQueue <- o
	}
}

func clientUnsubscribed(e Event, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.clientId
	o.packet = Unsuback(e.packet.packetIdentifier, e.packet.subscribedCount)
	outQueue <- o
}

func clientUnsubscription(subscriptions Subscriptions, e Event) {
	err := subscriptions.remSubscription(e.topic, e.clientId)
	if err != nil {
		log.Println("could not remove topic subscription")
	}
}

func clientPublish(subs Subscriptions, retains Retains, connections Connections, e Event, outQueue chan<- OutData) {
	if e.published.retain {
		saveRetain(retains, e)
	}
	dests := subs.findTopicSubscribers(e.published.topic)
	count := sendToDests(connections, dests, e.packet, outQueue)
	if e.published.qos == 1 {
		var res uint8
		if count == 0 {
			res = PUBACK_NO_MATCHING_SUBSCRIBERS
		} else {
			res = PUBACK_SUCCESS
		}
		log.Println("pub ack", e.packet.packetIdentifier, "being sent to", e.clientId)
		var o OutData
		o.clientId = e.clientId
		o.packet = Puback(e.packet.packetIdentifier, res)
		outQueue <- o
	}
}

func sendToDests(connections Connections, dests []Subscription, p Packet, outQueue chan<- OutData) int {
	count := 0
	for i := 0; i < len(dests); i++ {
		var o OutData
		o.clientId = dests[i].clientId
		o.packet = p
		outQueue <- o
	}
	return count
}

func clientPuback(e Event) {
	// find msg identifier sent
	// check reasoncode
}

func saveRetain(retains Retains, e Event) {
	var r Retain
	r.topic = e.published.topic
	r.applicationMessage = e.packet.remainingBytes[e.packet.applicationMessage:]
	r.createdAt = time.Now()
	err := retains.addRetain(r)
	if err != nil {
		log.Println("could not save retained message:", err)
	}
}

func clientPing(e Event) {
	p := PingResp()
	_, werr := e.connection.publish(p.toByteSlice())
	if werr != nil {
		log.Println("could not write to", e.clientId)
	}
}

func clientDisconnect(subscriptions Subscriptions, connections Connections, e Event) {
	subscriptions.disableClientSubscriptions(e.clientId)
	if toRem, ok := connections.findConn(e.clientId); ok {

		err0 := connections.remConn(toRem.clientId)
		if err0 != nil {
			log.Println("could not remove connection from connections")
		}
		err := toRem.conn.Close()
		if err != nil {
			log.Println("could not close conn", err)
		}
	}
}
