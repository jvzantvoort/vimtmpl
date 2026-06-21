package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWhich(t *testing.T) {

	_, err := Which("unknowncommand")
	if err == nil {
		t.Errorf("got %s did not want that", err)
	}

}

func TestWhichFound(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "mycommand")
	if err := os.WriteFile(target, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create test binary: %v", err)
	}
	t.Setenv("PATH", dir)

	got, err := Which("mycommand")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != target {
		t.Errorf("got %q want %q", got, target)
	}
}

func TestWhichSearchesMultiplePathEntries(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	target := filepath.Join(dir2, "mycommand")
	if err := os.WriteFile(target, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create test binary: %v", err)
	}
	t.Setenv("PATH", dir1+":"+dir2)

	got, err := Which("mycommand")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != target {
		t.Errorf("got %q want %q", got, target)
	}
}
