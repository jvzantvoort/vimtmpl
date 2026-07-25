// Package main is the entry point for vimtmplx, the graphical (terminal UI)
// version of vimtmpl.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	version  = "dev"
	revision = "unknown"
)

func main() {
	opts := parseFlags()

	if opts.Cwd != "" {
		if err := os.Chdir(opts.Cwd); err != nil {
			fmt.Fprintf(os.Stderr, "vimtmplx: cannot change to directory %s: %s\n", opts.Cwd, err)
			os.Exit(1)
		}
	}

	p := tea.NewProgram(newModel(opts.Filename), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "vimtmplx: %s\n", err)
		os.Exit(1)
	}
}
