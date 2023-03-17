package routers

import (
	"sync"

	"github.com/ilgianlu/tagyou/model"
)

type SimpleConnections struct {
	Conns map[string]model.TagyouConn
	Mu    sync.RWMutex
}

func (c *SimpleConnections) Add(clientId string, conn model.TagyouConn) {
	c.Mu.Lock()
	c.Conns[clientId] = conn
	c.Mu.Unlock()
}

func (c *SimpleConnections) Exists(clientId string) (model.TagyouConn, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	if conn, ok := c.Conns[clientId]; ok {
		return conn, true
	} else {
		return nil, false
	}
}

func (c *SimpleConnections) Close(clientId string) error {
	if connection, ok := c.Exists(clientId); ok {
		return connection.Close()
	}
	return nil
}

func (c *SimpleConnections) Remove(clientId string) {
	c.Mu.Lock()
	delete(c.Conns, clientId)
	c.Mu.Unlock()
}
