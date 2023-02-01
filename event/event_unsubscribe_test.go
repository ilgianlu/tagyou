package event

import (
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sqlrepository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestClientUnsubscription(t *testing.T) {
	os.Setenv("DEBUG", "1")
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	sqlrepository.Migrate(db)

	persistence.InitSqlRepositories(db)

	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM subscriptions")
	sess1 := model.Session{ID: 1, ClientId: "pippo", Connected: true}
	sub1 := model.Subscription{SessionID: sess1.ID, ClientId: "pippo", Topic: "topic1"}
	ssub1 := model.Subscription{SessionID: sess1.ID, ClientId: "pippo", ShareName: "share1", Topic: "sharedTopic1"}

	sess2 := model.Session{ID: 2, Connected: true, ClientId: "pluto"}
	sub2 := model.Subscription{SessionID: sess2.ID, ClientId: "pluto", Topic: "topic1"}
	ssub2 := model.Subscription{SessionID: sess2.ID, ClientId: "pluto", ShareName: "share2", Topic: "sharedTopic1"}

	subscriptions := []model.Subscription{sub1, ssub1, sub2, ssub2}
	sessions := []model.Session{sess1, sess2}
	db.Create(&sessions)
	db.Create(&subscriptions)

	res := clientUnsubscription("pippo", model.Subscription{Topic: "topic1"})
	if res != 0 {
		t.Error("unsuccessful subscription, expected success")
	}

	res = clientUnsubscription("pippo", model.Subscription{Topic: "topic2"})
	if res != 17 {
		t.Errorf("expecting 17 (no subscription to unsub), received %d", res)
	}

	res = clientUnsubscription("pluto", model.Subscription{Topic: "sharedTopic1", ShareName: "share2"})
	if res != 0 {
		t.Errorf("expecting 0 (success), received %d", res)
	}
	s := model.Subscription{}
	if err := db.Where("share_name = ? and topic = ?", "share2", "sharedTopic1").First(&s).Error; err == nil {
		t.Errorf("shared subscription was not removed!")
	}
}
