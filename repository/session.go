package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

type SessionRepository interface {
	PersistSession(running *model.RunningSession) (sessionId int64, err error)
	CleanSession(clientId string) error
	SessionExists(clientId string) (model.Session, bool)
	DisconnectSession(clientId string)
	GetById(sessionId uint) (model.Session, error)
	GetAll() []model.Session
	IsOnline(sessionId int64) bool
}
