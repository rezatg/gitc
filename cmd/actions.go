package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rezatg/gitc/internal/ai"
	"github.com/rezatg/gitc/internal/ai/openai"
	"github.com/rezatg/gitc/internal/git"
	"github.com/rezatg/gitc/pkg/config"
	"github.com/rezatg/gitc/pkg/utils"
	"github.com/urfave/cli/v2"
)

// App encapsulates the application logic and dependencies
type App struct {
	gitService git.GitService
	config     *config.Config
}

// NewApp creates a new App instance
func NewApp(gitService git.GitService, config *config.Config) *App {
	return &App{
		gitService: gitService,
		config:     config,
	}
}

// buildAIConfig constructs the AI configuration with proper validation
func (a *App) buildAIConfig(c *cli.Context) (*ai.Config, error) {
	cfg := &ai.Config{
		Provider:         c.String("provider"),
		Model:            c.String("model"),
		APIKey:           c.String("api-key"),
		Timeout:          time.Duration(c.Int("timeout")) * time.Second,
		MaxLength:        c.Int("max-length"),
		Language:         c.String("lang"),
		MaxRedirects:     c.Int("max-redirects"),
		Proxy:            c.String("proxy"),
		CommitType:       c.String("commit-type"),
		CustomConvention: c.String("custom-convention"),
		UseGitmoji:       c.Bool("emoji"),
	}

	// Apply config defaults
	a.applyConfigDefaults(cfg)

	// Validate configuration
	if err := a.validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (a *App) applyConfigDefaults(cfg *ai.Config) {
	if cfg.Provider == "" {
		cfg.Provider = a.config.Provider
	}
	if cfg.Model == "" {
		cfg.Model = a.config.OpenAI.Model
	}
	if cfg.APIKey == "" {
		cfg.APIKey = a.config.OpenAI.APIKey
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = time.Duration(a.config.Timeout) * time.Second
	}
	if cfg.MaxLength == 0 {
		cfg.MaxLength = a.config.MaxLength
	}
	if cfg.Language == "" {
		cfg.Language = a.config.Language
	}
	if cfg.MaxRedirects == 0 {
		cfg.MaxRedirects = a.config.MaxRedirects
	}
	if !cfg.UseGitmoji {
		cfg.UseGitmoji = a.config.UseGitmoji
	}
}

func (a *App) validateConfig(cfg *ai.Config) error {
	switch {
	case cfg.Provider == "":
		return errors.New("AI provider must be specified")
	case cfg.APIKey == "":
		return errors.New("API key is required")
	case cfg.Timeout <= 0:
		return errors.New("timeout must be positive")
	case cfg.MaxLength <= 0:
		return errors.New("max length must be positive")
	}
	return nil
}

// CommitAction handles the generation of commit messages
func (a *App) CommitAction(c *cli.Context) error {
	// Load git diff
	diff, err := a.gitService.GetDiff(context.Background())
	if err != nil {
		return fmt.Errorf("❌ failed to get git diff: %v", err)
	} else if diff == "" {
		return fmt.Errorf("❌ nothing staged for commit")
	}

	// Build AI configuration
	aiConfig, err := a.buildAIConfig(c)
	if err != nil {
		return fmt.Errorf("failed to build AI config: %w", err)
	}

	// Initialize AI provider
	var provider ai.AIProvider
	switch aiConfig.Provider {
	case "openai":
		provider, err = openai.NewOpenAIProvider(aiConfig.APIKey, aiConfig.Proxy, a.config.OpenAI.URL)
		if err != nil {
			return fmt.Errorf("failed to initialize OpenAI provider: %w", err)
		}
	default:
		return fmt.Errorf("unsupported AI provider: %s", aiConfig.Provider)
	}

	// Generate commit message
	ctx, cancel := context.WithTimeout(context.Background(), aiConfig.Timeout)
	defer cancel()

	msg, err := provider.GenerateCommitMessage(
		ctx, diff,
		ai.MessageOptions{
			Model:            aiConfig.Model,
			Language:         aiConfig.Language,
			CommitType:       c.String("commit-type"),
			CustomConvention: aiConfig.CustomConvention,
			MaxLength:        aiConfig.MaxLength,
			MaxRedirects:     aiConfig.MaxRedirects,
		})
	if err != nil {
		return fmt.Errorf("❌ failed to generate commit message: %v", err)
	}

	// Apply Gitmoji if enabled
	if c.Bool("emoji") || a.config.UseGitmoji {
		msg = utils.AddGitmojiToCommitMessage(msg)
	}

	fmt.Println("✅ Commit message generated. You can now run:")
	fmt.Printf("   git commit -m \"%s\"\n", strings.ReplaceAll(msg, "\n", ""))
	return nil
}

// ConfigAction handles configuration updates
func (a *App) ConfigAction(c *cli.Context) error {
	// Create a copy of the current config
	cfg := *a.config

	// Update fields based on CLI arguments
	if provider := c.String("provider"); provider != "" {
		cfg.Provider = provider
	}
	if model := c.String("model"); model != "" {
		cfg.OpenAI.Model = model
	}
	if apiKey := c.String("api-key"); apiKey != "" {
		cfg.OpenAI.APIKey = apiKey
	}
	if lang := c.String("lang"); lang != "" {
		cfg.Language = lang
	}
	if timeout := c.Int("timeout"); timeout != 0 {
		cfg.Timeout = timeout
	}
	if maxLength := c.Int("maxLength"); maxLength != 0 {
		cfg.MaxLength = maxLength
	}
	if proxy := c.String("proxy"); proxy != "" {
		cfg.Proxy = proxy
	}
	if commitType := c.String("commit-type"); commitType != "" {
		cfg.CommitType = commitType
	}
	if customConvention := c.String("custom-convention"); customConvention != "" {
		cfg.CustomConvention = customConvention
	}
	if c.IsSet("emoji") {
		cfg.UseGitmoji = c.Bool("emoji")
	}
	if maxRedirects := c.Int("max-redirects"); maxRedirects != 0 {
		cfg.MaxRedirects = maxRedirects
	}

	// Save updated config
	if err := config.Save(&cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Update the app's config
	a.config = &cfg
	fmt.Println("✅ Configuration updated successfully")
	return nil
}
