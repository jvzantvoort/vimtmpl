# vimtmpl

[![Build](https://github.com/jvzantvoort/vimtmpl/actions/workflows/build.yml/badge.svg)](https://github.com/jvzantvoort/vimtmpl/actions/workflows/build.yml)
[![Lint](https://github.com/jvzantvoort/vimtmpl/actions/workflows/lint.yml/badge.svg)](https://github.com/jvzantvoort/vimtmpl/actions/workflows/lint.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/jvzantvoort/vimtmpl.svg)](https://pkg.go.dev/github.com/jvzantvoort/vimtmpl)

[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)

Generate boilerplate scripts and source files from
[Go `text/template`][Go text/template] templates - from the shell
with **vimtmpl**, or from a small terminal UI with **vimtmplx**.

Both fill in author, company, license, and date automatically,
support optional per-template switches (safety flags, headers,
helper functions, ...), and open the result straight in your editor
once it's written.

:exclamation: This is for goofing around... be warned.

## Contents

- [Install](#install)
- [Quick start](#quick-start)
- [vimtmpl (CLI)](#vimtmpl-cli)
- [vimtmplx (terminal UI)](#vimtmplx-terminal-ui)
- [Configuration file](#configuration-file)
- [Templates](#templates)
- [Editor integration](#editor-integration)
- [Development](#development)
- [License](#license)


## Install


### Prebuilt release

You can get the latest releases from 
[GitHub Releases](https://github.com/jvzantvoort/vimtmpl/releases).

The bundled `vimtmpl-update` script detects your platform, downloads
the latest matching release, and installs both binaries into
`/usr/local/bin` (as root) or `$GOBIN`/`~/bin` otherwise. Grab it
and run it locally - worth a read before running, as with any
install script:

```shell
curl -fsSLO https://raw.githubusercontent.com/jvzantvoort/vimtmpl/master/vimtmpl-update
chmod +x vimtmpl-update
./vimtmpl-update
```
![download vimtmpl](docs/images/download.gif)

Re-run it any time to update to the latest release.

### Build from source

Requires Go 1.25+.

```shell
git clone https://github.com/jvzantvoort/vimtmpl.git
cd vimtmpl
./build.sh build      # build both binaries for the host platform into build/<goos>/<arch>/
./build.sh install    # install them into $GOBIN (or ~/go/bin)
```

See [Development](#development) for the rest of `build.sh`'s subcommands.

## Quick start

```shell
vimtmpl init # install default templates + skeleton config
```

![init vimtmpl](docs/images/init.gif)

**Note**: Update the ``~/.template.cfg`` file.

```shell
vimtmpl bash ~/bin/deploy.sh -d "Deploy application to production" # generate, then open in $EDITOR
```

or interactively:

```shell
vimtmplx
```

![create bash](docs/images/bash.gif)

## vimtmpl (CLI)

```
vimtmpl init
vimtmpl help
vimtmpl <template> <filename> [options]
```

`init` installs the bundled templates into `~/.templates.d/default/`
and writes a skeleton `~/.template.cfg` - safe to run repeatedly, it
never overwrites existing files. `help` prints the full built-in
reference (flags, template variables, config format, examples) - run
`vimtmpl help` for the authoritative, always-up-to-date version.

| Flag | Description |
|---|---|
| `-m, --mailaddress <addr>` | Author email address |
| `-c, --company <name>` | Company name |
| `-C, --copyright <holder>` | Copyright holder (defaults to company) |
| `-l, --license <id>` | License identifier, e.g. `MIT`, `Apache-2.0` |
| `-U, --user <account>` | Login/account name of the author |
| `-u, --username <name>` | Author's full name |
| `-t, --title <title>` | Title for templates that need one, e.g. a Python class name |
| `-d, --description <text>` | Short description of the file being created |
| `-f, --flags <name[,name...]>` | Enable one or more template switches (repeatable) |
| `-i, --info` | Print the chosen template's available switches and exit |
| `-v, --verbose` | Enable debug logging |

Any flag not given falls back to the config file, then to the
template's built-in default.

## vimtmplx (terminal UI)

```
vimtmplx [-c|--cwd <dir>] [-f|--filename <name>]
```

`--cwd` sets the directory relative filenames resolve against;
`--filename` pre-fills the filename field. Templates, switches, and
config are shared with `vimtmpl` - it writes the exact same output.

| Key | Action |
|---|---|
| `tab` / `shift+tab` | Move between fields |
| `↑` / `↓` | Move within the focused field |
| `←` / `→` | Change the selected template |
| `space` / `enter` | Toggle the highlighted switch |
| `ctrl+s` | Generate the file, then open it in your editor |
| `esc` / `ctrl+c` | Quit without writing |

## Configuration file

Read from `~/.template.cfg`, INI format. A `[DEFAULT]` section
applies to every template; a section named after a template
overrides individual keys just for that one.

```ini
[DEFAULT]
user        = dduck
username    = Donald Duck
company     = Ducktown
copyright   = Donald Duck
mailaddress = d.duck@example.com
license     = MIT
mode        = 0644

[bash]
description = Bash script
mode        = 0755

[go]
mode = 0644

[playbook]
mode = 0644

[python]
mode = 0755

[pythonclass]
mode = 0644
```

`mode` sets the octal permissions of the generated file (`vimtmpl
init` writes it unset - add it per section, or under `[DEFAULT]`, to
control it). `description` is used when `-d/--description` isn't
passed on the command line.

## Templates

Templates are `.gtmpl` files
([Go `text/template`](https://pkg.go.dev/text/template) syntax) with an
optional INI header between `---` markers declaring their
`[switches]`. Two locations are searched, in order:

1. `~/.templates.d/local/` - your own templates, checked first.
2. `~/.templates.d/default/` - the bundled defaults, installed by `vimtmpl init`.

A file in `local` with the same name as one in `default` wins, so
you can override a bundled template without touching it.

### Bundled templates

| Template | Switches |
|---|---|
| `bash` | `safe` (adds `-eu` to the shebang), `header` (comment block header), `pathmunge` (adds a `pathmunge()` helper), `staging` (adds a staging-dir helper) |
| `go` | none |
| `playbook` | `header` |
| `python` | `header`, `version` |
| `pythonclass` | none |

List a template's switches without writing the file (a filename is
still required as a positional argument, but nothing is created):

```shell
vimtmpl bash deploy.sh --info
```

### Template variables

Available inside every `.gtmpl` file, resolved as: command-line flag
> config file section for the chosen template > `[DEFAULT]` section
> built-in default.

| Variable | Source |
|---|---|
| `{{.ScriptName}}` | Basename of `<filename>` |
| `{{.FullPath}}` | Full path of the output file |
| `{{.Lang}}` | Template name used |
| `{{.Date}}` / `{{.Year}}` | Auto-populated creation date |
| `{{.User}}` / `{{.UserName}}` | `-U` / `user` key, `-u` / `username` key |
| `{{.MailAddress}}` | `-m` / `mailaddress` key |
| `{{.Company}}` / `{{.Copyright}}` | `-c` / `company` key, `-C` / `copyright` key |
| `{{.License}}` | `-l` / `license` key |
| `{{.Title}}` | `-t` flag |
| `{{.Description}}` | `-d` flag |
| `{{.Enabled "name"}}` / `{{index .Flags "name"}}` | Whether switch `name` was passed via `-f` |

### Writing your own

Drop a `.gtmpl` file into `~/.templates.d/local/`. To declare
switches, start the file with an INI header:

```
---
[switches]
safe = Enable safety switches
---
#!/bin/bash{{ if index .Flags "safe" }} -eu{{ end }}
echo {{.ScriptName}}
```

A template with no header simply has no switches.

## Editor integration

After writing the file, both tools open it in `$EDITOR` (falling
back to `vim`) by replacing their own process with the editor
(`execve`, like a shell's `exec`) - there's no parent process left
waiting once the editor starts. `vimtmplx` first suspends its
terminal UI so the editor gets a clean screen. On platforms where
process replacement isn't available (Windows), both tools fall back
to spawning the editor as a child process and waiting for it to exit
instead.

## Development

```shell
./build.sh build       # build both binaries for the host platform
./build.sh buildr <goos> <goarch> <outdir>   # cross-compile into outdir/<goos>.<goarch>/
./build.sh install     # install into $GOBIN
./build.sh check       # go vet + golangci-lint + staticcheck
./build.sh cleanup     # remove build/ and pkg/

go test ./...
go test -race -coverprofile=coverage.out ./...
```

CI runs the build matrix, `golangci-lint`, `go
vet`/`staticcheck`/race-tested coverage, and a weekly
dependency-update PR (see `.github/workflows/`). Pushing a
`vimtmpl-*` tag triggers `release.yml`, which builds and publishes
the release archives described in [Install](#install).

## License

[MIT](LICENSE) © John van Zantvoort
