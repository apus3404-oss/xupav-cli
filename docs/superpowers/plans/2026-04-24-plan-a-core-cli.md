# Core CLI + Configuration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the foundational CLI structure with Cobra, configuration management, secure API key storage, and basic commands.

**Architecture:** Go-based CLI using Cobra for command structure. Configuration stored in YAML with environment variable overrides. API keys secured in OS keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service) with encrypted fallback.

**Tech Stack:**
- **cobra** - CLI framework
- **viper** - Configuration management
- **zalando/go-keyring** - OS keychain integration
- **gopkg.in/yaml.v3** - YAML parsing

---

## File Structure

**New Files:**
```
cmd/mycli/
  main.go                    # Entry point
internal/
  cli/
    root.go                  # Root command
    chat.go                  # Chat command (stub)
    config.go                # Config commands
    version.go               # Version command
  config/
    config.go                # Config struct and loading
    defaults.go              # Default configuration
    validation.go            # Config validation
    keychain.go              # Keychain integration
    env.go                   # Environment variable overrides
go.mod                       # Go module definition
go.sum                       # Dependency checksums
.gitignore                   # Git ignore patterns
```

**Test Files:**
```
internal/config/
  config_test.go
  keychain_test.go
  validation_test.go
  env_test.go
```

---
## Task 1: Project Initialization

**Files:**
- Create: `go.mod`
- Create: `.gitignore`
- Create: `cmd/mycli/main.go`

- [ ] **Step 1: Initialize Go module**

```bash
go mod init github.com/yourusername/mycli
```

Expected output: `go: creating new go.mod: module github.com/yourusername/mycli`

- [ ] **Step 2: Create .gitignore**

```gitignore
# Binaries
bin/
*.exe
*.dll
*.so
*.dylib

# Test coverage
*.out
coverage.html

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Build artifacts
dist/
```

- [ ] **Step 3: Write minimal main.go**

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("mycli v0.1.0")
	os.Exit(0)
}
```

- [ ] **Step 4: Test build**

```bash
go build -o bin/mycli cmd/mycli/main.go
./bin/mycli
```

Expected output: `mycli v0.1.0`

- [ ] **Step 5: Commit**

```bash
git add go.mod .gitignore cmd/mycli/main.go
git commit -m "chore: initialize Go project structure"
```

---

## Task 2: Add Dependencies

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Add Cobra dependency**

```bash
go get github.com/spf13/cobra@latest
```

- [ ] **Step 2: Add Viper dependency**

```bash
go get github.com/spf13/viper@latest
```

- [ ] **Step 3: Add keyring dependency**

```bash
go get github.com/zalando/go-keyring@latest
```

- [ ] **Step 4: Add YAML dependency**

```bash
go get gopkg.in/yaml.v3@latest
```

- [ ] **Step 5: Verify dependencies**

```bash
go mod tidy
cat go.mod
```

Expected: All four dependencies listed in `require` section

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: add CLI and config dependencies"
```

---
## Task 3: Config Struct and Defaults

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/defaults.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write failing test for default config**

```go
// internal/config/config_test.go
package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	if cfg.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", cfg.Version)
	}
	
	if cfg.Providers.OpenRouter.DefaultModel != "deepseek/deepseek-r1" {
		t.Errorf("expected default model deepseek/deepseek-r1, got %s", 
			cfg.Providers.OpenRouter.DefaultModel)
	}
	
	if cfg.Behavior.CodeMode != "interactive" {
		t.Errorf("expected code mode interactive, got %s", cfg.Behavior.CodeMode)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/config -v
```

Expected: FAIL with "undefined: DefaultConfig"

- [ ] **Step 3: Write config struct**

