package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestEditorPrefix(t *testing.T) {
	e := Editor{}
	prefix := e.Prefix()
	if prefix != "Editor.TestEditorPrefix" {
		t.Errorf("got %q want %q", prefix, "Editor.TestEditorPrefix")
	}
}

func TestEditorExecute(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test relies on /bin/echo")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	e := Editor{Command: "/bin/echo", Cwd: cwd}

	stdout, stderr, err := e.Execute("hello", "world")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(stdout) != 1 || stdout[0] != "hello world" {
		t.Errorf("got stdout %v want [%q]", stdout, "hello world")
	}
	if len(stderr) != 0 {
		t.Errorf("got stderr %v want empty", stderr)
	}
}

func TestEditorExecuteCommandFails(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test relies on /bin/false")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	e := Editor{Command: "/bin/false", Cwd: cwd}

	_, _, err = e.Execute()
	if err == nil {
		t.Errorf("expected error from failing command")
	}
}

func TestEditorEdit(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test relies on /bin/echo")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	e := Editor{Command: "/bin/echo", Cwd: cwd}

	stdout, _, err := e.Edit("from-edit")
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if len(stdout) != 1 || stdout[0] != "from-edit" {
		t.Errorf("got stdout %v want [%q]", stdout, "from-edit")
	}
}

func TestNewEditor(t *testing.T) {
	dir := t.TempDir()
	vimPath := filepath.Join(dir, "vim")
	if err := os.WriteFile(vimPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create fake vim binary: %v", err)
	}
	t.Setenv("PATH", dir)

	e := NewEditor()
	if e.Path == nil {
		t.Fatalf("expected Path to be initialized")
	}
	if e.Cwd == "" {
		t.Errorf("expected Cwd to be populated")
	}
	if e.Command != vimPath {
		t.Errorf("got Command %q want %q", e.Command, vimPath)
	}
}

func TestNewEditorCommandNotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)

	e := NewEditor()
	if e.Command != "" {
		t.Errorf("expected empty Command when vim is not found, got %q", e.Command)
	}
}
