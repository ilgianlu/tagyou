package persistence

import (
	"github.com/ilgianlu/tagyou/repository"
)

var (
	AuthRepository         repository.AuthRepository
	SessionRepository      repository.SessionRepository
	SubscriptionRepository repository.SubscriptionRepository
	RetainRepository       repository.RetainRepository
	RetryRepository        repository.RetryRepository
)

type Persistence interface {
	Init() error
	Close()
}
