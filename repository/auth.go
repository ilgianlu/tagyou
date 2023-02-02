package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type AuthRepository interface {
	Create(auth model.Auth) error
	DeleteByClientIdUsername(clientId string, username string) error
	GetAll() []model.Auth
	GetByClientIdUsername(clientId string, username string) (model.Auth, error)
}
