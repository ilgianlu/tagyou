package mqtt

type Connections interface {
	addConn(string, Connection) error
	remConn(string) error
	findConn(string) (Connection, bool)
}

type inMemoryConnections map[string]Connection

func (is inMemoryConnections) addConn(clientId string, conn Connection) error {
	is[clientId] = conn
	return nil
}

func (is inMemoryConnections) remConn(clientId string) error {
	delete(is, clientId)
	return nil
}

func (is inMemoryConnections) findConn(clientId string) (Connection, bool) {
	v, ok := is[clientId]
	return v, ok
}
