// Package commands provides utilities for handling and manipulating PATH-like environment variables.
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

// Path represents a PATH-like environment variable, including its type, home directory, and a list of directories.
type Path struct {
	Type        string
	Home        string
	Directories []string
}

// Prefix returns a formatted string with the caller function name and the Path type.
func (p Path) Prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("%s[%s]", elements[len(elements)-1], p.Type)
}

// HavePath checks if the given directory exists in the Path's Directories slice.
func (p Path) HavePath(inputdir string) bool {
	for _, element := range p.Directories {
		if element == inputdir {
			return true
		}
	}
	return false
}

// AppendPath appends a directory to the Path's Directories slice if it is not already present.
// It expands the input directory and checks its existence.
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

// PrependPath prepends a directory to the Path's Directories slice if it is not already present.
// It expands the input directory.
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

// Import splits a PATH-like string and appends each directory to the Path's Directories slice.
func (p *Path) Import(path string) {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	log.Debugf("%s: path=%s", log_prefix, path)

	for _, dirn := range strings.Split(path, ":") {
		err := p.AppendPath(dirn)
		if err != nil {
			log.Errorf("Append failed: %v", err)
		}
	}
}

// IsEmpty returns true if the Path's Directories slice is empty.
func (p Path) IsEmpty() bool {
	if len(p.Directories) == 0 {
		return true
	} else {
		return false
	}
}

// ReturnExport returns a shell export command for the Path's type and directories.
func (p Path) ReturnExport() string {
	return fmt.Sprintf("export %s=\"%s\"", p.Type, strings.Join(p.Directories, ":"))

}

// targetExists returns true if the target path exists in the filesystem.
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

// Lookup searches for the target executable in the Path's Directories.
// Returns the full path if found, or an error if not found.
func (p Path) Lookup(target string) (string, error) {
	log_prefix := p.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	var retv string
	var err error
	err = fmt.Errorf("command %s not found", target)

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

// LookupMulti searches for multiple target executables in the Path's Directories.
// Returns the first found path, or an error if none are found.
func (p Path) LookupMulti(targets ...string) (string, error) {
	for _, target := range targets {
		if result, err := p.Lookup(target); err != nil {
			return result, nil
		}
	}
	return "", fmt.Errorf("targets not found")
}

// MapGetPlatform retrieves the platform-specific command from the provided map.
// Returns the command for the current OS or the default if available.
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

// LookupPlatform looks up the platform-specific command in the Path's Directories using a map of platform commands.
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

	return "", fmt.Errorf("target not found")
}

// NewPath creates a new Path instance for the given environment variable name.
// It initializes the Path type, home directory, and imports the environment variable value.
func NewPath(pathname string) *Path {
	retv := &Path{}
	retv.Type = pathname
	retv.Home, _ = homedir.Dir()
	retv.Import(os.Getenv(retv.Type))
	return retv
}
