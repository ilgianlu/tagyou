package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

type SessionRepository interface {
	PersistSession(running *model.RunningSession, connected bool) (sessionId uint, err error)
	CleanSession(clientId string) error
	SessionExists(clientId string) (model.Session, bool)
	DisconnectSession(clientId string)
	Save(*model.Session)
}
