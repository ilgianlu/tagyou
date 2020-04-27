package mqtt

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func authSeed() (SqliteAuths, error) {
	Seed("../testdata/sqlite.db3")
	db, err := sql.Open("sqlite3", "../testdata/sqlite.db3")
	return SqliteAuths{db: db}, err
}

func TestPassword(t *testing.T) {
	var a Auth
	err := a.setPassword("pippo123")
	if err != nil {
		t.Error("Did not expect any error")
	}
	err = a.checkPassword("pippo123")
	if err != nil {
		t.Error("Password not matching")
	}
}

func TestAuth(t *testing.T) {
	auths, err := authSeed()
	if err != nil {
		t.Error("Did not expect any error")
	}
	a := Auth{
		clientId:     "client1",
		username:     "user",
		subscribeAcl: "a/#",
		publishAcl:   "b/#",
		createdAt:    time.Now(),
	}
	z := Auth{
		clientId:     "client2",
		username:     "user0",
		subscribeAcl: "a0/#",
		publishAcl:   "b0/#",
		createdAt:    time.Now(),
	}
	err = auths.createAuth(z)
	if err != nil {
		t.Error("Did not expect any error")
	}
	err = a.setPassword("pippo123")
	if err != nil {
		t.Error("Did not expect any error")
	}
	err = auths.createAuth(a)
	if err != nil {
		t.Error("Did not expect any error")
	}
	b, ok := auths.findAuth(a.clientId)
	if !ok {
		t.Error("expected one auth!")
	}
	if b.username != "user" {
		t.Errorf("expected user, found %s", b.username)
	}
}
