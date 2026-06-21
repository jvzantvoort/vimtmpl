// Package commands provides utility functions for command execution and application control.
package commands

import (
	"errors"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Which searches for an executable named command in the system's PATH.
// It returns the full path to the executable if found, or an error if not found.
func Which(command string) (string, error) {
	for _, dirname := range strings.Split(os.Getenv("PATH"), ":") {
		fpath := path.Join(dirname, command)
		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			continue
		} else {
			return fpath, nil
		}
	}
	return command, errors.New("unable to find command " + command)
}

// ExitApplication logs the end of the application and exits with status code 0.
// The output parameter is currently unused.
func ExitApplication(output string) {
	log.Info("End")
	os.Exit(0)
}
