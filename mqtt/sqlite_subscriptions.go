package mqtt

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteSubscriptions struct {
	db *sql.DB
}

func (is sqliteSubscriptions) addSubscription(topic string, subscriber string) error {
	return nil
}

func (is sqliteSubscriptions) remSubscription(topic string, subscriber string) error {
	return nil
}

func (is sqliteSubscriptions) findSubscribers(topic string) []string {
	// only topics have topic separator in name
	topicSegments := strings.Split(topic, TOPIC_SEPARATOR)
	return topicSegments
}

func (is sqliteSubscriptions) findSubscribed(subscribed string) ([]string, bool) {
	return []string{}, false
}

func (is sqliteSubscriptions) remSubscribed(subscribed string) {
}
