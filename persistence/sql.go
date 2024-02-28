package persistence

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/sqlrepository"
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
		slog.Error("could not open DB", "err", err)
		return err
	}

	return p.InnerInit(db, conf.CLEAN_EXPIRED_SESSIONS, conf.CLEAN_EXPIRED_RETRIES, conf.INIT_ADMIN_PASSWORD)
}

func (p *SqlPersistence) InnerInit(db *gorm.DB, startSessionCleaner bool, startRetryCleaner bool, newAdminPassword []byte) error {
	sqlrepository.Migrate(db)

	ClientRepository = sqlrepository.ClientSqlRepository{Db: db}
	SessionRepository = sqlrepository.SessionSqlRepository{Db: db}
	SubscriptionRepository = sqlrepository.SubscriptionSqlRepository{Db: db}
	RetainRepository = sqlrepository.RetainSqlRepository{Db: db}
	UserRepository = sqlrepository.UserSqlRepository{Db: db}
	RetryRepository = sqlrepository.RetrySqlRepository{Db: db}

	if len(newAdminPassword) > 0 {
		sqlrepository.AdminPasswordReset(db, newAdminPassword)
	}
	if startSessionCleaner {
		sqlrepository.StartSessionCleaner(db)
	}
	if startRetryCleaner {
		sqlrepository.StartRetryCleaner(db)
	}
	p.db = db
	return nil
}

func openDb() (*gorm.DB, error) {
	logLevel := logger.Silent
	if os.Getenv("DEBUG") != "" {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // using log until can be subst with slog
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	return gorm.Open(sqlite.Open(conf.DB_PATH+conf.DB_NAME), &gorm.Config{
		Logger: newLogger,
	})
}
