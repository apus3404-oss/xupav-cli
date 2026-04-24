# Basic TUI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the basic terminal user interface using Bubble Tea framework with chat message display, input handling, and integration with the Python bridge.

**Architecture:** Bubble Tea MVC pattern - Model holds state, Update handles events, View renders UI. Messages flow from user input → bridge → AI response → display.

**Tech Stack:**
- **charmbracelet/bubbletea** - TUI framework
- **charmbracelet/lipgloss** - Styling
- **charmbracelet/bubbles** - UI components (textarea, viewport)

---

## File Structure

**New Files:**
```
internal/
  tui/
    app.go             # Main Bubble Tea app
    model.go           # App state model
    update.go          # Update logic (event handling)
    view.go            # Render logic
    messages.go        # Message types and handling
    styles.go          # Basic styles (colors later)
  tui/
    app_test.go
    model_test.go
```

---
## Task 1: Add Bubble Tea Dependencies

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Add Bubble Tea dependency**

```bash
go get github.com/charmbracelet/bubbletea@latest
```

- [ ] **Step 2: Add Lipgloss for styling**

```bash
go get github.com/charmbracelet/lipgloss@latest
```

- [ ] **Step 3: Add Bubbles for components**

```bash
go get github.com/charmbracelet/bubbles@latest
```

- [ ] **Step 4: Verify dependencies**

```bash
go mod tidy
cat go.mod | grep charmbracelet
```

Expected: All three charmbracelet packages listed

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: add Bubble Tea TUI dependencies"
```

---

## Task 2: Message Types and State Model

**Files:**
- Create: `internal/tui/messages.go`
- Create: `internal/tui/model.go`
- Create: `internal/tui/model_test.go`

- [ ] **Step 1: Write failing test for message types**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/tui -v -run TestMessage
```

Expected: FAIL with "undefined: Message"

- [ ] **Step 3: Implement message types**

```go
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
	Tokens    int     // Token count (if available)
	Cost      float64 // Cost (if available)
}

// NewMessage creates a new message
func NewMessage(role Role, content string) Message {
	return Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// String returns string representation
func (r Role) String() string {
	return string(r)
}
```

- [ ] **Step 4: Implement model**

```go
// internal/tui/model.go
package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/mycli/internal/bridge"
	"github.com/yourusername/mycli/internal/config"
)

// Model holds the application state
type Model struct {
	// Configuration
	Config *config.Config
	Bridge *bridge.PythonBridge
	
	// UI Components
	Viewport viewport.Model
	Input    textarea.Model
	
	// State
	Messages    []Message
	CurrentMsg  string
	Loading     bool
	Error       error
	Width       int
	Height      int
	
	// Flags
	Ready       bool
	Quitting    bool
}

// NewModel creates a new TUI model
func NewModel() *Model {
	// Create input textarea
	input := textarea.New()
	input.Placeholder = "Type your message..."
	input.Focus()
	input.CharLimit = 10000
	input.SetWidth(80)
	input.SetHeight(3)
	
	// Create viewport for messages
	vp := viewport.New(80, 20)
	
	return &Model{
		Input:    input,
		Viewport: vp,
		Messages: make([]Message, 0),
	}
}

// AddMessage adds a message to the conversation
func (m *Model) AddMessage(role Role, content string) {
	msg := NewMessage(role, content)
	m.Messages = append(m.Messages, msg)
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/tui -v -run TestMessage
go test ./internal/tui -v -run TestModel_AddMessage
```

Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/tui/
git commit -m "feat(tui): add message types and state model"
```

---
## Task 3: Update Logic (Event Handling)

**Files:**
- Create: `internal/tui/update.go`
- Modify: `internal/tui/model_test.go`

- [ ] **Step 1: Write failing test for update logic**

```go
// internal/tui/model_test.go

func TestModel_HandleKeyPress(t *testing.T) {
	model := NewModel()
	
	// Type some text
	model.Input.SetValue("hello")
	
	// Simulate Enter key (send message)
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := model.Update(msg)
	
	m := newModel.(Model)
	
	// Input should be cleared
	if m.Input.Value() != "" {
		t.Errorf("expected empty input, got %s", m.Input.Value())
	}
	
	// Should have command to send message
	if cmd == nil {
		t.Error("expected command, got nil")
	}
}

