package sender

import (
	"sort"
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
	sess1 := model.Session{ID: 1, Connected: true}
	sub1 := model.Subscription{SessionID: sess1.ID, ClientId: "pippo", ShareName: "share1"}
	sess2 := model.Session{ID: 2, Connected: true}
	sub2 := model.Subscription{SessionID: sess2.ID, ClientId: "pluto", ShareName: "share2"}
	sess3 := model.Session{ID: 3, Connected: true}
	sub3 := model.Subscription{SessionID: sess3.ID, ClientId: "minnie", ShareName: "share1"}
	sess4 := model.Session{ID: 4, Connected: true}
	sub4 := model.Subscription{SessionID: sess4.ID, ClientId: "topolino", ShareName: "share2"}
	sess5 := model.Session{ID: 5, Connected: false}
	sub5 := model.Subscription{SessionID: sess5.ID, ClientId: "paperino", ShareName: "share1"}
	sess6 := model.Session{ID: 5, Connected: false}
	sub6 := model.Subscription{SessionID: sess6.ID, ClientId: "paperina", ShareName: "share2"}
	ungrouped := []model.Subscription{sub1, sub2, sub3, sub4, sub5, sub6}
	sessions := []model.Session{sess1, sess2, sess3, sess4, sess5, sess6}
	db.Create(&ungrouped)
	db.Create(&sessions)
	groups := groupSubscribers(ungrouped)
	if group, ok := groups["share1"]; ok {
		res := []uint{sess1.ID, sess3.ID}
		data := []uint{}
		for _, cl := range group {
			data = append(data, cl.SessionID)
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i] > data[j]
		})
		for i, d := range data {
			if d != res[i] {
				t.Errorf("expecting %d received %d", res[i], d)
			}
		}
	} else {
		t.Errorf("no group share1 found!")
	}
	if group, ok := groups["share2"]; ok {
		res := []uint{sess2.ID, sess4.ID}
		data := []uint{}
		for _, cl := range group {
			data = append(data, cl.SessionID)
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i] > data[j]
		})
		for i, d := range data {
			if d != res[i] {
				t.Errorf("expecting %d received %d", res[i], d)
			}
		}
	} else {
		t.Errorf("no group share2 found!")
	}

}
