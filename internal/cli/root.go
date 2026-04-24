// internal/cli/root.go
package cli

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "Budget-friendly AI coding assistant",
	Long: `mycli is a terminal-based AI coding assistant that prioritizes
cost efficiency through smart token management.

Features:
  - Chat-based AI assistance with modern TUI
  - 60-70% token savings through intelligent optimization
  - Support for budget-friendly models (DeepSeek-R1, StepFun)
  - LSP-powered codebase intelligence
  - Multiple code modification modes`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default: ~/.mycli/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"enable debug mode")
}
