// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
	Enabled      bool    `yaml:"enabled"`
	DefaultModel string  `yaml:"default_model"`
	APIKeySource string  `yaml:"api_key_source"`
	BaseURL      string  `yaml:"base_url"`
	MaxTokens    int     `yaml:"max_tokens"`
	Temperature  float64 `yaml:"temperature"`
	Timeout      int     `yaml:"timeout"`
}

type OllamaConfig struct {
	Enabled      bool   `yaml:"enabled"`
	BaseURL      string `yaml:"base_url"`
	DefaultModel string `yaml:"default_model"`
	Fallback     bool   `yaml:"fallback"`
	Timeout      int    `yaml:"timeout"`
}

type BehaviorConfig struct {
	CodeMode           string `yaml:"code_mode"`
	AutoAddFiles       bool   `yaml:"auto_add_files"`
	MaxContextFiles    int    `yaml:"max_context_files"`
	TokenBudget        int    `yaml:"token_budget"`
	ConversationMemory int    `yaml:"conversation_memory"`
}

type UIConfig struct {
	Theme           string `yaml:"theme"`
	Animations      bool   `yaml:"animations"`
	SyntaxHighlight bool   `yaml:"syntax_highlight"`
	ShowTokenCount  bool   `yaml:"show_token_count"`
	ShowCost        bool   `yaml:"show_cost"`
	ColorScheme     string `yaml:"color_scheme"`
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
