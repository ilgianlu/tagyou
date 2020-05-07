package mqtt

import (
	"log"
	"time"
)

func rangeEvents(subscriptions Subscriptions, retains Retains, retries Retries, connections Connections, auths Auths, events <-chan Event, outQueue chan<- OutData) {
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
			clientPublish(subscriptions, retains, e, outQueue)
		case EVENT_PUBACKED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client acked message", e.packet.packetIdentifier)
			clientPuback(e, retries)
		case EVENT_PUBRECED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub received message", e.packet.packetIdentifier)
			clientPubrec(e, outQueue)
		case EVENT_PUBRELED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub releases message", e.packet.packetIdentifier)
			clientPubrel(e, outQueue)
		case EVENT_PUBCOMPED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "pub complete message", e.packet.packetIdentifier)
			clientPubcomp(e, retries)
		case EVENT_PING:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client ping")
			clientPing(e, outQueue)
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
	if c, ok := connections.findConn(e.clientId); ok {
		log.Println("session taken over")
		p := Connack(false, SESSION_TAKEN_OVER)
		c.publish(p.toByteSlice())
		closeRemoveClient(c, connections)
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
	if e.err != 0 {
		respond(e.clientId, Connack(false, e.err), outQueue)
	} else {
		respond(e.clientId, Connack(false, CONNECT_OK), outQueue)
	}
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
		respond(e.clientId, Publish(e.subscription.QoS, true, r.topic, r.applicationMessage), outQueue)
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

func clientPublish(subs Subscriptions, retains Retains, e Event, outQueue chan<- OutData) {
	if e.published.retain {
		saveRetain(retains, e)
	}
	dests := subs.findTopicSubscribers(e.published.topic)
	sendToDests(dests, e.packet, outQueue)
	if e.published.qos == 1 {
		respond(e.clientId, Puback(e.packet.packetIdentifier, PUBACK_SUCCESS), outQueue)
	}
	if e.published.qos == 2 {
		respond(e.clientId, Pubrec(e.packet.packetIdentifier, PUBREC_SUCCESS), outQueue)
	}
}

func sendToDests(dests []Subscription, p Packet, outQueue chan<- OutData) int {
	count := 0
	for i := 0; i < len(dests); i++ {
		respond(dests[i].clientId, p, outQueue)
	}
	return count
}

func respond(clientId string, p Packet, outQueue chan<- OutData) {
	var o OutData
	o.clientId = clientId
	o.packet = p
	outQueue <- o
}

// func forward(topic string, p Packet, outQueue chan<- OutData) {
// 	var o OutData
// 	o.clientId = clientId
// 	o.packet = p
// 	outQueue <- o
// }

func clientPuback(e Event, retries Retries) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok remove retry
	if _, ok := retries.findRetry(e.clientId, e.packet.packetIdentifier); ok {
		log.Println("retry found, removing...", e.clientId, e.packet.packetIdentifier)
		err := retries.remRetry(e.clientId, e.packet.packetIdentifier)
		if err != nil {
			log.Println("could not remove retry", e.clientId, e.packet.packetIdentifier)
		}
	} else {
		log.Println("ack for invalid retry", e.clientId, e.packet.packetIdentifier)
	}
}

func clientPubrec(e Event, outQueue chan<- OutData) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok
	// if retry in wait for pub rec -> send pub rel
	var o OutData
	o.clientId = e.clientId
	o.packet = Pubrel(e.packet.packetIdentifier, PUBREL_SUCCESS)
	outQueue <- o
	// change retry state to wait for pubcomp
}

func clientPubrel(e Event, outQueue chan<- OutData) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok
	// if retry in wait for pubrel -> send pub comp
	var o OutData
	o.clientId = e.clientId
	o.packet = Pubcomp(e.packet.packetIdentifier, PUBCOMP_SUCCESS)
	outQueue <- o
}

func clientPubcomp(e Event, retries Retries) {
	// find msg identifier sent
	// check reasoncode
	// if reasoncode ok
	// if retry in wait for pubcomp -> remove retry
	if _, ok := retries.findRetry(e.clientId, e.packet.packetIdentifier); ok {
		err := retries.remRetry(e.clientId, e.packet.packetIdentifier)
		if err != nil {
			log.Println("could not remove retry", e.clientId, e.packet.packetIdentifier)
		}
	} else {
		log.Println("pub complete for invalid retry", e.clientId, e.packet.packetIdentifier)
	}
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

func clientPing(e Event, outQueue chan<- OutData) {
	var o OutData
	o.clientId = e.clientId
	o.packet = PingResp()
	outQueue <- o
}

func clientDisconnect(subscriptions Subscriptions, connections Connections, e Event) {
	subscriptions.disableClientSubscriptions(e.clientId)
	if toRem, ok := connections.findConn(e.clientId); ok {
		closeRemoveClient(toRem, connections)
	}
}

func closeRemoveClient(connection Connection, connections Connections) {
	err0 := connections.remConn(connection.clientId)
	if err0 != nil {
		log.Println("could not remove connection from connections")
	}
	err := connection.conn.Close()
	if err != nil {
		log.Println("could not close conn", err)
	}
}
