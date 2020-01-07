package main

import (
	"flag"
	"github.com/jvzantvoort/vimtmpl"
	"strings"
	"os"
	log "github.com/sirupsen/logrus"
)

const (
	templateType = "go"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}


func main() {
	log.Info("Start as " + templateType)

	company := "companyname"
	copyright := "copyright holder"
	description := ""
	license := "undefined"
	mailaddress := "undefined"
	scriptname := "undefined"
	title := "Title"
	username := "undefined"
	user := "undefined"

	// try environment variables
	// mailaddress = vimtmpl.GetEnv("MAILADDRESS", mailaddress)
	// user = GetEnv("USER", user)

	// try command line options
	flags := flag.NewFlagSet("vimtmpl", flag.ExitOnError)

	flags.StringVar(&company, "company", company, "Company name")
	flags.StringVar(&company, "c", company, "Company name")

	flags.StringVar(&copyright, "copyright", copyright, "Copyright holder")

	flags.StringVar(&description, "description", description, "Script description")
	flags.StringVar(&description, "d", description, "Script description")

	flags.StringVar(&license, "license", license, "License")
	flags.StringVar(&license, "l", license, "License")

	flags.StringVar(&mailaddress, "mailaddress", mailaddress, "mailaddress")
	flags.StringVar(&mailaddress, "m", mailaddress, "mailaddress")

	flags.StringVar(&scriptname, "scriptname", scriptname, "Script name")
	flags.StringVar(&scriptname, "s", scriptname, "Script name")

	flags.StringVar(&title, "title", title, "Title (of e.g. python class)")
	flags.StringVar(&title, "t", title, "Title (of e.g. python class)")

	flags.StringVar(&username, "username", username, "Users full name")
	flags.StringVar(&username, "u", username, "Users full name")

	flags.StringVar(&user, "user", user, "User account name")
	flags.StringVar(&user, "U", user, "User account name")

	flags.Parse(os.Args[1:])

	// create the template with the options
	tmpl := vimtmpl.NewTemplate()
	tmpl.Set("company", company)
	tmpl.Set("copyright", copyright)
	tmpl.Set("description", description)
	tmpl.Set("license", license)
	tmpl.Set("mailaddress", mailaddress)
	tmpl.Set("scriptname", scriptname)
	tmpl.Set("title", title)
	tmpl.Set("username", username)
	tmpl.Set("user", user)

	var scriptcontent string
	dialect := tmpl.GetDialect(templateType)
	scriptcontent = tmpl.Parse(dialect.Template)
	if len(dialect.Extension) != 0 {
		log.Debug("extenstion: " + dialect.Extension)
		if !strings.HasSuffix(scriptname, dialect.Extension) {
			log.Warning("did not find extension")
			scriptname = scriptname + dialect.Extension
		} else {
			log.Info("found extension")
		}
	}
	log.Info("debug: scriptname %s", scriptname)
	file, err := os.Create(scriptname)
	_, err = file.WriteString(scriptcontent)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	mode := os.FileMode(dialect.Mode)
	if err := os.Chmod(scriptname, mode); err != nil {
		log.Fatal(err)
	}
	return
}
