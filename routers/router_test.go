package routers

import "testing"

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
		res := explodeFull(value.data)
		if !arrEq(res, value.result) {
			t.Error("expected", value.data, "to explode into", value.result, "received", res)
		}
	}
}

func TestExplodeFull(t *testing.T) {
	type test struct {
		data   string
		result []string
	}

	tests := []test{
		{"a", []string{"#", "a", "+"}},
		{"a/b", []string{"#", "a/#", "+/#", "a/b", "+/b", "a/+", "+/+"}},
		{"/a", []string{"#", "/#", "+/#", "/a", "+/a", "/+", "+/+"}},
		{
			"a/b/c",
			[]string{
				"#",
				"a/#", "+/#",
				"a/b/#", "+/b/#", "a/+/#", "+/+/#",
				"a/b/c", "+/b/c", "a/+/c", "+/+/c", "a/b/+", "+/b/+", "a/+/+", "+/+/+"},
		},
	}

	for _, value := range tests {
		res := explodeFull(value.data)
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
