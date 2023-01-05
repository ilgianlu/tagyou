package sqlrepository

import (
	"testing"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	Migrate(db)

	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM subscriptions")

	un := Subscription{ClientId: "uno", Topic: "uno"}
	if err := db.Create(&un).Error; err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	unBis := Subscription{ClientId: "uno", Topic: "uno"}
	if err := db.Create(&unBis).Error; err == nil {
		t.Error("subscription (duplicate client id and topic) create should throw err!")
	}

	du := Subscription{ClientId: "due", Topic: "uno"}
	if err := db.Create(&du).Error; err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}
}