```go
// internal/config/config.go
package config

type Config struct {
	Version   string              `yaml:"version"`
	Providers ProvidersConfig     `yaml:"providers"`
	Behavior  BehaviorConfig      `yaml:"behavior"`
	UI        UIConfig            `yaml:"ui"`
	LSP       LSPConfig           `yaml:"lsp"`
	Cache     CacheConfig         `yaml:"cache"`
	Logging   LoggingConfig       `yaml:"logging"`
}

type ProvidersConfig struct {
	OpenRouter OpenRouterConfig `yaml:"openrouter"`
	Ollama     OllamaConfig     `yaml:"ollama"`
}

type OpenRouterConfig struct {
	Enabled      bool   `yaml:"enabled"`
	DefaultModel string `yaml:"default_model"`
	APIKeySource string `yaml:"api_key_source"`
	BaseURL      string `yaml:"base_url"`
	MaxTokens    int    `yaml:"max_tokens"`
	Temperature  float64 `yaml:"temperature"`
	Timeout      int    `yaml:"timeout"`
}

type OllamaConfig struct {
	Enabled      bool   `yaml:"enabled"`
	BaseURL      string `yaml:"base_url"`
	DefaultModel string `yaml:"default_model"`
	Fallback     bool   `yaml:"fallback"`
	Timeout      int    `yaml:"timeout"`
}

type BehaviorConfig struct {
	CodeMode            string `yaml:"code_mode"`
	AutoAddFiles        bool   `yaml:"auto_add_files"`
	MaxContextFiles     int    `yaml:"max_context_files"`
	TokenBudget         int    `yaml:"token_budget"`
	ConversationMemory  int    `yaml:"conversation_memory"`
}

type UIConfig struct {
	Theme            string `yaml:"theme"`
	Animations       bool   `yaml:"animations"`
	SyntaxHighlight  bool   `yaml:"syntax_highlight"`
	ShowTokenCount   bool   `yaml:"show_token_count"`
	ShowCost         bool   `yaml:"show_cost"`
	ColorScheme      string `yaml:"color_scheme"`
}

type LSPConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Languages  []string `yaml:"languages"`
	Timeout    int      `yaml:"timeout"`
	MaxQueries int      `yaml:"max_queries"`
}

type CacheConfig struct {
	Enabled  bool   `yaml:"enabled"`
	TTL      int    `yaml:"ttl"`
	MaxSize  string `yaml:"max_size"`
	Location string `yaml:"location"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    string `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
}
```

- [ ] **Step 4: Write default config**

```go
// internal/config/defaults.go
package config

func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Providers: ProvidersConfig{
			OpenRouter: OpenRouterConfig{
				Enabled:      true,
				DefaultModel: "deepseek/deepseek-r1",
				APIKeySource: "keychain",
				BaseURL:      "https://openrouter.ai/api/v1",
				MaxTokens:    4096,
				Temperature:  0.7,
				Timeout:      60000,
			},
			Ollama: OllamaConfig{
				Enabled:      true,
				BaseURL:      "http://localhost:11434",
				DefaultModel: "codellama:13b",
				Fallback:     true,
				Timeout:      120000,
			},
		},
		Behavior: BehaviorConfig{
			CodeMode:           "interactive",
			AutoAddFiles:       true,
			MaxContextFiles:    10,
			TokenBudget:        60000,
			ConversationMemory: 10,
		},
		UI: UIConfig{
			Theme:           "gradient",
			Animations:      true,
			SyntaxHighlight: true,
			ShowTokenCount:  true,
			ShowCost:        true,
			ColorScheme:     "auto",
		},
		LSP: LSPConfig{
			Enabled:    true,
			Languages:  []string{"go", "python", "javascript", "typescript", "rust"},
			Timeout:    5000,
			MaxQueries: 5,
		},
		Cache: CacheConfig{
			Enabled:  true,
			TTL:      3600,
			MaxSize:  "100MB",
			Location: "~/.mycli/cache",
		},
		Logging: LoggingConfig{
			Level:      "info",
			File:       "~/.mycli/logs/mycli.log",
			MaxSize:    "10MB",
			MaxBackups: 3,
		},
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/config -v
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/config/
git commit -m "feat(config): add config struct and defaults"
```

---

## Task 4: Config Loading and Saving

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write failing test for load/save**

```go
// internal/config/config_test.go
func TestConfigLoadSave(t *testing.T) {
	// Create temp config file
	tmpFile := "/tmp/test-config.yaml"
	defer os.Remove(tmpFile)
	
	// Create config
	cfg := DefaultConfig()
	cfg.Providers.OpenRouter.DefaultModel = "test-model"
	
	// Save
	err := cfg.Save(tmpFile)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}
	
	// Load
	loaded, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	// Verify
	if loaded.Providers.OpenRouter.DefaultModel != "test-model" {
		t.Errorf("expected test-model, got %s", 
			loaded.Providers.OpenRouter.DefaultModel)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/config -v -run TestConfigLoadSave
```

