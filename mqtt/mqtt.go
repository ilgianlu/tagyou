package mqtt

import (
	"errors"
	"fmt"
	"net"
	"time"

	bolt "go.etcd.io/bbolt"
)

type MQTT struct {
	db *bolt.DB
}

func New(db *bolt.DB) MQTT {
	var m MQTT
	m.db = db
	return m
}

func (m MQTT) HandleConnection(conn net.Conn) {
	for {
		p, rerr := readPacket(conn)
		if rerr != nil {
			fmt.Printf("err %s\n", rerr)
			defer conn.Close()
			break
		}

		p.PrettyLog()

		resp, err := p.Respond(m.db)
		if err != nil {
			defer conn.Close()
			break
		}
		werr := writePacket(conn, resp)
		if werr != nil {
			fmt.Printf("err %s\n", werr)
			defer conn.Close()
			break
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
