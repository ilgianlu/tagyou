package routers

import "github.com/ilgianlu/tagyou/model"

func NewSimple() Router {
	simple := SimpleConnections{
		Conns: make(map[string]model.TagyouConn),
	}
	return SimpleRouter{Conns: &simple}
}
