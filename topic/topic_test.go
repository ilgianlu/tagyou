package topic

import (
	"log"
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
		{"a/b", []string{"#", "a/b", "a/#", "+/#", "+/b", "a/+"}},
		{"a/b/c", []string{"#", "a/b/c", "a/#", "+/#", "+/b/c", "a/b/#", "a/+/c", "a/+/#", "a/b/+"}},
	}

	for _, value := range tests {
		res := Explode(value.data)
		log.Println(value.data, "==>", res, "==", value.result)
		if !arrEq(res, value.result) {
			t.Error("expected", value.data, "to explode into", value.result, "received", res)
		}
	}
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
