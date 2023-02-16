package message

import (
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/messages"

type MessageController struct {
}

func New() *MessageController {
	return &MessageController{}
}

func (mc MessageController) RegisterRoutes(r *httprouter.Router) {
	r.POST(resourceName, mc.PostMessage)
}
