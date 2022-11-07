// PATH type handling.
package commands

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

type Path struct {
	Type        string
	Home        string
	Directories []string
}

func (p Path) Prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("%s[%s]", elements[len(elements)-1], p.Type)
}

func (p Path) HavePath(inputdir string) bool {
	for _, element := range p.Directories {
		if element == inputdir {
			return true
		}
	}
	return false
}

// AppendPath append a path to the list of Directories
func (p *Path) AppendPath(inputdir string) error {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	log.Debugf("%s: inputdir=%s", log_prefix, inputdir)

	if len(inputdir) == 0 {
		return nil
	}

	fullpath, err := homedir.Expand(inputdir)
	if err != nil {
		log.Errorf("error %s", err)
		return err
	}

	_, err = os.Stat(fullpath)
	if err != nil {
		return err
	}

	if !p.HavePath(fullpath) {
		p.Directories = append(p.Directories, fullpath)
	}

	return nil
}

func (p *Path) PrependPath(inputdir string) error {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	log.Debugf("%s: inputdir=%s", log_prefix, inputdir)

	if len(inputdir) == 0 {
		return nil
	}

	fullpath, err := homedir.Expand(inputdir)
	if err != nil {
		log.Errorf("error %s", err)
		return err
	}

	if !p.HavePath(fullpath) {
		p.Directories = append([]string{fullpath}, p.Directories...)
	}

	return nil
}

func (p *Path) Import(path string) {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	log.Debugf("%s: path=%s", log_prefix, path)

	for _, dirn := range strings.Split(path, ":") {
		p.AppendPath(dirn)
	}
}

func (p Path) IsEmpty() bool {
	if len(p.Directories) == 0 {
		return true
	} else {
		return false
	}
}

func (p Path) ReturnExport() string {
	return fmt.Sprintf("export %s=\"%s\"", p.Type, strings.Join(p.Directories, ":"))

}

// targetExists return true if target exists
func (p Path) targetExists(targetpath string) bool {
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (p Path) Lookup(target string) (string, error) {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	var retv string
	var err error
	err = fmt.Errorf("Command %s not found", target)

	for _, dirname := range p.Directories {
		fullpath := path.Join(dirname, target)
		if p.targetExists(fullpath) {
			retv = fullpath
			err = nil
			break
		}
	}
	return retv, err
}

func (p Path) LookupMulti(targets ...string) (string, error) {
	for _, target := range targets {
		if result, err := p.Lookup(target); err != nil {
			return result, nil
		}
	}
	return "", fmt.Errorf("Targets not found")
}

func (p Path) MapGetPlatform(pathmap map[string]string) (string, error) {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	goos := runtime.GOOS
	log.Debugf("%s: os=%s", log_prefix, goos)
	if target, ok := pathmap[goos]; ok {
		log.Debugf("%s: found key: %s -> %s", log_prefix, goos, target)
		return target, nil
	}
	if target, ok := pathmap["default"]; ok {
		log.Debugf("%s: found key: default -> %s", log_prefix, target)
		return target, nil
	}

	return "", fmt.Errorf("%s: map keys not found", log_prefix)
}

// LookupPlatform lookup paths based on platform
func (p Path) LookupPlatform(pathmap map[string]string) (string, error) {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	commandname, err := p.MapGetPlatform(pathmap)
	if err != nil {
		log.Errorf("Err: %s", err)
		return "", err
	}

	if result, err := p.Lookup(commandname); err == nil {
		log.Debugf("%s: found: %s -> %s", log_prefix, commandname, result)
		return result, nil
	}
	log.Errorf("%s: cannot find %s in path", log_prefix, commandname)

	return "", fmt.Errorf("Target not found")
}

func NewPath(pathname string) *Path {
	retv := &Path{}
	retv.Type = pathname
	retv.Home, _ = homedir.Dir()
	retv.Import(os.Getenv(retv.Type))
	return retv
}
