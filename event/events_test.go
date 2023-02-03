package event

import (
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	p := persistence.SqlPersistence{}
	p.InnerInit(db)

	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM subscriptions")
	sess1 := model.Session{ID: 1, ClientId: "pippo", Connected: true}
	sub1 := model.Subscription{SessionID: sess1.ID, ClientId: "pippo", ShareName: "share1"}
	sess2 := model.Session{ID: 2, ClientId: "pluto", Connected: true}
	sub2 := model.Subscription{SessionID: sess2.ID, ClientId: "pluto", ShareName: "share2"}
	sess3 := model.Session{ID: 3, ClientId: "minnie", Connected: true}
	sub3 := model.Subscription{SessionID: sess3.ID, ClientId: "minnie", ShareName: "share1"}
	sess4 := model.Session{ID: 4, ClientId: "topolino", Connected: true}
	sub4 := model.Subscription{SessionID: sess4.ID, ClientId: "topolino", ShareName: "share2"}
	sess5 := model.Session{ID: 5, ClientId: "paperino", Connected: false}
	sub5 := model.Subscription{SessionID: sess5.ID, ClientId: "paperino", ShareName: "share1"}
	sess6 := model.Session{ID: 5, ClientId: "paperina", Connected: false}
	sub6 := model.Subscription{SessionID: sess6.ID, ClientId: "paperina", ShareName: "share2"}
	ungrouped := []model.Subscription{sub1, sub2, sub3, sub4, sub5, sub6}
	sessions := []model.Session{sess1, sess2, sess3, sess4, sess5, sess6}
	db.Create(&ungrouped)
	db.Create(&sessions)
	groups := groupSubscribers(ungrouped)
	if group, ok := groups["share1"]; ok {
		if len(group) != 3 {
			t.Errorf("found %d subs in share1, expected %d!", len(group), 3)
		}
		for _, s := range group {
			if s.ClientId != "pippo" && s.ClientId != "minnie" && s.ClientId != "paperino" {
				t.Errorf("unexpected %s client in share1!", s.ClientId)
			}
		}
	} else {
		t.Errorf("no group share1 found!")
	}
	if group, ok := groups["share2"]; ok {
		if len(group) != 3 {
			t.Errorf("found %d subs in share2, expected %d!", len(group), 3)
		}
		for _, s := range group {
			if s.ClientId != "pluto" && s.ClientId != "topolino" && s.ClientId != "paperina" {
				t.Errorf("unexpected %s client in share1!", s.ClientId)
			}
		}
	} else {
		t.Errorf("no group share2 found!")
	}

}
