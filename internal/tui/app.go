// internal/tui/app.go
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/apus3404-oss/xupav-cli/internal/bridge"
	"github.com/apus3404-oss/xupav-cli/internal/config"
)

// App wraps the TUI application
type App struct {
	Config *config.Config
	Bridge *bridge.PythonBridge
}

// NewApp creates a new TUI app
func NewApp(cfg *config.Config, br *bridge.PythonBridge) *App {
	return &App{
		Config: cfg,
		Bridge: br,
	}
}

// Run starts the TUI
func (a *App) Run() error {
	// Create model
	model := NewModel()

	// Add welcome message
	model.AddMessage(RoleSystem, "Welcome to mycli! Type your message and press Enter to send.")

	// Create program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	// Check if there was an error in the model
	if m, ok := finalModel.(Model); ok {
		if m.Input.Err != nil {
			return m.Input.Err
		}
	}

	return nil
}
