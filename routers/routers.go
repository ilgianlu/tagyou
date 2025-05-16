package routers

import (
	"log/slog"
	"strings"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func NewSimple(connections model.Connections) model.Router {
	return SimpleRouter{Conns: connections}
}

func NewDebug(connections model.Connections) model.Router {
	return DebugRouter{Conns: connections}
}

func NewStandard(connections model.Connections) model.Router {
	return StandardRouter{Conns: connections}
}

func NewDefault(mode string, connections model.Connections) model.Router {
	slog.Info("default router mode", "mode", strings.ToUpper(mode))
	switch mode {
	case conf.ROUTER_MODE_DEBUG:
		return NewDebug(connections)
	case conf.ROUTER_MODE_SIMPLE:
		return NewSimple(connections)
	default:
		return NewStandard(connections)
	}
}

func ByClientId(clientId string, connections model.Connections) model.Router {
	if strings.Contains(conf.DEBUG_CLIENTS, clientId) {
	    slog.Debug("router selection", "sender", clientId, "mode", "debug")
		return NewDebug(connections)
	} else if strings.Contains(conf.SIMPLE_CLIENTS, clientId) {
	    slog.Debug("router selection", "sender", clientId, "mode", "simple")
		return NewSimple(connections)
	}
	slog.Debug("router selection", "sender", clientId, "mode", "standard")
	return NewStandard(connections)
}
