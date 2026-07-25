package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230"))

	labelStyle        = lipgloss.NewStyle().Bold(true)
	focusedLabelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))

	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	helpStyle = dimStyle
)

// fieldLabel renders a field name, highlighting it when focused.
func fieldLabel(focused bool, name string) string {
	if focused {
		return focusedLabelStyle.Render("▶ " + name)
	}
	return labelStyle.Render("  " + name)
}
