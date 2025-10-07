// Package main is the entry point for the vimtmpl command-line tool.
package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/templates"
	log "github.com/sirupsen/logrus"
)

// init configures the logger for the application.
func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// WriteFile writes the generated template content to a file.
// Returns an error if the target file already exists or on write failure.
func WriteFile(tmpl *config.TemplateConfig, content string) error {
	log.Debugf("WriteFile, start")
	defer log.Debugf("WriteFile, end")

	if config.TargetExists(tmpl.FullPath) {
		return fmt.Errorf("target already exists: %s", tmpl.FullPath)
	}

	file, _ := os.Create(tmpl.FullPath)
	defer file.Close()

	obj := tmpl.GetItem(tmpl.Lang)
	log.Debugf("language found: %s", obj.Name)
	log.Debugf("mode: %o", obj.Mode)

	_, err := file.WriteString(content)
	if err != nil {
		log.Errorf("file.WriteString: error: %s", err)
		log.Errorf("%#v", content)
		return err
	} else {
		log.Debugf("write file: %s", tmpl.FullPath)
	}

	mode := os.FileMode(obj.Mode)
	if err := os.Chmod(tmpl.FullPath, mode); err != nil {
		log.Error(err)
		return err
	} else {
		log.Debugf("mode: %d\n", obj.Mode)
	}
	return nil
}

func main() {

	cfg, err := ArgParse()
	if err != nil {
		log.Errorf("Failed: %s", err)
		os.Exit(1)
	}

	if cfg.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Setup description
	if len(cfg.Description) == 0 {
		cfg.Description = cfg.GetKeyAsString("description")
	}

	// get template content
	templatestring, _ := templates.GetTemplateContent(cfg.Lang)
	text_template, err := template.New("tmpl").Parse(templatestring)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	err = text_template.Execute(buf, *cfg)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	content := buf.String()

	// Print to stdout if needed
	if cfg.Stdout {
		fmt.Print(content)
		return
	}

	err = WriteFile(cfg, content)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
