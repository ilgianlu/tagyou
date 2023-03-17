package mqtt

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/routers"
)

func Start(mqttPort string, wsPort string) {
	// prepare connections repo
	connections := model.Connections{}
	connections.Conns = make(map[string]model.TagyouConn)

	router := routers.SimpleRouter{Connections: &connections}

	go StartWebSocket(wsPort, router)
	StartMQTT(mqttPort, router)
}
