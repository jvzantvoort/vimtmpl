package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/content"
)

func setupTempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return dir
}

func TestRunInitCreatesDirectories(t *testing.T) {
	home := setupTempHome(t)

	if err := runInit(); err != nil {
		t.Fatalf("runInit returned error: %v", err)
	}

	for _, sub := range []string{
		filepath.Join(".templates.d", "default"),
		filepath.Join(".templates.d", "local"),
	} {
		path := filepath.Join(home, sub)
		if !config.TargetExists(path) {
			t.Errorf("expected directory %q to exist after init", path)
		}
	}
}

func TestRunInitInstallsAllTemplates(t *testing.T) {
	home := setupTempHome(t)

	if err := runInit(); err != nil {
		t.Fatalf("runInit returned error: %v", err)
	}

	defaultDir := filepath.Join(home, ".templates.d", "default")

	entries, err := fs.ReadDir(content.Templates, ".")
	if err != nil {
		t.Fatalf("cannot read embedded templates: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		dest := filepath.Join(defaultDir, entry.Name())
		if !config.TargetExists(dest) {
			t.Errorf("expected installed template %q to exist", dest)
		}
	}
}

func TestRunInitCreatesConfig(t *testing.T) {
	home := setupTempHome(t)

	if err := runInit(); err != nil {
		t.Fatalf("runInit returned error: %v", err)
	}

	cfgPath := filepath.Join(home, config.ConfigFilename)
	if !config.TargetExists(cfgPath) {
		t.Errorf("expected configuration file %q to exist after init", cfgPath)
	}
}

func TestRunInitIsIdempotent(t *testing.T) {
	setupTempHome(t)

	if err := runInit(); err != nil {
		t.Fatalf("first runInit returned error: %v", err)
	}
	if err := runInit(); err != nil {
		t.Fatalf("second runInit returned error: %v", err)
	}
}

func TestRunInitPreservesExistingTemplate(t *testing.T) {
	home := setupTempHome(t)

	defaultDir := filepath.Join(home, ".templates.d", "default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatalf("failed to create default dir: %v", err)
	}

	sentinel := []byte("# sentinel content")
	existing := filepath.Join(defaultDir, "bash.gtmpl")
	if err := os.WriteFile(existing, sentinel, 0644); err != nil {
		t.Fatalf("failed to write sentinel file: %v", err)
	}

	if err := runInit(); err != nil {
		t.Fatalf("runInit returned error: %v", err)
	}

	got, err := os.ReadFile(existing)
	if err != nil {
		t.Fatalf("failed to read template after init: %v", err)
	}
	if string(got) != string(sentinel) {
		t.Errorf("existing template was overwritten: got %q, want %q", got, sentinel)
	}
}

func TestRunInitPreservesExistingConfig(t *testing.T) {
	home := setupTempHome(t)

	sentinel := []byte("[default]\ncompany = sentinel\n")
	cfgPath := filepath.Join(home, config.ConfigFilename)
	if err := os.WriteFile(cfgPath, sentinel, 0644); err != nil {
		t.Fatalf("failed to write sentinel config: %v", err)
	}

	if err := runInit(); err != nil {
		t.Fatalf("runInit returned error: %v", err)
	}

	got, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("failed to read config after init: %v", err)
	}
	if string(got) != string(sentinel) {
		t.Errorf("existing config was overwritten: got %q, want %q", got, sentinel)
	}
}
