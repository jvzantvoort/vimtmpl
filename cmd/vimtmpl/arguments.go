package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
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

	set1 := make(map[string]string)
	set2 := make(map[string]string)

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

	set1["verbose"] = strconv.FormatBool(false)
	f.BoolVar(&cfg.Verbose, "v", false, "Verbose logging")

	set1["company"] = cfg.Company
	f.StringVar(&cfg.Company, "company", cfg.Company, "Company name")
	f.StringVar(&cfg.Company, "c", cfg.Company, "Company name")

	set1["copyright"] = cfg.Copyright
	f.StringVar(&cfg.Copyright, "copyright", cfg.Copyright, "Copyright holder")

	set1["description"] = cfg.Description
	f.StringVar(&cfg.Description, "description", "", "Script description")
	f.StringVar(&cfg.Description, "d", "", "Script description")

	set1["license"] = cfg.License
	f.StringVar(&cfg.License, "license", cfg.License, "License")
	f.StringVar(&cfg.License, "l", cfg.License, "License")

	set1["mailaddress"] = cfg.MailAddress
	f.StringVar(&cfg.MailAddress, "mailaddress", cfg.MailAddress, "mailaddress")
	f.StringVar(&cfg.MailAddress, "m", cfg.MailAddress, "mailaddress")

	set1["fullpath"] = cfg.FullPath
	f.StringVar(&cfg.FullPath, "scriptname", "", "Script name")
	f.StringVar(&cfg.FullPath, "s", "", "Script name")

	set1["title"] = cfg.Title
	f.StringVar(&cfg.Title, "title", "", "Title (of e.g. python class)")
	f.StringVar(&cfg.Title, "t", "", "Title (of e.g. python class)")

	set1["username"] = cfg.UserName
	f.StringVar(&cfg.UserName, "username", cfg.UserName, "Users full name")
	f.StringVar(&cfg.UserName, "u", cfg.UserName, "Users full name")

	set1["user"] = cfg.User
	f.StringVar(&cfg.User, "user", cfg.User, "User account name")
	f.StringVar(&cfg.User, "U", cfg.User, "User account name")

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s", Usage(cfg.Lang))
		f.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n\n")
	}

	f.Parse(args)

	if len(cfg.FullPath) == 0 {
		args := f.Args()
		if len(args) > 0 {
			cfg.FullPath = args[0]
		}
	}

	if len(cfg.FullPath) == 0 {
		cfg.Stdout = true
		cfg.FullPath = "stdout"
	} else {
		cfg.ScriptName = path.Base(cfg.FullPath)

	}

	set2["company"] = cfg.Company
	set2["copyright"] = cfg.Copyright
	set2["description"] = cfg.Description
	set2["license"] = cfg.License
	set2["mailaddress"] = cfg.MailAddress
	set2["fullpath"] = cfg.FullPath
	set2["title"] = cfg.Title
	set2["user"] = cfg.User
	set2["username"] = cfg.UserName
	set1["verbose"] = strconv.FormatBool(cfg.Verbose)

	for keyn, keyv := range set1 {
		if keyv == set2[keyn] {
			continue
		}
		log.Debugf("   changed %s from %s to %s", keyn, keyv, set2[keyn])
	}

	return cfg, nil
}
