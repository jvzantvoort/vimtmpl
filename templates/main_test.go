package templates

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func TestListTemplateFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "bash.gtmpl"), "bash content")
	writeFile(t, filepath.Join(dir, "python.gtmpl"), "python content")
	writeFile(t, filepath.Join(dir, "README.md"), "not a template")
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	got, err := ListTemplateFile(dir)
	if err != nil {
		t.Fatalf("ListTemplateFile returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("got %d entries, want 2: %v", len(got), got)
	}

	wantBash := filepath.Join(dir, "bash.gtmpl")
	if got["bash"] != wantBash {
		t.Errorf("got bash path %q want %q", got["bash"], wantBash)
	}
	wantPython := filepath.Join(dir, "python.gtmpl")
	if got["python"] != wantPython {
		t.Errorf("got python path %q want %q", got["python"], wantPython)
	}
	if _, ok := got["README"]; ok {
		t.Errorf("did not expect README to be listed")
	}
}

func setupTemplatesHome(t *testing.T) (home, defaultDir, localDir string) {
	t.Helper()
	home = t.TempDir()
	t.Setenv("HOME", home)
	defaultDir = filepath.Join(home, ".templates.d", "default")
	localDir = filepath.Join(home, ".templates.d", "local")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatalf("failed to create default dir: %v", err)
	}
	if err := os.MkdirAll(localDir, 0755); err != nil {
		t.Fatalf("failed to create local dir: %v", err)
	}
	return
}

func TestListTemplatesNoTemplatesDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	_, err := ListTemplates()
	if err == nil {
		t.Errorf("expected error when .templates.d does not exist")
	}
}

func TestListTemplatesMergesDefaultAndLocal(t *testing.T) {
	_, defaultDir, localDir := setupTemplatesHome(t)

	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), "default bash")
	writeFile(t, filepath.Join(defaultDir, "python.gtmpl"), "default python")
	writeFile(t, filepath.Join(localDir, "python.gtmpl"), "local python override")
	writeFile(t, filepath.Join(localDir, "go.gtmpl"), "local go")

	got, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates returned error: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("got %d entries, want 3: %v", len(got), got)
	}

	wantPython := filepath.Join(localDir, "python.gtmpl")
	if got["python"] != wantPython {
		t.Errorf("expected local to override default for python: got %q want %q", got["python"], wantPython)
	}
}

func TestListTemplatesOnlyDefaultDirExists(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	defaultDir := filepath.Join(home, ".templates.d", "default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatalf("failed to create default dir: %v", err)
	}
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), "default bash")

	got, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1: %v", len(got), got)
	}
}

func TestListTemplateNames(t *testing.T) {
	_, defaultDir, _ := setupTemplatesHome(t)
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), "bash content")
	writeFile(t, filepath.Join(defaultDir, "python.gtmpl"), "python content")

	names := ListTemplateNames()
	if len(names) != 2 {
		t.Fatalf("got %d names, want 2: %v", len(names), names)
	}
}

func TestListTemplateNamesNoTemplatesDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	names := ListTemplateNames()
	if len(names) != 0 {
		t.Errorf("expected no names, got %v", names)
	}
}

func TestGetTemplateFile(t *testing.T) {
	_, defaultDir, _ := setupTemplatesHome(t)
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), "bash content")

	got, err := GetTemplateFile("bash")
	if err != nil {
		t.Fatalf("GetTemplateFile returned error: %v", err)
	}
	want := filepath.Join(defaultDir, "bash.gtmpl")
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}

	_, err = GetTemplateFile("doesnotexist")
	if err == nil {
		t.Errorf("expected error for unknown language")
	}
}

func TestGetTemplateFileTemplatesDirMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	_, err := GetTemplateFile("bash")
	if err == nil {
		t.Errorf("expected error when templates dir is missing")
	}
}

func TestParseLangWithHeader(t *testing.T) {
	_, defaultDir, _ := setupTemplatesHome(t)
	content := "---\n[switches]\nfoo = bar description\n---\nbody content here\n"
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), content)

	header, body, err := ParseLang("bash")
	if err != nil {
		t.Fatalf("ParseLang returned error: %v", err)
	}
	if body != "body content here\n" {
		t.Errorf("got body %q want %q", body, "body content here\n")
	}
	if !header.HasSection("switches") {
		t.Errorf("expected header to have switches section")
	}
}

func TestParseLangWithoutHeader(t *testing.T) {
	_, defaultDir, _ := setupTemplatesHome(t)
	content := "plain body without header\n"
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), content)

	_, body, err := ParseLang("bash")
	if err != nil {
		t.Fatalf("ParseLang returned error: %v", err)
	}
	if body != content {
		t.Errorf("got body %q want %q", body, content)
	}
}

func TestParseLangUnclosedHeader(t *testing.T) {
	_, defaultDir, _ := setupTemplatesHome(t)
	content := "---\nno closing marker here\n"
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), content)

	_, body, err := ParseLang("bash")
	if err != nil {
		t.Fatalf("ParseLang returned error: %v", err)
	}
	if body != content {
		t.Errorf("got body %q want %q", body, content)
	}
}

func TestParseLangUnknownLanguage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	_, _, err := ParseLang("doesnotexist")
	if err == nil {
		t.Errorf("expected error for unknown language")
	}
}

func TestSplitDocument(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantHeader string
		wantBody   string
	}{
		{"no header", "just body\n", "", "just body\n"},
		{"with header", "---\nkey = value\n---\nbody text\n", "key = value\n", "body text\n"},
		{"unclosed header", "---\nkey = value\n", "", "---\nkey = value\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header, body := splitDocument(tt.content)
			if header != tt.wantHeader {
				t.Errorf("got header %q want %q", header, tt.wantHeader)
			}
			if body != tt.wantBody {
				t.Errorf("got body %q want %q", body, tt.wantBody)
			}
		})
	}
}

func TestGetTemplateContent(t *testing.T) {
	_, defaultDir, _ := setupTemplatesHome(t)
	content := "---\nkey = value\n---\nthe body\n"
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), content)

	got, err := GetTemplateContent("bash")
	if err != nil {
		t.Fatalf("GetTemplateContent returned error: %v", err)
	}
	if got != "the body\n" {
		t.Errorf("got %q want %q", got, "the body\n")
	}
}

func TestGetTemplateContentUnknownLanguage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	_, err := GetTemplateContent("doesnotexist")
	if err == nil {
		t.Errorf("expected error for unknown language")
	}
}

func TestGetTemplateConfig(t *testing.T) {
	_, defaultDir, _ := setupTemplatesHome(t)
	content := "---\nkey = value\n---\nthe body\n"
	writeFile(t, filepath.Join(defaultDir, "bash.gtmpl"), content)

	got, err := GetTemplateConfig("bash")
	if err != nil {
		t.Fatalf("GetTemplateConfig returned error: %v", err)
	}
	if got != "key = value\n" {
		t.Errorf("got %q want %q", got, "key = value\n")
	}
}

func TestGetTemplateConfigUnknownLanguage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	_, err := GetTemplateConfig("doesnotexist")
	if err == nil {
		t.Errorf("expected error for unknown language")
	}
}
