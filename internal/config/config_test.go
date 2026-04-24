// internal/config/config_test.go
package config

import (
	"os"
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
