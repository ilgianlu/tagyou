package mqtt

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const SUBSCRIPTION_INSERT = `insert into subscriptions(
	topic, clientid, enabled, created_at,
	qos, retain_handling, retain_as_published, no_local)
	values(?, ?, ?, ?, ?, ?, ?, ?)`
const SUBSCRIPTION_DELETE = "delete from subscriptions where topic = ? and clientid = ?"
const SUBSCRIPTION_DELETE_TOPIC = "delete from subscriptions where topic = ?"
const SUBSCRIPTION_DELETE_CLIENTID = "delete from subscriptions where clientid = ?"
const SUBSCRIPTION_SELECT_TOPIC = "select topic, clientid, enabled from subscriptions where topic = ? and enabled = 1"
const SUBSCRIPTION_SELECT_CLIENTID = "select topic, clientid, enabled from subscriptions where clientid = ? and enabled = 1"
const SUBSCRIPTION_UPDATE_CLIENTID = "update subscriptions set enabled = ? where clientid = ?"

type SqliteSubscriptions struct {
	db *sql.DB
}

func (is SqliteSubscriptions) addSubscription(s Subscription) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(SUBSCRIPTION_INSERT)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		s.topic, s.clientId, s.enabled, s.createdAt.Unix(),
		s.QoS, s.RetainHandling, s.RetainAsPublished, s.NoLocal,
	)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteSubscriptions) remSubscription(topic string, clientId string) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare(SUBSCRIPTION_DELETE)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(topic, clientId)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteSubscriptions) findTopicSubscribers(topic string) []Subscription {
	topicSegments := strings.Split(topic, TOPIC_SEPARATOR)
	if len(topicSegments) == 1 {
		return is.findSubscriptionsByTopic(topic)
	} else {
		return is.multiSegmentSubs(topicSegments)
	}
}

func (is SqliteSubscriptions) findSubscriptionsByTopic(topic string) []Subscription {
	subscribers := []Subscription{}
	rows, err := is.db.Query(SUBSCRIPTION_SELECT_TOPIC, topic)
	if err != nil {
		log.Fatal(err)
		return subscribers
	}
	defer rows.Close()
	for rows.Next() {
		var s Subscription
		var enabled int
		err = rows.Scan(&s.topic, &s.clientId, &enabled)
		if enabled == 1 {
			s.enabled = true
		}
		if err != nil {
			log.Println(err)
			return subscribers
		}
		subscribers = append(subscribers, s)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return subscribers
	}

	return subscribers
}

func (is SqliteSubscriptions) findSubscriptionsByClientId(clientId string) []Subscription {
	subscribers := []Subscription{}
	rows, err := is.db.Query(SUBSCRIPTION_SELECT_CLIENTID, clientId)
	if err != nil {
		log.Fatal(err)
		return subscribers
	}
	defer rows.Close()
	for rows.Next() {
		var s Subscription
		var enabled int
		err = rows.Scan(&s.topic, &s.clientId, &enabled)
		if enabled == 1 {
			s.enabled = true
		}
		if err != nil {
			log.Fatal(err)
			return subscribers
		}
		subscribers = append(subscribers, s)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return subscribers
	}

	return subscribers
}

func (is SqliteSubscriptions) multiSegmentSubs(topicSegments []string) []Subscription {
	subs := []Subscription{}
	for i := 1; i <= len(topicSegments); i++ {
		subT := append(make([]string, 0), topicSegments[:i]...)
		if len(subT) < len(topicSegments) {
			subT = append(subT, TOPIC_WILDCARD)
		}
		t := strings.Join(subT, TOPIC_SEPARATOR)
		ss := is.findSubscriptionsByTopic(t)
		subs = append(subs, ss...)
	}
	return subs
}

func (is SqliteSubscriptions) remSubscriptionsByTopic(topic string) {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(SUBSCRIPTION_DELETE_TOPIC)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(topic)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	_ = tx.Commit()
}

func (is SqliteSubscriptions) remSubscriptionsByClientId(clientId string) {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(SUBSCRIPTION_DELETE_CLIENTID)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientId)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	_ = tx.Commit()
}

func (is SqliteSubscriptions) disableClientSubscriptions(clientId string) {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(SUBSCRIPTION_UPDATE_CLIENTID)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(0, clientId)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	_ = tx.Commit()
}

func (is SqliteSubscriptions) enableClientSubscriptions(clientId string) {
	tx, err := is.db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(SUBSCRIPTION_UPDATE_CLIENTID)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(1, clientId)
	if err != nil {
		log.Fatal(err)
		_ = tx.Rollback()
		return
	}
	_ = tx.Commit()
}
