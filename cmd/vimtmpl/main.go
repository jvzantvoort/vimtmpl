package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jvzantvoort/vimtmpl"
	"github.com/jvzantvoort/vimtmpl/config"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	Verbose     bool
	Company     string
	Copyright   string
	Description string
	License     string
	MailAddress string
	ScriptName  string
	Title       string
	User        string
	UserName    string
	Lang        string
}

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

func Usage() {

	langs := vimtmpl.ListTemplateNames()
	fmt.Printf("%s [%s] <options>\n", os.Args[0], strings.Join(langs, "|"))

}

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		Usage()
		return

	}

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

	tmplfile, err := vimtmpl.GetTemplateFile(lang)
	if err != nil {
		fmt.Printf("unable to get template for %s: %s\n\n", lang, err)
		Usage()
		return
	}

	log.Debugf("template: %s", tmplfile)

	options := &Options{}
	options.Lang = lang

	cfg := config.NewTemplateConfig(options.Lang)
	cfg.Load()

	f := flag.NewFlagSet("prompt", flag.ExitOnError)

	f.BoolVar(&options.Verbose, "v", false, "Verbose logging")

	f.StringVar(&options.Company, "company", cfg.Company, "Company name")
	f.StringVar(&options.Company, "c", cfg.Company, "Company name")

	f.StringVar(&options.Copyright, "copyright", cfg.Copyright, "Copyright holder")

	f.StringVar(&options.Description, "description", "", "Script description")
	f.StringVar(&options.Description, "d", "", "Script description")

	f.StringVar(&options.License, "license", cfg.License, "License")
	f.StringVar(&options.License, "l", cfg.License, "License")

	f.StringVar(&options.MailAddress, "mailaddress", cfg.MailAdress, "mailaddress")
	f.StringVar(&options.MailAddress, "m", cfg.MailAdress, "mailaddress")

	f.StringVar(&options.ScriptName, "scriptname", "", "Script name")
	f.StringVar(&options.ScriptName, "s", "", "Script name")

	f.StringVar(&options.Title, "title", "", "Title (of e.g. python class)")
	f.StringVar(&options.Title, "t", "", "Title (of e.g. python class)")

	f.StringVar(&options.UserName, "username", cfg.Username, "Users full name")
	f.StringVar(&options.UserName, "u", cfg.Username, "Users full name")

	f.StringVar(&options.User, "user", cfg.User, "User account name")
	f.StringVar(&options.User, "U", cfg.User, "User account name")

	f.Parse(args)

	if len(options.ScriptName) != 0 {
		args := f.Args()
		if len(args) > 0 {
			options.ScriptName = args[0]
		}
	}

	if options.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	for _, indx := range cfg.Object.Sections() {
		fmt.Printf("section: %s\n", indx.Name())
	}
}
