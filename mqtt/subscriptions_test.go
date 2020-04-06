package mqtt

import "testing"

func TestMultiSegmentSubs(t *testing.T) {
	var is inMemorySubscriptions = make(map[string][]string)
	is["anna/#"] = []string{"client0", "client1"}
	is["anna/barbara/#"] = []string{"client2", "client3"}
	is["anna/barbara/carlo"] = []string{"client4"}
	subs := is.multiSegmentSubs([]string{"anna", "barbara", "carlo"})
	if len(subs) != 5 {
		t.Errorf("retrieved less than 5 subscribers %d", len(subs))
	}
}

func TestFindSubscriber(t *testing.T) {
	subs := []string{"client0", "client1"}
	a := findSubscriber(subs, "client0")
	if a != 0 {
		t.Errorf("expected %d to be 0", a)
	}
	b := findSubscriber(subs, "client1")
	if b != 1 {
		t.Errorf("expected %d to be 1", b)
	}
	c := findSubscriber(subs, "client2")
	if c != -1 {
		t.Errorf("expected %d to be -1", c)
	}
	emptys := []string{}
	d := findSubscriber(emptys, "client0")
	if d != -1 {
		t.Errorf("expected %d to be -1", a)
	}
}

func TestAddSub(t *testing.T) {
	var is inMemorySubscriptions = make(map[string][]string)
	is["anna/#"] = []string{"client0", "client1"}
	err := is.addSub("anna/#", "client0")
	if err != nil {
		t.Error("did not expect any error")
	}
	if len(is["anna/#"]) != 2 {
		t.Errorf("expected 2 subscribers, %d found", len(is["anna/#"]))
	}
	err0 := is.addSub("anna/#", "client2")
	if err0 != nil {
		t.Error("did not expect any error")
	}
	i := findSubscriber(is["anna/#"], "client2")
	if i != 2 {
		t.Errorf("expected client2 in subscribers position 2, %d found", i)
	}
	is["barbara"] = []string{}
	err1 := is.addSub("barbara", "client3")
	if err1 != nil {
		t.Error("did not expect any error")
	}
	j := findSubscriber(is["barbara"], "client3")
	if j != 0 {
		t.Errorf("expected client3 in subscribers position 0, %d found", j)
	}
}

func TestRemSub(t *testing.T) {
	var is0 inMemorySubscriptions = make(map[string][]string)
	err0 := is0.remSub("anna/#", "client0")
	if err0 != nil {
		t.Error("did not expect any error")
	}
	var is inMemorySubscriptions = make(map[string][]string)
	is["anna/#"] = []string{"client0", "client1"}
	is["barbara/#"] = []string{"client0", "client1"}
	err := is.remSub("anna/#", "client0")
	if err != nil {
		t.Error("did not expect any error")
	}
	if len(is["anna/#"]) != 1 {
		t.Errorf("expected 1 subscribers, %d found", len(is["anna/#"]))
	}

}
