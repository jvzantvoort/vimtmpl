package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/templates"
	log "github.com/sirupsen/logrus"
)

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

func Usage(lang string) string {
	if lang == "" {
		langs := templates.ListTemplateNames()
		lang = fmt.Sprintf("[%s]", strings.Join(langs, "|"))
	}
	return fmt.Sprintf("USAGE:\n\n\t%s %s [<filename>] <options>\n\n", os.Args[0], lang)
}

func WriteFile(tmpl *config.TemplateConfig, content string) error {
	log.Debugf("WriteFile, start")
	defer log.Debugf("WriteFile, end")

	if config.TargetExists(tmpl.ScriptName) {
		return fmt.Errorf("target already exists: %s", tmpl.ScriptName)
	}

	file, err := os.Create(tmpl.ScriptName)
	defer file.Close()

	obj := tmpl.GetItem(tmpl.Lang)
	log.Debugf("language found: %s", obj.Name)
	log.Debugf("mode: %o", obj.Mode)

	_, err = file.WriteString(content)
	if err != nil {
		log.Errorf("file.WriteString: error: %s", err)
		return err
	} else {
		log.Info("write file")
	}

	mode := os.FileMode(obj.Mode)
	if err := os.Chmod(tmpl.ScriptName, mode); err != nil {
		log.Error(err)
		return err
	} else {
		log.Debugf("mode: %d\n", obj.Mode)
	}
	return nil
}

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Print(Usage(""))
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

	tmplfile, err := templates.GetTemplateFile(lang)
	if err != nil {
		fmt.Printf("unable to get template for %s: %s\n\n", lang, err)
		Usage("")
		return
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

	if cfg.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Setup description
	cfg.Description = cfg.GetKeyAsString("description")

	// get template content
	templatestring, _ := templates.GetTemplateContent(cfg.Lang)
	text_template, err := template.New("tmpl").Parse(templatestring)
	if err != nil {
		log.Error(err)
		return
	}

	buf := new(bytes.Buffer)
	err = text_template.Execute(buf, *cfg)

	if err != nil {
		log.Error(err)
		return
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
	}

}
