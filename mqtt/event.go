package mqtt

const EVENT_CONNECT = 0
const EVENT_PUBLISH = 2
const EVENT_SUBSCRIBED = 10
const EVENT_SUBSCRIPTION = 11
const EVENT_PING = 12
const EVENT_UNSUBSCRIBED = 20
const EVENT_UNSUBSCRIPTION = 21
const EVENT_DISCONNECT = 100

type Event struct {
	eventType  int
	clientId   string
	topic      string
	connection *Connection
	packet     *Packet
	err        uint8
}
