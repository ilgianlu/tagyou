package persistence

import "github.com/ilgianlu/tagyou/sqlrepository"

var (
	ClientRepository       sqlrepository.ClientSqlRepository
	SessionRepository      sqlrepository.SessionSqlRepository
	SubscriptionRepository sqlrepository.SubscriptionSqlRepository
	RetainRepository       sqlrepository.RetainSqlRepository
	RetryRepository        sqlrepository.RetrySqlRepository
	UserRepository         sqlrepository.UserSqlRepository
)
