package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	mqtt "github.com/ilgianlu/tagyou/mqtt"
	dotenv "github.com/joho/godotenv"
	bolt "go.etcd.io/bbolt"
)

func handleConnection(conn net.Conn) {
	for {
		p, rerr := readPacket(conn)
		if rerr != nil {
			fmt.Printf("err %s\n", rerr)
			defer conn.Close()
			break
		}

		p.PrettyLog()

		resp, err := p.Respond()
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

func writePacket(conn net.Conn, p mqtt.Packet) error {
	n, err := conn.Write(p)
	if err != nil {
		fmt.Printf("err %s\n", err)
		return err
	} else {
		fmt.Printf("wrote %d bytes\n", n)
		return nil
	}
}

func readPacket(conn net.Conn) (mqtt.Packet, error) {
	p := make(mqtt.Packet, 255)
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

func main() {
	// load env vars
	berr := dotenv.Load()
	if berr != nil {
		fmt.Println("error loading env", berr)
	}

	// open k/v store
	db, derr := bolt.Open(os.Getenv("DB_PATH"), 0666, nil)
	if derr != nil {
		fmt.Println("error opening bbolt", derr)
	}
	defer db.Close()

	// start tcp socket
	ln, err := net.Listen("tcp", os.Getenv("LISTEN_PORT"))
	if err != nil {
		// handle error
		fmt.Println("error", err)
	}
	fmt.Println("listen", os.Getenv("LISTEN_PORT"))
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("error ", err)
		}
		go handleConnection(conn)
	}
}
