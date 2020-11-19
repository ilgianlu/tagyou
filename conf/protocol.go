package conf

const MQTT_V5 = 5
const MQTT_V3_11 = 4

const LEVEL_SEPARATOR = "/"
const WILDCARD_MULTI_LEVEL = "#"
const WILDCARD_SINGLE_LEVEL = "+"

const TOPIC_SHARED = "$share"

const MINIMUM_SUPPORTED_PROTOCOL = 4

const MAX_WAIT_FOR_ACK = 4

// client can't subscribe more than N topics on a single "subscribe" command
const MAX_TOPIC_SINGLE_SUBSCRIBE = 10

const QOS0 = 0
const QOS1 = 1
const QOS2 = 2

// GENERIC REASON CODE
const SUCCESS = 0

// SUBSCRIBE REASON CODE
const SUB_TOPIC_FILTER_INVALID = 143

// UNSUBSCRIBE REASON CODE
const UNSUB_NO_SUB_EXISTED = 17
const UNSUB_UNSPECIFIED_ERROR = 128
