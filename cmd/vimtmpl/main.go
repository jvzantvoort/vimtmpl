package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&ConfigSubCmd{}, "config")
	subcommands.Register(&BashSubCmd{}, "template")
	subcommands.Register(&BashLibSubCmd{}, "template")
	subcommands.Register(&GoSubCmd{}, "template")
	subcommands.Register(&PlaybookSubCmd{}, "template")
	subcommands.Register(&PythonSubCmd{}, "template")
	subcommands.Register(&PythonLibSubCmd{}, "template")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))

}
