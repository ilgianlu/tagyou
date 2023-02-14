package persistence

import (
	"os"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/sqlrepository"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqlPersistence struct {
	db *gorm.DB
}

func (p SqlPersistence) Init() error {
	db, err := openDb()
	if err != nil {
		log.Error().Err(err).Msg("could not open DB")
		return err
	}

	return p.InnerInit(db)
}

func (p *SqlPersistence) InnerInit(db *gorm.DB) error {
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
	p.db = db
	return nil
}

func (p SqlPersistence) Close() {
	sql, err := p.db.DB()
	if err != nil {
		log.Error().Err(err).Msg("could not close DB")
		return
	}
	sql.Close()
}

func openDb() (*gorm.DB, error) {
	logLevel := logger.Silent
	if os.Getenv("DEBUG") != "" {
		logLevel = logger.Info
	}
	return gorm.Open(sqlite.Open(conf.DB_PATH+conf.DB_NAME), &gorm.Config{
		Logger: logger.New(
			&log.Logger,
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logLevel,
				Colorful:      true,
			},
		),
	})
}
