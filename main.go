package vimtmpl

import (
	"bytes"
	"fmt"
	"os"
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

func GetEnv(envvar string, defaultvar string) string {
	var retv string
	val, ok := os.LookupEnv(envvar)
	if !ok {
		retv = defaultvar
	} else {
		retv = val
	}
	return retv
}

func (t TemplateFields) GetDialect(ttype string) Dialect {
	for _, dialect := range dialectSets {
		if dialect.Name == ttype {
			return dialect
		}
	}
	return Dialect{}
}
