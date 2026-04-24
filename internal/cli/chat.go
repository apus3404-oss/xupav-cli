// internal/cli/chat.go
package cli

import (
	"fmt"

	"github.com/apus3404-oss/xupav-cli/internal/config"
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