Expected: FAIL with "undefined: Config.Save"

- [ ] **Step 3: Add imports to config.go**

```go
// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	
	"gopkg.in/yaml.v3"
)

// ... existing structs ...
```

- [ ] **Step 4: Implement Save method**

```go
// internal/config/config.go

func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}
```

- [ ] **Step 5: Implement Load function**

```go
// internal/config/config.go

func Load(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}
	
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Unmarshal YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	return &cfg, nil
}
```

- [ ] **Step 6: Add missing import to test**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"testing"
)

// ... existing tests ...
```

- [ ] **Step 7: Run test to verify it passes**

```bash
go test ./internal/config -v -run TestConfigLoadSave
```

Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/config/
git commit -m "feat(config): add load and save functionality"
```

---
## Task 5: Config Validation

**Files:**
- Create: `internal/config/validation.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write failing test for validation**

```go
// internal/config/config_test.go

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name: "token budget too low",
			cfg: &Config{
				Behavior: BehaviorConfig{
					TokenBudget: 500,
				},
			},
			wantErr: true,
			errMsg:  "token_budget must be >= 1000",
		},
		{
			name: "invalid code mode",
			cfg: &Config{
				Behavior: BehaviorConfig{
					CodeMode:    "invalid",
					TokenBudget: 60000,
				},
			},
			wantErr: true,
			errMsg:  "code_mode must be one of: safe, interactive, auto",
		},
		{
			name: "max context files too low",
			cfg: &Config{
				Behavior: BehaviorConfig{
					CodeMode:        "interactive",
					TokenBudget:     60000,
					MaxContextFiles: 0,
				},
			},
			wantErr: true,
			errMsg:  "max_context_files must be >= 1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
				} else if !contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		len(s) > len(substr)+1 && s[1:len(substr)+1] == substr))
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/config -v -run TestConfigValidation
```

Expected: FAIL with "undefined: Config.Validate"

- [ ] **Step 3: Implement validation**

```go
// internal/config/validation.go
package config

import (
	"fmt"
	"strings"
)

func (c *Config) Validate() error {
	// Validate behavior
	if err := c.validateBehavior(); err != nil {
		return err
	}
	
	// Validate providers
	if err := c.validateProviders(); err != nil {
		return err
	}
	
	// Validate LSP
	if err := c.validateLSP(); err != nil {
		return err
	}
	
	return nil
}

func (c *Config) validateBehavior() error {
	// Code mode
	validModes := []string{"safe", "interactive", "auto"}
	if !contains(validModes, c.Behavior.CodeMode) {
		return fmt.Errorf("code_mode must be one of: %s", strings.Join(validModes, ", "))
	}
	
	// Token budget
	if c.Behavior.TokenBudget < 1000 {
		return fmt.Errorf("token_budget must be >= 1000, got %d", c.Behavior.TokenBudget)
	}
	
	// Max context files
	if c.Behavior.MaxContextFiles < 1 {
		return fmt.Errorf("max_context_files must be >= 1, got %d", c.Behavior.MaxContextFiles)
	}
	
	// Conversation memory
	if c.Behavior.ConversationMemory < 1 {
		return fmt.Errorf("conversation_memory must be >= 1, got %d", c.Behavior.ConversationMemory)
	}
	
	return nil
}

func (c *Config) validateProviders() error {
	// At least one provider must be enabled
	if !c.Providers.OpenRouter.Enabled && !c.Providers.Ollama.Enabled {
		return fmt.Errorf("at least one provider must be enabled")
	}
	
	// OpenRouter validation
	if c.Providers.OpenRouter.Enabled {
		if c.Providers.OpenRouter.DefaultModel == "" {
			return fmt.Errorf("openrouter.default_model cannot be empty")
		}
		if c.Providers.OpenRouter.MaxTokens < 100 {
			return fmt.Errorf("openrouter.max_tokens must be >= 100")
		}
	}
	
	// Ollama validation
	if c.Providers.Ollama.Enabled {
		if c.Providers.Ollama.BaseURL == "" {
			return fmt.Errorf("ollama.base_url cannot be empty")
		}
	}
	
	return nil
}

