// internal/tui/model_test.go
package tui

import (
	"testing"
)

func TestMessage(t *testing.T) {
	msg := Message{
		Role:    RoleUser,
		Content: "test message",
	}

	if msg.Role != RoleUser {
		t.Errorf("expected role user, got %s", msg.Role)
	}

	if msg.Content != "test message" {
		t.Errorf("expected 'test message', got %s", msg.Content)
	}
}

func TestModel_AddMessage(t *testing.T) {
	model := NewModel()

	model.AddMessage(RoleUser, "hello")
	model.AddMessage(RoleAssistant, "hi there")

	if len(model.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(model.Messages))
	}

	if model.Messages[0].Role != RoleUser {
		t.Error("first message should be from user")
	}

	if model.Messages[1].Role != RoleAssistant {
		t.Error("second message should be from assistant")
	}
}

func TestModel_GetLastMessage(t *testing.T) {
	model := NewModel()

	// Empty model
	if msg := model.GetLastMessage(); msg != nil {
		t.Error("expected nil for empty model")
	}

	// Add messages
	model.AddMessage(RoleUser, "first")
	model.AddMessage(RoleAssistant, "second")

	last := model.GetLastMessage()
	if last == nil {
		t.Fatal("expected last message, got nil")
	}

	if last.Content != "second" {
		t.Errorf("expected 'second', got %s", last.Content)
	}
}
