// internal/tui/view.go
package tui

import (
	"fmt"
	"strings"
)

// View renders the UI
func (m Model) View() string {
	if !m.Ready {
		return loadingStyle.Render("Initializing...")
	}

	var sections []string

	// Header
	sections = append(sections, m.renderHeader())

	// Messages
	sections = append(sections, m.renderMessages())

	// Input
	sections = append(sections, inputStyle.Render(m.Input.View()))

	// Help
	sections = append(sections, m.renderHelp())

	return strings.Join(sections, "\n")
}

// renderHeader renders the header
func (m Model) renderHeader() string {
	title := "🤖 mycli"
	title += fmt.Sprintf(" | Messages: %d", len(m.Messages))
	return headerStyle.Render(title)
}

// renderMessages renders the message history
func (m Model) renderMessages() string {
	var b strings.Builder

	for _, msg := range m.Messages {
		var prefix string
		var style = messageAssistantStyle

		switch msg.Role {
		case RoleUser:
			prefix = "You: "
			style = messageUserStyle
		case RoleAssistant:
			prefix = "AI: "
			style = messageAssistantStyle
		case RoleSystem:
			prefix = "System: "
			style = messageSystemStyle
		}

		b.WriteString(style.Render(prefix) + msg.Content + "\n\n")
	}

	// Set viewport content
	m.Viewport.SetContent(b.String())

	return m.Viewport.View()
}

// renderHelp renders help text
func (m Model) renderHelp() string {
	help := "Enter: send | Ctrl+C: quit"
	return helpStyle.Render(help)
}
