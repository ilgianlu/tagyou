package badgerrepository

import (
	"encoding/gob"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/model"
)

func TestFindRetains(t *testing.T) {
	gob.Register(model.Retain{})
	dbe, _ := badger.Open(badger.DefaultOptions("ret.db"))
	// defer dbe.Close()
	defer os.RemoveAll("ret.db")

	retainRepository := RetainBadgerRepository{Db: dbe}

	un := model.Retain{Topic: "topic1", ApplicationMessage: []byte("message1")}
	du := model.Retain{Topic: "topic2", ApplicationMessage: []byte("message2")}
	retainRepository.Create(un)
	retainRepository.Create(du)

	retains := retainRepository.FindRetains("topic1")
	if len(retains) != 1 {
		t.Errorf("expected %d subscriptions, found %d", 1, len(retains))
	}
	if retains[0].Topic != "topic1" {
		t.Errorf("expected %s topic retain, found %s", "topic1", retains[0].Topic)
	}
}
