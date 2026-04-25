// internal/cli/chat.go
package cli

import (
	"fmt"

	"github.com/apus3404-oss/xupav-cli/internal/bridge"
	"github.com/apus3404-oss/xupav-cli/internal/config"
	"github.com/apus3404-oss/xupav-cli/internal/tui"
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

		// Create Python bridge
		pythonBridge := bridge.NewPythonBridge(bridge.PythonConfig{
			PythonPath: "python",
			Config: map[string]interface{}{
				"provider": "openrouter",
				"api_key":  "test-key", // Will be loaded from keychain in future
			},
		})

		// Start bridge
		if err := pythonBridge.Start(); err != nil {
			return fmt.Errorf("failed to start Python bridge: %w", err)
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
