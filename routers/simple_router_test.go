package routers

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

func TestPickDest(t *testing.T) {
	sub1 := model.Subscription{
		ClientId: "pippo",
	}
	sub2 := model.Subscription{
		ClientId: "pluto",
	}
	subGroup := []model.Subscription{sub1, sub2}
	dest := pickDest(subGroup, 0)

	if dest.ClientId != "pippo" {
		t.Errorf("expecting %s received %s", "pippo", dest.ClientId)
	}
}

func TestGroupSubscribers(t *testing.T) {
	p := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	db, err := p.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	sess1, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{Connected: sql.NullInt64{Int64: 1, Valid: true}})
	sub1 := model.Subscription{SessionID: sess1.ID, ClientId: "pippo", ShareName: "share1"}
	sess2, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{Connected: sql.NullInt64{Int64: 1, Valid: true}})
	sub2 := model.Subscription{SessionID: sess2.ID, ClientId: "pluto", ShareName: "share2"}
	sess3, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{Connected: sql.NullInt64{Int64: 1, Valid: true}})
	sub3 := model.Subscription{SessionID: sess3.ID, ClientId: "minnie", ShareName: "share1"}
	sess4, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{Connected: sql.NullInt64{Int64: 1, Valid: true}})
	sub4 := model.Subscription{SessionID: sess4.ID, ClientId: "topolino", ShareName: "share2"}
	sess5, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{Connected: sql.NullInt64{Int64: 0, Valid: true}})
	sub5 := model.Subscription{SessionID: sess5.ID, ClientId: "paperino", ShareName: "share1"}
	sess6, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{Connected: sql.NullInt64{Int64: 0, Valid: true}})
	sub6 := model.Subscription{SessionID: sess6.ID, ClientId: "paperina", ShareName: "share2"}
	ungrouped := []model.Subscription{sub1, sub2, sub3, sub4, sub5, sub6}

	groups := groupSubscribers(ungrouped)
	if group, ok := groups["share1"]; ok {
		if len(group) != 2 {
			t.Errorf("expecting %d subs for share1, received %d", 2, len(group))
		}
		if group[0].SessionID != sess1.ID {
			t.Errorf("expecting %d sub id position %d, received %d", sess1.ID, 0, group[0].SessionID)
		}
		if group[1].SessionID != sess3.ID {
			t.Errorf("expecting %d sub id position %d, received %d", sess3.ID, 1, group[1].SessionID)
		}
	} else {
		t.Errorf("no group share1 found!")
	}
	if group, ok := groups["share2"]; ok {
		if len(group) != 2 {
			t.Errorf("expecting %d subs for share2, received %d", 2, len(group))
		}
		if group[0].SessionID != sess2.ID {
			t.Errorf("expecting %d sub id position %d, received %d", sess2.ID, 0, group[0].SessionID)
		}
		if group[1].SessionID != sess4.ID {
			t.Errorf("expecting %d sub id position %d, received %d", sess4.ID, 1, group[1].SessionID)
		}
	} else {
		t.Errorf("no group share2 found!")
	}
}
