package mqtt

import (
	"net"
)

type Connection struct {
	clientId        string
	protocolVersion uint8
	connectFlags    uint8
	keepAlive       int
	conn            net.Conn
}

func (c Connection) publish(msg []byte) (int, error) {
	n, err := c.conn.Write(msg)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (c Connection) reservedBit() bool {
	return (c.connectFlags & 0x01) == 0
}

func (c Connection) cleanStart() bool {
	return (c.connectFlags & 0x02) > 0
}

func (c Connection) willFlag() bool {
	return (c.connectFlags & 0x04) > 0
}

func (c Connection) willQoS() uint8 {
	return (c.connectFlags & 0x18 >> 3)
}

func (c Connection) willRetain() bool {
	return (c.connectFlags & 0x20) > 0
}

func (c Connection) havePass() bool {
	return (c.connectFlags & 0x40) > 0
}

func (c Connection) haveUser() bool {
	return (c.connectFlags & 0x80) > 0
}
