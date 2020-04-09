package mqtt

import (
	"net"
)

type Connection struct {
	clientId        string
	protocolVersion uint8
	connectFlags    uint8
	keepAlive       []byte
	conn            net.Conn
}

func (c Connection) publish(msg []byte) (int, error) {
	n, err := c.conn.Write(msg)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (c Connection) cleanSession() bool {
	return (c.connectFlags & 0x02 >> 1) == 1
}

func (c Connection) willFlag() bool {
	return (c.connectFlags & 0x04 >> 2) == 1
}

func (c Connection) willQoS() bool {
	return (c.connectFlags & 0x18 >> 3) == 1
}
