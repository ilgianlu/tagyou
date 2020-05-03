package mqtt

const MINIMUM_SUPPORTED_PROTOCOL = 4

const DEFAULT_KEEPALIVE = 30

const MQTT_V5 = 5
const MQTT_V3_11 = 4

const MAX_TOPIC_SINGLE_SUBSCRIBE = 10

var DISALLOW_ANONYMOUS_LOGIN bool = true

const PACKET_TYPE_CONNECT = 1
const PACKET_TYPE_CONNACK = 2
const PACKET_TYPE_PUBLISH = 3
const PACKET_TYPE_PUBACK = 4
const PACKET_TYPE_PUBREC = 5
const PACKET_TYPE_PUBREL = 6
const PACKET_TYPE_PUBCOMP = 7
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

// publish ack in QoS 1
const PUBACK_SUCCESS = 0x00
const PUBACK_NO_MATCHING_SUBSCRIBERS = 0x10

// publish in QoS 2
const PUBCOMP_SUCCESS = 0x00
const PUBREL_SUCCESS = 0x00
const PUBREC_SUCCESS = 0x00
