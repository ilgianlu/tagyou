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

func (c ConnStatus) cleanSession() bool {
	return (c.connectFlags & 0x02 >> 1) == 1
}

func (c ConnStatus) willFlag() bool {
	return (c.connectFlags & 0x04 >> 2) == 1
}

func (c ConnStatus) willQoS() bool {
	return (c.connectFlags & 0x18 >> 3) == 1
}
