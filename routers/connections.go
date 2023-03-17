package routers

import "github.com/ilgianlu/tagyou/model"

type Connections interface {
	Add(clientId string, conn model.TagyouConn)
	Exists(clientId string) (model.TagyouConn, bool)
	Close(clientId string) error
	Remove(clientId string)
}
