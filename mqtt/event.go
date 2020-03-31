package mqtt

import (
	"net"
	"time"
)

type Event struct {
	eventType       int
	clientId        string
	protocolVersion uint8
	connectFlags    uint8
	keepAlive       []byte
	topic           string
	message         string
	conn            net.Conn
	timestamp       time.Time
}
