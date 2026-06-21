package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestUserHomeDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Setenv("HOMEDRIVE", "C:")
		t.Setenv("HOMEPATH", `\Users\test`)
		t.Setenv("USERPROFILE", "")
		got := UserHomeDir()
		want := `C:\Users\test`
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		return
	}

	t.Setenv("HOME", "/tmp/somehome")
	got := UserHomeDir()
	if got != "/tmp/somehome" {
		t.Errorf("got %q want %q", got, "/tmp/somehome")
	}
}

func TestUserHomeDirWindowsFallsBackToUserProfile(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only behaviour")
	}
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")
	t.Setenv("USERPROFILE", `C:\Users\fallback`)
	got := UserHomeDir()
	if got != `C:\Users\fallback` {
		t.Errorf("got %q want %q", got, `C:\Users\fallback`)
	}
}

func TestUserName(t *testing.T) {
	got := UserName()
	if got == "" {
		t.Errorf("expected non-empty username")
	}
}

func TestTargetExists(t *testing.T) {
	dir := t.TempDir()

	if TargetExists(filepath.Join(dir, "nope")) {
		t.Errorf("expected non-existent path to report false")
	}

	existing := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(existing, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if !TargetExists(existing) {
		t.Errorf("expected existing file to report true")
	}

	if !TargetExists(dir) {
		t.Errorf("expected existing directory to report true")
	}
}
