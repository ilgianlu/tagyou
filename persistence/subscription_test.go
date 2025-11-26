package persistence

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/ilgianlu/tagyou/sqlc"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreate(t *testing.T) {
	os.Remove("test.db3")

	dbConn, err := sql.Open("sqlite3", "test.db3")
	if err != nil {
		t.Errorf("[API] failed to connect database")
	}

	dbConn.ExecContext(context.Background(), sqlc.DBSchema)

	db := dbaccess.New(dbConn)

	un := dbaccess.CreateSubscriptionParams{ClientID: sql.NullString{String: "uno", Valid: true}, Topic: sql.NullString{String: "uno", Valid: true}}
	if _, err := db.CreateSubscription(context.Background(), un); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	if _, err := db.CreateSubscription(context.Background(), un); err != nil {
		t.Error("subscription (duplicate client id and topic) create should throw err!")
	}

	du := dbaccess.CreateSubscriptionParams{ClientID: sql.NullString{String: "due", Valid: true}, Topic: sql.NullString{String: "uno", Valid: true}}
	if _, err := db.CreateSubscription(context.Background(), du); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}
}
