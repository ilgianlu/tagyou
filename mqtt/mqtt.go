package mqtt

import (
	"errors"
	"fmt"
	"net"
	"time"

	bolt "go.etcd.io/bbolt"
)

type MQTT struct {
	db    *bolt.DB
	e     chan Event
	conns map[string]net.Conn
}

func New(db *bolt.DB) MQTT {
	var m MQTT
	m.db = db
	m.e = make(chan Event)
	m.conns = make(map[string]net.Conn)
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(SUBSCRIPTION_BUCKET))
		tx.CreateBucketIfNotExists([]byte(CLIENTS_BUCKET))
		return nil
	})
	return m
}

func (m MQTT) Start(port string) {
	go func(events <-chan Event) {
		for e := range events {
			fmt.Println("///////////// EVENT START")
			fmt.Println("new event type:", e.eventType)
			fmt.Println("clientId:", e.clientId)
			fmt.Println("topic:", e.topic)
			fmt.Println("message:", e.message)
			fmt.Println("/////////////")
			switch e.eventType {
			case EVENT_CONNECT:
				fmt.Println("new conn for :", e.clientId)
				m.conns[e.clientId] = e.conn
			case EVENT_SUBSCRIBE:
				var s Subscription
				s.topic = e.topic
				s.clientId = e.clientId
				err := s.persist(m.db)
				if err != nil {
					fmt.Println(err)
				}
			case EVENT_PUBLISH:
				dests := findSubs(m.db, e.topic)
				for i := 0; i < len(dests); i++ {
					c := m.conns[dests[i]]
					if c != nil {
						n, err := c.Write(*e.packt)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println("published", n, "bytes to", dests[i])
					}
					i++
				}
			}
		}
	}(m.e)

	// start tcp socket
	ln, err := net.Listen("tcp", port)
	if err != nil {
		// handle error
		fmt.Println("error", err)
		return
	}
	fmt.Println("mqtt listening on", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("error", err)
		}
		go m.handleConnection(conn)
	}
}

func (m MQTT) handleConnection(conn net.Conn) {
	var connStatus ConnStatus
	for {
		p, rerr := readPacket(conn)
		if rerr != nil {
			fmt.Printf("err %s\n", rerr)
			defer conn.Close()
			break
		}

		resp, err := p.Respond(m.db, m.e, &connStatus)
		if err != nil {
			defer conn.Close()
			break
		}

		if resp != nil {
			werr := writePacket(conn, resp)
			if werr != nil {
				fmt.Printf("err %s\n", werr)
				defer conn.Close()
				break
			}
		}

		derr := conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if derr != nil {
			fmt.Printf("err %s\n", derr)
			defer conn.Close()
			break
		}
	}
}

func writePacket(conn net.Conn, p Packet) error {
	n, err := conn.Write(p)
	if err != nil {
		fmt.Printf("err %s\n", err)
		return err
	} else {
		fmt.Printf("wrote %d bytes\n", n)
		return nil
	}
}

func readPacket(conn net.Conn) (Packet, error) {
	p := make(Packet, 255)
	n, err := conn.Read(p)
	if err != nil {
		fmt.Printf("err %s\n", err)
		return nil, err
	}
	if n < 2 {
		fmt.Printf("reading fewer bytes: %d\n", n)
		return nil, errors.New("read fewer than expected bytes")
	}
	return p, nil
}
