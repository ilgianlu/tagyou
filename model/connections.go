package model

import (
	"net"
	"sync"
)

type Connections struct {
	Conns map[string]net.Conn
	Mu    sync.RWMutex
}

func (c *Connections) Add(clientId string, conn net.Conn) {
	c.Mu.Lock()
	c.Conns[clientId] = conn
	c.Mu.Unlock()
}

func (c *Connections) Exists(clientId string) (net.Conn, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	if conn, ok := c.Conns[clientId]; ok {
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
	c.Mu.Lock()
	delete(c.Conns, clientId)
	c.Mu.Unlock()
}
