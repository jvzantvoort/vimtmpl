package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-ini/ini"
)

func TestNewTemplateConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	tc := NewTemplateConfig("bash")

	if tc.Lang != "bash" {
		t.Errorf("got Lang %q want %q", tc.Lang, "bash")
	}
	if tc.Stdout {
		t.Errorf("expected Stdout to default to false")
	}
	if tc.Company != "company" {
		t.Errorf("got Company %q want %q", tc.Company, "company")
	}
	if tc.Copyright != "copyright" {
		t.Errorf("got Copyright %q want %q", tc.Copyright, "copyright")
	}
	if tc.License != "license" {
		t.Errorf("got License %q want %q", tc.License, "license")
	}
	if tc.MailAddress != "mailaddress" {
		t.Errorf("got MailAddress %q want %q", tc.MailAddress, "mailaddress")
	}
	if tc.UserName != "username" {
		t.Errorf("got UserName %q want %q", tc.UserName, "username")
	}
	if tc.Date == "" {
		t.Errorf("expected Date to be populated")
	}
	if tc.Year == "" {
		t.Errorf("expected Year to be populated")
	}
	if tc.Homedir != home {
		t.Errorf("got Homedir %q want %q", tc.Homedir, home)
	}
	wantFilepath := filepath.Join(home, ConfigFilename)
	if tc.Filepath != wantFilepath {
		t.Errorf("got Filepath %q want %q", tc.Filepath, wantFilepath)
	}
	if tc.Flags == nil {
		t.Errorf("expected Flags map to be initialized")
	}
}

func TestEnabled(t *testing.T) {
	tc := TemplateConfig{Flags: map[string]bool{"on": true, "off": false}}

	if !tc.Enabled("on") {
		t.Errorf("expected 'on' to be enabled")
	}
	if tc.Enabled("off") {
		t.Errorf("expected 'off' to be disabled")
	}
	if tc.Enabled("missing") {
		t.Errorf("expected missing key to be disabled")
	}
}

func newTestConfigFile(t *testing.T, content string) *ini.File {
	t.Helper()
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, []byte(content))
	if err != nil {
		t.Fatalf("failed to load test ini content: %v", err)
	}
	return cfg
}

func TestGetKeyAsString(t *testing.T) {
	content := `
[DEFAULT]
company = DefaultCorp
description = default description

[bash]
company = BashCorp
`
	cfg := newTestConfigFile(t, content)

	tc := TemplateConfig{Lang: "bash", Object: cfg}
	if got := tc.GetKeyAsString("company"); got != "BashCorp" {
		t.Errorf("got %q want %q", got, "BashCorp")
	}

	// falls back to DEFAULT section when lang-specific key is absent
	if got := tc.GetKeyAsString("description"); got != "default description" {
		t.Errorf("got %q want %q", got, "default description")
	}

	// known empty-default keys
	tcEmpty := TemplateConfig{Lang: "missinglang", Object: cfg}
	if got := tcEmpty.GetKeyAsString("description"); got != "default description" {
		// description exists in DEFAULT so falls back there even for unknown lang
		t.Errorf("got %q want %q", got, "default description")
	}
	if got := tcEmpty.GetKeyAsString("extension"); got != "" {
		t.Errorf("got %q want empty string", got)
	}
	if got := tcEmpty.GetKeyAsString("nonexistentkey"); got != "" {
		t.Errorf("got %q want empty string", got)
	}
}

