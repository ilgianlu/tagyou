package event

import "testing"

func TestCheckAcls(t *testing.T) {
	type test struct {
		data   []string
		result bool
	}

	tests := []test{
		{[]string{"/a/b", "{}"}, false},
		{[]string{"/a/b", "[{\"pattern\": \"/a/#\"}]"}, true},
	}

	for _, value := range tests {
		if CheckAcl(value.data[0], value.data[1]) != value.result {
			t.Errorf("expected %s to acl check %s, %t", value.data[0], value.data[1], value.result)
		}
	}
}
