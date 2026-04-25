// internal/tui/view.go
package tui

import (
	"fmt"
	"strings"
)

// View renders the UI
func (m Model) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Render messages
	for _, msg := range m.Messages {
		prefix := "You: "
		if msg.Role == RoleAssistant {
			prefix = "AI: "
		}
		b.WriteString(fmt.Sprintf("%s%s\n\n", prefix, msg.Content))
	}

	// Set viewport content
	m.Viewport.SetContent(b.String())

	// Render viewport
	view := m.Viewport.View() + "\n\n"

	// Render input
	view += m.Input.View()

	return view
}
