package main

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/gobuffalo/packr"
)

type TemplateFields struct {
	Company     string
	Copyright   string
	Date        string
	Description string
	Mailaddress string
	License     string
	Scriptname  string
	Username    string
	Title       string
	User        string
	Year        string
}

type Dialect struct {
	Name      string
	Template  string
	Mode      int
	Extension string
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

var Box = packr.NewBox("./templates")

var dialectSets = []Dialect{
	Dialect{Name: "bash",
		Template:  Box.String("bash.template"),
		Extension: "",
		Mode:      0755},
	Dialect{Name: "bashlib",
		Template:  Box.String("bash.template"),
		Extension: ".sh",
		Mode:      0644},
	Dialect{Name: "go",
		Template:  Box.String("go.template"),
		Extension: ".go",
		Mode:      0644},
	Dialect{Name: "playbook",
		Template:  Box.String("playbook.template"),
		Extension: ".yml",
		Mode:      0644},
	Dialect{Name: "pythonlib",
		Template:  Box.String("pythonclass.template"),
		Extension: ".py",
		Mode:      0644},
	Dialect{Name: "python",
		Template:  Box.String("python.template"),
		Extension: "",
		Mode:      0755},
}

func NewTemplate() *TemplateFields {
	timest := time.Now()
	t := &TemplateFields{
		Date: fmt.Sprintf("%4d-%02d-%02d", timest.Year(), timest.Month(), timest.Day()),
		Year: fmt.Sprintf("%04d", timest.Year()),
	}
	return t
}

func (t *TemplateFields) Set(iname string, ivar string) {
	switch iname {
	case "company":
		t.Company = ivar
	case "copyright":
		t.Copyright = ivar
	case "date":
		t.Date = ivar
	case "description":
		t.Description = ivar
	case "mailaddress":
		t.Mailaddress = ivar
	case "license":
		t.License = ivar
	case "scriptname":
		t.Scriptname = ivar
	case "username":
		t.Username = ivar
	case "user":
		t.User = ivar
	case "year":
		t.Year = ivar
	}
}

// buildConfig contruct the text from the template definition and arguments.
func (t TemplateFields) Parse(templatestring string) string {
	tmpl, err := template.New("prompt").Parse(templatestring)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func getEnv(envvar string, defaultvar string) string {
	var retv string
	val, ok := os.LookupEnv(envvar)
	if !ok {
		retv = defaultvar
	} else {
		retv = val
	}
	return retv
}

func main() {

	filename := os.Args[0]
	filename = path.Base(filename)
	templateType := strings.Replace(filename, "vimtmpl_", "", -1)
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
	mailaddress = getEnv("MAILADDRESS", mailaddress)
	user = getEnv("USER", user)

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
	tmpl := NewTemplate()
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
	log.Info("provide filename " + filename)
	for _, dialect := range dialectSets {
		if dialect.Name == templateType {
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
			break
		}
	}
}
