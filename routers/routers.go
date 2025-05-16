package routers

import (
	"log/slog"
	"strings"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func NewSimple() model.Router {
	simple := SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	return SimpleRouter{Conns: &simple}
}

func NewDebug() model.Router {
	debug := SimpleConnections{
		Conns: make(map[string]model.TagyouConn),
	}
	return DebugRouter{Conns: &debug}
}

func NewStandard() model.Router {
	conns := SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	return StandardRouter{Conns: &conns}
}

func NewDefault(mode string) model.Router {
	slog.Info("default router mode", "mode", strings.ToUpper(mode))
	switch mode {
	case conf.ROUTER_MODE_DEBUG:
		return NewDebug()
	case conf.ROUTER_MODE_SIMPLE:
		return NewSimple()
	default:
		return NewStandard()
	}
}

func ByClientId(clientId string) model.Router {
	if strings.Contains(conf.DEBUG_CLIENTS, clientId) {
	    slog.Debug("router selection", "sender", clientId, "mode", "debug")
		return NewDebug()
	} else if strings.Contains(conf.SIMPLE_CLIENTS, clientId) {
	    slog.Debug("router selection", "sender", clientId, "mode", "simple")
		return NewSimple()
	}
	slog.Debug("router selection", "sender", clientId, "mode", "standard")
	return NewStandard()
}
