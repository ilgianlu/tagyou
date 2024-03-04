package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type ClientRepository interface {
	Create(client model.Client) error
	GetAll() []model.Client
	GetByClientIdUsername(clientId string, username string) (model.Client, error)
	GetById(id int64) (model.Client, error)
	DeleteById(id int64) error
}
