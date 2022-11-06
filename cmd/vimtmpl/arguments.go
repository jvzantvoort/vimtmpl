package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/templates"
	log "github.com/sirupsen/logrus"
)

func Usage(lang string) string {
	if lang == "" {
		langs := templates.ListTemplateNames()
		lang = fmt.Sprintf("[%s]", strings.Join(langs, "|"))
	}
	return fmt.Sprintf("USAGE:\n\n\t%s %s [<filename>] <options>\n\n", os.Args[0], lang)
}

func ArgParse(args ...string) (*config.TemplateConfig, error) {

	if len(args) == 0 {
		fmt.Print(Usage(""))
		return &config.TemplateConfig{}, fmt.Errorf("Not enough arguments")
	}

	// take first argument to lang
	lang := args[0]
	args = args[1:]

	// force the setting of verbose before handling
	if len(args) >= 1 {
		for _, opt := range args {
			if opt == "-v" {
				log.SetLevel(log.DebugLevel)
			}
			if opt == "--verbose" {
				log.SetLevel(log.DebugLevel)
			}
		}
	}

	// Start logging
	log.Debugf("ArgParse, start")
	defer log.Debugf("ArgParse, end")

	log.Debugf("lang: %s", lang)

	tmplfile, err := templates.GetTemplateFile(lang)
	if err != nil {
		Usage("")
		return &config.TemplateConfig{}, fmt.Errorf("unable to get template for %s: %s\n\n", lang, err)
	}

	log.Debugf("template: %s", tmplfile)

	cfg := config.NewTemplateConfig(lang)
	cfg.Load()

	f := flag.NewFlagSet("prompt", flag.ExitOnError)

	f.BoolVar(&cfg.Verbose, "v", false, "Verbose logging")

	f.StringVar(&cfg.Company, "company", cfg.Company, "Company name")
	f.StringVar(&cfg.Company, "c", cfg.Company, "Company name")

	f.StringVar(&cfg.Copyright, "copyright", cfg.Copyright, "Copyright holder")

	f.StringVar(&cfg.Description, "description", "", "Script description")
	f.StringVar(&cfg.Description, "d", "", "Script description")

	f.StringVar(&cfg.License, "license", cfg.License, "License")
	f.StringVar(&cfg.License, "l", cfg.License, "License")

	f.StringVar(&cfg.MailAddress, "mailaddress", cfg.MailAdress, "mailaddress")
	f.StringVar(&cfg.MailAddress, "m", cfg.MailAdress, "mailaddress")

	f.StringVar(&cfg.ScriptName, "scriptname", "", "Script name")
	f.StringVar(&cfg.ScriptName, "s", "", "Script name")

	f.StringVar(&cfg.Title, "title", "", "Title (of e.g. python class)")
	f.StringVar(&cfg.Title, "t", "", "Title (of e.g. python class)")

	f.StringVar(&cfg.UserName, "username", cfg.UserName, "Users full name")
	f.StringVar(&cfg.UserName, "u", cfg.UserName, "Users full name")

	f.StringVar(&cfg.User, "user", cfg.User, "User account name")
	f.StringVar(&cfg.User, "U", cfg.User, "User account name")

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s", Usage(cfg.Lang))
		f.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n\n")
	}

	f.Parse(args)

	if len(cfg.ScriptName) == 0 {
		args := f.Args()
		if len(args) > 0 {
			cfg.ScriptName = args[0]
		}
	}
	if len(cfg.ScriptName) == 0 {
		cfg.Stdout = true
		cfg.ScriptName = "stdout"
	}
	return cfg, nil
}
