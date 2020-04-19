package mqtt

const MINIMUM_SUPPORTED_PROTOCOL = 4

const DEFAULT_KEEPALIVE = 30

const MAX_TOPIC_SINGLE_SUBSCRIBE = 10

const PACKET_TYPE_CONNECT = 1
const PACKET_TYPE_CONNACK = 2
const PACKET_TYPE_PUBLISH = 3
const PACKET_TYPE_SUBSCRIBE = 8
const PACKET_TYPE_SUBACK = 9
const PACKET_TYPE_UNSUBSCRIBE = 10
const PACKET_TYPE_UNSUBACK = 11
const PACKET_TYPE_PINGREQ = 12
const PACKET_TYPE_PINGRES = 13
const PACKET_TYPE_DISCONNECT = 14
const PACKET_MAX_SIZE = 65000

const TOPIC_SEPARATOR = "/"
const TOPIC_WILDCARD = "#"

// 0x00
// connect OK
const CONNECT_OK = 0

// 0x80
// don't know what to do
const UNSPECIFIED_ERROR = 0x80

const MALFORMED_PACKET = 0x81

// 0x84
// Unsupported Protocol Version
// The Server does not support the version of the MQTT protocol requested by the Client
const UNSUPPORTED_PROTOCOL_VERSION = 132