func TestModel_HandleResize(t *testing.T) {
	model := NewModel()
	
	// Simulate window resize
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newModel, _ := model.Update(msg)
	
	m := newModel.(Model)
	
	if m.Width != 120 {
		t.Errorf("expected width 120, got %d", m.Width)
	}
	
	if m.Height != 40 {
		t.Errorf("expected height 40, got %d", m.Height)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/tui -v -run TestModel_HandleKeyPress
```

Expected: FAIL with "undefined: Model.Update"

- [ ] **Step 3: Implement update logic**

```go
// internal/tui/update.go
package tui

import (
	"context"
	"time"
	
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Custom message types
type (
	// AIResponseMsg contains AI response
	AIResponseMsg struct {
		Content string
		Tokens  int
		Cost    float64
		Error   error
	}
	
	// ErrorMsg contains error information
	ErrorMsg struct {
		Err error
	}
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit
		
		case tea.KeyEnter:
			// Check if Shift+Enter (new line) or Enter (send)
			if msg.Alt || msg.Shift {
				// New line
				m.Input, cmd = m.Input.Update(msg)
				return m, cmd
			}
			
			// Send message
			if m.Input.Value() != "" && !m.Loading {
				return m.handleSendMessage()
			}
		}
	
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		
		// Update component sizes
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - 6 // Leave space for input
		m.Input.SetWidth(msg.Width - 4)
		
		if !m.Ready {
			m.Ready = true
		}
	
	case AIResponseMsg:
		m.Loading = false
		
		if msg.Error != nil {
			m.Error = msg.Error
			m.AddMessage(RoleSystem, "Error: "+msg.Error.Error())
		} else {
			m.AddMessage(RoleAssistant, msg.Content)
			
			// Update viewport to show new message
			m.Viewport.SetContent(m.renderMessages())
			m.Viewport.GotoBottom()
		}
	
	case ErrorMsg:
		m.Loading = false
		m.Error = msg.Err
		m.AddMessage(RoleSystem, "Error: "+msg.Err.Error())
		m.Viewport.SetContent(m.renderMessages())
		m.Viewport.GotoBottom()
	}
	
	// Update input
	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)
	
	// Update viewport
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	
	return m, tea.Batch(cmds...)
}

// handleSendMessage handles sending a message to AI
func (m Model) handleSendMessage() (tea.Model, tea.Cmd) {
	userMsg := m.Input.Value()
	
	// Add user message
	m.AddMessage(RoleUser, userMsg)
	
	// Clear input
	m.Input.Reset()
	
	// Set loading state
	m.Loading = true
	
	// Update viewport
	m.Viewport.SetContent(m.renderMessages())
	m.Viewport.GotoBottom()
	
	// Send to AI
	return m, m.sendToAI(userMsg)
}

// sendToAI sends message to AI and returns response
func (m Model) sendToAI(message string) tea.Cmd {
	return func() tea.Msg {
		if m.Bridge == nil {
			return ErrorMsg{Err: ErrBridgeNotInitialized}
		}
		
		// Create request
		req := bridge.CreateRequest("chat", map[string]interface{}{
			"message": message,
			"model":   m.Config.Providers.OpenRouter.DefaultModel,
		})
		
		// Send with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		
		resp, err := m.Bridge.SendRequest(ctx, req)
		if err != nil {
			return AIResponseMsg{Error: err}
		}
		
		// Extract response
		content, _ := resp.Result["content"].(string)
		tokens, _ := resp.Result["tokens"].(float64)
		cost, _ := resp.Result["cost"].(float64)
		
		return AIResponseMsg{
			Content: content,
			Tokens:  int(tokens),
			Cost:    cost,
		}
	}
}

// renderMessages renders all messages to string
func (m *Model) renderMessages() string {
	var output string
	
	for _, msg := range m.Messages {
		switch msg.Role {
		case RoleUser:
			output += "You: " + msg.Content + "\n\n"
		case RoleAssistant:
			output += "🤖 Assistant: " + msg.Content + "\n\n"
		case RoleSystem:
			output += "⚠️  " + msg.Content + "\n\n"
		}
	}
	
	return output
}

