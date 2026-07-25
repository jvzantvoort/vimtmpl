package main

import (
	"os"

	"github.com/spf13/pflag"
)

// Options holds the command-line options accepted by vimtmplx.
type Options struct {
	Cwd      string
	Filename string
}

func parseFlags() *Options {
	opts := &Options{}

	pflag.Usage = func() {
		printHelp()
		os.Exit(0)
	}

	pflag.StringVarP(&opts.Cwd, "cwd", "c", "", "directory to start from (default is the current working directory)")
	pflag.StringVarP(&opts.Filename, "filename", "f", "", "pre-fill the filename field")

	pflag.Parse()

	return opts
}
