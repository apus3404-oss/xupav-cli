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
					CodeMode:           "interactive",
					TokenBudget:        500,
					MaxContextFiles:    10,
					ConversationMemory: 10,
				},
				Providers: ProvidersConfig{
					OpenRouter: OpenRouterConfig{
						Enabled:      true,
						DefaultModel: "test",
						MaxTokens:    1000,
					},
				},
			},
			wantErr: true,
			errMsg:  "token_budget must be >= 1000",
		},
		{
			name: "invalid code mode",
			cfg: &Config{
				Behavior: BehaviorConfig{
					CodeMode:           "invalid",
					TokenBudget:        60000,
					MaxContextFiles:    10,
					ConversationMemory: 10,
				},
				Providers: ProvidersConfig{
					OpenRouter: OpenRouterConfig{
						Enabled:      true,
						DefaultModel: "test",
						MaxTokens:    1000,
					},
				},
			},
			wantErr: true,
			errMsg:  "code_mode must be one of: safe, interactive, auto",
		},
		{
			name: "max context files too low",
			cfg: &Config{
				Behavior: BehaviorConfig{
					CodeMode:           "interactive",
					TokenBudget:        60000,
					MaxContextFiles:    0,
					ConversationMemory: 10,
				},
				Providers: ProvidersConfig{
					OpenRouter: OpenRouterConfig{
						Enabled:      true,
						DefaultModel: "test",
						MaxTokens:    1000,
					},
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
				} else if !containsSubstring(err.Error(), tt.errMsg) {
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

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)+1 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

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
