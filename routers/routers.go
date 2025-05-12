package routers

import (
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func NewSimple() Router {
	simple := SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	return SimpleRouter{Conns: &simple}
}

func NewDebug(debugClients string) Router {
	debug := SimpleConnections{
		Conns: make(map[string]model.TagyouConn),
	}
	return DebugRouter{Conns: &debug, DebugClients: debugClients}
}

func NewStandard() Router {
	conns := SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	return StandardRouter{Conns: &conns}
}
