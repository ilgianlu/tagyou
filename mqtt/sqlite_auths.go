package mqtt

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const AUTH_INSERT = `insert into auths(
	clientid, username, password, subscribe_acl, publish_acl, created_at)
	values(?, ?, ?, ?, ?, ?)`
const AUTH_DELETE = "delete from auths where username = ? and clientid = ?"
const AUTH_FIND_CLIENTID = "select * from auths where clientid = ?"
const AUTH_VERIFY = "select * from auths where clientid = ? and username = ? and password = ?"

type SqliteAuths struct {
	db *sql.DB
}

func (is SqliteAuths) createAuth(a Auth) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(AUTH_INSERT)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		a.clientId, a.username, a.password, a.subscribeAcl, a.publishAcl, a.createdAt.Unix(),
	)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteAuths) remAuth(username string, clientId string) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(AUTH_DELETE)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, clientId)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteAuths) findAuth(clientId string) (Auth, bool) {
	var a Auth
	rows, err := is.db.Query(AUTH_FIND_CLIENTID, clientId)
	if err != nil {
		log.Println(err)
		return a, false
	}
	defer rows.Close()
	if rows.Next() {
		var createdAt int64
		err = rows.Scan(&a.clientId, &a.username, &a.password, &a.subscribeAcl, &a.publishAcl, &createdAt)
		a.createdAt = time.Unix(createdAt, 0)
		if err != nil {
			log.Println(err)
			return a, false
		}

	}
	return a, true
}
