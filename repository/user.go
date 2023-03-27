package repository

import (
	"github.com/ilgianlu/tagyou/model"
)

type UserRepository interface {
	GetAll() []model.User
	GetByUsername(username string) (model.User, error)
}
