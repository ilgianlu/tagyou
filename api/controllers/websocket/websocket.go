package websocket

import (
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/ws"

type WebsocketController struct {
}

func New() *WebsocketController {
	return &WebsocketController{}
}

func (wc WebsocketController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, wc.UpgradeConnection)
}
