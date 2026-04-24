// tests/integration_test.go
package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apus3404-oss/xupav-cli/internal/config"
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
