package routers

import (
	"os"

	"github.com/ilgianlu/tagyou/model"
)

func NewSimple() Router {
	simple := SimpleConnections{
		Conns: make(map[string]model.TagyouConn),
	}
	return SimpleRouter{Conns: &simple}
}

func NewDebug(debugFile *os.File) Router {
	debug := SimpleConnections{
		Conns: make(map[string]model.TagyouConn),
	}
	return DebugRouter{Conns: &debug, DebugFile: debugFile}
}

func NewStandard() Router {
	conns := SimpleConnections{
		Conns: make(map[string]model.TagyouConn),
	}
	return StandardRouter{Conns: &conns}
}
