package mqtt

import "net"

type Connections map[string]net.Conn

func (c *Connections) Add(clientId string, conn net.Conn) {
	(*c)[clientId] = conn
}

func (c *Connections) Exists(clientId string) (net.Conn, bool) {
	if conn, ok := (*c)[clientId]; ok {
		return conn, true
	} else {
		return nil, false
	}
}

func (c *Connections) Close(clientId string) error {
	if connection, ok := c.Exists(clientId); ok {
		return connection.Close()
	}
	return nil
}

func (c *Connections) Remove(clientId string) {
	delete(*c, clientId)
}
