// Package main provides command-line argument parsing and usage for the vimtmpl tool.
package main

import (
	"fmt"
	"os"
	"path"

	//	"strconv"
	"strings"

	"github.com/spf13/pflag"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/templates"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Company     string
	Copyright   string
	Description string
	FullPath    string
	License     string
	MailAddress string
	Title       string
	UserName    string
	User        string
	Info        bool
	Verbose     bool
	Flags       []string
	Lang        string
}

func parseFlags() *Config {
	retv := &Config{}

	pflag.Usage = func() {
		printHelp()
		os.Exit(0)
	}

	pflag.StringVarP(&retv.MailAddress, "mailaddress", "m", "", "mailaddress")
	pflag.StringVarP(&retv.Company, "company", "c", "", "Company name")
	pflag.StringVarP(&retv.Copyright, "copyright", "C", "", "Copyright holder")
	pflag.StringVarP(&retv.License, "license", "l", "", "License")
	pflag.StringVarP(&retv.User, "user", "U", "", "User account name")
	pflag.StringVarP(&retv.UserName, "username", "u", "", "Users full name")

	pflag.BoolVarP(&retv.Verbose, "verbose", "v", false, "Enable verbose output")
	pflag.BoolVarP(&retv.Info, "info", "i", false, "Print Information about the template")

	pflag.StringVarP(&retv.FullPath, "scriptname", "s", "", "Script name")
	pflag.StringVarP(&retv.Title, "title", "t", "", "Title (of e.g. python class)")
	pflag.StringVarP(&retv.Description, "description", "d", "", "Script description")

	// pflag.StringToStringVarP(&retv.Flags, "flags", "f", nil, "Custom template switchesi as key=value pairs")
	pflag.StringArrayVarP(&retv.Flags, "flags", "f", nil, "Custom template switchesi as key=value pairs")

	pflag.Parse()

	// Positional arguments
	args := pflag.Args()
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, Usage(""))
		fmt.Fprintln(os.Stderr, "Error: a template name is required as the first argument.")
		os.Exit(1)
	}

	if len(args) < 2 {
		fmt.Fprint(os.Stderr, Usage(args[0]))
		fmt.Fprintln(os.Stderr, "Error: an output filename is required as the second argument.")
		os.Exit(1)
	}

	retv.Lang = args[0]
	retv.FullPath = args[1]

	if retv.FullPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --output is required.")
		pflag.Usage()
		os.Exit(1)
	}

	return retv
}

// Usage returns a formatted usage string for the given language.
// When lang is empty all available templates are listed, along with the named subcommands.
func Usage(lang string) string {
	if lang == "" {
		langs := templates.ListTemplateNames()
		lang = fmt.Sprintf("[%s]", strings.Join(langs, "|"))
		return fmt.Sprintf("USAGE:\n\n\t%[1]s init\n\t%[1]s help\n\t%[1]s %[2]s <filename> <options>\n\n",
			os.Args[0], lang)
	}
	return fmt.Sprintf("USAGE:\n\n\t%s %s <filename> <options>\n\n", os.Args[0], lang)
}

// ArgParse parses command-line arguments and returns a TemplateConfig or an error.
func ArgParse() (*config.TemplateConfig, error) {

	// Start logging
	log.Debugf("ArgParse, start")
	defer log.Debugf("ArgParse, end")

	flags := parseFlags()

	if flags.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	cfg := config.NewTemplateConfig(flags.Lang)
	cfg.Load()

	cfg.Info = flags.Info

	if flags.MailAddress != "" {
		cfg.MailAddress = flags.MailAddress
	}

	if flags.Company != "" {
		cfg.Company = flags.Company
	}

	if flags.Copyright != "" {
		cfg.Copyright = flags.Copyright
	}

	if flags.License != "" {
		cfg.License = flags.License
	}

	if flags.User != "" {
		cfg.User = flags.User
	}

	if flags.UserName != "" {
		cfg.UserName = flags.UserName
	}

	if flags.FullPath != "" {
		cfg.FullPath = flags.FullPath
	}

	if flags.Title != "" {
		cfg.Title = flags.Title
	}

	if flags.Description != "" {
		cfg.Description = flags.Description
	}

	// switches := make(map[string]bool)

	for _, indx := range flags.Flags {
		if strings.Contains(indx, ",") {
			for _, sindx := range strings.Split(indx, ",") {
				cfg.Flags[sindx] = true

			}

		} else {
			cfg.Flags[indx] = true
		}
	}

	cfg.ScriptName = path.Base(cfg.FullPath)

	_, err := templates.GetTemplateFile(cfg.Lang)
	if err != nil {
		return &config.TemplateConfig{}, fmt.Errorf("unable to get template for %s: %s", cfg.Lang, err)
	}

	return cfg, nil
}
