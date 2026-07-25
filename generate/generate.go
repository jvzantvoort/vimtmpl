// Package generate renders template configs into file content and writes the
// result to disk. It is shared by the vimtmpl CLI and the vimtmplx TUI so both
// front-ends produce identical output.
package generate

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/templates"
	log "github.com/sirupsen/logrus"
)

// Render parses the template selected by cfg.Lang and executes it against cfg,
// returning the rendered content.
func Render(cfg *config.TemplateConfig) (string, error) {
	_, templatestring, err := templates.ParseLang(cfg.Lang)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("tmpl").Parse(templatestring)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, *cfg); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WriteFile writes the rendered content to cfg.FullPath, applying the file
// mode configured for cfg.Lang. Returns an error if the target file already
// exists or on write failure.
func WriteFile(cfg *config.TemplateConfig, content string) error {
	log.Debugf("WriteFile, start")
	defer log.Debugf("WriteFile, end")

	if config.TargetExists(cfg.FullPath) {
		return fmt.Errorf("target already exists: %s", cfg.FullPath)
	}

	file, err := os.Create(cfg.FullPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	obj := cfg.GetItem(cfg.Lang)
	log.Debugf("language found: %s", obj.Name)
	log.Debugf("mode: %o", obj.Mode)

	if _, err = file.WriteString(content); err != nil {
		log.Errorf("file.WriteString: error: %s", err)
		log.Errorf("%#v", content)
		return err
	}
	log.Debugf("write file: %s", cfg.FullPath)

	mode := os.FileMode(obj.Mode)
	if err := os.Chmod(cfg.FullPath, mode); err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("mode: %d\n", obj.Mode)
	return nil
}
