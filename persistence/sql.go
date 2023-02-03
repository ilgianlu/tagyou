package persistence

import (
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/sqlrepository"
	"gorm.io/gorm"
)

func InitSqlRepositories(db *gorm.DB) {
	sqlrepository.Migrate(db)

	AuthRepository = sqlrepository.AuthSqlRepository{Db: db}
	SessionRepository = sqlrepository.SessionSqlRepository{Db: db}
	SubscriptionRepository = sqlrepository.SubscriptionSqlRepository{Db: db}
	RetainRepository = sqlrepository.RetainSqlRepository{Db: db}
	RetryRepository = sqlrepository.RetrySqlRepository{Db: db}

	if conf.CLEAN_EXPIRED_SESSIONS {
		sqlrepository.StartSessionCleaner(db)
	}
	if conf.CLEAN_EXPIRED_RETRIES {
		sqlrepository.StartRetryCleaner(db)
	}
}
