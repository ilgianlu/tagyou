package mqtt

import (
	"fmt"
	"io"
	"net"
	"time"
)

func StartMQTT(port string) {
	connections := make(inMemoryConnections)
	topicSubscriptions := make(inMemorySubscriptions)
	clientSubscriptions := make(inMemorySubscriptions)
	events := make(chan Event, 1)

	go rangeEvents(topicSubscriptions, clientSubscriptions, connections, events)

	startTCP(topicSubscriptions, events, port)
}

func rangeEvents(topicSubs Subscriptions, clientSubs Subscriptions, connections Connections, events <-chan Event) {
	for e := range events {
		switch e.eventType {
		case EVENT_CONNECT:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client connect")
			clientConnection(connections, topicSubs, clientSubs, e)
		case EVENT_SUBSCRIBED:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client subscribed")
			clientSubscribed(e)
		case EVENT_SUBSCRIPTION:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client subscription", e.topic)
			clientSubscription(topicSubs, clientSubs, e)
		case EVENT_UNSUBSCRIBED:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscribed")
			clientUnsubscribed(e)
		case EVENT_UNSUBSCRIPTION:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscription", e.topic)
			clientUnsubscription(topicSubs, clientSubs, e)
		case EVENT_PUBLISH:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client published to", e.topic)
			clientPublish(topicSubs, connections, e)
		case EVENT_PING:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client ping")
			clientPing(e)
		case EVENT_DISCONNECT:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client disconnect")
			clientDisconnect(connections, e)
		}
	}
}

func clientConnection(connections Connections, topicSubs Subscriptions, clientSubs Subscriptions, e Event) {
	aerr := connections.addConn(e.clientId, *e.connection)
	if aerr != nil {
		fmt.Println("could not add connection", aerr)
	}
	if e.connection.cleanStart() {
		clientSubs.remSubscribed(e.clientId)
	} else {
		if s, ok := clientSubs.findSubscribed(e.clientId); ok {
			for i := 0; i < len(s); i++ {
				topicSubs.remSubscription(e.clientId, s[i])
			}
		}
	}
	if e.err != 0 {
		_, werr := e.connection.conn.Write(Connack(e.err))
		if werr != nil {
			fmt.Println("could not write to", e.clientId)
		}
	} else {
		_, werr := e.connection.conn.Write(Connack(0))
		if werr != nil {
			fmt.Println("could not write to", e.clientId)
		}
	}
}

func clientSubscribed(e Event) {
	p := Suback(e.packet.packetIdentifier, e.packet.subscribedCount)
	_, werr := e.connection.conn.Write(p)
	if werr != nil {
		fmt.Println("could not write to", e.clientId)
	}
}

func clientSubscription(topicSubs Subscriptions, clientSubs Subscriptions, e Event) {
	err := topicSubs.addSubscription(e.topic, e.clientId)
	if err != nil {
		fmt.Println("cannot persist subscription:", err)
	}
	err0 := clientSubs.addSubscription(e.clientId, e.topic)
	if err0 != nil {
		fmt.Println("cannot persist subscription:", err0)
	}
}

func clientUnsubscribed(e Event) {
	p := Unsuback(e.packet.packetIdentifier, e.packet.subscribedCount)
	_, werr := e.connection.conn.Write(p)
	if werr != nil {
		fmt.Println("could not write to", e.clientId)
	}
}

func clientUnsubscription(topicSubs Subscriptions, clientSubs Subscriptions, e Event) {
	pos0 := topicSubs.remSubscription(e.topic, e.clientId)
	if pos0 < 0 {
		fmt.Println("could not remove topic subscription")
	}
	pos1 := clientSubs.remSubscription(e.clientId, e.topic)
	if pos1 < 0 {
		fmt.Println("could not remove client subscription")
	}
}

func clientPublish(subs Subscriptions, connections Connections, e Event) {
	dests := subs.findSubscribers(e.topic)
	for i := 0; i < len(dests); i++ {
		if c, ok := connections.findConn(dests[i]); ok {
			n, err := c.publish(append(e.packet.header, e.packet.remainingBytes...))
			if err != nil {
				fmt.Println("cannot write to", dests[i], ":", err)
			}
			fmt.Println("published", n, "bytes to", dests[i])
		} else {
			fmt.Println(dests[i], "is not connected")
		}
	}
}

func clientPing(e Event) {
	_, werr := e.connection.conn.Write(PingResp())
	if werr != nil {
		fmt.Println("could not write to", e.clientId)
	}
}

func clientDisconnect(connections Connections, e Event) {
	if toRem, ok := connections.findConn(e.clientId); ok {
		err0 := connections.remConn(toRem.clientId)
		if err0 != nil {
			fmt.Println("could not remove connection from connections")
		}
		err := toRem.conn.Close()
		if err != nil {
			fmt.Println("could not close conn", err)
		}
	}
}

func startTCP(subs Subscriptions, events chan<- Event, port string) {
	// start tcp socket
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("tcp listen error", err)
		return
	}
	fmt.Println("mqtt listening on", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("tcp accept error", err)
		}
		go handleConnection(events, conn)
	}
}

func handleConnection(events chan<- Event, conn net.Conn) {
	var connection Connection
	connection.conn = conn
	connection.keepAlive = DEFAULT_KEEPALIVE
	for {
		_, rErr := ReadPacket(conn, &connection, events)
		if rErr != nil {
			fmt.Println("could not elaborate packet ", rErr)
			if rErr == io.EOF {
				fmt.Println("connection closed!")
			}
			defer conn.Close()
			break
		}

		derr := conn.SetReadDeadline(time.Now().Add(time.Duration(connection.keepAlive*2) * time.Second))
		if derr != nil {
			fmt.Println("cannot set read deadline", derr)
			defer conn.Close()
			break
		}
	}
	fmt.Println("abandon closed connection!")
}
