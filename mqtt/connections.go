package mqtt

import "net"

type Connections interface {
	addConn(string, net.Conn) error
	remConn(string) error
	findConn(string) net.Conn
}

type inMemoryConnections map[string]net.Conn

func (is inMemoryConnections) addConn(clientId string, conn net.Conn) error {
	is[clientId] = conn
	return nil
}

func (is inMemoryConnections) remConn(clientId string) error {
	delete(is, clientId)
	return nil
}

func (is inMemoryConnections) findConn(clientId string) net.Conn {
	if v, ok := is[clientId]; ok {
		return v
	} else {
		return nil
	}
}
