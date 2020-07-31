package session

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
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
