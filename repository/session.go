package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

type SessionRepository interface {
	PersistSession(running *model.RunningSession, connected bool) (sessionId uint, err error)
	CleanSession(clientId string) error
	SessionExists(clientId string) (model.RunningSession, bool)
	DisconnectSession(clientId string)
	GetById(sessionId uint) (model.RunningSession, error)
	GetAll() []model.RunningSession
	Save(*model.RunningSession)
	IsOnline(sessionId uint) bool
}
