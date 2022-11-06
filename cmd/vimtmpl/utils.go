package main

import (
	"os"

	"github.com/jvzantvoort/vimtmpl"
	log "github.com/sirupsen/logrus"
)

func WriteFile(filename, content string, dialect vimtmpl.Dialect) error {

	file, err := os.Create(filename)
	defer file.Close()

	_, err = file.WriteString(content)

	if err != nil {
		return err
	}

	mode := os.FileMode(dialect.Mode)
	if err := os.Chmod(filename, mode); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
