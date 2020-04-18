package mqtt

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func testSeed() (sqliteSubscriptions, error) {
	db, err := sql.Open("sqlite3", "../testdata/sqlite.db3")
	return sqliteSubscriptions{db: db}, err
}

func TestSqliteAddRemoveSubsciption(t *testing.T) {
	subscriptions, err := testSeed()
	if err != nil {
		t.Errorf("error opening test db %s", err)
	}
	err = subscriptions.addSubscription("topic", "gianluca")
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	err = subscriptions.addSubscription("topic", "carlo")
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	subs := subscriptions.findSubscribers("topic")
	if len(subs) != 2 {
		t.Errorf("expected 2 subscribers, got %d", len(subs))
	}
	err = subscriptions.remSubscription("topic", "gianluca")
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
}
