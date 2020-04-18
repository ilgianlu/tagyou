package mqtt

import (
	"database/sql"
	"strings"
	"log"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteSubscriptions struct {
	db *sql.DB
}

func (is sqliteSubscriptions) addSubscription(topic string, subscriber string) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("insert into %s(topic, clientid) values(?, ?)", TABLE_SUBSCRIPTIONS))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(topic, subscriber)
	if err != nil {
		log.Fatal(err)
		return err
	}
	tx.Commit()
	return nil
}

func (is sqliteSubscriptions) remSubscription(topic string, subscriber string) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("delete from %s where topic = ? and clientid = ?", TABLE_SUBSCRIPTIONS))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(topic, subscriber)
	if err != nil {
		log.Fatal(err)
		return err
	}
	tx.Commit()
	return nil
}

func (is sqliteSubscriptions) findSubscribers(topic string) []string {
	topicSegments := strings.Split(topic, TOPIC_SEPARATOR)
	if len(topicSegments) == 1 {
		if s, ok := is.findSubscribed(topic); ok {
			return s
		} else {
			return []string{}
		}
	} else {
		return is.multiSegmentSubs(topicSegments)
	}
}

func (is sqliteSubscriptions) findSubscribed(topic string) ([]string, bool) {
	subscribers := []string{}
	rows, err := is.db.Query(fmt.Sprintf("select topic, clientid from %s where topic = ?", TABLE_SUBSCRIPTIONS), topic)
	if err != nil {
		log.Fatal(err)
		return subscribers, false
	}
	defer rows.Close()
	for rows.Next() {
		var topic string
		var clientid string
		err = rows.Scan(&topic, &clientid)
		if err != nil {
			log.Fatal(err)
			return subscribers, false
		}
		subscribers = append(subscribers, clientid)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return subscribers, false
	}

	return subscribers, true
}

func (is sqliteSubscriptions) multiSegmentSubs(topicSegments []string) []string {
	subs := make([]string, 0)
	for i := 1; i <= len(topicSegments); i++ {
		subT := append(make([]string, 0), topicSegments[:i]...)
		if len(subT) < len(topicSegments) {
			subT = append(subT, TOPIC_WILDCARD)
		}
		t := strings.Join(subT, TOPIC_SEPARATOR)
		if ss, ok := is.findSubscribed(t); ok {
			subs = append(subs, ss...)
		}
	}
	return subs
}

func (is sqliteSubscriptions) remSubscribed(topic string) {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(fmt.Sprintf("delete from %s where topic = ?", TABLE_SUBSCRIPTIONS))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(topic)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
