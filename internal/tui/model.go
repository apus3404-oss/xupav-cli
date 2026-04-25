// internal/tui/model.go
package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
)

// Model represents the TUI application state
type Model struct {
	Messages []Message
	Input    textarea.Model
	Viewport viewport.Model
	Width    int
	Height   int
	Ready    bool
}

// NewModel creates a new TUI model
func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Type your message..."
	ta.Focus()

	vp := viewport.New(80, 20)

	return Model{
		Messages: make([]Message, 0),
		Input:    ta,
		Viewport: vp,
	}
}

// AddMessage adds a message to the chat history
func (m *Model) AddMessage(role Role, content string) {
	msg := NewMessage(role, content)
	m.Messages = append(m.Messages, msg)
}

// GetLastMessage returns the last message or nil if empty
func (m *Model) GetLastMessage() *Message {
	if len(m.Messages) == 0 {
		return nil
	}
	return &m.Messages[len(m.Messages)-1]
}

// ClearInput clears the input textarea
func (m *Model) ClearInput() {
	m.Input.Reset()
}
