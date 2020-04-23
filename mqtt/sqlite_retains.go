package mqtt

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const RETAIN_INSERT = `insert into retains(topic, application_message, created_at)
	values(?, ?, ?)`
const RETAIN_DELETE = "delete from retains where topic = ?"
const RETAIN_SELECT_TOPIC = "select topic, application_message from retains where topic = ?"
const RETAIN_SELECT_LIKE_TOPIC = "select topic, application_message from retains where topic like ?"

type SqliteRetains struct {
	db *sql.DB
}

func (is SqliteRetains) addRetain(r Retain) error {
	rErr := is.remRetain(r.topic)
	if rErr != nil {
		log.Println(rErr)
		return rErr
	}
	if len(r.applicationMessage) == 0 {
		return nil
	}
	tx, err := is.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare(RETAIN_INSERT)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.topic, r.applicationMessage, r.createdAt.Unix())
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteRetains) remRetain(topic string) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare(RETAIN_DELETE)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(topic)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteRetains) findRetainsByTopic(topic string) []Retain {
	retains := []Retain{}
	var rows *sql.Rows
	var err error
	if withWildCard(topic) {
		rows, err = is.db.Query(RETAIN_SELECT_LIKE_TOPIC, fmt.Sprint(topic[:len(topic)-1], "%"))
	} else {
		rows, err = is.db.Query(RETAIN_SELECT_TOPIC, topic)
	}
	if err != nil {
		log.Println(err)
		return retains
	}
	defer rows.Close()
	for rows.Next() {
		var r Retain
		err = rows.Scan(&r.topic, &r.applicationMessage)
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

func withWildCard(topic string) bool {
	return topic[len(topic)-1:] == TOPIC_WILDCARD
}
