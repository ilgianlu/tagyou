package mqtt

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sender"
)

func Start(mqttPort string, wsPort string) {
	// prepare connections repo
	connections := model.Connections{}
	connections.Conns = make(map[string]model.TagyouConn)

	sender := sender.SimpleSender{Connections: &connections}

	go StartWebSocket(wsPort, sender, &connections)
	StartMQTT(mqttPort, sender, &connections)
}
