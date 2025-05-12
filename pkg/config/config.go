package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytedance/sonic"
)

// Config holds application configuration
type Config struct {
	Provider         string `json:"provider"`
	MaxLength        int    `json:"max_length"`
	Proxy            string `json:"proxy"`
	Language         string `json:"language"`
	Timeout          int    `json:"timeout"`
	CommitType       string `json:"commit_type"`
	CustomConvention string `json:"custom-convention"`
	UseGitmoji       bool   `json:"use_gitmoji"`
	MaxRedirects     int    `json:"max_redirects"`

	OpenAI OpenAI `json:"open_ai"`
}

// OpenAI holds OpenAI-specific configuration
type OpenAI struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
	URL    string `json:"url"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Provider:         "openai",
		MaxLength:        200,
		Proxy:            "",
		Language:         "en",
		Timeout:          10,
		CommitType:       "",
		CustomConvention: "",
		UseGitmoji:       false,
		MaxRedirects:     5,

		OpenAI: OpenAI{
			APIKey: os.Getenv("AI_API_KEY"), // Load from env if available
			Model:  "gpt-4o-mini",
			URL:    "",
		},
	}
}

// configPath is set to config.json in the project root
var configPath = "./config.json"

// SetConfigPath updates the configuration file path
func SetConfigPath(path string) {
	configPath = path
}

// Load loads the configuration from file or creates a default one if it doesn't exist
func Load() (*Config, error) {
	// Resolve absolute path to ensure consistency
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if os.IsNotExist(err) {
		// Create default config file if it doesn't exist
		cfg := DefaultConfig()
		if err := Save(cfg); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}

		return cfg, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := sonic.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply defaults for unset fields
	defaults := DefaultConfig()
	if cfg.Provider == "" {
		cfg.Provider = defaults.Provider
	}
	if cfg.MaxLength == 0 {
		cfg.MaxLength = defaults.MaxLength
	}
	if cfg.Language == "" {
		cfg.Language = defaults.Language
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaults.Timeout
	}
	if cfg.OpenAI.Model == "" {
		cfg.OpenAI.Model = defaults.OpenAI.Model
	}
	if cfg.OpenAI.APIKey == "" {
		cfg.OpenAI.APIKey = defaults.OpenAI.APIKey
	}

	return &cfg, nil
}

// Save saves the configuration to file
func Save(cfg *Config) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := sonic.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	} else if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func Reset() error {
	defaultConfig := DefaultConfig()
	return Save(defaultConfig)
}
