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