func (c *Config) validateLSP() error {
	if c.LSP.Enabled {
		if c.LSP.Timeout < 1000 {
			return fmt.Errorf("lsp.timeout must be >= 1000ms, got %d", c.LSP.Timeout)
		}
		if c.LSP.MaxQueries < 1 {
			return fmt.Errorf("lsp.max_queries must be >= 1, got %d", c.LSP.MaxQueries)
		}
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/config -v -run TestConfigValidation
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat(config): add validation logic"
```

---

## Task 6: Environment Variable Overrides

**Files:**
- Create: `internal/config/env.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write failing test for env overrides**

```go
// internal/config/config_test.go

func TestEnvOverrides(t *testing.T) {
	// Set env vars
	os.Setenv("MYCLI_MODEL", "test-model")
	os.Setenv("MYCLI_CODE_MODE", "auto")
	os.Setenv("MYCLI_TOKEN_BUDGET", "80000")
	defer func() {
		os.Unsetenv("MYCLI_MODEL")
		os.Unsetenv("MYCLI_CODE_MODE")
		os.Unsetenv("MYCLI_TOKEN_BUDGET")
	}()
	
	// Load default config
	cfg := DefaultConfig()
	
	// Apply env overrides
	ApplyEnvOverrides(cfg)
	
	// Verify overrides
	if cfg.Providers.OpenRouter.DefaultModel != "test-model" {
		t.Errorf("expected test-model, got %s", cfg.Providers.OpenRouter.DefaultModel)
	}
	
	if cfg.Behavior.CodeMode != "auto" {
		t.Errorf("expected auto, got %s", cfg.Behavior.CodeMode)
	}
	
	if cfg.Behavior.TokenBudget != 80000 {
		t.Errorf("expected 80000, got %d", cfg.Behavior.TokenBudget)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/config -v -run TestEnvOverrides
```

Expected: FAIL with "undefined: ApplyEnvOverrides"

- [ ] **Step 3: Implement env overrides**

```go
// internal/config/env.go
package config

import (
	"os"
	"strconv"
	"strings"
)

func ApplyEnvOverrides(cfg *Config) {
	// Provider selection
	if v := os.Getenv("MYCLI_PROVIDER"); v != "" {
		switch v {
		case "openrouter":
			cfg.Providers.OpenRouter.Enabled = true
			cfg.Providers.Ollama.Enabled = false
		case "ollama":
			cfg.Providers.OpenRouter.Enabled = false
			cfg.Providers.Ollama.Enabled = true
		}
	}
	
	// Model
	if v := os.Getenv("MYCLI_MODEL"); v != "" {
		cfg.Providers.OpenRouter.DefaultModel = v
	}
	
	// Behavior
	if v := os.Getenv("MYCLI_CODE_MODE"); v != "" {
		cfg.Behavior.CodeMode = v
	}
	
	if v := os.Getenv("MYCLI_TOKEN_BUDGET"); v != "" {
		if budget, err := strconv.Atoi(v); err == nil {
			cfg.Behavior.TokenBudget = budget
		}
	}
	
	// UI
	if v := os.Getenv("MYCLI_THEME"); v != "" {
		cfg.UI.Theme = v
	}
	
	if v := os.Getenv("MYCLI_NO_ANIMATIONS"); v != "" {
		cfg.UI.Animations = !parseBool(v)
	}
	
	if v := os.Getenv("MYCLI_NO_COLOR"); v != "" {
		if parseBool(v) {
			cfg.UI.ColorScheme = "none"
		}
	}
	
	// Debug
	if v := os.Getenv("MYCLI_DEBUG"); v != "" {
		if parseBool(v) {
			cfg.Logging.Level = "debug"
		}
	}
	
	if v := os.Getenv("MYCLI_LOG_LEVEL"); v != "" {
		cfg.Logging.Level = v
	}
}

func parseBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "1" || s == "yes" || s == "on"
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/config -v -run TestEnvOverrides
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat(config): add environment variable overrides"
```

---
## Task 7: Keychain Integration

**Files:**
- Create: `internal/config/keychain.go`
- Modify: `internal/config/keychain_test.go`

- [ ] **Step 1: Write failing test for keychain**

```go
// internal/config/keychain_test.go
package config

import (
	"testing"
)

func TestKeychainSetGet(t *testing.T) {
	service := "mycli-test"
	account := "openrouter"
	secret := "test-api-key-12345"
	
	// Set key
	err := SetAPIKey(service, account, secret)
	if err != nil {
		t.Skipf("keychain not available: %v", err)
	}
	
	// Get key
	retrieved, err := GetAPIKey(service, account)
	if err != nil {
		t.Fatalf("failed to get API key: %v", err)
	}
	
	if retrieved != secret {
		t.Errorf("expected %s, got %s", secret, retrieved)
	}
	
	// Clean up
	DeleteAPIKey(service, account)
}

func TestKeychainDelete(t *testing.T) {
	service := "mycli-test"
	account := "test-delete"
	secret := "test-key"
	
	// Set key
	err := SetAPIKey(service, account, secret)
	if err != nil {
		t.Skipf("keychain not available: %v", err)
	}
	
	// Delete key
	err = DeleteAPIKey(service, account)
	if err != nil {
		t.Fatalf("failed to delete API key: %v", err)
	}
	
	// Verify deleted
	_, err = GetAPIKey(service, account)
	if err == nil {
		t.Error("expected error for deleted key, got nil")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/config -v -run TestKeychain
```

Expected: FAIL with "undefined: SetAPIKey"

- [ ] **Step 3: Implement keychain functions**

```go
// internal/config/keychain.go
package config

import (
	"fmt"
	
	"github.com/zalando/go-keyring"
)

const (
	KeychainService = "mycli"
)

// SetAPIKey stores an API key in the system keychain
func SetAPIKey(service, account, secret string) error {
	err := keyring.Set(service, account, secret)
	if err != nil {
		return fmt.Errorf("failed to store API key in keychain: %w", err)
	}
	return nil
}

// GetAPIKey retrieves an API key from the system keychain
func GetAPIKey(service, account string) (string, error) {
	secret, err := keyring.Get(service, account)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve API key from keychain: %w", err)
	}
	return secret, nil
}

// DeleteAPIKey removes an API key from the system keychain
func DeleteAPIKey(service, account string) error {
	err := keyring.Delete(service, account)
	if err != nil {
		return fmt.Errorf("failed to delete API key from keychain: %w", err)
	}
	return nil
}

// GetOpenRouterKey retrieves the OpenRouter API key
func GetOpenRouterKey() (string, error) {
	return GetAPIKey(KeychainService, "openrouter")
}

// SetOpenRouterKey stores the OpenRouter API key
func SetOpenRouterKey(key string) error {
	return SetAPIKey(KeychainService, "openrouter", key)
}

// GetOllamaKey retrieves the Ollama API key (if needed)
func GetOllamaKey() (string, error) {
	return GetAPIKey(KeychainService, "ollama")
}

// SetOllamaKey stores the Ollama API key
func SetOllamaKey(key string) error {
	return SetAPIKey(KeychainService, "ollama", key)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/config -v -run TestKeychain
```

Expected: PASS (or SKIP if keychain not available)

- [ ] **Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat(config): add keychain integration for secure API key storage"
```

---

## Task 8: Cobra CLI Setup

**Files:**
- Create: `internal/cli/root.go`
- Create: `internal/cli/version.go`
- Modify: `cmd/mycli/main.go`

- [ ] **Step 1: Write root command**

```go
// internal/cli/root.go
package cli

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "Budget-friendly AI coding assistant",
	Long: `mycli is a terminal-based AI coding assistant that prioritizes 
cost efficiency through smart token management.

Features:
  - Chat-based AI assistance with modern TUI
  - 60-70% token savings through intelligent optimization
  - Support for budget-friendly models (DeepSeek-R1, StepFun)
  - LSP-powered codebase intelligence
  - Multiple code modification modes`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", 
		"config file (default: ~/.mycli/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, 
		"enable debug mode")
}
```

- [ ] **Step 2: Write version command**

```go
// internal/cli/version.go
package cli

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	GitCommit = "dev"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mycli version %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

