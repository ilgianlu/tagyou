package model

import (
	"net"
	"sync"
)

type TagyouConn interface {
	Write(b []byte) (n int, err error)
	Close() error
	RemoteAddr() net.Addr
}

type Connections struct {
	Conns map[string]TagyouConn
	Mu    sync.RWMutex
}

func (c *Connections) Add(clientId string, conn TagyouConn) {
	c.Mu.Lock()
	c.Conns[clientId] = conn
	c.Mu.Unlock()
}

func (c *Connections) Exists(clientId string) (TagyouConn, bool) {
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
