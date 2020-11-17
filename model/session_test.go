package model

import (
	"testing"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSessionDelete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	Migrate(db)

	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM subscriptions")

	s1 := Session{ClientId: "sessionOne"}
	if err := db.Create(&s1).Error; err != nil {
		t.Errorf("session create should not throw err: %s", err)
	}

	sId := s1.ID

	un := Subscription{ClientId: "uno", Topic: "uno", SessionID: s1.ID}
	if err := db.Create(&un).Error; err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	du := Subscription{ClientId: "due", Topic: "uno", SessionID: s1.ID}
	if err := db.Create(&du).Error; err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	db.Delete(&s1)

	sess := []Session{}
	db.Where("client_id = ?", "sessionOne").Find(&sess)
	if len(sess) > 0 {
		t.Errorf("session find should find nothing")
	}

	subs := []Subscription{}
	db.Where("session_id = ?", sId).Find(&subs)
	if len(subs) > 0 {
		t.Errorf("subscription find should find nothing")
	}
}
