package mqtt

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const RETAIN_INSERT = `insert into retains(topic, application_message, created_at)
	values(?, ?, ?)`
const RETAIN_DELETE = "delete from retains where topic = ?"
const RETAIN_SELECT_TOPIC = "select topic, application_message, created_at from retains where topic = ?"

type SqliteRetains struct {
	db *sql.DB
}

func (is SqliteRetains) addRetain(r Retain) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(RETAIN_INSERT)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.topic, r.applicationMessage, r.createdAt.Unix())
	if err != nil {
		log.Println(err)
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteRetains) remRetain(topic string) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(RETAIN_DELETE)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(topic)
	if err != nil {
		log.Println(err)
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteRetains) findRetainByTopic(topic string) []Retain {
	retains := []Retain{}
	rows, err := is.db.Query(RETAIN_SELECT_TOPIC, topic)
	if err != nil {
		log.Fatal(err)
		return retains
	}
	defer rows.Close()
	for rows.Next() {
		var r Retain
		err = rows.Scan(&r.topic, &r.applicationMessage, &r.createdAt)
		if err != nil {
			log.Println(err)
			return retains
		}
		retains = append(retains, r)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return retains
	}

	return retains
}
