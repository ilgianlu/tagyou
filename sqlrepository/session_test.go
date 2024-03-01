package sqlrepository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/ilgianlu/tagyou/sqlc"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

func TestSessionDelete(t *testing.T) {
	os.Remove("test.db")

	dbConn, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Errorf("[API] failed to connect database")
	}

	dbConn.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")
	dbConn.ExecContext(context.Background(), sqlc.DBSchema)

	db := dbaccess.New(dbConn)

	s1 := dbaccess.CreateSessionParams{ClientID: sql.NullString{String: "sessionOne", Valid: true}}
	s, err := db.CreateSession(context.Background(), s1)
	if err != nil {
		t.Errorf("session create should not throw err: %s", err)
	}

	un := dbaccess.CreateSubscriptionParams{
		ClientID:  sql.NullString{String: "uno", Valid: true},
		Topic:     sql.NullString{String: "uno", Valid: true},
		SessionID: sql.NullInt64{Int64: s.ID, Valid: true},
	}
	if _, err := db.CreateSubscription(context.Background(), un); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	du := dbaccess.CreateSubscriptionParams{
		ClientID:  sql.NullString{String: "due", Valid: true},
		Topic:     sql.NullString{String: "uno", Valid: true},
		SessionID: sql.NullInt64{Int64: s.ID, Valid: true},
	}
	if _, err := db.CreateSubscription(context.Background(), du); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	db.DeleteSessionByClientId(context.Background(), s1.ClientID)

	_, err = db.GetSessionByClientId(context.Background(), sql.NullString{String: "sessionOne", Valid: true})
	if err != sql.ErrNoRows {
		t.Errorf("session find should find nothing, err: %s", err)
	}

	subs, _ := db.GetSubscriptionsBySessionId(context.Background(), sql.NullInt64{Int64: s.ID, Valid: true})
	if len(subs) > 0 {
		t.Errorf("subscription find should find nothing, found %d", len(subs))
	}
}
