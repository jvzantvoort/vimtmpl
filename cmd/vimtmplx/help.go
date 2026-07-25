package main

import (
	"fmt"
	"os"
)

const helpTemplate = `NAME
    vimtmplx - graphical (terminal UI) version of vimtmpl

SYNOPSIS
    %[1]s [-c|--cwd <dir>] [-f|--filename <name>]
    %[1]s -h | --help

DESCRIPTION
    vimtmplx is a terminal UI front-end for vimtmpl. It lets you pick the
    output filename, the template to use, and the template's switches
    interactively, then writes the file exactly as "vimtmpl" would.

    Templates and configuration are shared with vimtmpl. Run
    "vimtmpl init" once to install the bundled default templates and
    create a skeleton ~/.template.cfg.

OPTIONS
    -c, --cwd <dir>
            Directory to start from (default: current working directory).
            Relative filenames are resolved against this directory.

    -f, --filename <name>
            Pre-fill the filename field with <name>.

    -h, --help
            Print this help text.

KEYS
    tab / shift+tab   move between fields
    up / down         move within the focused field
    left / right      change the selected template
    space / enter     toggle the highlighted switch
    ctrl+s            generate the file
    esc / ctrl+c      quit without writing

`

func printHelp() {
	fmt.Printf(helpTemplate, os.Args[0])
}
