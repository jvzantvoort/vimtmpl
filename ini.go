package vimtmpl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const configfilename string = ".vimtmplrc"

type UserConfig struct {
	mailaddress string
	company     string
}

func (uc UserConfig) UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func (uc UserConfig) Configfile() string {
	return filepath.Join(uc.UserHomeDir(), configfilename)
}

func (uc UserConfig) Load() map[string]string {
	retv := make(map[string]string)
	retv["mailaddress"] = "mailaddress"
	retv["company"] = "company"
	retv["license"] = "license"
	retv["copyright"] = "John Doe"
	retv["user"] = "John Doe"
	retv["username"] = "John Doe"

	cfg, err := ini.Load(uc.Configfile())
	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		return retv
	}
	retv["mailaddress"] = cfg.Section("user").Key("mailaddress").String()
	retv["company"] = cfg.Section("user").Key("company").String()

	return retv
}

func (uc UserConfig) LoadFile() *ini.File {
	cfg, err := ini.Load(uc.Configfile())
	if err != nil {
		log.Errorf("Failed to load %s", configfilename)
		cfg = ini.Empty()
	}
	return cfg
}

func (uc UserConfig) GetSection(sectionname string) ini.Section {
	data := uc.LoadFile()

	if _, err := data.GetSection(sectionname); err != nil {
		log.Errorf("Failed to load %s: %s", sectionname, err)
		data.NewSection(sectionname)
	}
	sect, _ := data.GetSection(sectionname)
	return *sect
}

func (uc UserConfig) Set(parameter string, value string) {
	log.Debugf("Set %s: %s, start", parameter, value)
	defer log.Debugf("Set %s: %s, end", parameter, value)

	sectionname := "user"
	cfg := uc.LoadFile()
	if _, err := cfg.GetSection(sectionname); err != nil {
		log.Errorf("Failed to load %s: %s", sectionname, err)
		cfg.NewSection(sectionname)
	}
	log.Debugf("Set %s/%s = %s", sectionname, parameter, value)
	cfg.Section(sectionname).Key(parameter).SetValue(value)
	cfg.SaveTo(uc.Configfile())
}

func (uc UserConfig) Get(parameter string, defaultstr string) string {

	log.Debugf("Get %s[%s], start", parameter, defaultstr)
	defer log.Debugf("Get %s[%s], end", parameter, defaultstr)

	sectionname := "user"

	configfilename := uc.Configfile()
	cfg, err := ini.Load(configfilename)
	if err != nil {
		log.Errorf("Failed to load %s", configfilename)
		cfg = ini.Empty()
	}

	for _, indx := range cfg.Sections() {
		log.Debugf("%s/%s", sectionname, indx.Name())
	}

	retv := cfg.Section(sectionname).Key(parameter).String()
	if retv == "" {
		log.Errorf("Get failed for %s in %s", parameter, sectionname)
		retv = defaultstr
	}
	log.Debugf("Get %s/%s = %s", sectionname, parameter, retv)
	return retv
}

func NewUserConfig() *UserConfig {
	log.Debugf("UserConfig, start")
	defer log.Debugf("UserConfig, end")

	retv := &UserConfig{}
	retv.Load()

	return retv


}
