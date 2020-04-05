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
