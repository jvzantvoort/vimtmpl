package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jvzantvoort/vimtmpl/config"
)

func TestWriteFileCreatesFileWithContent(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "out.sh")

	cfg := &config.TemplateConfig{
		FullPath: target,
		Lang:     "bash",
		Items: []*config.TemplateItem{
			{Name: "bash", Mode: 0755},
		},
	}

	if err := WriteFile(cfg, "#!/bin/bash\necho hi\n"); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(data) != "#!/bin/bash\necho hi\n" {
		t.Errorf("got content %q, want %q", string(data), "#!/bin/bash\necho hi\n")
	}

	if runtime.GOOS != "windows" {
		info, err := os.Stat(target)
		if err != nil {
			t.Fatalf("failed to stat written file: %v", err)
		}
		if info.Mode().Perm() != 0755 {
			t.Errorf("got mode %o want %o", info.Mode().Perm(), 0755)
		}
	}
}

func TestWriteFileTargetAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "out.sh")
	if err := os.WriteFile(target, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	cfg := &config.TemplateConfig{FullPath: target, Lang: "bash"}

	err := WriteFile(cfg, "new content")
	if err == nil {
		t.Errorf("expected error when target already exists")
	}

	data, _ := os.ReadFile(target)
	if string(data) != "existing" {
		t.Errorf("expected existing file to remain untouched, got %q", string(data))
	}
}

func TestWriteFileUnknownLangUsesDefaultMode(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "out.sh")

	cfg := &config.TemplateConfig{FullPath: target, Lang: "doesnotexist"}

	if err := WriteFile(cfg, "content"); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if _, err := os.Stat(target); err != nil {
		t.Fatalf("expected file to be created: %v", err)
	}
}
