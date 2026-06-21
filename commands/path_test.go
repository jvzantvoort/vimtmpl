package commands

import (
	"os"
	"path/filepath"
	"runtime"
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

func TestAppendPathEmptyInput(t *testing.T) {
	p := Path{Directories: []string{"/tmp"}}
	if err := p.AppendPath(""); err != nil {
		t.Errorf("expected nil error for empty input, got %v", err)
	}
	if len(p.Directories) != 1 {
		t.Errorf("expected no directories to be added, got %v", p.Directories)
	}
}

func TestAppendPathDoesNotDuplicate(t *testing.T) {
	p := Path{Directories: []string{"/tmp"}}
	if err := p.AppendPath("/tmp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Directories) != 1 {
		t.Errorf("expected directory to not be duplicated, got %v", p.Directories)
	}
}

//
// PrependPath tests
//

func TestPrependPathEmptyInput(t *testing.T) {
	p := Path{Directories: []string{"/tmp"}}
	if err := p.PrependPath(""); err != nil {
		t.Errorf("expected nil error for empty input, got %v", err)
	}
	if len(p.Directories) != 1 {
		t.Errorf("expected no directories to be added, got %v", p.Directories)
	}
}

func TestPrependPathAddsToFront(t *testing.T) {
	p := Path{Directories: []string{"/tmp"}}
	if err := p.PrependPath("/var"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Directories) != 2 || p.Directories[0] != "/var" {
		t.Errorf("expected /var to be prepended, got %v", p.Directories)
	}
}

func TestPrependPathDoesNotDuplicate(t *testing.T) {
	p := Path{Directories: []string{"/tmp"}}
	if err := p.PrependPath("/tmp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Directories) != 1 {
		t.Errorf("expected directory to not be duplicated, got %v", p.Directories)
	}
}

//
// Import tests
//

func TestImport(t *testing.T) {
	p := Path{Type: "PATH"}
	p.Import("/tmp:/var:/etc")
	want := []string{"/tmp", "/var", "/etc"}
	if len(p.Directories) != len(want) {
		t.Fatalf("got %v want %v", p.Directories, want)
	}
	for i, dir := range want {
		if p.Directories[i] != dir {
			t.Errorf("got %v want %v", p.Directories, want)
			break
		}
	}
}

func TestImportSkipsMissingDirs(t *testing.T) {
	p := Path{Type: "PATH"}
	p.Import("/tmp:/this-dir-should-not-exist-xyz")
	if len(p.Directories) != 1 || p.Directories[0] != "/tmp" {
		t.Errorf("got %v want only /tmp", p.Directories)
	}
}

//
// IsEmpty tests
//

func TestIsEmpty(t *testing.T) {
	p := Path{}
	if !p.IsEmpty() {
		t.Errorf("expected new Path to be empty")
	}

	p.Directories = []string{"/tmp"}
	if p.IsEmpty() {
		t.Errorf("expected Path with directories to not be empty")
	}
}

//
// ReturnExport tests
//

func TestReturnExport(t *testing.T) {
	p := Path{Type: "PATH", Directories: []string{"/tmp", "/var"}}
	got := p.ReturnExport()
	want := `export PATH="/tmp:/var"`
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

//
// targetExists tests
//

func TestTargetExists(t *testing.T) {
	p := Path{}
	dir := t.TempDir()

	if !p.targetExists(dir) {
		t.Errorf("expected existing dir to report true")
	}

	if p.targetExists(filepath.Join(dir, "doesnotexist")) {
		t.Errorf("expected missing path to report false")
	}
}

//
// Lookup tests
//

func TestLookup(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "mybinary")
	if err := os.WriteFile(target, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create test binary: %v", err)
	}

	p := Path{Directories: []string{dir}}

	got, err := p.Lookup("mybinary")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != target {
		t.Errorf("got %q want %q", got, target)
	}

	_, err = p.Lookup("doesnotexist")
	if err == nil {
		t.Errorf("expected error for missing target")
	}
}

//
// LookupMulti tests
//

func TestLookupMultiNotFoundReturnsNilError(t *testing.T) {
	// NOTE: LookupMulti's current implementation returns a nil error and an
	// empty string as soon as a target is *not* found in the first
	// directory checked, rather than continuing to search. This test
	// documents that existing behaviour.
	p := Path{Directories: []string{t.TempDir()}}

	got, err := p.LookupMulti("doesnotexist")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != "" {
		t.Errorf("got %q want empty string", got)
	}
}

func TestLookupMultiNoTargets(t *testing.T) {
	p := Path{Directories: []string{t.TempDir()}}

	_, err := p.LookupMulti()
	if err == nil {
		t.Errorf("expected error when no targets are given")
	}
}

//
// MapGetPlatform tests
//

func TestMapGetPlatform(t *testing.T) {
	p := Path{}

	m := map[string]string{runtime.GOOS: "platform-specific", "default": "fallback"}
	got, err := p.MapGetPlatform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "platform-specific" {
		t.Errorf("got %q want %q", got, "platform-specific")
	}
}

func TestMapGetPlatformFallsBackToDefault(t *testing.T) {
	p := Path{}

	m := map[string]string{"some-other-os": "other", "default": "fallback"}
	got, err := p.MapGetPlatform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "fallback" {
		t.Errorf("got %q want %q", got, "fallback")
	}
}

func TestMapGetPlatformNoMatch(t *testing.T) {
	p := Path{}

	m := map[string]string{"some-other-os": "other"}
	_, err := p.MapGetPlatform(m)
	if err == nil {
		t.Errorf("expected error when no matching key exists")
	}
}

//
// LookupPlatform tests
//

func TestLookupPlatform(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "mybinary")
	if err := os.WriteFile(target, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create test binary: %v", err)
	}

	p := Path{Directories: []string{dir}}

	m := map[string]string{runtime.GOOS: "mybinary"}
	got, err := p.LookupPlatform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != target {
		t.Errorf("got %q want %q", got, target)
	}
}

func TestLookupPlatformMapMiss(t *testing.T) {
	p := Path{Directories: []string{t.TempDir()}}

	m := map[string]string{"some-other-os": "mybinary"}
	_, err := p.LookupPlatform(m)
	if err == nil {
		t.Errorf("expected error when platform map has no match")
	}
}

func TestLookupPlatformLookupMiss(t *testing.T) {
	p := Path{Directories: []string{t.TempDir()}}

	m := map[string]string{runtime.GOOS: "doesnotexist"}
	_, err := p.LookupPlatform(m)
	if err == nil {
		t.Errorf("expected error when target is not found in directories")
	}
}

//
// NewPath tests
//

func TestNewPath(t *testing.T) {
	t.Setenv("VIMTMPL_TEST_PATH", "/tmp:/var")

	p := NewPath("VIMTMPL_TEST_PATH")
	if p.Type != "VIMTMPL_TEST_PATH" {
		t.Errorf("got Type %q want %q", p.Type, "VIMTMPL_TEST_PATH")
	}
	if len(p.Directories) != 2 {
		t.Errorf("got Directories %v want 2 entries", p.Directories)
	}
}
