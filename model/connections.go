package model

type Connections interface {
	Add(clientId string, conn TagyouConn)
	Exists(clientId string) (TagyouConn, bool)
	Close(clientId string) error
	Remove(clientId string)
}
