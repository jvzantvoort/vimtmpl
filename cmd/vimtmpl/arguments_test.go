package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// resetFlags gives each test a fresh pflag.CommandLine, since parseFlags
// registers flags on the package-level singleton and pflag panics if a flag
// is registered twice.
func resetFlags(t *testing.T) {
	t.Helper()
	old := pflag.CommandLine
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	t.Cleanup(func() { pflag.CommandLine = old })
}

func setupTemplatesHomeForArgs(t *testing.T) (home, defaultDir string) {
	t.Helper()
	home = t.TempDir()
	t.Setenv("HOME", home)
	defaultDir = filepath.Join(home, ".templates.d", "default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatalf("failed to create default dir: %v", err)
	}
	return
}

func TestUsageWithLang(t *testing.T) {
	got := Usage("bash")
	if !strings.Contains(got, "USAGE:") {
		t.Errorf("expected usage text to contain USAGE:, got %q", got)
	}
	if !strings.Contains(got, "bash") {
		t.Errorf("expected usage text to mention lang %q, got %q", "bash", got)
	}
}

func TestUsageWithoutLang(t *testing.T) {
	_, defaultDir := setupTemplatesHomeForArgs(t)
	if err := os.WriteFile(filepath.Join(defaultDir, "bash.gtmpl"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	got := Usage("")
	if !strings.Contains(got, "bash") {
		t.Errorf("expected usage text to list available templates, got %q", got)
	}
}

func TestArgParseSuccess(t *testing.T) {
	_, defaultDir := setupTemplatesHomeForArgs(t)
	if err := os.WriteFile(filepath.Join(defaultDir, "bash.gtmpl"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	resetFlags(t)

	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	outFile := filepath.Join(t.TempDir(), "out.sh")
	os.Args = []string{"vimtmpl", "-m", "dev@example.com", "-c", "Acme", "-t", "MyTitle", "bash", outFile}

	cfg, err := ArgParse()
	if err != nil {
		t.Fatalf("ArgParse returned error: %v", err)
	}

	if cfg.Lang != "bash" {
		t.Errorf("got Lang %q want %q", cfg.Lang, "bash")
	}
	if cfg.FullPath != outFile {
		t.Errorf("got FullPath %q want %q", cfg.FullPath, outFile)
	}
	if cfg.MailAddress != "dev@example.com" {
		t.Errorf("got MailAddress %q want %q", cfg.MailAddress, "dev@example.com")
	}
	if cfg.Company != "Acme" {
		t.Errorf("got Company %q want %q", cfg.Company, "Acme")
	}
	if cfg.Title != "MyTitle" {
		t.Errorf("got Title %q want %q", cfg.Title, "MyTitle")
	}
	if cfg.ScriptName != "out.sh" {
		t.Errorf("got ScriptName %q want %q", cfg.ScriptName, "out.sh")
	}
}

func TestArgParseUnknownTemplate(t *testing.T) {
	setupTemplatesHomeForArgs(t)

	resetFlags(t)

	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	os.Args = []string{"vimtmpl", "doesnotexist", filepath.Join(t.TempDir(), "out.sh")}

	_, err := ArgParse()
	if err == nil {
		t.Errorf("expected error for unknown template")
	}
}

func TestArgParseFlagsWithCommaList(t *testing.T) {
	_, defaultDir := setupTemplatesHomeForArgs(t)
	if err := os.WriteFile(filepath.Join(defaultDir, "bash.gtmpl"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	resetFlags(t)

	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	outFile := filepath.Join(t.TempDir(), "out.sh")
	os.Args = []string{"vimtmpl", "-f", "one,two", "-f", "three", "bash", outFile}

	cfg, err := ArgParse()
	if err != nil {
		t.Fatalf("ArgParse returned error: %v", err)
	}

	for _, want := range []string{"one", "two", "three"} {
		if !cfg.Flags[want] {
			t.Errorf("expected flag %q to be enabled, got %v", want, cfg.Flags)
		}
	}
}
