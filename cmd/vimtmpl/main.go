// Package main is the entry point for the vimtmpl command-line tool.
package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/generate"
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
	return generate.WriteFile(tmpl, content)
}

func main() {

	if len(os.Args) >= 2 && os.Args[1] == "help" {
		printHelp()
		return
	}

	if len(os.Args) >= 2 && os.Args[1] == "init" {
		if err := runInit(); err != nil {
			log.Errorf("init failed: %s", err)
			os.Exit(1)
		}
		return
	}

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

	// get template header (used for --info)
	tmplcfg, _, err := templates.ParseLang(cfg.Lang)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if cfg.Info {
		sections := []string{}
		sections = append(sections, tmplcfg.SectionStrings()...)
		if slices.Contains(sections, "switches") {
			fmt.Println("Template switches:")
			sec, err := tmplcfg.GetSection("switches")
			if err != nil {
				fmt.Printf("CRAP\n")
			}
			for name, description := range sec.KeysHash() {
				fmt.Printf("%-15s %s\n", name, description)
			}
		}
		os.Exit(0)
	}

	content, err := generate.Render(cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

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

	// ExecEditor replaces this process with the editor outright (like a
	// shell's "exec"); it only returns on platforms where that isn't
	// possible (e.g. Windows), in which case we fall back to spawning it.
	if err := generate.ExecEditor(cfg.FullPath); err != nil {
		log.Debugf("exec editor failed, falling back to subprocess: %s", err)
		if err := generate.OpenEditor(cfg.FullPath); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

}
