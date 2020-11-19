package event

import (
	"testing"
)

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

func TestReadAcls(t *testing.T) {
	var acls []Acl
	var err error
	acls, err = readAcls("")
	if err == nil {
		t.Errorf("expected to err! %s", acls)
	}
	acls, err = readAcls("{}")
	if err == nil {
		t.Errorf("expected to err! %s", acls)
	}
	acls, err = readAcls("[]")
	if err != nil || len(acls) != 0 {
		t.Errorf("expected no error and empty acl array! %s", acls)
	}
	acls, err = readAcls("[{\"pattern\": \"/a/#\"}]")
	if err != nil || len(acls) != 1 {
		t.Errorf("expected no error and one element acl array! %s", acls)
	}
}
