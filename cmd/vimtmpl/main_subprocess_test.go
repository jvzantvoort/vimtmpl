package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMainHelperProcess is not a real test. It is invoked as a subprocess by
// the Test*ViaSubprocess tests below so that main(), which calls os.Exit,
// can be exercised end-to-end without killing the real test binary.
func TestMainHelperProcess(t *testing.T) {
	if os.Getenv("VIMTMPL_BE_MAIN") != "1" {
		t.Skip("not running as main helper process")
	}
	args := strings.Split(os.Getenv("VIMTMPL_MAIN_ARGS"), "\x1f")
	os.Args = append([]string{"vimtmpl"}, args...)
	main()
}

func runMainSubprocess(t *testing.T, home string, args ...string) (stdout string, exitCode int) {
	t.Helper()

	out, err := runSelfAsSubprocess(t, home, args)
	stdout = out

	if err == nil {
		return stdout, 0
	}
	if exitErr, ok := asExitError(err); ok {
		return stdout, exitErr
	}
	t.Fatalf("subprocess failed unexpectedly: %v\noutput:\n%s", err, out)
	return
}

func TestMainHelpViaSubprocess(t *testing.T) {
	home := t.TempDir()

	stdout, code := runMainSubprocess(t, home, "help")
	if code != 0 {
		t.Fatalf("got exit code %d want 0, output:\n%s", code, stdout)
	}
	if !strings.Contains(stdout, "NAME") {
		t.Errorf("expected help output to contain NAME, got %q", stdout)
	}
}

func TestMainInitViaSubprocess(t *testing.T) {
	home := t.TempDir()

	_, code := runMainSubprocess(t, home, "init")
	if code != 0 {
		t.Fatalf("got exit code %d want 0", code)
	}

	if _, err := os.Stat(filepath.Join(home, ".templates.d", "default", "bash.gtmpl")); err != nil {
		t.Errorf("expected bash.gtmpl to be installed: %v", err)
	}
}

func TestMainGenerateViaSubprocess(t *testing.T) {
	home := t.TempDir()
	defaultDir := filepath.Join(home, ".templates.d", "default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatalf("failed to create default dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(defaultDir, "bash.gtmpl"), []byte("#!/bin/bash\necho {{.ScriptName}}\n"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	out := filepath.Join(home, "out.sh")
	_, code := runMainSubprocess(t, home, "bash", out)
	if code != 0 {
		t.Fatalf("got exit code %d want 0", code)
	}

	// The bundled bash.gtmpl template has no matching config section here,
	// so WriteFile chmods the output to mode 0 (TemplateConfig.GetItem's
	// zero-value default). Restore read access before inspecting it.
	if err := os.Chmod(out, 0644); err != nil {
		t.Fatalf("failed to chmod output file: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("expected output file to be created: %v", err)
	}
	if !strings.Contains(string(data), "echo out.sh") {
		t.Errorf("got content %q, expected it to mention out.sh", string(data))
	}
}

func TestMainGenerateUnknownTemplateViaSubprocess(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".templates.d", "default"), 0755); err != nil {
		t.Fatalf("failed to create default dir: %v", err)
	}

	out := filepath.Join(home, "out.sh")
	_, code := runMainSubprocess(t, home, "doesnotexist", out)
	if code == 0 {
		t.Errorf("expected non-zero exit code for unknown template")
	}
}

func TestMainMissingArgsViaSubprocess(t *testing.T) {
	home := t.TempDir()

	_, code := runMainSubprocess(t, home)
	if code == 0 {
		t.Errorf("expected non-zero exit code when no arguments are given")
	}
}
