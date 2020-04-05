package mqtt

import (
	"fmt"
	"io"
	"net"
	"time"
)

func StartMQTT(port string) {
	connections := make(map[string]net.Conn)
	subs := make(inMemorySubscriptions)
	events := make(chan Event, 2)

	go rangeEvents(subs, events, connections)

	startTCP(subs, events, port)
}

func rangeEvents(subs Subscriptions, events <-chan Event, connections map[string]net.Conn) {
	for e := range events {
		switch e.eventType {
		case EVENT_CONNECT:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "connected clientId")
			clientConnection(connections, e)
		case EVENT_SUBSCRIBED:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client subscribed")
			clientSubscribed(e)
		case EVENT_SUBSCRIPTION:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client subscribed", e.topic)
			clientSubscription(subs, e)
		case EVENT_PUBLISH:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client published", e.topic)
			clientPublish(subs, connections, e)
		case EVENT_DISCONNECT:
			fmt.Println("//!! EVENT type", e.eventType, e.clientId, "client disconnects")
			clientDisconnect(connections, e)
		}
	}
}

func clientConnection(connections map[string]net.Conn, e Event) {
	connections[e.clientId] = e.conn
	if e.err != 0 {
		_, werr := e.conn.Write(Connack(e.err))
		if werr != nil {
			fmt.Println("could not write to", e.clientId)
		}
	} else {
		_, werr := e.conn.Write(Connack(0))
		if werr != nil {
			fmt.Println("could not write to", e.clientId)
		}
	}
}

func clientSubscribed(e Event) {
	_, werr := e.conn.Write(Suback(e.packetIdentifier, e.subscribedCount))
	if werr != nil {
		fmt.Println("could not write to", e.clientId)
	}
}

func clientSubscription(subs Subscriptions, e Event) {
	err := subs.addSub(e.topic, e.clientId)
	if err != nil {
		fmt.Println("cannot persist subscription:", err)
	}
}

func clientPublish(subs Subscriptions, connections map[string]net.Conn, e Event) {
	dests := subs.findSubs(e.topic)
	for i := 0; i < len(dests); i++ {
		if c, ok := connections[dests[i]]; ok {
			n, err := c.Write(append(e.header, e.remainingBytes...))
			if err != nil {
				fmt.Println("cannot write to", dests[i], ":", err)
			}
			fmt.Println("published", n, "bytes to", dests[i])
		} else {
			fmt.Println(dests[i], "is not connected")
			err := subs.remSub(e.topic, dests[i])
			if err != nil {
				fmt.Println("could not remove subscription :", err)
			}
		}
	}
}

func clientDisconnect(connections map[string]net.Conn, e Event) {
	if c, ok := connections[e.clientId]; ok {
		delete(connections, e.clientId)
		err := c.Close()
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
	var connStatus ConnStatus
	for {
		var event Event
		event.conn = conn
		event.clientId = connStatus.clientId
		event.timestamp = time.Now()

		rErr := readHeader(conn, &event)
		if rErr != nil {
			fmt.Println("read header error ", rErr)
			if rErr == io.EOF {
				fmt.Println("connection closed!")
			}
			defer conn.Close()
			break
		}

		rbErr := readRemainingBytes(conn, &event)
		if rbErr != nil {
			fmt.Println("read remaining bytes error ", rbErr)
			if rbErr == io.EOF {
				fmt.Println("connection closed!")
			}
			defer conn.Close()
			break
		}

		mErr := manageEvent(events, &connStatus, &event)
		if mErr != nil {
			fmt.Println("error managing event", mErr)
			defer conn.Close()
			break
		}

		derr := conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if derr != nil {
			fmt.Println("cannot set read deadline", derr)
			defer conn.Close()
			break
		}
	}
	fmt.Println("abandon closed connection!")
}
