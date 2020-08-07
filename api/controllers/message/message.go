package message

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/messages"

type MessageController struct {
	mqttClient mqtt.Client
}

func New(mqttClient mqtt.Client) *MessageController {
	return &MessageController{mqttClient}
}

func (mc MessageController) RegisterRoutes(r *httprouter.Router) {
	r.POST(resourceName, mc.PostMessage)
}
