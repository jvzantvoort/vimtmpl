package templates

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/jvzantvoort/vimtmpl/config"

	log "github.com/sirupsen/logrus"
)

func ListTemplateFile(tmpldir string) (map[string]string, error) {
	retv := make(map[string]string)

	files, err := ioutil.ReadDir(tmpldir)
	if err != nil {
		log.Fatal(err)
		return retv, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		fileExtension := filepath.Ext(name)
		if fileExtension != ".gtmpl" {
			continue
		}

		fpath := path.Join(tmpldir, name)
		fpath, _ = filepath.Abs(fpath)

		name = name[:len(name)-len(".gtmpl")]

		retv[name] = fpath
	}
	return retv, nil
}

func ListTemplates() (map[string]string, error) {
	retv := make(map[string]string)

	homedir := config.UserHomeDir()

	templatesdir := path.Join(homedir, ".templates.d")

	if !config.TargetExists(templatesdir) {
		return retv, fmt.Errorf("%s does not exist", templatesdir)
	}

	defaulttemplates := path.Join(templatesdir, "default")
	localtemplates := path.Join(templatesdir, "local")

	if config.TargetExists(defaulttemplates) {
		data, err := ListTemplateFile(defaulttemplates)
		if err != nil {
			log.Error(err)
		}
		for key, element := range data {
			retv[key] = element
		}
	} else {
		log.Warnf("Directory %s does not exist", defaulttemplates)
	}

	if config.TargetExists(localtemplates) {
		data, err := ListTemplateFile(localtemplates)
		if err != nil {
			log.Error(err)
		}
		for key, element := range data {
			retv[key] = element
		}
	} else {
		log.Warnf("Directory %s does not exist", localtemplates)
	}

	return retv, nil

}

func ListTemplateNames() []string {
	retv := []string{}
	data, err := ListTemplates()
	if err != nil {
		return retv
	}
	for keyn, _ := range data {
		retv = append(retv, keyn)
	}
	return retv
}

func GetTemplateFile(lang string) (string, error) {
	data, err := ListTemplates()
	if err != nil {
		return "", err
	}
	if rev, ok := data[lang]; ok {
		return rev, nil
	}
	return "", fmt.Errorf("Language not found: %s", lang)
}

func GetTemplateContent(lang string) (string, error) {
	target, err := GetTemplateFile(lang)
	retv := ""
	if err != nil {
		return string(retv), err
	}

	content, err := ioutil.ReadFile(target)
	retv = string(content)

	if err != nil {
		return retv, err
	}
	return retv, nil

}
