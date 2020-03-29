package mqtt

const MINIMUM_SUPPORTED_PROTOCOL = 4

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
