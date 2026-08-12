package main

import (
	"os"
	"path/filepath"
	"testing"
)

func installTemplateFixture(t *testing.T, home, name, content string) {
	t.Helper()
	dir := filepath.Join(home, ".templates.d", "default")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".gtmpl"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write template fixture: %v", err)
	}
}

func TestLoadSwitchesReturnsOrderedEntries(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	installTemplateFixture(t, home, "bash", "---\n[switches]\nsafe = Enable safety switches\nheader = Add basic header\n---\necho hi\n")

	switches := loadSwitches("bash")
	if len(switches) != 2 {
		t.Fatalf("expected 2 switches, got %d: %+v", len(switches), switches)
	}
	if switches[0].Name != "safe" || switches[0].Description != "Enable safety switches" {
		t.Errorf("unexpected first switch: %+v", switches[0])
	}
	if switches[1].Name != "header" || switches[1].Description != "Add basic header" {
		t.Errorf("unexpected second switch: %+v", switches[1])
	}
}

func TestLoadSwitchesNoHeader(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	installTemplateFixture(t, home, "plain", "echo hi\n")

	if switches := loadSwitches("plain"); switches != nil {
		t.Errorf("expected nil switches for headerless template, got %+v", switches)
	}
}

func TestModelGenerateWritesFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	installTemplateFixture(t, home, "bash",
		"---\n[switches]\nsafe = Enable safety switches\n---\n#!/bin/bash {{ if index .Flags \"safe\" }}-eu{{ end }}\necho {{.ScriptName}}\n")

	cfgFixture := "[DEFAULT]\nuser = dduck\n\n[bash]\nmode = 0644\n"
	if err := os.WriteFile(filepath.Join(home, ".template.cfg"), []byte(cfgFixture), 0644); err != nil {
		t.Fatalf("failed to write config fixture: %v", err)
	}

	workdir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(workdir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWd) })

	m := newModel("")
	if len(m.langs) != 1 || m.langs[0] != "bash" {
		t.Fatalf("expected langs=[bash], got %+v", m.langs)
	}
	if len(m.switches) != 1 {
		t.Fatalf("expected 1 switch loaded, got %+v", m.switches)
	}

	m.filename.SetValue("deploy.sh")
	m.switches[0].Enabled = true

	updated, cmd := m.generate()
	got, ok := updated.(model)
	if !ok {
		t.Fatalf("generate() did not return a model")
	}
	if cmd == nil {
		t.Errorf("expected a tea.Cmd from generate() to open the editor, got nil")
	}
	if got.status.isErr {
		t.Fatalf("expected success status, got error: %s", got.status.text)
	}

	target := filepath.Join(workdir, "deploy.sh")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}
	want := "#!/bin/bash -eu\necho deploy.sh\n"
	if string(data) != want {
		t.Errorf("got content %q, want %q", string(data), want)
	}
}

func TestModelGenerateRequiresFilename(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	installTemplateFixture(t, home, "bash", "echo hi\n")

	m := newModel("")
	updated, _ := m.generate()
	got := updated.(model)
	if !got.status.isErr {
		t.Errorf("expected error status when filename is empty")
	}
}

func TestCycleLangWrapsAndReloadsSwitches(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	installTemplateFixture(t, home, "aaa", "---\n[switches]\nfoo = Foo switch\n---\nbody\n")
	installTemplateFixture(t, home, "bbb", "body\n")

	m := newModel("")
	if len(m.langs) != 2 {
		t.Fatalf("expected 2 langs, got %+v", m.langs)
	}

	m.cycleLang(-1)
	if m.langIdx != len(m.langs)-1 {
		t.Errorf("expected wrap to last index, got %d", m.langIdx)
	}

	m.cycleLang(1)
	if m.langIdx != 0 {
		t.Errorf("expected wrap back to 0, got %d", m.langIdx)
	}
}
