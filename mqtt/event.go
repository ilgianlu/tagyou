package mqtt

import (
	"net"
	"time"
)

const EVENT_CONNECT = 0
const EVENT_SUBSCRIBED = 10
const EVENT_SUBSCRIPTION = 11
const EVENT_PUBLISH = 2
const EVENT_DISCONNECT = 100

type Event struct {
	eventType        int
	clientId         string
	topic            string
	header           []byte
	packetType       uint8
	flags            uint8
	remainingLength  uint8
	remainingBytes   []byte
	packetIdentifier int
	subscribedCount  int
	err              uint8
	conn             net.Conn
	timestamp        time.Time
}
