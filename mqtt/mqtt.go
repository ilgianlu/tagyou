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
	DISALLOW_ANONYMOUS_LOGIN = os.Getenv("DISALLOW_ANONYMOUS_LOGIN") == "true"
	db, err := openDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	connections := make(inMemoryConnections)
	subscriptions := SqliteSubscriptions{db: db}
	retains := SqliteRetains{db: db}
	retries := SqliteRetries{db: db}
	auths := SqliteAuths{db: db}
	events := make(chan Event, 1)
	outQueue := make(chan OutData, 1)

	go rangeEvents(subscriptions, retains, connections, auths, events, outQueue)
	go rangeOutQueue(connections, retries, outQueue)

	startTCP(events, port)
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
	var connection Connection
	connection.conn = conn
	connection.keepAlive = DEFAULT_KEEPALIVE
	buffers := make(chan []byte, 2)
	packets := make(chan Packet)
	for {
		buffer := make([]byte, 1024)
		bytesCount, err := conn.Read(buffer)
		if err != nil {
			log.Println("could read packet", err)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("keepalive not respected!")
				sendWill(conn, &connection, events)
				disconnectClient(&connection, events)
				break
			}
			if err == io.EOF {
				log.Println("connection closed!")
				sendWill(conn, &connection, events)
				disconnectClient(&connection, events)
				break
			}
		}
		buffers <- buffer[:bytesCount]

		go rangeBuffers(buffers, packets)
		go rangePackets(&connection, packets, events)

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
		event.packet = willPacket
		e <- event
	}
}

func disconnectClient(connection *Connection, e chan<- Event) {
	var event Event
	event.eventType = EVENT_DISCONNECT
	event.clientId = connection.clientId
	event.connection = connection
	e <- event
}
