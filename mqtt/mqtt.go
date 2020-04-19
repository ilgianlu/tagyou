package mqtt

import (
	"database/sql"
	"io"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func StartMQTT(port string) {
	db, err := openDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	connections := make(inMemoryConnections)
	subscriptions := SqliteSubscriptions{db: db}
	events := make(chan Event, 1)

	go rangeEvents(subscriptions, connections, events)

	startTCP(subscriptions, events, port)
}

func openDb() (*sql.DB, error) {
	if _, err := os.Stat(os.Getenv("DB_FILE")); err != nil {
		if os.IsNotExist(err) {
			Seed(os.Getenv("DB_FILE"))
		} else {
			return nil, err
		}
	}
	return sql.Open("sqlite3", os.Getenv("DB_FILE"))
}

func rangeEvents(subscriptions Subscriptions, connections Connections, events <-chan Event) {
	for e := range events {
		switch e.eventType {
		case EVENT_CONNECT:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client connect")
			clientConnection(connections, subscriptions, e)
		case EVENT_SUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client subscribed")
			clientSubscribed(e)
		case EVENT_SUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client subscription", e.topic)
			clientSubscription(subscriptions, e)
		case EVENT_UNSUBSCRIBED:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscribed")
			clientUnsubscribed(e)
		case EVENT_UNSUBSCRIPTION:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client unsubscription", e.topic)
			clientUnsubscription(subscriptions, e)
		case EVENT_PUBLISH:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client published to", e.topic)
			clientPublish(subscriptions, connections, e)
		case EVENT_PING:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client ping")
			clientPing(e)
		case EVENT_DISCONNECT:
			log.Println("//!! EVENT type", e.eventType, e.clientId, "client disconnect")
			clientDisconnect(connections, e)
		}
	}
}

func clientConnection(connections Connections, subscriptions Subscriptions, e Event) {
	aerr := connections.addConn(e.clientId, *e.connection)
	if aerr != nil {
		log.Println("could not add connection", aerr)
	}
	if e.connection.cleanStart() {
		subscriptions.remSubscriptionsByClientId(e.clientId)
	} else {
		subscriptions.enableClientSubscriptions(e.clientId)
	}
	if e.err != 0 {
		_, werr := e.connection.conn.Write(Connack(e.err))
		if werr != nil {
			log.Println("could not write to", e.clientId)
		}
	} else {
		_, werr := e.connection.conn.Write(Connack(0))
		if werr != nil {
			log.Println("could not write to", e.clientId)
		}
	}
}

func clientSubscribed(e Event) {
	p := Suback(e.packet.packetIdentifier, e.packet.subscribedCount)
	_, werr := e.connection.conn.Write(p)
	if werr != nil {
		log.Println("could not write to", e.clientId)
	}
}

func clientSubscription(subscriptions Subscriptions, e Event) {
	err := subscriptions.addSubscription(e.subscription)
	if err != nil {
		log.Println("cannot persist subscription:", err)
	}
}

func clientUnsubscribed(e Event) {
	p := Unsuback(e.packet.packetIdentifier, e.packet.subscribedCount)
	_, werr := e.connection.conn.Write(p)
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

func clientPublish(subs Subscriptions, connections Connections, e Event) {
	dests := subs.findSubscriptionsByTopic(e.topic)
	for i := 0; i < len(dests); i++ {
		if c, ok := connections.findConn(dests[i].clientId); ok {
			n, err := c.publish(append(e.packet.header, e.packet.remainingBytes...))
			if err != nil {
				log.Println("cannot write to", dests[i].clientId, ":", err)
			}
			log.Println("published", n, "bytes to", dests[i].clientId)
		} else {
			log.Println(dests[i], "is not connected")
		}
	}
}

func clientPing(e Event) {
	_, werr := e.connection.conn.Write(PingResp())
	if werr != nil {
		log.Println("could not write to", e.clientId)
	}
}

func clientDisconnect(connections Connections, e Event) {
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

func startTCP(subs Subscriptions, events chan<- Event, port string) {
	// start tcp socket
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("tcp listen error", err)
		return
	}
	log.Println("mqtt listening on", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("tcp accept error", err)
		}
		go handleConnection(events, conn)
	}
}

func handleConnection(events chan<- Event, conn net.Conn) {
	defer conn.Close()
	var connection Connection
	connection.conn = conn
	connection.keepAlive = DEFAULT_KEEPALIVE
	for {
		_, rErr := ReadPacket(conn, &connection, events)
		if rErr != nil {
			log.Println("could read packet", rErr)
			if err, ok := rErr.(net.Error); ok && err.Timeout() {
				log.Println("keepalive not respected!")
				sendWill(conn, &connection, events)
				break
			}
			if rErr == io.EOF {
				log.Println("connection closed!")
				sendWill(conn, &connection, events)
				break
			}
		}

		derr := conn.SetReadDeadline(time.Now().Add(time.Duration(connection.keepAlive*2) * time.Second))
		if derr != nil {
			log.Println("cannot set read deadline", derr)
			defer conn.Close()
			break
		}
	}
	log.Println("abandon closed connection!")
}

func sendWill(conn net.Conn, connection *Connection, e chan<- Event) {
	// publish will message event
	if connection.willTopic != "" {
		willPacket := Publish(connection.willQoS(), connection.willRetain(), connection.willTopic, connection.willMessage)
		var event Event
		event.eventType = EVENT_PUBLISH
		event.topic = connection.willTopic
		event.packet = &willPacket
		e <- event
	}
}
