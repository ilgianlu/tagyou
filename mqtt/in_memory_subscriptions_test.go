package mqtt

import (
	"testing"
)

func TestInMemoryMultiSegmentSubs(t *testing.T) {
	var is inMemorySubscriptions
	is.topicSubscriptions = make(map[string][]string)
	is.topicSubscriptions["anna/#"] = []string{"client0", "client1"}
	is.topicSubscriptions["anna/barbara/#"] = []string{"client2", "client3"}
	is.topicSubscriptions["anna/barbara/carlo"] = []string{"client4"}
	subs := is.multiSegmentSubs([]string{"anna", "barbara", "carlo"})
	if len(subs) != 5 {
		t.Errorf("retrieved less than 5 subscribers %d", len(subs))
	}
}

func TestInMemoryFindSubscriber(t *testing.T) {
	subs := []string{"client0", "client1"}
	a := findIn(subs, "client0")
	if a != 0 {
		t.Errorf("expected %d to be 0", a)
	}
	b := findIn(subs, "client1")
	if b != 1 {
		t.Errorf("expected %d to be 1", b)
	}
	c := findIn(subs, "client2")
	if c != -1 {
		t.Errorf("expected %d to be -1", c)
	}
	emptys := []string{}
	d := findIn(emptys, "client0")
	if d != -1 {
		t.Errorf("expected %d to be -1", a)
	}
}

func TestInMemoryAddSubscription(t *testing.T) {
	var is inMemorySubscriptions
	is.clientSubscriptions = make(map[string][]string)
	is.topicSubscriptions = make(map[string][]string)
	is.topicSubscriptions["anna/#"] = []string{"client0", "client1"}
	is.clientSubscriptions["client0"] = []string{"anna/#"}
	is.clientSubscriptions["client1"] = []string{"anna/#"}
	err := is.addSubscription("anna/#", "client0")
	if err != nil {
		t.Error("did not expect any error")
	}
	if len(is.topicSubscriptions["anna/#"]) != 2 {
		t.Errorf("expected 2 subscribers, %d found", len(is.topicSubscriptions["anna/#"]))
	}
	err0 := is.addSubscription("anna/#", "client2")
	if err0 != nil {
		t.Error("did not expect any error")
	}
	i := findIn(is.topicSubscriptions["anna/#"], "client2")
	if i != 2 {
		t.Errorf("expected client2 in subscribers position 2, %d found", i)
	}
	is.topicSubscriptions["barbara"] = []string{}
	err1 := is.addSubscription("barbara", "client3")
	if err1 != nil {
		t.Error("did not expect any error")
	}
	j := findIn(is.topicSubscriptions["barbara"], "client3")
	if j != 0 {
		t.Errorf("expected client3 in subscribers position 0, %d found", j)
	}
}

func TestInMemoryRemSubscription(t *testing.T) {
	var is inMemorySubscriptions
	is.clientSubscriptions = make(map[string][]string)
	is.topicSubscriptions = make(map[string][]string)
	is.topicSubscriptions["anna/#"] = []string{"client0", "client1"}
	is.topicSubscriptions["barbara/#"] = []string{"client0", "client1"}
	is.clientSubscriptions["client0"] = []string{"anna/#", "barbara/#"}
	is.clientSubscriptions["client1"] = []string{"anna/#", "barbara/#"}

	err1 := is.remSubscription("anna/#", "client0")
	if err1 != nil {
		t.Errorf("not expect any error: %s\n", err1)
	}
	pos1 := findIn(is.topicSubscriptions["anna/#"], "client0")
	if pos1 != -1 {
		t.Errorf("expecting -1, found %d\n", pos1)
	}
	if len(is.topicSubscriptions["anna/#"]) != 1 {
		t.Errorf("expected 1 subscribers, %d found", len(is.topicSubscriptions["anna/#"]))
	}
	err2 := is.remSubscription("anna/#", "client2")
	if err2 == nil {
		t.Errorf("expected an error, clientid does not exist\n")
	}
	pos2 := findIn(is.topicSubscriptions["anna/#"], "client2")
	if pos2 != -1 {
		t.Errorf("expecting -1, found %d\n", pos2)
	}
	err3 := is.remSubscription("carlo", "client0")
	if err3 == nil {
		t.Errorf("expected an error, topic does not exist\n")
	}
	pos3 := findIn(is.topicSubscriptions["carlo"], "client0")
	if pos3 != -1 {
		t.Errorf("expecting -1, found %d\n", pos3)
	}
	err4 := is.remSubscription("barbara/#", "client1")
	if err4 != nil {
		t.Errorf("not expect any error: %s\n", err4)
	}
	pos4 := findIn(is.topicSubscriptions["barbara/#"], "client1")
	if pos4 != -1 {
		t.Errorf("expecting -1, found %d\n", pos4)
	}
	pos5 := findIn(is.topicSubscriptions["barbara/#"], "client0")
	if pos5 != 0 {
		t.Errorf("expecting 0, found %d\n", pos5)
	}
}
