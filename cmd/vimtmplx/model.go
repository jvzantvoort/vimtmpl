package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jvzantvoort/vimtmpl/config"
	"github.com/jvzantvoort/vimtmpl/generate"
	"github.com/jvzantvoort/vimtmpl/templates"
)

const (
	fieldFilename = iota
	fieldTemplate
	fieldSwitches
	numFields
)

// switchItem represents a single toggleable template switch, as declared in
// a template's [switches] header section.
type switchItem struct {
	Name        string
	Description string
	Enabled     bool
}

// statusMsg is a short-lived message shown at the bottom of the UI.
type statusMsg struct {
	text  string
	isErr bool
}

// editorFinishedMsg reports the outcome of trying to replace this process
// with the external editor after a successful write. err is only non-nil
// when the replacement itself couldn't happen (e.g. unsupported on
// Windows) — on success the process is gone and this is never sent.
type editorFinishedMsg struct {
	path string
	err  error
}

// editorFallbackFinishedMsg reports the outcome of running the editor as an
// ordinary child process, used when process replacement isn't available.
type editorFallbackFinishedMsg struct {
	path string
	err  error
}

// execReplaceCommand adapts generate.ExecEditor to bubbletea's ExecCommand
// interface so tea.Exec can run it. Since ExecEditor replaces the process
// outright, the terminal streams bubbletea would otherwise wire up are
// irrelevant — the new process inherits the real fds directly.
type execReplaceCommand struct {
	path string
}

func (c execReplaceCommand) Run() error          { return generate.ExecEditor(c.path) }
func (c execReplaceCommand) SetStdin(io.Reader)  {}
func (c execReplaceCommand) SetStdout(io.Writer) {}
func (c execReplaceCommand) SetStderr(io.Writer) {}

type model struct {
	filename textinput.Model

	langs   []string
	langIdx int

	switches  []switchItem
	switchIdx int

	focus  int
	status statusMsg

	width, height int
	quitting      bool
}

func newModel(prefill string) model {
	ti := textinput.New()
	ti.Placeholder = "output filename, e.g. deploy.sh"
	ti.CharLimit = 512
	ti.Width = 40
	if prefill != "" {
		ti.SetValue(prefill)
	}
	ti.Focus()

	langs := templates.ListTemplateNames()
	sort.Strings(langs)

	m := model{
		filename: ti,
		langs:    langs,
		focus:    fieldFilename,
	}
	if len(langs) > 0 {
		m.switches = loadSwitches(langs[0])
	}
	return m
}

// loadSwitches returns the ordered list of switches declared in the header
// of the template named lang, or nil if it declares none.
func loadSwitches(lang string) []switchItem {
	header, _, err := templates.ParseLang(lang)
	if err != nil {
		return nil
	}

	sec, err := header.GetSection("switches")
	if err != nil {
		return nil
	}

	keys := sec.Keys()
	items := make([]switchItem, 0, len(keys))
	for _, k := range keys {
		items = append(items, switchItem{Name: k.Name(), Description: k.Value()})
	}
	return items
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % numFields
			m.syncFocus()
			return m, nil
		case "shift+tab":
			m.focus = (m.focus + numFields - 1) % numFields
			m.syncFocus()
			return m, nil
		case "ctrl+s":
			return m.generate()
		}

		switch m.focus {
		case fieldFilename:
			var cmd tea.Cmd
			m.filename, cmd = m.filename.Update(msg)
			return m, cmd

		case fieldTemplate:
			switch msg.String() {
			case "left", "up":
				m.cycleLang(-1)
			case "right", "down":
				m.cycleLang(1)
			}
			return m, nil

		case fieldSwitches:
			switch msg.String() {
			case "up":
				if m.switchIdx > 0 {
					m.switchIdx--
				}
			case "down":
				if m.switchIdx < len(m.switches)-1 {
					m.switchIdx++
				}
			case " ", "enter":
				if len(m.switches) > 0 {
					m.switches[m.switchIdx].Enabled = !m.switches[m.switchIdx].Enabled
				}
			}
			return m, nil
		}

	case editorFinishedMsg:
		// Process replacement failed to even start (e.g. unsupported on
		// this platform) — fall back to running the editor as a child
		// process and resuming the UI afterward.
		m.status = statusMsg{text: fmt.Sprintf("opening editor via exec failed (%s), falling back...", msg.err)}
		editorCmd := generate.EditorCommand(msg.path)
		return m, tea.ExecProcess(editorCmd, func(err error) tea.Msg {
			return editorFallbackFinishedMsg{path: msg.path, err: err}
		})

	case editorFallbackFinishedMsg:
		if msg.err != nil {
			m.status = statusMsg{text: fmt.Sprintf("editor exited with error: %s", msg.err), isErr: true}
		} else {
			m.status = statusMsg{text: fmt.Sprintf("wrote %s", msg.path)}
		}
		return m, nil
	}

	return m, nil
}

