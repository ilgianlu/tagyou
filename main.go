package main

import (
	"fmt"
	"net"

	"github.com/ilgianlu/tagyou/mqtt"
)

func handleConnection(conn net.Conn) {
	p := make(mqtt.Packet, 255)
	n, err := conn.Read(p)
	if err != nil {
		fmt.Printf("err %s\n", err)
	}
	if n < 2 {
		fmt.Printf("reading fewer bytes: %d\n", n)
	}
	fmt.Printf("read %d bytes\n", n)
	// fmt.Println(bs)
	fmt.Printf("packet type %d\n", p.PacketType())
	fmt.Printf("flags %d\n", p.Flags())
	fmt.Printf("remaining length %d\n", p.RemainingLength())
	fmt.Println("payload", p.Payload())
	if p.PacketType() == 1 {
		fmt.Println(p.ProtocolName())
	}
}

func main() {
	ln, err := net.Listen("tcp", ":3000")
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
