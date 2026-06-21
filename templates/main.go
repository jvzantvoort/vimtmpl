// Package templates provides functions for listing and managing template files.
package templates

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jvzantvoort/vimtmpl/config"
	"gopkg.in/ini.v1"

	log "github.com/sirupsen/logrus"
)

// ListTemplateFile returns a map of template names to their absolute file paths in the given directory.
func ListTemplateFile(tmpldir string) (map[string]string, error) {
	retv := make(map[string]string)

	files, err := os.ReadDir(tmpldir)
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
	for keyn := range data {
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
	return "", fmt.Errorf("language not found: %s", lang)
}

func ParseLang(lang string) (*ini.File, string, error) {

	const marker = "---\n"

	// load file
	// -------------------------------------
	target, err := GetTemplateFile(lang)
	if err != nil {
		return ini.Empty(), "", err
	}

	contentb, err := os.ReadFile(target)
	if err != nil {
		return ini.Empty(), "", err
	}
	content := string(contentb)
	// -------------------------------------

	if !strings.HasPrefix(content, marker) {
		return ini.Empty(), content, nil
	}

	// Add header config
	// -------------------------------------
	rest := content[len(marker):]

	idx := strings.Index(rest, marker)
	if idx == -1 {
		// Opening marker exists but no closing marker.
		return ini.Empty(), content, nil
	}

	headerBlob := rest[:idx]
	body := rest[idx+len(marker):]

	header, err := ini.LoadSources(
		ini.LoadOptions{
			IgnoreInlineComment: true,
		},
		[]byte(headerBlob),
	)
	if err != nil {
		return nil, "", err
	}

	return header, body, nil

}

func splitDocument(content string) (header string, body string) {
	const marker = "---\n"

	if !strings.HasPrefix(content, marker) {
		return "", content
	}

	rest := content[len(marker):]

	idx := strings.Index(rest, marker)
	if idx == -1 {
		// Opening marker exists but no closing marker.
		return "", content
	}

	header = rest[:idx]
	body = rest[idx+len(marker):]

	return header, body
}

func GetTemplateContent(lang string) (string, error) {
	target, err := GetTemplateFile(lang)
	retv := ""
	if err != nil {
		return string(retv), err
	}

	content, err := os.ReadFile(target)
	_, retv = splitDocument(string(content))

	if err != nil {
		return retv, err
	}
	return retv, nil

}

func GetTemplateConfig(lang string) (string, error) {
	target, err := GetTemplateFile(lang)
	retv := ""
	if err != nil {
		return string(retv), err
	}

	content, err := os.ReadFile(target)
	retv, _ = splitDocument(string(content))

	if err != nil {
		return retv, err
	}
	return retv, nil

}