func (m *model) syncFocus() {
	if m.focus == fieldFilename {
		m.filename.Focus()
	} else {
		m.filename.Blur()
	}
}

func (m *model) cycleLang(delta int) {
	n := len(m.langs)
	if n == 0 {
		return
	}
	m.langIdx = ((m.langIdx+delta)%n + n) % n
	m.switches = loadSwitches(m.langs[m.langIdx])
	m.switchIdx = 0
	m.status = statusMsg{}
}

// generate renders the selected template with the current field values and
// writes it to disk, reporting the outcome in the status line.
func (m model) generate() (tea.Model, tea.Cmd) {
	if len(m.langs) == 0 {
		m.status = statusMsg{text: `no templates found - run "vimtmpl init" first`, isErr: true}
		return m, nil
	}

	filename := strings.TrimSpace(m.filename.Value())
	if filename == "" {
		m.status = statusMsg{text: "filename is required", isErr: true}
		return m, nil
	}

	lang := m.langs[m.langIdx]

	cfg := config.NewTemplateConfig(lang)
	cfg.Load()

	fullPath := filename
	if !filepath.IsAbs(fullPath) {
		if wd, err := os.Getwd(); err == nil {
			fullPath = filepath.Join(wd, fullPath)
		}
	}
	cfg.FullPath = fullPath
	cfg.ScriptName = filepath.Base(fullPath)

	if cfg.Description == "" {
		cfg.Description = cfg.GetKeyAsString("description")
	}

	for _, sw := range m.switches {
		cfg.Flags[sw.Name] = sw.Enabled
	}

	content, err := generate.Render(cfg)
	if err != nil {
		m.status = statusMsg{text: fmt.Sprintf("render failed: %s", err), isErr: true}
		return m, nil
	}

	if err := generate.WriteFile(cfg, content); err != nil {
		m.status = statusMsg{text: fmt.Sprintf("write failed: %s", err), isErr: true}
		return m, nil
	}

	m.status = statusMsg{text: fmt.Sprintf("wrote %s, opening editor...", fullPath)}
	return m, tea.Exec(execReplaceCommand{path: fullPath}, func(err error) tea.Msg {
		return editorFinishedMsg{path: fullPath, err: err}
	})
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("vimtmplx %s (%s)", version, revision)))
	b.WriteString("\n\n")

	// Filename
	b.WriteString(fieldLabel(m.focus == fieldFilename, "Filename"))
	b.WriteString("\n  ")
	b.WriteString(m.filename.View())
	b.WriteString("\n\n")

	// Template
	b.WriteString(fieldLabel(m.focus == fieldTemplate, "Template"))
	b.WriteString("\n  ")
	if len(m.langs) == 0 {
		b.WriteString(dimStyle.Render(`(none found - run "vimtmpl init" first)`))
	} else {
		fmt.Fprintf(&b, "<  %s  >", selectedStyle.Render(m.langs[m.langIdx]))
	}
	b.WriteString("\n\n")

	// Switches
	b.WriteString(fieldLabel(m.focus == fieldSwitches, "Switches"))
	b.WriteString("\n")
	if len(m.switches) == 0 {
		b.WriteString(dimStyle.Render("  (this template has no switches)"))
		b.WriteString("\n")
	} else {
		for i, sw := range m.switches {
			box := "[ ]"
			if sw.Enabled {
				box = "[x]"
			}
			line := fmt.Sprintf("%s %-14s %s", box, sw.Name, sw.Description)
			if m.focus == fieldSwitches && i == m.switchIdx {
				line = cursorStyle.Render("> " + line)
			} else {
				line = "  " + line
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	if m.status.text != "" {
		style := okStyle
		if m.status.isErr {
			style = errStyle
		}
		b.WriteString(style.Render(m.status.text))
		b.WriteString("\n\n")
	}

	b.WriteString(helpStyle.Render("tab: next field  ·  ←/→: template  ·  space: toggle switch  ·  ctrl+s: generate  ·  esc: quit"))

	return b.String()
}
