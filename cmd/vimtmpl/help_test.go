package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAvailableTemplatesSectionEmpty(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got := availableTemplatesSection()
	if !strings.Contains(got, "none found") {
		t.Errorf("expected 'none found' message, got %q", got)
	}
}

func TestAvailableTemplatesSectionListsTemplates(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	defaultDir := filepath.Join(home, ".templates.d", "default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatalf("failed to create default dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(defaultDir, "bash.gtmpl"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	got := availableTemplatesSection()
	if !strings.Contains(got, "bash") {
		t.Errorf("expected templates list to contain bash, got %q", got)
	}
}

func TestPrintHelp(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	oldStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = oldStdout })

	printHelp()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}

	got := string(out)
	if !strings.Contains(got, "NAME") {
		t.Errorf("expected help text to contain NAME section, got %q", got)
	}
	if !strings.Contains(got, "SUBCOMMANDS") {
		t.Errorf("expected help text to contain SUBCOMMANDS section, got %q", got)
	}
}
