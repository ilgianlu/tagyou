package mqtt

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func testSeed() (SqliteSubscriptions, error) {
	Seed("../testdata/sqlite.db3")
	db, err := sql.Open("sqlite3", "../testdata/sqlite.db3")
	return SqliteSubscriptions{db: db}, err
}

func TestSqliteAddRemoveSubscription(t *testing.T) {
	subscriptions, err := testSeed()
	if err != nil {
		t.Errorf("error opening test db %s", err)
	}
	s0 := Subscription{
		topic:     "topic",
		clientId:  "gianluca",
		enabled:   true,
		createdAt: time.Now(),
	}
	err = subscriptions.addSubscription(s0)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	s1 := Subscription{
		topic:     "topic",
		clientId:  "carlo",
		enabled:   true,
		createdAt: time.Now(),
	}
	err = subscriptions.addSubscription(s1)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	subs := subscriptions.findSubscriptionsByTopic("topic")
	if len(subs) != 2 {
		t.Errorf("expected 2 subscribers, got %d", len(subs))
	}
	err = subscriptions.remSubscription("topic", "gianluca")
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
}

func TestSqliteEnableDisableSubscription(t *testing.T) {
	subscriptions, err := testSeed()
	if err != nil {
		t.Errorf("error opening test db %s", err)
	}
	s0 := Subscription{
		topic:    "topic",
		clientId: "gianluca",
		enabled:  true,
	}
	err = subscriptions.addSubscription(s0)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	s1 := Subscription{
		topic:     "topic",
		clientId:  "carlo",
		enabled:   true,
		createdAt: time.Now(),
	}
	err = subscriptions.addSubscription(s1)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	subs := subscriptions.findSubscriptionsByTopic("topic")
	if len(subs) != 2 {
		t.Errorf("expected 2 subscribers, got %d", len(subs))
	}
	subscriptions.disableClientSubscriptions("gianluca")
	subs = subscriptions.findSubscriptionsByTopic("topic")
	if len(subs) != 1 {
		t.Errorf("expected 1 subscribers, got %d", len(subs))
	}
	subscriptions.enableClientSubscriptions("gianluca")
	subs = subscriptions.findSubscriptionsByTopic("topic")
	if len(subs) != 2 {
		t.Errorf("expected 2 subscribers, got %d", len(subs))
	}
}

func TestSqliteCannotDuplicateSubscription(t *testing.T) {
	subscriptions, err := testSeed()
	if err != nil {
		t.Errorf("error opening test db %s", err)
	}
	s0 := Subscription{
		topic:     "topic",
		clientId:  "gianluca",
		enabled:   true,
		createdAt: time.Now(),
	}
	err = subscriptions.addSubscription(s0)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	s1 := Subscription{
		topic:     "topic",
		clientId:  "gianluca",
		enabled:   true,
		createdAt: time.Now(),
	}
	err = subscriptions.addSubscription(s1)
	if err == nil {
		t.Errorf("expected an error!")
	}
	subs := subscriptions.findSubscriptionsByTopic("topic")
	if len(subs) != 1 {
		t.Errorf("duplicated subscription!, got %d subscriptions", len(subs))
	}
}
