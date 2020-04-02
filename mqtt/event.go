package mqtt

import (
	"net"
	"time"
)

const EVENT_CONNECT = 0
const EVENT_SUBSCRIBE = 1
const EVENT_PUBLISH = 2
const EVENT_DISCONNECT = 100

type Event struct {
	eventType int
	clientId  string
	topic     string
	message   string
	packt     *Packet
	conn      net.Conn
	timestamp time.Time
}
