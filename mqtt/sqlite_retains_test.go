package mqtt

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func retainsSeed() (SqliteRetains, error) {
	Seed("../testdata/sqlite.db3")
	db, err := sql.Open("sqlite3", "../testdata/sqlite.db3")
	return SqliteRetains{db: db}, err
}

func TestSqliteAddRemoveRetains(t *testing.T) {
	retains, err := retainsSeed()
	if err != nil {
		t.Errorf("error opening test db %s", err)
	}
	r0 := Retain{
		topic:              "topic1",
		applicationMessage: []byte("hello"),
		createdAt:          time.Now(),
	}
	err = retains.addRetain(r0)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	r1 := Retain{
		topic:              "topic2",
		applicationMessage: []byte("goodbye"),
		createdAt:          time.Now(),
	}
	err = retains.addRetain(r1)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	rets := retains.findRetainsByTopic("topic1")
	if len(rets) != 1 {
		t.Errorf("expected 1 retain, got %d", len(rets))
	}
	if string(rets[0].applicationMessage) != "hello" {
		t.Errorf("expected hello , got %s", string(rets[0].applicationMessage))
	}
	r0d := Retain{
		topic:              "topic1",
		applicationMessage: []byte{},
		createdAt:          time.Now(),
	}
	err = retains.addRetain(r0d)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	rets = retains.findRetainsByTopic("topic1")
	if len(rets) != 0 {
		t.Errorf("expected 0 retain, got %d", len(rets))
	}
}

func TestSqliteRetainsWithWildcard(t *testing.T) {
	retains, err := retainsSeed()
	if err != nil {
		t.Errorf("error opening test db %s", err)
	}
	r0 := Retain{
		topic:              "topic/a/b",
		applicationMessage: []byte("hello"),
		createdAt:          time.Now(),
	}
	err = retains.addRetain(r0)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	r1 := Retain{
		topic:              "topic/a",
		applicationMessage: []byte("goodbye"),
		createdAt:          time.Now(),
	}
	err = retains.addRetain(r1)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	r2 := Retain{
		topic:              "topic",
		applicationMessage: []byte("goodby"),
		createdAt:          time.Now(),
	}
	err = retains.addRetain(r2)
	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}
	rets := retains.findRetainsByTopic("topic/#")
	if len(rets) != 2 {
		t.Errorf("expected 2 retain, got %d", len(rets))
	}
	am := string(rets[0].applicationMessage)
	if am != "hello" && am != "goodbye" {
		t.Errorf("expected hello or goodbye, got %s", am)
	}
}

func TestWithWildcard(t *testing.T) {
	topic := "topic/#"
	if !withWildCard(topic) {
		t.Errorf("%s expected true", topic)
	}
	topic = "#"
	if !withWildCard(topic) {
		t.Errorf("%s expected true", topic)
	}
	topic = "#abc"
	if withWildCard(topic) {
		t.Errorf("%s expected false", topic)
	}
	topic = "a/b/c"
	if withWildCard(topic) {
		t.Errorf("%s expected false", topic)
	}
}
