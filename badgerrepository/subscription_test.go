package badgerrepository

import (
	"encoding/gob"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/model"
)

func TestCreate(t *testing.T) {
	gob.Register(model.Subscription{})
	dbe, _ := badger.Open(badger.DefaultOptions("sub.db"))
	// defer dbe.Close()
	defer os.RemoveAll("sub.db")

	subscriptionRepository := SubscriptionBadgerRepository{Db: dbe}

	un := model.Subscription{ClientId: "uno", Topic: "uno"}
	if err := subscriptionRepository.CreateOne(un); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	unBis := model.Subscription{ClientId: "uno", Topic: "uno"}
	if err := subscriptionRepository.CreateOne(unBis); err == nil {
		t.Error("subscription (duplicate client id and topic) create should throw err!")
	}

	du := model.Subscription{ClientId: "due", Topic: "uno"}
	if err := subscriptionRepository.CreateOne(du); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}
}
