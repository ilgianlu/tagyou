package mqtt

import (
	"github.com/ilgianlu/tagyou/routers"
)

func Start(mqttPort string, wsPort string) {
	router := routers.NewSimple()

	go StartWebSocket(wsPort, router)
	StartMQTT(mqttPort, router)
}
