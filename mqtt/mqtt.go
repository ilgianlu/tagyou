package mqtt

import (
	"net"

	"github.com/ilgianlu/tagyou/model"
)

func Start(mqttPort string, wsPort string) {
	// prepare connections repo
	connections := model.Connections{}
	connections.Conns = make(map[string]net.Conn)

	go StartWebSocket(wsPort, &connections)
	StartMQTT(mqttPort, &connections)
}
