package mqtt

import (
	"errors"
	"fmt"
	"io"
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
					fmt.Println("cannot persist subscription:", err)
				}
			case EVENT_PUBLISH:
				dests := findSubs(m.db, e.topic)
				fmt.Println("clients subscribed :", dests)
				for i := 0; i < len(dests); i++ {
					if c, ok := m.conns[dests[i]]; ok {
						fmt.Println(dests[i], "is connected", c)
						n, err := c.Write(*e.packt)
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
				// if c, ok := m.conns[e.clientId]; ok {
				fmt.Println(e.clientId, "wants to disconnect")
				delete(m.conns, e.clientId)
				fmt.Println(e.clientId, "cleaned")
				// c.Close()
				// }
			}
		}
	}(m.e)

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
		go m.handleConnection(conn)
	}
}

func (m MQTT) handleConnection(conn net.Conn) {
	var connStatus ConnStatus
	for {
		p, rerr := readPacket(conn)
		if rerr != nil {
			fmt.Println("read packet error ", rerr)
			if rerr == io.EOF {
				fmt.Println("connection closed!")
			}
			defer conn.Close()
			break
		}

		var event Event
		event.conn = conn
		event.timestamp = time.Now()

		resp, err := p.Respond(m.db, m.e, &connStatus, &event)
		if err != nil {
			defer conn.Close()
			break
		}

		if resp != nil {
			werr := writePacket(conn, resp)
			if werr != nil {
				fmt.Println("cannot respond", werr)
				defer conn.Close()
				break
			}
		}

		derr := conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if derr != nil {
			fmt.Println("cannot set read deadline", derr)
			defer conn.Close()
			break
		}
	}
}

func writePacket(conn net.Conn, p Packet) error {
	n, err := conn.Write(p)
	if err != nil {
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
		return nil, err
	}
	if n < 2 {
		fmt.Printf("reading fewer bytes: %d\n", n)
		return nil, errors.New("read fewer than expected bytes")
	}
	return p, nil
}
