package conf

const MQTT_V5 = 5
const MQTT_V3_11 = 4

const TOPIC_SEPARATOR = "/"
const TOPIC_WILDCARD = "#"

const MINIMUM_SUPPORTED_PROTOCOL = 4

const MAX_WAIT_FOR_ACK = 4

// client can't subscribe more than N topics on a single "subscribe" command
const MAX_TOPIC_SINGLE_SUBSCRIBE = 10

const QOS0 = 0
const QOS1 = 1
const QOS2 = 2

// SUBSCRIBE REASON CODE
const SUB_TOPIC_FILTER_INVALID = 143
