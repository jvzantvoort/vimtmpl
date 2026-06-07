package content

import (
	"io/fs"
	"testing"
)

var expectedTemplates = []string{
	"bash.gtmpl",
	"go.gtmpl",
	"playbook.gtmpl",
	"python.gtmpl",
	"pythonclass.gtmpl",
}

func TestTemplatesHasExpectedFiles(t *testing.T) {
	for _, name := range expectedTemplates {
		_, err := Templates.Open(name)
		if err != nil {
			t.Errorf("expected embedded template %q not found: %v", name, err)
		}
	}
}

func TestTemplatesFilesAreNonEmpty(t *testing.T) {
	for _, name := range expectedTemplates {
		data, err := fs.ReadFile(Templates, name)
		if err != nil {
			t.Errorf("cannot read embedded template %q: %v", name, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("embedded template %q is empty", name)
		}
	}
}
