package sqlrepository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"

	_ "github.com/mattn/go-sqlite3"
)

func TestRetryCleaner(t *testing.T) {
	dbConn, err := sql.Open("", conf.DB_PATH+conf.DB_NAME)
	if err != nil {
		t.Errorf("cannot connect to db")
	}

	db := dbaccess.New(dbConn)

	dbConn.Exec("DELETE FROM sessions")
	dbConn.Exec("DELETE FROM retries")

	s1 := dbaccess.CreateSessionParams{ClientID: sql.NullString{String: "sessionOne", Valid: true}}
	insertedSession, err := db.CreateSession(context.Background(), s1)
	if err != nil {
		t.Errorf("session create should not throw err: %s", err)
	}

	sId := insertedSession.ID

	un := dbaccess.CreateRetryParams{
		ClientID:           sql.NullString{String: "uno", Valid: true},
		ApplicationMessage: []byte{1, 2, 3},
		PacketIdentifier:   sql.NullInt64{Int64: 50, Valid: true},
		Qos:                sql.NullInt64{Int64: 1, Valid: true},
		Dup:                sql.NullInt64{Int64: 0, Valid: true},
		Retries:            sql.NullInt64{Int64: 3, Valid: true},
		AckStatus:          sql.NullInt64{Int64: 0, Valid: true},
		CreatedAt:          sql.NullInt64{Int64: time.Now().Unix() - 30, Valid: true},
		SessionID:          sql.NullInt64{Int64: sId, Valid: true},
	}
	db.CreateRetry(context.Background(), un)

	du := dbaccess.CreateRetryParams{
		ClientID:           sql.NullString{String: "due", Valid: true},
		ApplicationMessage: []byte{1, 2, 3},
		PacketIdentifier:   sql.NullInt64{Int64: 50, Valid: true},
		Qos:                sql.NullInt64{Int64: 1, Valid: true},
		Dup:                sql.NullInt64{Int64: 0, Valid: true},
		Retries:            sql.NullInt64{Int64: 3, Valid: true},
		AckStatus:          sql.NullInt64{Int64: 0, Valid: true},
		CreatedAt:          sql.NullInt64{Int64: time.Now().Unix() - 30, Valid: true},
		SessionID:          sql.NullInt64{Int64: sId, Valid: true},
	}
	db.CreateRetry(context.Background(), du)

	before, err := db.GetAllRetries(context.Background())

	fmt.Println(before)

	if len(before) != 2 {
		t.Errorf("expected 2 retry, found: %d", len(before))
	}

	cleanRetries(db)

	after, err := db.GetAllRetries(context.Background())

	if len(after) != 1 {
		t.Errorf("expected 1 retry (expiration %d secs), found: %d", conf.RETRY_EXPIRATION, len(after))
	}
}
