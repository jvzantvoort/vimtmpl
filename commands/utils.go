package commands

import (
	"errors"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

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

func ExitApplication(output string) {
	log.Info("End")
	os.Exit(0)
}
