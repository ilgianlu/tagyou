package session

import (
	"net/http"

	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
)

type SessionController struct {
}

func NewController() *SessionController {
	return &SessionController{}
}

func (sc SessionController) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /sessions", middlewares.Authenticated(sc.GetSessions))
}
