package topic

import (
	"testing"
)

func TestMatch(t *testing.T) {
	type test struct {
		data   []string
		result bool
	}

	tests := []test{
		{[]string{"/a/b/c", "/a/b/c"}, true},
		{[]string{"/a/b/c", "/a/b/#"}, true},
		{[]string{"/a/b/c", "/a/#"}, true},
		{[]string{"/a/b/c", "/#"}, true},
		{[]string{"", ""}, true},
		{[]string{"/a", "/b"}, false},
		{[]string{"/a/b/c", "/a/+/c"}, true},
		{[]string{"/a/b/c", "/+/b/c"}, true},
		{[]string{"/a/b/d", "/a/+/c"}, false},
		{[]string{"/a/b/c", "/a/+/d"}, false},
		{[]string{"/a/b/c/d", "/a/+/d"}, false},
		{[]string{"/a/b/c/d", "/a/+/+/d"}, true},
		{[]string{"/a/b/c/d", "/a/b/c"}, false},
	}

	for _, value := range tests {
		if Match(value.data[0], value.data[1]) != value.result {
			t.Errorf("expected %s to match %s, %t", value.data[0], value.data[1], value.result)
		}
	}
}

func TestMatcherSubset(t *testing.T) {
	type test struct {
		data   []string
		result bool
	}

	tests := []test{
		{[]string{"/a/b/c", "/a/b/c"}, true},
		{[]string{"/a/b/c", "/a/b/#"}, true},
		{[]string{"/a/b/c", "/a/#"}, true},
		{[]string{"/a/b/c", "/#"}, true},
		{[]string{"/a/b/#", "/a/#"}, true},
		{[]string{"/a/#", "/a/b/#"}, false},
		{[]string{"/a/b/#", "/a/b/c"}, false},
		{[]string{"/a/b/c/d", "/a/b/c"}, false},
		{[]string{"", ""}, false},
		{[]string{"/a", "/b"}, false},
		{[]string{"/a/b/c", "/a/+/c"}, true},
		{[]string{"/a/+/c", "/a/b/c"}, false},
	}

	for _, value := range tests {
		if MatcherSubset(value.data[0], value.data[1]) != value.result {
			t.Errorf("expected %s to match %s, %t", value.data[0], value.data[1], value.result)
		}
	}
}

func TestSharedSubscription(t *testing.T) {
	type test struct {
		data   string
		result bool
	}

	tests := []test{
		{"/a/b/c", false},
		{"$share//a/b/c", false},
		{"$share/a/b/c", true},
		{"$share/pippo/a/b/c", true},
	}

	for _, value := range tests {
		if SharedSubscription(value.data) != value.result {
			t.Errorf("expected %s to match %t", value.data, value.result)
		}
	}
}

func TestSharedSubscriptionTopicParse(t *testing.T) {
	type ResultData struct {
		ShareName  string
		ShareTopic string
	}

	type test struct {
		data   string
		result ResultData
	}

	tests := []test{
		{data: "$share/shareName/a/b/c", result: ResultData{ShareName: "shareName", ShareTopic: "a/b/c"}},
		{data: "$share/a/b/c", result: ResultData{ShareName: "a", ShareTopic: "b/c"}},
	}

	for _, value := range tests {
		shareName, sharedTopic := SharedSubscriptionTopicParse(value.data)
		if shareName != value.result.ShareName || sharedTopic != value.result.ShareTopic {
			t.Errorf("shareName expected %s to match %s && sharedTopic expected %s to match %s", shareName, value.result.ShareName, sharedTopic, value.result.ShareTopic)
		}
	}
}
