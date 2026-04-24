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
