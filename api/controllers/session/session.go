package session

import (
	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/sessions"

type SessionController struct {
}

func New() *SessionController {
	return &SessionController{}
}

func (sc SessionController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, middlewares.Authenticated(sc.GetSessions))
}
