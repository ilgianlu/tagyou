package event

import (
	"testing"

	"github.com/ilgianlu/tagyou/model"
)

// func TestCheckAcls(t *testing.T) {
// 	type test struct {
// 		data   []string
// 		result bool
// 	}

// 	tests := []test{
// 		{[]string{"/a/b", "{}"}, false},
// 		{[]string{"/a/b", "[{\"pattern\": \"/a/#\"}]"}, true},
// 	}

// 	for _, value := range tests {
// 		if CheckAcl(value.data[0], value.data[1]) != value.result {
// 			t.Errorf("expected %s to acl check %s, %t", value.data[0], value.data[1], value.result)
// 		}
// 	}
// }

func TestPickDest(t *testing.T) {
	sub1 := model.Subscription{
		ClientId: "pippo",
	}
	sub2 := model.Subscription{
		ClientId: "pluto",
	}
	subGroup := []model.Subscription{sub1, sub2}
	dest := pickDest(subGroup)

	if dest.ClientId != "pippo" {
		t.Errorf("expecting %s received %s", "pippo", dest.ClientId)
	}
}

func TestGroupSubscribers(t *testing.T) {
	sub1 := model.Subscription{
		ClientId:  "pippo",
		ShareName: "share1",
	}
	sub2 := model.Subscription{
		ClientId:  "pluto",
		ShareName: "share2",
	}
	sub3 := model.Subscription{
		ClientId:  "minnie",
		ShareName: "share1",
	}
	sub4 := model.Subscription{
		ClientId:  "topolino",
		ShareName: "share2",
	}
	ungrouped := []model.Subscription{sub1, sub2, sub3, sub4}
	groups := groupSubscribers(ungrouped)
	if group, ok := groups["share1"]; ok {
		if group[0].ClientId != "pippo" {
			t.Errorf("expecting %s received %s", "pippo", group[0].ClientId)
		}
	} else {
		t.Errorf("no group share1 found!")
	}
	if group, ok := groups["share2"]; ok {
		if group[0].ClientId != "pluto" {
			t.Errorf("expecting %s received %s", "pluto", group[0].ClientId)
		}
	} else {
		t.Errorf("no group share2 found!")
	}
}
