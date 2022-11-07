package commands

import (
	"strings"
	"testing"
)

// ErrorContains checks if the error message in out contains the text in
// want.
//
// This is safe when out is nil. Use an empty string for want if you want to
// test that err is nil.
func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

// Prefix tests
func TestPrefix(t *testing.T) {
	p := Path{}
	p.Type = "lala"
	log_prefix := p.Prefix()
	if log_prefix != "TestPrefix[lala]" {
		t.Errorf("got %s expected %s", log_prefix, "lala")
	}
}

// Path tests
type pathTest struct {
	inputstr string
	path     []string
	result   bool
}

var pathTests = []pathTest{
	pathTest{"a", []string{"a", "b", "c"}, true},
	pathTest{"b", []string{"a", "b", "c"}, true},
	pathTest{"c", []string{"a", "b", "c"}, true},
	pathTest{"d", []string{"a", "b", "c"}, false},
}

func TestHavePath(t *testing.T) {
	for _, test := range pathTests {
		p := Path{}
		p.Directories = test.path
		result := p.HavePath(test.inputstr)
		if result != test.result {
			t.Errorf("expected %s in %s (%v)", test.inputstr, strings.Join(test.path, ", "), result)
		}
	}
}

//
// AppendPath tests
//

type appendPathTest struct {
	inputstr string
	path     []string
	errstr   string
}

var appendPathTests = []appendPathTest{
	appendPathTest{"/etc", []string{"/tmp", "/var", "/etc"}, ""},
	appendPathTest{"/NOEXIST", []string{"/tmp", "/var", "/etc"}, "stat /NOEXIST: no such file or directory"},
	appendPathTest{"/bin", []string{"/tmp", "/var", "/etc"}, ""},
	appendPathTest{"/usr", []string{"/tmp", "/var", "/etc"}, ""},
}

func TestAppendPath(t *testing.T) {
	for _, test := range appendPathTests {
		p := Path{}
		p.Directories = test.path

		result := p.AppendPath(test.inputstr)
		if !ErrorContains(result, test.errstr) {
			t.Errorf("expected %s in %s (%q)", test.inputstr, strings.Join(test.path, ", "), result)
		}
	}
}
