package config

import (
	"path/filepath"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const ConfigFilename string = ".template.cfg"

type TemplateItem struct {
	Name        string
	Description string
}

type TemplateConfig struct {
	Filepath   string
	MailAdress string
	Company    string
	Copyright  string
	License    string
	User       string
	Username   string
	Lang       string
	Homedir    string
	Object     *ini.File
	Items      []*TemplateItem
}

func NewTemplateConfig(lang string) *TemplateConfig {
	log.Debugf("TemplateConfig, start")
	defer log.Debugf("TemplateConfig, end")

	retv := &TemplateConfig{}
	retv.Lang = lang

	retv.Company = "company"
	retv.Copyright = "copyright"
	retv.License = "license"
	retv.MailAdress = "mailaddress"
	retv.Username = "username"

	retv.User = UserName()
	retv.Homedir = UserHomeDir()

	retv.Filepath = filepath.Join(retv.Homedir, ConfigFilename)

	return retv

}

func (tc TemplateConfig) GetKeyAsString(keyname string) string {
	log.Debugf("GetKeyAsString: %s/%s, start", tc.Lang, keyname)
	defer log.Debugf("GetKeyAsString: %s/%s, end", tc.Lang, keyname)

	result, err := tc.Object.Section(tc.Lang).GetKey(keyname)
	if err == nil {
		return result.String()
	}

	result, err = tc.Object.Section(ini.DefaultSection).GetKey(keyname)
	if err == nil {
		return result.String()
	}

	log.Errorf("Error: %s", err)
	return ""
	// return tc.Object.Section(tc.Lang).GetKey(keyname).String()
}

func (tc TemplateConfig) LoadFile() *ini.File {
	cfg, err := ini.Load(tc.Filepath)
	if err != nil {
		log.Errorf("Failed to load %s", ConfigFilename)
		cfg = ini.Empty()
	}
	return cfg
}

func (tc *TemplateConfig) Load() {

	tc.Object = tc.LoadFile()

	tc.Company = tc.GetKeyAsString("company")
	tc.Copyright = tc.GetKeyAsString("copyright")
	tc.License = tc.GetKeyAsString("license")
	tc.MailAdress = tc.GetKeyAsString("mailaddress")
	tc.Username = tc.GetKeyAsString("username")
	tc.User = tc.GetKeyAsString("user")

	for _, indx := range tc.Object.Sections() {
		ti := &TemplateItem{}

		ti.Name = indx.Name()
		ti.Description = tc.GetKeyAsString("description")

		tc.Items = append(tc.Items, ti)
	}
}

func (tc TemplateConfig) SaveTo(filename string) error {

	cfg := ini.Empty()
	cfg.Section("").Key("company").SetValue(tc.Company)
	cfg.Section("").Key("copyright").SetValue(tc.Copyright)
	cfg.Section("").Key("license").SetValue(tc.License)
	cfg.Section("").Key("mailaddress").SetValue(tc.MailAdress)
	cfg.Section("").Key("username").SetValue(tc.Username)
	cfg.Section("").Key("user").SetValue(tc.User)

	for _, obj := range tc.Items {
		sec, err := cfg.NewSection(obj.Name)
		if err != nil {
			log.Error(err)
		}
		sec.Key("description").SetValue(obj.Description)
	}

	return cfg.SaveTo(filename)

}