func TestGetKeyAsInt(t *testing.T) {
	content := `
[DEFAULT]
mode = 644

[bash]
mode = 755
`
	cfg := newTestConfigFile(t, content)

	tc := TemplateConfig{Lang: "bash", Object: cfg}
	if got := tc.GetKeyAsInt("mode"); got != 755 {
		t.Errorf("got %d want %d", got, 755)
	}

	tcDefault := TemplateConfig{Lang: "missinglang", Object: cfg}
	if got := tcDefault.GetKeyAsInt("mode"); got != 644 {
		t.Errorf("got %d want %d", got, 644)
	}

	if got := tcDefault.GetKeyAsInt("nonexistentkey"); got != 0 {
		t.Errorf("got %d want 0", got)
	}
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.cfg")
	if err := os.WriteFile(cfgPath, []byte("[DEFAULT]\ncompany = Acme\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	tc := TemplateConfig{Filepath: cfgPath}
	cfg := tc.LoadFile()
	if cfg == nil {
		t.Fatalf("expected non-nil ini file")
	}
	key, err := cfg.Section(ini.DefaultSection).GetKey("company")
	if err != nil {
		t.Fatalf("expected company key: %v", err)
	}
	if key.String() != "Acme" {
		t.Errorf("got %q want %q", key.String(), "Acme")
	}
}

func TestLoadFileMissing(t *testing.T) {
	tc := TemplateConfig{Filepath: filepath.Join(t.TempDir(), "doesnotexist.cfg")}
	cfg := tc.LoadFile()
	if cfg == nil {
		t.Fatalf("expected an empty ini file rather than nil")
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.cfg")
	content := `
[DEFAULT]
company = Acme
copyright = Acme Corp
license = MIT
mailaddress = dev@acme.com
username = Dev User
user = devuser

[bash]
description = bash script
extension = sh
mode = 755
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	tc := &TemplateConfig{Filepath: cfgPath, Lang: "bash"}
	tc.Load()

	if tc.Company != "Acme" {
		t.Errorf("got Company %q want %q", tc.Company, "Acme")
	}
	if tc.Copyright != "Acme Corp" {
		t.Errorf("got Copyright %q want %q", tc.Copyright, "Acme Corp")
	}
	if tc.License != "MIT" {
		t.Errorf("got License %q want %q", tc.License, "MIT")
	}
	if tc.MailAddress != "dev@acme.com" {
		t.Errorf("got MailAddress %q want %q", tc.MailAddress, "dev@acme.com")
	}
	if tc.UserName != "Dev User" {
		t.Errorf("got UserName %q want %q", tc.UserName, "Dev User")
	}
	if tc.User != "devuser" {
		t.Errorf("got User %q want %q", tc.User, "devuser")
	}

	if len(tc.Items) == 0 {
		t.Fatalf("expected items to be populated")
	}

	found := false
	for _, item := range tc.Items {
		if item.Name == "bash" {
			found = true
			if item.Description != "bash script" {
				t.Errorf("got Description %q want %q", item.Description, "bash script")
			}
			if item.Extension != "sh" {
				t.Errorf("got Extension %q want %q", item.Extension, "sh")
			}
			if item.Mode != 755 {
				t.Errorf("got Mode %d want %d", item.Mode, 755)
			}
		}
	}
	if !found {
		t.Errorf("expected to find bash item in Items")
	}
}

func TestGetItem(t *testing.T) {
	tc := TemplateConfig{
		Items: []*TemplateItem{
			{Name: "bash", Description: "bash script"},
			{Name: "python", Description: "python script"},
		},
	}

	got := tc.GetItem("bash")
	if got.Description != "bash script" {
		t.Errorf("got %q want %q", got.Description, "bash script")
	}

	missing := tc.GetItem("nonexistent")
	if missing.Name != "nonexistent" {
		t.Errorf("got Name %q want %q", missing.Name, "nonexistent")
	}
	if missing.Description != "" {
		t.Errorf("expected empty description for unknown item, got %q", missing.Description)
	}
}

func TestSaveTo(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.cfg")

	tc := TemplateConfig{
		Company:     "Acme",
		Copyright:   "Acme Corp",
		License:     "MIT",
		MailAddress: "dev@acme.com",
		UserName:    "Dev User",
		User:        "devuser",
		Items: []*TemplateItem{
			{Name: "bash", Description: "bash script"},
		},
	}

	if err := tc.SaveTo(dest); err != nil {
		t.Fatalf("SaveTo returned error: %v", err)
	}

	saved, err := ini.Load(dest)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	key, err := saved.Section("").GetKey("company")
	if err != nil {
		t.Fatalf("expected company key: %v", err)
	}
	if key.String() != "Acme" {
		t.Errorf("got %q want %q", key.String(), "Acme")
	}

	if !saved.HasSection("bash") {
		t.Errorf("expected bash section to exist in saved config")
	}
	descKey, err := saved.Section("bash").GetKey("description")
	if err != nil {
		t.Fatalf("expected description key: %v", err)
	}
	if descKey.String() != "bash script" {
		t.Errorf("got %q want %q", descKey.String(), "bash script")
	}
}

func TestSaveToInvalidPath(t *testing.T) {
	tc := TemplateConfig{}
	err := tc.SaveTo(filepath.Join(t.TempDir(), "nonexistent-dir", "out.cfg"))
	if err == nil {
		t.Errorf("expected error when saving to a non-existent directory")
	}
}
