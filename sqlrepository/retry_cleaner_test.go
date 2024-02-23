package sqlrepository

import (
	"fmt"
	"testing"
	"time"

	"github.com/ilgianlu/tagyou/conf"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRetryCleaner(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		t.Errorf("[API] failed to connect database")
	}

	Migrate(db)

	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM retries")

	s1 := Session{ClientId: "sessionOne"}
	if err := db.Create(&s1).Error; err != nil {
		t.Errorf("session create should not throw err: %s", err)
	}

	sId := s1.ID

	un := Retry{
		ClientId:           "uno",
		ApplicationMessage: []byte{1, 2, 3},
		PacketIdentifier:   50,
		Qos:                1,
		Dup:                false,
		Retries:            3,
		AckStatus:          0,
		CreatedAt:          time.Now().Unix() - 30,
		SessionID:          sId,
	}
	db.Create(&un)

	du := Retry{
		ClientId:           "due",
		ApplicationMessage: []byte{4, 5, 6},
		PacketIdentifier:   51,
		Qos:                1,
		Dup:                false,
		Retries:            3,
		AckStatus:          0,
		CreatedAt:          time.Now().Unix() - 90,
		SessionID:          sId,
	}
	db.Create(&du)

	before := []Retry{}
	db.Find(&before)

	fmt.Println(before)

	if len(before) != 2 {
		t.Errorf("expected 2 retry, found: %d", len(before))
	}

	cleanRetries(db)

	after := []Retry{}
	db.Find(&after)

	if len(after) != 1 {
		t.Errorf("expected 1 retry (expiration %d secs), found: %d", conf.RETRY_EXPIRATION, len(after))
	}
}
