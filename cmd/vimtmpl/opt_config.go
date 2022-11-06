package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

type ConfigSubCmd struct {
	Verbose bool
}

func (*ConfigSubCmd) Name() string {
	return "config"
}

func (c *ConfigSubCmd) Synopsis() string {
	return fmt.Sprintf("Update the configuration")
}

func (*ConfigSubCmd) Usage() string {
	return "oke"
}

func (c *ConfigSubCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.Verbose, "v", false, "Verbose logging")
}

func (c *ConfigSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	return subcommands.ExitSuccess
}
