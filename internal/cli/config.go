// internal/cli/config.go
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/apus3404-oss/xupav-cli/internal/config"
	"github.com/spf13/cobra"
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
