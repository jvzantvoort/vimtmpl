package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/jvzantvoort/vimtmpl"
	log "github.com/sirupsen/logrus"
)

type BashLibSubCmd struct {
	CommonSubCmd
}

func (*BashLibSubCmd) Name() string {
	return "bashlib"
}

func (c *BashLibSubCmd) Synopsis() string {
	return fmt.Sprintf("Create a new %s script", c.Name())
}

func (c *BashLibSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	c.SubName = c.Name()

	var scriptcontent string

	if c.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")
	defer log.Debugln("End")

	// create the template with the options
	tmpl := vimtmpl.NewTemplate()
	tmpl.Company = c.Company
	tmpl.Copyright = c.Copyright
	tmpl.Description = c.Description
	tmpl.Mailaddress = c.MailAddress
	tmpl.License = c.License
	tmpl.Scriptname = c.ScriptName
	tmpl.Username = c.UserName
	tmpl.Title = c.Title
	tmpl.User = c.User

	dialect := tmpl.GetDialect(c.SubName)

	scriptcontent = tmpl.Parse(dialect.Template)

	if len(c.ScriptName) == 0 {
		fmt.Print(scriptcontent)
		return subcommands.ExitSuccess
	}

	c.ScriptName = dialect.AddExtension(c.ScriptName)

	err := WriteFile(c.ScriptName, scriptcontent, dialect)

	if err != nil {
		return subcommands.ExitFailure
	} else {
		return subcommands.ExitSuccess
	}
}
