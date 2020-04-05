package mqtt

const MINIMUM_SUPPORTED_PROTOCOL = 4

const PACKET_TYPE_CONNECT = 1
const PACKET_TYPE_PUBLISH = 3
const PACKET_TYPE_SUBSCRIBE = 8
const PACKET_TYPE_DISCONNECT = 14

const TOPIC_SEPARATOR = "/"
const TOPIC_WILDCARD = "#"

// 0x00
// connect OK
const CONNECT_OK = 0

// 0x80
// don't know what to do
const CONNECT_UNSPECIFIED_ERROR = 128

// 0x84
// Unsupported Protocol Version
// The Server does not support the version of the MQTT protocol requested by the Client
const CONNECT_UNSUPPORTED_PROTOCOL_VERSION = 132
