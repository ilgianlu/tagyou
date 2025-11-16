package persistence

import (
	"github.com/ilgianlu/tagyou/repository"
)

var (
	ClientRepository       repository.ClientRepository
	SessionRepository      repository.SessionRepository
	SubscriptionRepository repository.SubscriptionRepository
	RetainRepository       repository.RetainRepository
	RetryRepository        repository.RetryRepository
	UserRepository         repository.UserRepository
)