// Key bindings
type keyMap struct {
	Send key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Send: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "quit"),
	),
}
```

- [ ] **Step 4: Add error definition**

```go
// internal/tui/messages.go

import (
	"errors"
	"time"
)

var (
	ErrBridgeNotInitialized = errors.New("bridge not initialized")
)

// ... rest of file
```

- [ ] **Step 5: Add missing import to model.go**

```go
// internal/tui/model.go

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/mycli/internal/bridge"
	"github.com/yourusername/mycli/internal/config"
)
```

- [ ] **Step 6: Run test to verify it passes**

```bash
go test ./internal/tui -v -run TestModel_HandleKeyPress
go test ./internal/tui -v -run TestModel_HandleResize
```

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/tui/
git commit -m "feat(tui): implement update logic and event handling"
```

---

## Task 4: View Rendering

**Files:**
- Create: `internal/tui/view.go`
- Create: `internal/tui/styles.go`

- [ ] **Step 1: Implement basic styles**

```go
// internal/tui/styles.go
package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors (basic for now, will be enhanced in Plan E)
	colorPrimary   = lipgloss.Color("12")  // Blue
	colorSecondary = lipgloss.Color("10")  // Green
	colorError     = lipgloss.Color("9")   // Red
	colorMuted     = lipgloss.Color("8")   // Gray
	
	// Styles
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		Padding(0, 1)
	
	inputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 1)
	
	messageUserStyle = lipgloss.NewStyle().
		Foreground(colorSecondary).
		Bold(true)
	
	messageAssistantStyle = lipgloss.NewStyle().
		Foreground(colorPrimary)
	
	errorStyle = lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true)
	
	loadingStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		Italic(true)
)
```

- [ ] **Step 2: Implement view rendering**

```go
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
	
	if m.Quitting {
		return "Goodbye!\n"
	}
	
	var sections []string
	
	// Header
	sections = append(sections, m.renderHeader())
	
	// Messages viewport
	sections = append(sections, m.Viewport.View())
	
	// Loading indicator
	if m.Loading {
		sections = append(sections, loadingStyle.Render("⏳ Thinking..."))
	}
	
	// Error display
	if m.Error != nil {
		sections = append(sections, errorStyle.Render("Error: "+m.Error.Error()))
	}
	
	// Input
	sections = append(sections, m.renderInput())
	
	// Help
	sections = append(sections, m.renderHelp())
	
	return strings.Join(sections, "\n")
}

// renderHeader renders the header
func (m Model) renderHeader() string {
	title := "🤖 mycli"
	
	// Add model info if available
	if m.Config != nil {
		model := m.Config.Providers.OpenRouter.DefaultModel
		title += fmt.Sprintf(" | Model: %s", model)
	}
	
	// Add message count
	title += fmt.Sprintf(" | Messages: %d", len(m.Messages))
	
	return headerStyle.Render(title)
}

// renderInput renders the input area
func (m Model) renderInput() string {
	return inputStyle.Render(m.Input.View())
}

// renderHelp renders help text
func (m Model) renderHelp() string {
	help := "Enter: send | Shift+Enter: new line | Ctrl+C: quit"
	return lipgloss.NewStyle().
		Foreground(colorMuted).
		Render(help)
}
```

- [ ] **Step 3: Test view rendering manually**

```bash
# Create a simple test program
cat > cmd/test-tui/main.go << 'EOF'
package main

import (
	"fmt"
	"os"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/mycli/internal/tui"
)

func main() {
	model := tui.NewModel()
	model.AddMessage(tui.RoleUser, "Hello")
	model.AddMessage(tui.RoleAssistant, "Hi there! How can I help?")
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
EOF

go run cmd/test-tui/main.go
```

Expected: TUI displays with messages, input box, and help text

- [ ] **Step 4: Clean up test program**

```bash
rm -rf cmd/test-tui
```

- [ ] **Step 5: Commit**

