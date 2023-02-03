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
	if err := subscriptionRepository.Create(un); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}

	unBis := model.Subscription{ClientId: "uno", Topic: "uno"}
	if err := subscriptionRepository.Create(unBis); err == nil {
		t.Error("subscription (duplicate client id and topic) create should throw err!")
	}

	du := model.Subscription{ClientId: "due", Topic: "uno"}
	if err := subscriptionRepository.Create(du); err != nil {
		t.Errorf("subscription create should not throw err: %s", err)
	}
}

func TestFindToUnsubscribe(t *testing.T) {
	gob.Register(model.Subscription{})
	dbe, _ := badger.Open(badger.DefaultOptions("sub.db"))
	// defer dbe.Close()
	defer os.RemoveAll("sub.db")

	subscriptionRepository := SubscriptionBadgerRepository{Db: dbe}

	un := model.Subscription{ClientId: "client1", Topic: "topic1"}
	du := model.Subscription{ClientId: "client2", Topic: "topic2"}
	subscriptionRepository.Create(un)
	subscriptionRepository.Create(du)

	sub, err := subscriptionRepository.FindToUnsubscribe("", "topic1", "client1")
	if err != nil {
		t.Errorf("FindToSubscribe should not throw err: %s", err)
	}
	if sub.ClientId != "client1" {
		t.Errorf("expected client id %s , found %s", "client1", sub.ClientId)
	}
}

func TestFindSubscriptions(t *testing.T) {
	gob.Register(model.Subscription{})
	dbe, _ := badger.Open(badger.DefaultOptions("sub.db"))
	// defer dbe.Close()
	defer os.RemoveAll("sub.db")

	subscriptionRepository := SubscriptionBadgerRepository{Db: dbe}

	ze := model.Subscription{ClientId: "clientOut", Topic: "outoftime"}
	un := model.Subscription{ClientId: "client1", Topic: "topic1"}
	du := model.Subscription{ClientId: "client1", Topic: "topic2"}
	tr := model.Subscription{ClientId: "client2", Topic: "topic1"}
	qu := model.Subscription{ClientId: "client3", Topic: "topic3"}
	ci := model.Subscription{ClientId: "client4", Topic: "topic/#"}
	se := model.Subscription{ClientId: "client5", Topic: "topic/+/message"}
	subscriptionRepository.Create(ze)
	subscriptionRepository.Create(un)
	subscriptionRepository.Create(du)
	subscriptionRepository.Create(tr)
	subscriptionRepository.Create(qu)
	subscriptionRepository.Create(ci)
	subscriptionRepository.Create(se)

	subs := subscriptionRepository.FindSubscriptions([]string{"topic1"}, false)
	if len(subs) != 2 {
		t.Errorf("expected %d subscriptions, found %d", 2, len(subs))
	}
	if subs[0].ClientId != "client1" {
		t.Errorf("expected client id %s , found %s", "client1", subs[0].ClientId)
	}
	if subs[1].ClientId != "client2" {
		t.Errorf("expected client id %s , found %s", "client2", subs[1].ClientId)
	}
	newSubs := subscriptionRepository.FindSubscriptions([]string{"topic/#"}, false)
	if len(newSubs) != 1 {
		t.Errorf("expected %d subscriptions, found %d", 1, len(newSubs))
	}
	if newSubs[0].ClientId != "client4" {
		t.Errorf("expected client id %s , found %s", "client4", newSubs[0].ClientId)
	}
}
