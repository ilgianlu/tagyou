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

func TestExplode(t *testing.T) {
	type test struct {
		data   string
		result []string
	}

	tests := []test{
		{"a", []string{"#", "a", "+"}},
		// {"a/b", []string{"#", "a/b", "a/#", "+/#", "+/b", "a/+"}},
		// {"a/b/c", []string{"#", "a/b/c", "a/#", "+/#", "+/b/c", "a/b/#", "a/+/c", "a/+/#", "a/b/+"}},
		// {"/a", []string{"#", "/a", "/#", "+/#", "+/a", "/+"}},
		// {"/b/c", []string{"#", "/b/c", "/#", "+/#", "+/b/c", "/b/#", "/+/c", "/+/#", "/b/+"}},
	}

	for _, value := range tests {
		res := Explode(value.data)
		// log.Println(value.data, "==>", res, "==", value.result)
		if !arrEq(res, value.result) {
			t.Error("expected", value.data, "to explode into", value.result, "received", res)
		}
	}
}

func TestExplodeMultiLevel(t *testing.T) {
	type test struct {
		data   []string
		result []string
	}

	tests := []test{
		{[]string{"a"}, []string{"#"}},
		{[]string{"", "a"}, []string{"#", "/#"}},
		{[]string{"a", "b"}, []string{"#", "a/#"}},
		{[]string{"a", "b", "c"}, []string{"#", "a/#", "a/b/#"}},
	}

	for _, value := range tests {
		res := explodeMultiLevel(value.data)
		if !arrEq(res, value.result) {
			t.Error("expected", value.data, "to explode into", value.result, "received", res)
		}
	}
}

func TestExplodeFull(t *testing.T) {
	type test struct {
		data   []string
		result []string
	}

	tests := []test{
		{[]string{"a"}, []string{"#", "a", "+"}},
		{[]string{"a", "b"}, []string{"#", "a/#", "+/#", "a/b", "+/b", "a/+", "+/+"}},
		{[]string{"", "a"}, []string{"#", "/#", "+/#", "/a", "+/a", "/+", "+/+"}},
		{
			[]string{"a", "b", "c"},
			[]string{
				"#",
				"a/#", "+/#",
				"a/b/#", "+/b/#", "a/+/#", "+/+/#",
				"a/b/c", "+/b/c", "a/+/c", "+/+/c", "a/b/+", "+/b/+", "a/+/+", "+/+/+"},
		},
	}

	for _, value := range tests {
		res := explodeFull(value.data)
		// log.Println(value.data, "==>", res, "==", value.result)
		if !arrEq(res, value.result) {
			t.Error("expected", value.data, "to explode into", value.result, "received", res)
		}
	}
}

func TestExplodeSingleLevel(t *testing.T) {
	type test struct {
		data   []string
		result [][]string
	}

	tests := []test{
		{[]string{"a"}, [][]string{{"a"}, {"+"}}},
		{[]string{"", "a"}, [][]string{{"", "a"}, {"+", "a"}, {"", "+"}, {"+", "+"}}},
		{[]string{"a", "b"}, [][]string{{"a", "b"}, {"+", "b"}, {"a", "+"}, {"+", "+"}}},
		{[]string{"a", "b", "c"}, [][]string{{"a", "b", "c"}, {"+", "b", "c"}, {"a", "+", "c"}, {"+", "+", "c"}, {"a", "b", "+"}, {"+", "b", "+"}, {"a", "+", "+"}, {"+", "+", "+"}}},
	}

	for _, value := range tests {
		res := explodeSingleLevel(value.data)
		// log.Println(value.data, "==>", res, "==", value.result)
		if !arrArrEq(res, value.result) {
			t.Error("expected", value.data, "to explode into", value.result, "received", res)
		}
	}
}

func arrArrEq(a [][]string, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !arrEq(a[i], b[i]) {
			return false
		}
	}
	return true
}

func arrEq(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