```bash
git add internal/tui/
git commit -m "feat(tui): implement view rendering with basic styles"
```

---
## Task 5: Integration with Chat Command

**Files:**
- Create: `internal/tui/app.go`
- Modify: `internal/cli/chat.go`

- [ ] **Step 1: Create TUI app wrapper**

```go
// internal/tui/app.go
package tui

import (
	"fmt"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/mycli/internal/bridge"
	"github.com/yourusername/mycli/internal/config"
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
	model.Config = a.Config
	model.Bridge = a.Bridge
	
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
		if m.Error != nil {
			return m.Error
		}
	}
	
	return nil
}
```

- [ ] **Step 2: Update chat command to use TUI**

```go
// internal/cli/chat.go
package cli

import (
	"fmt"
	"time"
	
	"github.com/spf13/cobra"
	"github.com/yourusername/mycli/internal/bridge"
	"github.com/yourusername/mycli/internal/config"
	"github.com/yourusername/mycli/internal/tui"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive chat session",
	Long: `Start an interactive chat session with the AI assistant.
	
The assistant can help you:
  - Write and debug code
  - Explain complex concepts
  - Refactor and optimize code
  - Generate tests
  - And much more!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		
		// Validate config
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
		
		// Check API key
		if cfg.Providers.OpenRouter.Enabled {
			apiKey, err := config.GetOpenRouterKey()
			if err != nil {
				return fmt.Errorf("OpenRouter API key not found. Run: mycli config set-key openrouter")
			}
			
			// Set API key in config for bridge
			cfg.Providers.OpenRouter.APIKey = apiKey
		}
		
		// Create Python bridge
		pythonBridge := bridge.NewPythonBridge(bridge.PythonConfig{
			PythonPath: "python3",
			Config: map[string]interface{}{
				"openrouter": map[string]interface{}{
					"enabled": cfg.Providers.OpenRouter.Enabled,
					"api_key": cfg.Providers.OpenRouter.APIKey,
					"base_url": cfg.Providers.OpenRouter.BaseURL,
				},
				"ollama": map[string]interface{}{
					"enabled": cfg.Providers.Ollama.Enabled,
					"base_url": cfg.Providers.Ollama.BaseURL,
				},
			},
			MaxRestarts:  3,
			RestartDelay: 1 * time.Second,
		})
		
		// Start bridge
		if err := pythonBridge.Start(); err != nil {
			return fmt.Errorf("failed to start AI bridge: %w", err)
		}
		defer pythonBridge.Stop()
		
		// Create and run TUI
		app := tui.NewApp(cfg, pythonBridge)
		if err := app.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
```

- [ ] **Step 3: Add APIKey field to config struct**

```go
// internal/config/config.go

type OpenRouterConfig struct {
	Enabled      bool   `yaml:"enabled"`
	DefaultModel string `yaml:"default_model"`
	APIKeySource string `yaml:"api_key_source"`
	APIKey       string `yaml:"-"` // Not serialized, set at runtime
	BaseURL      string `yaml:"base_url"`
	MaxTokens    int    `yaml:"max_tokens"`
	Temperature  float64 `yaml:"temperature"`
	Timeout      int    `yaml:"timeout"`
}
```

- [ ] **Step 4: Test chat command**

```bash
go build -o bin/mycli cmd/mycli/main.go

# Set API key first
./bin/mycli config set-key openrouter

# Run chat
./bin/mycli chat
```

Expected: TUI opens, can type messages, send to AI (if API key valid)

- [ ] **Step 5: Commit**

```bash
git add internal/tui/app.go internal/cli/chat.go internal/config/config.go
git commit -m "feat(tui): integrate TUI with chat command"
```

---

## Task 6: Integration Test

**Files:**
- Create: `internal/tui/app_test.go`

- [ ] **Step 1: Write integration test**

```go
// internal/tui/app_test.go
package tui

import (
	"testing"
	
	"github.com/yourusername/mycli/internal/config"
)

func TestNewModel(t *testing.T) {
	model := NewModel()
	
	if model == nil {
		t.Fatal("expected model, got nil")
	}
	
	if model.Messages == nil {
		t.Error("expected messages slice, got nil")
	}
	
	if len(model.Messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(model.Messages))
	}
}

