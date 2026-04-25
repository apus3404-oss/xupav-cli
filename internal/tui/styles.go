// internal/tui/styles.go
package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	colorPrimary   = lipgloss.Color("12")  // Blue
	colorSecondary = lipgloss.Color("10")  // Green
	colorError     = lipgloss.Color("9")   // Red
	colorMuted     = lipgloss.Color("8")   // Gray
	colorAccent    = lipgloss.Color("13")  // Magenta

	// Header style
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Padding(0, 1)

	// Input style
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	// Message styles
	messageUserStyle = lipgloss.NewStyle().
				Foreground(colorSecondary).
				Bold(true)

	messageAssistantStyle = lipgloss.NewStyle().
				Foreground(colorPrimary)

	messageSystemStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Italic(true)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	// Loading style
	loadingStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Italic(true)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted)
)
