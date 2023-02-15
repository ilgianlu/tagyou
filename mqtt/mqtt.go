package mqtt

import (
	"github.com/ilgianlu/tagyou/model"
)

func Start(mqttPort string, wsPort string) {
	// prepare connections repo
	connections := model.Connections{}
	connections.Conns = make(map[string]model.TagyouConn)

	go StartWebSocket(wsPort, &connections)
	StartMQTT(mqttPort, &connections)
}
