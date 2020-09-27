package session

import (
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

const resourceName string = "/sessions"

type SessionController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *SessionController {
	return &SessionController{db}
}

func (sc SessionController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, sc.GetSessions)
}
