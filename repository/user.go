package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

type UserRepository interface {
	GetAll() []model.User
	GetById(id int64) (model.User, error)
	GetByUsername(username string) (model.User, error)
	Create(user model.User) error
	DeleteById(id int64) error
}