- [ ] **Step 3: Update main.go to use Cobra**

```go
// cmd/mycli/main.go
package main

import (
	"fmt"
	"os"
	
	"github.com/yourusername/mycli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 4: Test CLI**

```bash
go build -o bin/mycli cmd/mycli/main.go
./bin/mycli --help
```

Expected output: Help text with description and available commands

- [ ] **Step 5: Test version command**

```bash
./bin/mycli version
```

Expected output:
```
mycli version 0.1.0
Git commit: dev
Build date: unknown
```

- [ ] **Step 6: Commit**

```bash
git add internal/cli/ cmd/mycli/
git commit -m "feat(cli): add Cobra CLI structure with root and version commands"
```

---
## Task 9: Config Commands

**Files:**
- Create: `internal/cli/config.go`

- [ ] **Step 1: Write config command structure**

```go
// internal/cli/config.go
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/yourusername/mycli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "View and modify mycli configuration settings",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		
		fmt.Println("Configuration:")
		fmt.Printf("  Provider: %s\n", getActiveProvider(cfg))
		fmt.Printf("  Model: %s\n", cfg.Providers.OpenRouter.DefaultModel)
		fmt.Printf("  Code Mode: %s\n", cfg.Behavior.CodeMode)
		fmt.Printf("  LSP: %v\n", cfg.LSP.Enabled)
		fmt.Printf("  Token Budget: %d\n", cfg.Behavior.TokenBudget)
		
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]
		
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		
		// Update config based on key
		if err := setConfigValue(cfg, key, value); err != nil {
			return err
		}
		
		// Save config
		cfgPath := getConfigPath()
		if err := cfg.Save(cfgPath); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		
		fmt.Printf("✓ %s updated to %s\n", key, value)
		return nil
	},
}

