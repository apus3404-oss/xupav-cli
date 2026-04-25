// internal/tui/update.go
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			// Get input value
			value := m.Input.Value()
			if value == "" {
				return m, nil
			}

			// Add user message
			m.AddMessage(RoleUser, value)

			// Clear input
			m.ClearInput()

			// Return with command to send to AI (placeholder for now)
			return m, func() tea.Msg {
				return nil
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		if !m.Ready {
			// Set up viewport
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - 5 // Leave space for input
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - 5
		}

		// Update input width
		m.Input.SetWidth(msg.Width - 4)
	}

	// Update input textarea
	m.Input, cmd = m.Input.Update(msg)

	return m, cmd
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}