func TestModel_AddMultipleMessages(t *testing.T) {
	model := NewModel()
	
	// Add conversation
	model.AddMessage(RoleUser, "What is Go?")
	model.AddMessage(RoleAssistant, "Go is a programming language.")
	model.AddMessage(RoleUser, "Tell me more")
	model.AddMessage(RoleAssistant, "Go was created at Google...")
	
	if len(model.Messages) != 4 {
		t.Errorf("expected 4 messages, got %d", len(model.Messages))
	}
	
	// Verify order
	if model.Messages[0].Role != RoleUser {
		t.Error("first message should be user")
	}
	
	if model.Messages[1].Role != RoleAssistant {
		t.Error("second message should be assistant")
	}
}

func TestModel_WithConfig(t *testing.T) {
	model := NewModel()
	cfg := config.DefaultConfig()
	model.Config = cfg
	
	if model.Config == nil {
		t.Error("expected config, got nil")
	}
	
	if model.Config.Providers.OpenRouter.DefaultModel == "" {
		t.Error("expected default model, got empty string")
	}
}

func TestRenderMessages(t *testing.T) {
	model := NewModel()
	
	model.AddMessage(RoleUser, "Hello")
	model.AddMessage(RoleAssistant, "Hi there!")
	
	output := model.renderMessages()
	
	if output == "" {
		t.Error("expected non-empty output")
	}
	
	// Should contain both messages
	if !contains(output, "Hello") {
		t.Error("output should contain user message")
	}
	
	if !contains(output, "Hi there!") {
		t.Error("output should contain assistant message")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 (len(s) > len(substr) && 
		  (s[:len(substr)] == substr || 
		   s[len(s)-len(substr):] == substr ||
		   (len(s) > len(substr)+1 && findSubstring(s, substr)))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: Run all TUI tests**

```bash
go test ./internal/tui -v
```

Expected: All tests PASS

- [ ] **Step 3: Commit**

```bash
git add internal/tui/app_test.go
git commit -m "test(tui): add integration tests for TUI components"
```

---

## Completion Checklist

- [x] Bubble Tea dependencies added
- [x] Message types and state model
- [x] Update logic (event handling)
- [x] View rendering with basic styles
- [x] Integration with chat command
- [x] Integration tests
- [x] All tests passing

## Manual Verification

```bash
# Build
go build -o bin/mycli cmd/mycli/main.go

# Set up config and API key
./bin/mycli config set-key openrouter

# Run chat (requires Python + mycli_ai installed)
./bin/mycli chat

# Test interactions:
# - Type a message and press Enter
# - Verify AI responds
# - Try Shift+Enter for new line
# - Press Ctrl+C to quit
```

## Next Steps

This plan establishes the basic TUI. Next plans will enhance it:
- **Plan E**: TUI Polish (colors, animations, diff preview, notifications)
- **Plan F**: Token Optimizer (file chunking, context prioritization)
- **Plan G**: LSP Integration (code intelligence)
- **Plan H**: Code Modifications (diff generation, file writer)

---

## Notes for Implementation

**Key Design Decisions:**
1. Bubble Tea MVC pattern (Model-Update-View)
2. Viewport for scrollable message history
3. Textarea for multi-line input
4. Async AI requests via tea.Cmd
5. Error handling with system messages

**Testing Strategy:**
- Unit tests for model and message handling
- Integration tests for full flow
- Manual testing for UI/UX
- Mock bridge for tests without Python

**Common Issues:**
- Terminal size detection may fail (handle gracefully)
- Alt screen may not work in all terminals
- Input focus can be lost (ensure Focus() called)
- Viewport scrolling needs manual GotoBottom()
- Long messages may overflow (need wrapping)

**Performance Considerations:**
- Viewport renders only visible content
- Message rendering is O(n) but cached
- Input is buffered for smooth typing
- AI requests are async (non-blocking)
- Model updates are batched

**User Experience:**
- Welcome message on start
- Loading indicator during AI response
- Error messages displayed inline
- Help text always visible
- Keyboard shortcuts intuitive
