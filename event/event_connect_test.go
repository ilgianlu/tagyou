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

func TestClientGoodConnection(t *testing.T) {
	os.Setenv("DEBUG", "1")
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	persistence := persistence.SqlPersistence{}
	persistence.InnerInit(db)

	db.Exec("DELETE FROM auths")

	auth1 := model.Auth{ClientId: "client-1", Username: "user1", SubscribeAcl: "[]", PublishAcl: "[]"}
	auth1.SetPassword("password")
	rAuth1 := sqlrepository.Auth{ClientId: "client-1", Username: "user1", SubscribeAcl: "[]", PublishAcl: "[]", Password: auth1.Password}
	db.Create(&rAuth1)

	session := model.RunningSession{ClientId: "client-1", Username: "user1", Password: "password"}

	res := doAuth(&session)

	if res != true {
		t.Errorf("expecting true (success), received false")
	}

}

func TestClientBadConnectionWrongPasswordConnection(t *testing.T) {
	os.Setenv("DEBUG", "1")
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	persistence := persistence.SqlPersistence{}
	persistence.InnerInit(db)

	db.Exec("DELETE FROM auths")

	auth1 := model.Auth{ClientId: "client-1", Username: "user1", SubscribeAcl: "[]", PublishAcl: "[]"}
	auth1.SetPassword("password")
	rAuth1 := sqlrepository.Auth{ClientId: "client-1", Username: "user1", SubscribeAcl: "[]", PublishAcl: "[]", Password: auth1.Password}
	db.Create(&rAuth1)

	session := model.RunningSession{ClientId: "client-1", Username: "user1", Password: "wrong"}

	res := doAuth(&session)

	if res != false {
		t.Errorf("expecting false (no access), received true")
	}
}

func TestClientBadConnectionWrongUsernameConnection(t *testing.T) {
	os.Setenv("DEBUG", "1")
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	persistence := persistence.SqlPersistence{}
	persistence.InnerInit(db)

	db.Exec("DELETE FROM auths")

	auth1 := model.Auth{ClientId: "client-1", Username: "user1", SubscribeAcl: "[]", PublishAcl: "[]"}
	auth1.SetPassword("password")
	rAuth1 := sqlrepository.Auth{ClientId: "client-1", Username: "user1", SubscribeAcl: "[]", PublishAcl: "[]", Password: auth1.Password}
	db.Create(&rAuth1)

	session := model.RunningSession{ClientId: "client-1", Username: "wrong", Password: "password"}

	res := doAuth(&session)

	if res != false {
		t.Errorf("expecting false (no access), received true")
	}
}