var configSetKeyCmd = &cobra.Command{
	Use:   "set-key <provider>",
	Short: "Set API key for a provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		
		fmt.Printf("Enter API key for %s: ", provider)
		var key string
		fmt.Scanln(&key)
		
		switch provider {
		case "openrouter":
			if err := config.SetOpenRouterKey(key); err != nil {
				return fmt.Errorf("failed to store API key: %w", err)
			}
		case "ollama":
			if err := config.SetOllamaKey(key); err != nil {
				return fmt.Errorf("failed to store API key: %w", err)
			}
		default:
			return fmt.Errorf("unknown provider: %s", provider)
		}
		
		fmt.Printf("✓ API key stored securely for %s\n", provider)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configSetKeyCmd)
}

func loadConfig() (*config.Config, error) {
	cfgPath := getConfigPath()
	
	// If config doesn't exist, return default
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return config.DefaultConfig(), nil
	}
	
	// Load config
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	// Apply env overrides
	config.ApplyEnvOverrides(cfg)
	
	return cfg, nil
}

func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	
	home, err := os.UserHomeDir()
	if err != nil {
		return ".mycli/config.yaml"
	}
	
	return filepath.Join(home, ".mycli", "config.yaml")
}

func getActiveProvider(cfg *config.Config) string {
	if cfg.Providers.OpenRouter.Enabled {
		return "openrouter"
	}
	if cfg.Providers.Ollama.Enabled {
		return "ollama"
	}
	return "none"
}

func setConfigValue(cfg *config.Config, key, value string) error {
	switch key {
	case "model":
		cfg.Providers.OpenRouter.DefaultModel = value
	case "code_mode":
		if value != "safe" && value != "interactive" && value != "auto" {
			return fmt.Errorf("invalid code_mode: %s (must be safe, interactive, or auto)", value)
		}
		cfg.Behavior.CodeMode = value
	case "theme":
		cfg.UI.Theme = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}
```

- [ ] **Step 2: Test config show command**

```bash
go build -o bin/mycli cmd/mycli/main.go
./bin/mycli config show
```

Expected output: Configuration values displayed

- [ ] **Step 3: Test config set command**

```bash
./bin/mycli config set model "stepfun/step-2-16k"
./bin/mycli config show
```

Expected: Model updated in output

- [ ] **Step 4: Commit**

```bash
git add internal/cli/config.go
git commit -m "feat(cli): add config commands (show, set, set-key)"
```

---

## Task 10: Chat Command Stub

**Files:**
- Create: `internal/cli/chat.go`

- [ ] **Step 1: Write chat command stub**

```go
// internal/cli/chat.go
package cli

import (
	"fmt"
	
	"github.com/spf13/cobra"
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
			_, err := config.GetOpenRouterKey()
			if err != nil {
				return fmt.Errorf("OpenRouter API key not found. Run: mycli config set-key openrouter")
			}
		}
		
		// TODO: Launch TUI (will be implemented in Plan D)
		fmt.Println("🤖 mycli chat")
		fmt.Println("Chat functionality will be implemented in the next phase.")
		fmt.Println("For now, configuration is ready!")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
```

- [ ] **Step 2: Add missing import to config.go**

```go
// internal/cli/config.go
import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/yourusername/mycli/internal/config"
)
```

- [ ] **Step 3: Test chat command**

```bash
go build -o bin/mycli cmd/mycli/main.go
./bin/mycli chat
```

Expected output: Stub message about chat functionality

- [ ] **Step 4: Commit**

```bash
git add internal/cli/chat.go
git commit -m "feat(cli): add chat command stub"
```

---
## Task 11: Integration Test

**Files:**
- Create: `tests/integration_test.go`

- [ ] **Step 1: Write integration test**

```go
// tests/integration_test.go
package tests

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/yourusername/mycli/internal/config"
)

