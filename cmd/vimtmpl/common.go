package main

import (
	"flag"

	"github.com/jvzantvoort/vimtmpl"
	log "github.com/sirupsen/logrus"
)

type CommonSubCmd struct {
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
	SubName        string
}

func (c *CommonSubCmd) SetFlags(f *flag.FlagSet) {

	userconfig := vimtmpl.NewUserConfig()

	company := userconfig.Get("company", "-company-")
	copyright := userconfig.Get("copyright", "-copyright-")
	username := userconfig.Get("username", "-username-")
	user := userconfig.Get("user", "-user-")
	license := userconfig.Get("license", "-license-")
	mailaddress := userconfig.Get("mailaddress", "-mailaddress-")

	f.BoolVar(&c.Verbose, "v", false, "Verbose logging")

	f.StringVar(&c.Company, "company", company, "Company name")
	f.StringVar(&c.Company, "c", company, "Company name")

	f.StringVar(&c.Copyright, "copyright", copyright, "Copyright holder")

	f.StringVar(&c.Description, "description", "", "Script description")
	f.StringVar(&c.Description, "d", "", "Script description")

	f.StringVar(&c.License, "license", license, "License")
	f.StringVar(&c.License, "l", license, "License")

	f.StringVar(&c.MailAddress, "mailaddress", mailaddress, "mailaddress")
	f.StringVar(&c.MailAddress, "m", mailaddress, "mailaddress")

	f.StringVar(&c.ScriptName, "scriptname", "", "Script name")
	f.StringVar(&c.ScriptName, "s", "", "Script name")

	f.StringVar(&c.Title, "title", "", "Title (of e.g. python class)")
	f.StringVar(&c.Title, "t", "", "Title (of e.g. python class)")

	f.StringVar(&c.UserName, "username", username, "Users full name")
	f.StringVar(&c.UserName, "u", username, "Users full name")

	f.StringVar(&c.User, "user", user, "User account name")
	f.StringVar(&c.User, "U", user, "User account name")

}

func (*CommonSubCmd) Usage() string {
	msgstr, err := vimtmpl.Asset("messages/usage_create")
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}


