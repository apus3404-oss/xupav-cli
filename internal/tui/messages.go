// internal/tui/messages.go
package tui

import (
	"time"
)

// Role represents message sender
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// Message represents a chat message
type Message struct {
	Role      Role
	Content   string
	Timestamp time.Time
}

// NewMessage creates a new message
func NewMessage(role Role, content string) Message {
	return Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
}
