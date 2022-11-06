package config

import (
	"fmt"
	"path/filepath"

	"time"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const ConfigFilename string = ".template.cfg"

type TemplateItem struct {
	Name        string
	Description string
	Mode        int
	Extension   string
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

	Date       string
	Year       string

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

	// add timestamps
	timest := time.Now()
	retv.Date = fmt.Sprintf("%4d-%02d-%02d", timest.Year(), timest.Month(), timest.Day())
	retv.Year = fmt.Sprintf("%04d", timest.Year())

	// add local parameters
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

	if keyname == "description" {
		return ""
	}

	log.Errorf("Error: %s", err)
	return ""
}

func (tc TemplateConfig) GetKeyAsInt(keyname string) int {
	log.Debugf("GetKeyAsInt: %s/%s, start", tc.Lang, keyname)
	defer log.Debugf("GetKeyAsInt: %s/%s, end", tc.Lang, keyname)

	result, err := tc.Object.Section(tc.Lang).GetKey(keyname)
	if err == nil {
		intval, ok := result.Int()
		if ok == nil {
			return intval
		}
	}

	result, err = tc.Object.Section(ini.DefaultSection).GetKey(keyname)
	if err == nil {
		intval, ok := result.Int()
		if ok == nil {
			return intval
		}
	}

	log.Errorf("Error: %s", err)
	return 0
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
		ti.Extension = tc.GetKeyAsString("extension")
		ti.Mode = tc.GetKeyAsInt("mode")

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
