package mqtt

import (
	"fmt"
	"io"
	"net"
	"time"

	bolt "go.etcd.io/bbolt"
)

func StartMQTT(port string, db *bolt.DB) {
	connections := make(map[string]net.Conn)
	events := make(chan Event, 2)
	dberr := initDb(db)
	if dberr != nil {
		fmt.Println("init db error", dberr)
	}

	go rangeEvents(db, events, connections)

	startTCP(db, events, port)
}

func initDb(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(SUBSCRIPTION_BUCKET))
		tx.CreateBucketIfNotExists([]byte(CLIENTS_BUCKET))
		return nil
	})
}

func rangeEvents(db *bolt.DB, events <-chan Event, connections map[string]net.Conn) {
	for e := range events {
		fmt.Println("//!! EVENT type", e.eventType, "clientId", e.clientId)
		switch e.eventType {
		case EVENT_CONNECT:
			connections[e.clientId] = e.conn
			if e.err != 0 {
				e.conn.Write(Connack(e.err))
			}
			e.conn.Write(Connack(0))
		case EVENT_SUBSCRIBED:
			e.conn.Write(Suback(e.packetIdentifier, e.subscribedCount))
		case EVENT_SUBSCRIPTION:
			var s Subscription
			s.topic = e.topic
			s.clientId = e.clientId
			err := s.persist(db)
			if err != nil {
				fmt.Println("cannot persist subscription:", err)
			}
		case EVENT_PUBLISH:
			dests := findSubs(db, e.topic)
			for i := 0; i < len(dests); i++ {
				if c, ok := connections[dests[i]]; ok {
					n, err := c.Write(append(e.header, e.remainingBytes...))
					if err != nil {
						fmt.Println("cannot write to", dests[i], ":", err)
					}
					fmt.Println("published", n, "bytes to", dests[i])
				} else {
					fmt.Println(dests[i], "is not connected")
					// clear subs
				}
			}
		case EVENT_DISCONNECT:
			if c, ok := connections[e.clientId]; ok {
				delete(connections, e.clientId)
				err := c.Close()
				if err != nil {
					fmt.Println("could not close conn", err)
				}
			}
		}
	}
}

func startTCP(db *bolt.DB, events chan<- Event, port string) {
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
