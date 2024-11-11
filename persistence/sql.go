package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"os"

	"github.com/ilgianlu/tagyou/sqlc"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
	"github.com/ilgianlu/tagyou/sqlrepository"

	_ "github.com/mattn/go-sqlite3"
)

type SqlPersistence struct {
	DbFile       string
	InitDatabase bool
	db           *sql.DB
}

func (p *SqlPersistence) Init(cleanExpiredSessions bool, cleanExpiredRetries bool, initAdminPassword []byte) (*dbaccess.Queries, error) {
	if p.InitDatabase {
		err := p.resetDb()
		if err != nil {
			slog.Error("could not reset DB", "err", err)
			return nil, err
		}
	}

	dbConn, err := p.openDb()
	if err != nil {
		slog.Error("could not open DB", "err", err)
		return nil, err
	}

	db := dbaccess.New(dbConn)

	ClientRepository = sqlrepository.ClientSqlRepository{Db: db}
	SessionRepository = sqlrepository.SessionSqlRepository{Db: db, SqlConn: dbConn}
	SubscriptionRepository = sqlrepository.SubscriptionSqlRepository{Db: db}
	RetainRepository = sqlrepository.RetainSqlRepository{Db: db}
	UserRepository = sqlrepository.UserSqlRepository{Db: db}
	RetryRepository = sqlrepository.RetrySqlRepository{Db: db}

	if len(initAdminPassword) > 0 {
		sqlrepository.AdminPasswordReset(db, initAdminPassword)
		os.Exit(0)
	}

	if cleanExpiredSessions {
		sqlrepository.StartSessionCleaner(db)
	}

	if cleanExpiredSessions {
		sqlrepository.StartRetryCleaner(db)
	}

	return db, nil
}

func (p *SqlPersistence) resetDb() error {
	os.Remove(p.DbFile)

	dbConn, err := p.openDb()
	if err != nil {
		slog.Error("could not open DB", "err", err)
		return err
	}

	if _, err := dbConn.ExecContext(context.Background(), sqlc.DBSchema); err != nil {
		return err
	}

	return nil
}

func (p *SqlPersistence) openDb() (*sql.DB, error) {
	if p.DbFile == "" {
		p.DbFile = "sqlite.db3"
	}

	db, err := sql.Open("sqlite3", p.DbFile)
	if err != nil {
		return db, err
	}
	db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")

	p.db = db

	return db, nil
}
