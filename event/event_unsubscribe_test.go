package event

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sqlc"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

func TestClientUnsubscription(t *testing.T) {
	os.Setenv("DEBUG", "1")
	os.Remove("test.db3")

	dbConn, err := sql.Open("sqlite3", "test.db3")
	if err != nil {
		t.Errorf("[API] failed to connect database")
	}

	dbConn.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")
	dbConn.ExecContext(context.Background(), sqlc.DBSchema)

	db := dbaccess.New(dbConn)

	p := persistence.SqlPersistence{}
	p.InnerInit(db, false, false, []byte(""))

	sess1, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{
		ClientID:  sql.NullString{String: "pippo", Valid: true},
		Connected: sql.NullInt64{Int64: 1, Valid: true},
	})
	db.CreateSubscription(context.Background(), dbaccess.CreateSubscriptionParams{
		Topic:     sql.NullString{String: "topic1", Valid: true},
		ClientID:  sql.NullString{String: "pippo", Valid: true},
		SessionID: sql.NullInt64{Int64: sess1.ID, Valid: true},
		Shared:    sql.NullInt64{Int64: 0, Valid: true},
	})
	db.CreateSubscription(context.Background(), dbaccess.CreateSubscriptionParams{
		Topic:     sql.NullString{String: "sharedTopic1", Valid: true},
		ShareName: sql.NullString{String: "share1", Valid: true},
		Shared:    sql.NullInt64{Int64: 1, Valid: true},
		ClientID:  sql.NullString{String: "pippo", Valid: true},
		SessionID: sql.NullInt64{Int64: sess1.ID, Valid: true},
	})

	sess2, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{
		ClientID:  sql.NullString{String: "pluto", Valid: true},
		Connected: sql.NullInt64{Int64: 1, Valid: true},
	})
	db.CreateSubscription(context.Background(), dbaccess.CreateSubscriptionParams{
		Topic:     sql.NullString{String: "topic1", Valid: true},
		Shared:    sql.NullInt64{Int64: 0, Valid: true},
		ClientID:  sql.NullString{String: "pluto", Valid: true},
		SessionID: sql.NullInt64{Int64: sess2.ID, Valid: true},
	})
	db.CreateSubscription(context.Background(), dbaccess.CreateSubscriptionParams{
		Topic:     sql.NullString{String: "sharedTopic1", Valid: true},
		ShareName: sql.NullString{String: "share2", Valid: true},
		Shared:    sql.NullInt64{Int64: 1, Valid: true},
		ClientID:  sql.NullString{String: "pluto", Valid: true},
		SessionID: sql.NullInt64{Int64: sess2.ID, Valid: true},
	})

	res := clientUnsubscription("pippo", model.Subscription{Topic: "topic1", Shared: false})
	if res != 0 {
		t.Error("unsuccessful subscription, expected success")
	}

	res = clientUnsubscription("pippo", model.Subscription{Topic: "topic2"})
	if res != 17 {
		t.Errorf("expecting 17 (no subscription to unsub), received %d", res)
	}

	res = clientUnsubscription("pluto", model.Subscription{Topic: "sharedTopic1", ShareName: "share2", Shared: true})
	if res != 0 {
		t.Errorf("expecting 0 (success), received %d", res)
	}

	ss, _ := db.GetSubscriptionsBySessionId(context.Background(), sql.NullInt64{Int64: sess2.ID, Valid: true})
	if len(ss) != 1 {
		t.Errorf("shared subscription was not removed!")
	}
	if ss[0].Topic.String != "topic1" {
		t.Errorf("shared subscription was not removed!")
	}
}
