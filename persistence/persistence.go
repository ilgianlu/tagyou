package persistence

import (
	"github.com/ilgianlu/tagyou/repository"
	"github.com/ilgianlu/tagyou/sqlrepository"
	"gorm.io/gorm"
)

var (
	AuthRepository         repository.AuthRepository
	SessionRepository      repository.SessionRepository
	SubscriptionRepository repository.SubscriptionRepository
	RetainRepository       repository.RetainRepository
	RetryRepository        repository.RetryRepository
)

func InitSqlRepositories(db *gorm.DB) {
	AuthRepository = sqlrepository.AuthSqlRepository{Db: db}
	SessionRepository = sqlrepository.SessionSqlRepository{Db: db}
	SubscriptionRepository = sqlrepository.SubscriptionSqlRepository{Db: db}
	RetainRepository = sqlrepository.RetainSqlRepository{Db: db}
	RetryRepository = sqlrepository.RetrySqlRepository{Db: db}
}

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&sqlrepository.Auth{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&sqlrepository.Retry{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&sqlrepository.Retain{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&sqlrepository.Session{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&sqlrepository.Subscription{}); err != nil {
		return
	}
}
