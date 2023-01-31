package repository

import "github.com/ilgianlu/tagyou/model"

type RetryRepository interface {
	SaveOne(retry model.Retry) error
}
