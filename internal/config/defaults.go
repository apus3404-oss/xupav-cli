// internal/config/defaults.go
package config

func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Providers: ProvidersConfig{
			OpenRouter: OpenRouterConfig{
				Enabled:      true,
				DefaultModel: "deepseek/deepseek-r1",
				APIKeySource: "keychain",
				BaseURL:      "https://openrouter.ai/api/v1",
				MaxTokens:    4096,
				Temperature:  0.7,
				Timeout:      60000,
			},
			Ollama: OllamaConfig{
				Enabled:      true,
				BaseURL:      "http://localhost:11434",
				DefaultModel: "codellama:13b",
				Fallback:     true,
				Timeout:      120000,
			},
		},
		Behavior: BehaviorConfig{
			CodeMode:           "interactive",
			AutoAddFiles:       true,
			MaxContextFiles:    10,
			TokenBudget:        60000,
			ConversationMemory: 10,
		},
		UI: UIConfig{
			Theme:           "gradient",
			Animations:      true,
			SyntaxHighlight: true,
			ShowTokenCount:  true,
			ShowCost:        true,
			ColorScheme:     "auto",
		},
		LSP: LSPConfig{
			Enabled:    true,
			Languages:  []string{"go", "python", "javascript", "typescript", "rust"},
			Timeout:    5000,
			MaxQueries: 5,
		},
		Cache: CacheConfig{
			Enabled:  true,
			TTL:      3600,
			MaxSize:  "100MB",
			Location: "~/.mycli/cache",
		},
		Logging: LoggingConfig{
			Level:      "info",
			File:       "~/.mycli/logs/mycli.log",
			MaxSize:    "10MB",
			MaxBackups: 3,
		},
	}
}
