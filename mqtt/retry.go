package mqtt

import "time"

// qos 1
const WAIT_FOR_PUB_ACK = 10

// qos 2
const WAIT_FOR_PUB_REC = 20
const WAIT_FOR_PUB_REL = 20
const WAIT_FOR_PUB_COMP = 21

type Retry struct {
	clientId           string
	applicationMessage []byte
	packetIdentifier   int
	qos                uint8
	retries            uint8
	ackStatus          uint8
	createdAt          time.Time
}
