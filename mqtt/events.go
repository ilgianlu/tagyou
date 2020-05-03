package mqtt

import (
	"log"
	"time"
)

func rangeEvents(subscriptions Subscriptions, retains Retains, connections Connections, auths Auths, events <-chan Event) {
	for e := range events {
		switch e.eventType {
		case EVENT_CONNECT:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client connect")
			clientConnection(connections, subscriptions, auths, e)
		case EVENT_SUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client subscribed")
			clientSubscribed(e)
		case EVENT_SUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.subscription.clientId, "client subscription", e.subscription.topic)
			clientSubscription(subscriptions, retains, e)
		case EVENT_UNSUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscribed")
			clientUnsubscribed(e)
		case EVENT_UNSUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscription", e.topic)
			clientUnsubscription(subscriptions, e)
		case EVENT_PUBLISH:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client published to", e.published.topic)
			clientPublish(subscriptions, retains, connections, e)
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

func clientConnection(connections Connections, subscriptions Subscriptions, auths Auths, e Event) {
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
	if e.err != 0 {
		p := Connack(false, e.err)
		_, werr := e.connection.publish(p.toByteSlice())
		if werr != nil {
			log.Println("could not write to", e.clientId)
		}
	} else {
		p := Connack(false, 0)
		_, werr := e.connection.publish(p.toByteSlice())
		if werr != nil {
			log.Println("could not write to", e.clientId)
		}
	}
}

func clientSubscribed(e Event) {
	p := Suback(e.packet.packetIdentifier, e.packet.subscribedCount)
	log.Println(p.toByteSlice())
	_, werr := e.connection.conn.Write(p.toByteSlice())
	if werr != nil {
		log.Println("could not write to", e.clientId)
	}
}

func clientSubscription(subscriptions Subscriptions, retains Retains, e Event) {
	err := subscriptions.addSubscription(e.subscription)
	if err != nil {
		log.Println("cannot persist subscription:", err)
	}
	sendRetain(retains, e)
}

func sendRetain(retains Retains, e Event) {
	rs := retains.findRetainsByTopic(e.subscription.topic)
	if len(rs) == 0 {
		return
	}
	for _, r := range rs {
		p := Publish(e.subscription.QoS, true, r.topic, r.applicationMessage)
		_, werr := e.connection.publish(p.toByteSlice())
		if werr != nil {
			log.Println("could not write to", e.clientId)
		}
	}
}

func clientUnsubscribed(e Event) {
	p := Unsuback(e.packet.packetIdentifier, e.packet.subscribedCount)
	_, werr := e.connection.publish(p.toByteSlice())
	if werr != nil {
		log.Println("could not write to", e.clientId)
	}
}

func clientUnsubscription(subscriptions Subscriptions, e Event) {
	err := subscriptions.remSubscription(e.topic, e.clientId)
	if err != nil {
		log.Println("could not remove topic subscription")
	}
}

func clientPublish(subs Subscriptions, retains Retains, connections Connections, e Event) {
	if e.published.retain {
		saveRetain(retains, e)
	}
	dests := subs.findTopicSubscribers(e.published.topic)
	for i := 0; i < len(dests); i++ {
		if c, ok := connections.findConn(dests[i].clientId); ok {
			n, err := c.publish(e.packet.toByteSlice())
			if err != nil {
				log.Println("cannot write to", dests[i].clientId, ":", err)
			}
			log.Println("published", n, "bytes to", dests[i].clientId)
		} else {
			log.Println(dests[i].clientId, "is not connected")
		}
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
