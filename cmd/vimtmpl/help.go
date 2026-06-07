package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/templates"
)

const helpTemplate = `NAME
    vimtmpl - create files from predefined templates

SYNOPSIS
    %[1]s init
    %[1]s help
    %[1]s <template> <filename> [options]

DESCRIPTION
    vimtmpl generates files from templates stored in ~/.templates.d/.
    Each template is a Go text/template file (.gtmpl) that expands
    variables such as author, company, and license. Values are read
    from ~/.template.cfg and can be overridden per invocation via flags.

    Run "%[1]s init" once to install the bundled default templates and
    create a skeleton configuration file.

SUBCOMMANDS
    init    Create ~/.templates.d/{default,local}/ and install the bundled
            default templates into the default directory. Also writes a
            skeleton ~/.template.cfg if none already exists.
            Safe to run more than once — existing files are never overwritten.

    help    Print this help text.

OPTIONS
    The following flags override values from the configuration file for a
    single invocation. Each entry shows the flag, the corresponding
    configuration file key, and the template variable it populates.

    -m, --mailaddress <addr>
            Author email address.
            Config key: mailaddress    Template variable: {{.MailAddress}}

    -c, --company <name>
            Company name.
            Config key: company        Template variable: {{.Company}}

    -C, --copyright <holder>
            Copyright holder (defaults to the company name when unset).
            Config key: copyright      Template variable: {{.Copyright}}

    -l, --license <id>
            License identifier, e.g. MIT, Apache-2.0, GPL-3.0.
            Config key: license        Template variable: {{.License}}

    -U, --user <account>
            Login / account name of the author.
            Config key: user           Template variable: {{.User}}

    -u, --username <name>
            Author full name.
            Config key: username       Template variable: {{.UserName}}

    -t, --title <title>
            Title used by templates that need one, e.g. a Python class name.
            No config key.             Template variable: {{.Title}}

    -d, --description <text>
            Short description of the file being created.
            No config key.             Template variable: {{.Description}}

    -v, --verbose
            Enable debug logging.

TEMPLATE VARIABLES
    All variables below are available inside .gtmpl files.  Values are
    resolved in this order: command-line flag > configuration file section
    for the chosen template > [DEFAULT] section > built-in default.

    {{.ScriptName}}    Basename of <filename> (derived automatically).
    {{.FullPath}}      Full path of the output file (the <filename> argument).
    {{.Lang}}          Name of the template used (the <template> argument).
    {{.Date}}          Creation date in YYYY-MM-DD format (auto-populated).
    {{.Year}}          Current year in YYYY format (auto-populated).
    {{.User}}          Login name of the author.        flag -U / key user
    {{.UserName}}      Full name of the author.         flag -u / key username
    {{.MailAddress}}   Author email address.            flag -m / key mailaddress
    {{.Company}}       Company name.                    flag -c / key company
    {{.Copyright}}     Copyright holder.                flag -C / key copyright
    {{.License}}       License identifier.              flag -l / key license
    {{.Title}}         Title (e.g. class name).         flag -t
    {{.Description}}   Short description of the file.  flag -d

CONFIGURATION FILE
    Location: %[2]s

    The file uses INI format.  A [DEFAULT] section provides values for all
    templates; per-template sections override individual keys.

    Skeleton (written by "vimtmpl init"):

        [DEFAULT]
        company     = Example Corp
        copyright   = Example Corp
        license     = MIT
        mailaddress = user@example.com
        username    = Full Name
        user        = loginname

    Per-template override example:

        [bash]
        license = GPL-2.0

AVAILABLE TEMPLATES
%[3]s
FILES
    %[2]s
            User configuration file.

    ~/.templates.d/default/
            Default templates, installed by "vimtmpl init".

    ~/.templates.d/local/
            Local templates; files here take precedence over default.

EXAMPLES
    First-time setup:
        vimtmpl init

    Create a bash script with a description:
        vimtmpl bash ~/bin/deploy.sh -d "Deploy application to production"

    Create a Python class, overriding the author email:
        vimtmpl pythonclass mymodule.py -t MyClass -m dev@example.com

    Create an Ansible playbook:
        vimtmpl playbook site.yml -d "Configure webservers"

`

func availableTemplatesSection() string {
	names := templates.ListTemplateNames()
	if len(names) == 0 {
		return "    (none found — run \"" + os.Args[0] + " init\" to install default templates)\n\n"
	}
	sort.Strings(names)
	var sb strings.Builder
	for _, n := range names {
		fmt.Fprintf(&sb, "    %s\n", n)
	}
	sb.WriteString("\n")
	return sb.String()
}

func printHelp() {
	cfgPath := config.UserHomeDir() + "/" + config.ConfigFilename
	fmt.Printf(helpTemplate, os.Args[0], cfgPath, availableTemplatesSection())
}
