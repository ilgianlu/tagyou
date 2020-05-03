package mqtt

const EVENT_CONNECT = 0
const EVENT_PUBLISH = 2
const EVENT_PUBACKED = 3
const EVENT_PUBRECED = 4
const EVENT_PUBRELED = 5
const EVENT_PUBCOMPED = 6
const EVENT_SUBSCRIBED = 10
const EVENT_SUBSCRIPTION = 11
const EVENT_PING = 12
const EVENT_UNSUBSCRIBED = 20
const EVENT_UNSUBSCRIPTION = 21
const EVENT_DISCONNECT = 100
const EVENT_KEEPALIVE_TIMEOUT = 101

const EVENT_PACKET_ERR = 1000

type Event struct {
	eventType    int
	clientId     string
	topic        string
	subscription Subscription
	published    Published
	connection   *Connection
	packet       Packet
	err          uint8
}
