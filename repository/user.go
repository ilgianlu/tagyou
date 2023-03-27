package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

type UserRepository interface {
	GetAll() []model.User
	GetById(id uint) (model.User, error)
	Create(user model.User) error
	DeleteById(id uint) error
}
