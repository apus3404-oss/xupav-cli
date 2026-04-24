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
