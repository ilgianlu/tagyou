package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type AuthRepository interface {
	GetByClientIdUsername(clientId string, username string) (model.Auth, error)
}
