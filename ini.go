package vimtmpl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ini/ini"
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

func (uc UserConfig) Set(parameter string, value string) {
	sectionname := "user"
	cfg, err := ini.Load(uc.Configfile())
	if err != nil {
		cfg = ini.Empty()
	}
	if _, err := cfg.GetSection(sectionname); err != nil {
		cfg.NewSection(sectionname)
	}
	cfg.Section(sectionname).Key(parameter).SetValue(value)
	cfg.SaveTo(uc.Configfile())
}

func (uc UserConfig) Get(parameter string, defaultstr string) string {
	sectionname := "user"
	cfg, err := ini.Load(uc.Configfile())
	if err != nil {
		cfg = ini.Empty()
	}
	retv := cfg.Section(sectionname).Key(parameter).String()
	if retv != "" {
		retv = defaultstr
	}
	return retv
}
