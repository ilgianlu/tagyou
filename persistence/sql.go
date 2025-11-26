package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"os"

	"github.com/ilgianlu/tagyou/sqlc"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"

	_ "github.com/mattn/go-sqlite3"
)

type SqlPersistence struct {
	DbFile       string
	InitDatabase bool
	db           *sql.DB
}

func (p *SqlPersistence) Init(cleanExpiredSessions bool, cleanExpiredRetries bool, initAdminPassword []byte) (*dbaccess.Queries, error) {
	dbConn, err := p.openDB()
	if err != nil {
		slog.Error("could not open DB", "err", err)
		return nil, err
	}

	db := dbaccess.New(dbConn)

	ClientRepository = ClientSqlRepository{Db: db}
	SessionRepository = SessionSqlRepository{Db: db, SqlConn: dbConn}
	SubscriptionRepository = SubscriptionSqlRepository{Db: db}
	RetainRepository = RetainSqlRepository{Db: db}
	UserRepository = UserSqlRepository{Db: db}
	RetryRepository = RetrySqlRepository{Db: db}

	if len(initAdminPassword) > 0 {
		AdminPasswordReset(db, initAdminPassword)
	}

	if cleanExpiredSessions {
		StartSessionCleaner(db)
	}

	if cleanExpiredSessions {
		StartRetryCleaner(db)
	}

	return db, nil
}

func (p *SqlPersistence) resetDB() error {
	return os.Remove(p.DbFile)
}

func (p *SqlPersistence) initTables(dbConn *sql.DB) error {
	if _, err := dbConn.ExecContext(context.Background(), sqlc.DBSchema); err != nil {
		os.Remove(p.DbFile)
		return err
	}

	return nil
}

func (p *SqlPersistence) needReset() bool {
	if p.InitDatabase {
		return true
	}
	_, err := os.Open(p.DbFile)
	if err != nil && os.IsNotExist(err) {
		return true
	}
	return false
}

func (p *SqlPersistence) openDB() (*sql.DB, error) {
	if p.DbFile == "" {
		p.DbFile = "sqlite.db3"
	}

	resetted := false
	if p.needReset() {
		p.resetDB()
		resetted = true
	}

	db, err := sql.Open("sqlite3", p.DbFile)
	if err != nil {
		return db, err
	}
	db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")

	p.db = db

	if resetted {
		p.initTables(db)
	}

	return db, nil
}

func (p *SqlPersistence) Close() {
	if p.db != nil {
		p.db.Close()
	}
}
