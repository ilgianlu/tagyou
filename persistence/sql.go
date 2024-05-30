package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"os"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/sqlc"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
	"github.com/ilgianlu/tagyou/sqlrepository"

	_ "github.com/mattn/go-sqlite3"
)

type SqlPersistence struct {
}

var dbFile = conf.DB_PATH + conf.DB_NAME

func (p SqlPersistence) Init() error {
	if conf.INIT_DB {
		err := resetDb()
		if err != nil {
			slog.Error("could not reset DB", "err", err)
			return err
		}
	}

	dbConn, err := openDb()
	if err != nil {
		slog.Error("could not open DB", "err", err)
		return err
	}

	db := dbaccess.New(dbConn)

	return p.InnerInit(db, conf.CLEAN_EXPIRED_SESSIONS, conf.CLEAN_EXPIRED_RETRIES, conf.INIT_ADMIN_PASSWORD)
}

func resetDb() error {
	os.Remove(dbFile)

	dbConn, err := openDb()
	if err != nil {
		slog.Error("could not open DB", "err", err)
		return err
	}

	if _, err := dbConn.ExecContext(context.Background(), sqlc.DBSchema); err != nil {
		return err
	}

	return nil
}

func (p *SqlPersistence) InnerInit(db *dbaccess.Queries, startSessionCleaner bool, startRetryCleaner bool, newAdminPassword []byte) error {
	ClientRepository = sqlrepository.ClientSqlRepository{Db: db}
	SessionRepository = sqlrepository.SessionSqlRepository{Db: db}
	SubscriptionRepository = sqlrepository.SubscriptionSqlRepository{Db: db}
	RetainRepository = sqlrepository.RetainSqlRepository{Db: db}
	UserRepository = sqlrepository.UserSqlRepository{Db: db}
	RetryRepository = sqlrepository.RetrySqlRepository{Db: db}

	if len(newAdminPassword) > 0 {
		sqlrepository.AdminPasswordReset(db, newAdminPassword)
		os.Exit(0)
	}
	if startSessionCleaner {
		sqlrepository.StartSessionCleaner(db)
	}
	if startRetryCleaner {
		sqlrepository.StartRetryCleaner(db)
	}
	return nil
}

func openDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return db, err
	}
	db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")
	return db, nil
}
