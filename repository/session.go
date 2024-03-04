package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

type SessionRepository interface {
	PersistSession(running *model.RunningSession) (int64, error)
	UpdateSession(sessionId int64, running *model.RunningSession) (int64, error)
	CleanSession(clientId string) error
	SessionExists(clientId string) (model.Session, bool)
	DisconnectSession(clientId string)
	GetById(sessionId int64) (model.Session, error)
	GetAll() []model.Session
	IsOnline(sessionId int64) bool
}
