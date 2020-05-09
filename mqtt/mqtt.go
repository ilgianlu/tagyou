package mqtt

import (
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func StartMQTT(port string) {
	DISALLOW_ANONYMOUS_LOGIN = os.Getenv("DISALLOW_ANONYMOUS_LOGIN") == "true"
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("failed to connect database")
	}
	defer db.Close()

	var connections Connections
	events := make(chan Event, 1)
	outQueue := make(chan OutData, 1)

	go rangeEvents(connections, db, events, outQueue)
	go rangeOutQueue(connections, db, outQueue)

	startTCP(events, port)
}

func startTCP(events chan<- Event, port string) {
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
		go handleConnection(conn, events)
	}
}

func handleConnection(conn net.Conn, events chan<- Event) {
	defer conn.Close()
	client := clientConn{keepAlive: DEFAULT_KEEPALIVE}
	buffers := make(chan []byte, 2)
	packets := make(chan Packet)
	connectOk := make(chan clientConn)
	for {
		buffer := make([]byte, 1024)
		bytesCount, err := conn.Read(buffer)
		if err != nil {
			log.Println("could read packet", err)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("keepalive not respected!")
				willEvent(client.clientId, events)
				disconnectClient(client.clientId, events)
				break
			}
			if err == io.EOF {
				log.Println("connection closed!")
				willEvent(client.clientId, events)
				disconnectClient(client.clientId, events)
				break
			}
		}
		buffers <- buffer[:bytesCount]

		go rangeBuffers(buffers, packets)
		go rangePackets(packets, events)

		derr := conn.SetReadDeadline(time.Now().Add(time.Duration(client.keepAlive*2) * time.Second))
		if derr != nil {
			log.Println("cannot set read deadline", derr)
			defer conn.Close()
			break
		}
		client = <-connectOk
	}
	log.Println("abandon closed connection!")
}

func willEvent(clientId string, e chan<- Event) {
	var event Event
	event.eventType = EVENT_WILL_SEND
	event.clientId = clientId
	e <- event
}

func disconnectClient(clientId string, e chan<- Event) {
	var event Event
	event.eventType = EVENT_DISCONNECT
	event.clientId = clientId
	e <- event
}
