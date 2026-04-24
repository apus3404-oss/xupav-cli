// internal/config/config.go
package config

type Config struct {
	Version   string              `yaml:"version"`
	Providers ProvidersConfig     `yaml:"providers"`
	Behavior  BehaviorConfig      `yaml:"behavior"`
	UI        UIConfig            `yaml:"ui"`
	LSP       LSPConfig           `yaml:"lsp"`
	Cache     CacheConfig         `yaml:"cache"`
	Logging   LoggingConfig       `yaml:"logging"`
}

type ProvidersConfig struct {
	OpenRouter OpenRouterConfig `yaml:"openrouter"`
	Ollama     OllamaConfig     `yaml:"ollama"`
}

type OpenRouterConfig struct {
	Enabled      bool    `yaml:"enabled"`
	DefaultModel string  `yaml:"default_model"`
	APIKeySource string  `yaml:"api_key_source"`
	BaseURL      string  `yaml:"base_url"`
	MaxTokens    int     `yaml:"max_tokens"`
	Temperature  float64 `yaml:"temperature"`
	Timeout      int     `yaml:"timeout"`
}

type OllamaConfig struct {
	Enabled      bool   `yaml:"enabled"`
	BaseURL      string `yaml:"base_url"`
	DefaultModel string `yaml:"default_model"`
	Fallback     bool   `yaml:"fallback"`
	Timeout      int    `yaml:"timeout"`
}

type BehaviorConfig struct {
	CodeMode           string `yaml:"code_mode"`
	AutoAddFiles       bool   `yaml:"auto_add_files"`
	MaxContextFiles    int    `yaml:"max_context_files"`
	TokenBudget        int    `yaml:"token_budget"`
	ConversationMemory int    `yaml:"conversation_memory"`
}

type UIConfig struct {
	Theme           string `yaml:"theme"`
	Animations      bool   `yaml:"animations"`
	SyntaxHighlight bool   `yaml:"syntax_highlight"`
	ShowTokenCount  bool   `yaml:"show_token_count"`
	ShowCost        bool   `yaml:"show_cost"`
	ColorScheme     string `yaml:"color_scheme"`
}

type LSPConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Languages  []string `yaml:"languages"`
	Timeout    int      `yaml:"timeout"`
	MaxQueries int      `yaml:"max_queries"`
}

type CacheConfig struct {
	Enabled  bool   `yaml:"enabled"`
	TTL      int    `yaml:"ttl"`
	MaxSize  string `yaml:"max_size"`
	Location string `yaml:"location"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    string `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
}
