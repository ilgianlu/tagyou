package main

import (
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
		p := make(mqtt.Packet, 255)
		n, err := conn.Read(p)
		if err != nil {
			fmt.Printf("err %s\n", err)
			defer conn.Close()
			break
		}
		if n < 2 {
			fmt.Printf("reading fewer bytes: %d\n", n)
		}
		fmt.Printf("read %d bytes\n", n)
		// fmt.Println(bs)
		p.PrettyLog()

		if p.PacketType() == 1 {
			if p.ProtocolVersion() < 4 {
				fmt.Println("unsupported protocol version err", p.ProtocolVersion())
				werr := writePacket(conn, mqtt.Connack(mqtt.CONNECT_UNSUPPORTED_PROTOCOL_VERSION))
				if werr != nil {
					fmt.Printf("err %s\n", err)
				}
				defer conn.Close()
				break
			}
			werr := writePacket(conn, mqtt.Connack(mqtt.CONNECT_OK))
			if werr != nil {
				fmt.Printf("err %s\n", err)
				defer conn.Close()
				break
			}
		}

		derr := conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if derr != nil {
			fmt.Printf("err %s\n", err)
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

func main() {
	berr := dotenv.Load()
	if berr != nil {
		fmt.Println("error loading env", berr)
	}
	db, derr := bolt.Open(os.Getenv("DB_PATH"), 0666, nil)
	if derr != nil {
		fmt.Println("error opening bbolt", derr)
	}
	defer db.Close()
	ln, err := net.Listen("tcp", os.Getenv("LISTEN_PORT"))
	if err != nil {
		// handle error
		fmt.Println("error", err)
	}
	fmt.Println("listen :3000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("error ", err)
		}
		go handleConnection(conn)
	}
}
