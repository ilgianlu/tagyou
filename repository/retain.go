package repository

import "github.com/ilgianlu/tagyou/model"

type RetainRepository interface {
	FindRetains(topics []string) []model.Retain
	Create(r model.Retain) error
	Delete(r model.Retain) error
}
