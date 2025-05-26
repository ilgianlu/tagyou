package engine

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"

	_ "github.com/mattn/go-sqlite3"
)

func TestClientGoodConnection(t *testing.T) {
	os.Setenv("DEBUG", "1")

	p := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	db, err := p.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	pwd, _ := password.EncodePassword([]byte("password"))

	db.CreateClient(
		context.Background(),
		dbaccess.CreateClientParams{
			ClientID:     sql.NullString{String: "client-1", Valid: true},
			Username:     sql.NullString{String: "user1", Valid: true},
			SubscribeAcl: sql.NullString{String: "[]", Valid: true},
			PublishAcl:   sql.NullString{String: "[]", Valid: true},
			Password:     pwd,
		},
	)

	session := model.RunningSession{ClientId: "client-1", Username: "user1", Password: "password"}

	res := doAuth(&session)

	if res != true {
		t.Errorf("expecting true (success), received false")
	}

}

func TestClientBadConnectionWrongPasswordConnection(t *testing.T) {
	os.Setenv("DEBUG", "1")

	p := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	db, err := p.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	pwd, _ := password.EncodePassword([]byte("password"))
	db.CreateClient(
		context.Background(),
		dbaccess.CreateClientParams{
			ClientID:     sql.NullString{String: "client-1", Valid: true},
			Username:     sql.NullString{String: "user1", Valid: true},
			SubscribeAcl: sql.NullString{String: "[]", Valid: true},
			PublishAcl:   sql.NullString{String: "[]", Valid: true},
			Password:     pwd,
		},
	)

	session := model.RunningSession{ClientId: "client-1", Username: "user1", Password: "wrong"}

	res := doAuth(&session)

	if res != false {
		t.Errorf("expecting false (no access), received true")
	}
}

func TestClientBadConnectionWrongUsernameConnection(t *testing.T) {
	os.Setenv("DEBUG", "1")

	p := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	db, err := p.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	pwd, _ := password.EncodePassword([]byte("password"))
	db.CreateClient(
		context.Background(),
		dbaccess.CreateClientParams{
			ClientID:     sql.NullString{String: "client-1", Valid: true},
			Username:     sql.NullString{String: "user1", Valid: true},
			SubscribeAcl: sql.NullString{String: "[]", Valid: true},
			PublishAcl:   sql.NullString{String: "[]", Valid: true},
			Password:     pwd,
		},
	)

	session := model.RunningSession{ClientId: "client-1", Username: "wrong", Password: "password"}

	res := doAuth(&session)

	if res != false {
		t.Errorf("expecting false (no access), received true")
	}
}
