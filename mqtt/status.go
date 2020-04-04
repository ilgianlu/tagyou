package mqtt

const CLIENTS_BUCKET = "clients"
const CONNECT_TIME = "connectTime"
const CLIENTID = "clientId"
const PROTOCOL_VERSION = "protocolVersion"
const CONNECT_FLAGS = "connectFlags"
const KEEP_ALIVE = "keepAlive"

type ConnStatus struct {
	clientId        string
	protocolVersion uint8
	connectFlags    uint8
	keepAlive       []byte
}
