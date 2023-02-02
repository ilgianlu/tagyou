package session

import (
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/sessions"

type SessionController struct {
}

func New() *SessionController {
	return &SessionController{}
}

func (sc SessionController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, sc.GetSessions)
}