func TestFullConfigFlow(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")
	
	// Step 1: Create default config
	cfg := config.DefaultConfig()
	
	// Step 2: Validate
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config validation failed: %v", err)
	}
	
	// Step 3: Modify config
	cfg.Providers.OpenRouter.DefaultModel = "test-model"
	cfg.Behavior.CodeMode = "auto"
	cfg.Behavior.TokenBudget = 80000
	
	// Step 4: Save config
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}
	
	// Step 5: Load config
	loaded, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	// Step 6: Verify loaded values
	if loaded.Providers.OpenRouter.DefaultModel != "test-model" {
		t.Errorf("expected test-model, got %s", loaded.Providers.OpenRouter.DefaultModel)
	}
	
	if loaded.Behavior.CodeMode != "auto" {
		t.Errorf("expected auto, got %s", loaded.Behavior.CodeMode)
	}
	
	if loaded.Behavior.TokenBudget != 80000 {
		t.Errorf("expected 80000, got %d", loaded.Behavior.TokenBudget)
	}
	
	// Step 7: Apply env overrides
	os.Setenv("MYCLI_MODEL", "env-model")
	defer os.Unsetenv("MYCLI_MODEL")
	
	config.ApplyEnvOverrides(loaded)
	
	if loaded.Providers.OpenRouter.DefaultModel != "env-model" {
		t.Errorf("env override failed: expected env-model, got %s", 
			loaded.Providers.OpenRouter.DefaultModel)
	}
	
	// Step 8: Validate modified config
	if err := loaded.Validate(); err != nil {
		t.Fatalf("modified config validation failed: %v", err)
	}
}
```

- [ ] **Step 2: Run integration test**

```bash
go test ./tests -v -run TestFullConfigFlow
```

Expected: PASS

- [ ] **Step 3: Run all tests**

```bash
go test ./... -v
```

Expected: All tests PASS

- [ ] **Step 4: Commit**

```bash
git add tests/
git commit -m "test: add integration test for config flow"
```

---

## Completion Checklist

- [x] Project initialized with Go modules
- [x] Dependencies added (Cobra, Viper, keyring, YAML)
- [x] Config struct with all required fields
- [x] Default configuration
- [x] Config load/save functionality
- [x] Config validation
- [x] Environment variable overrides
- [x] Keychain integration for secure API key storage
- [x] Cobra CLI structure (root, version commands)
- [x] Config commands (show, set, set-key)
- [x] Chat command stub
- [x] Integration tests
- [x] All tests passing

## Manual Verification

```bash
# Build
go build -o bin/mycli cmd/mycli/main.go

# Test commands
./bin/mycli --help
./bin/mycli version
./bin/mycli config show
./bin/mycli config set model "deepseek/deepseek-r1"
./bin/mycli config set-key openrouter
./bin/mycli chat

# Run tests
go test ./... -v -cover
```

## Next Steps

This plan establishes the foundation. Next plans will build on this:
- **Plan B**: Python AI Layer (OpenRouter/Ollama clients)
- **Plan C**: Go-Python Bridge (subprocess communication)
- **Plan D**: Basic TUI (Bubble Tea interface)

---

## Notes for Implementation

**Key Design Decisions:**
1. Config stored in `~/.mycli/config.yaml` by default
2. API keys stored in OS keychain, never in plain text
3. Environment variables override config file values
4. Validation happens on load to catch errors early
5. All config operations are testable without external dependencies

**Testing Strategy:**
- Unit tests for each component
- Integration test for full config flow
- Keychain tests skip if keychain unavailable
- Use temp directories for file operations

**Common Issues:**
- Keychain may not be available in CI/CD (tests will skip)
- Windows paths use backslashes (use `filepath.Join`)
- YAML indentation matters (use consistent spaces)
