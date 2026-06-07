package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/content"
	log "github.com/sirupsen/logrus"
)

func runInit() error {
	homedir := config.UserHomeDir()
	templatesDir := filepath.Join(homedir, ".templates.d")
	defaultDir := filepath.Join(templatesDir, "default")
	localDir := filepath.Join(templatesDir, "local")

	for _, dir := range []string{defaultDir, localDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Infof("directory ready: %s", dir)
	}

	entries, err := fs.ReadDir(content.Templates, ".")
	if err != nil {
		return fmt.Errorf("failed to read embedded templates: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		dest := filepath.Join(defaultDir, entry.Name())
		if config.TargetExists(dest) {
			log.Infof("skipping existing template: %s", entry.Name())
			continue
		}
		if err := installTemplate(entry.Name(), dest); err != nil {
			return fmt.Errorf("failed to install template %s: %w", entry.Name(), err)
		}
		log.Infof("installed template: %s", entry.Name())
	}

	cfgPath := filepath.Join(homedir, config.ConfigFilename)
	if config.TargetExists(cfgPath) {
		log.Infof("configuration already exists: %s", cfgPath)
		return nil
	}

	cfg := config.NewTemplateConfig("")
	if err := cfg.SaveTo(cfgPath); err != nil {
		return fmt.Errorf("failed to write configuration %s: %w", cfgPath, err)
	}
	log.Infof("created configuration: %s", cfgPath)
	return nil
}

func installTemplate(name, dest string) error {
	src, err := content.Templates.Open(name)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, src)
	return err
}
