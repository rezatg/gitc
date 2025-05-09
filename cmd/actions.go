package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
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

// buildAIConfig constructs the AI configuration from CLI flags and defaults
func (a *App) buildAIConfig(c *cli.Context) (ai.Config, error) {
	cfg := ai.Config{
		Ai:           c.String("ai"),
		Model:        c.String("model"),
		APIKey:       c.String("api-key"),
		Timeout:      time.Duration(c.Int("timeout")) * time.Second,
		MaxLength:    c.Int("maxLength"),
		Language:     c.String("lang"),
		MaxRedirects: c.Int("max-redirects"),
	}

	// Apply defaults from config if not provided
	if cfg.Ai == "" {
		cfg.Ai = a.config.AI
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

	// Handle proxy
	if proxy := c.String("proxy"); proxy != "" {
		cfg.Proxy = proxy
	}

	// Handle custom convention
	if customConvention := c.String("custom-convention"); customConvention != "" {
		var tmp string
		if err := sonic.Unmarshal([]byte(customConvention), &tmp); err != nil {
			return ai.Config{}, fmt.Errorf("invalid JSON for custom-convention: %w", err)
		}
		cfg.CustomConvention = tmp
	}

	// Validate required fields
	if cfg.Ai == "" {
		return ai.Config{}, fmt.Errorf("AI provider must be specified (use --ai or set in config)")
	}
	if cfg.APIKey == "" {
		return ai.Config{}, fmt.Errorf("API key must be specified (use --api-key or set in config)")
	}

	return cfg, nil
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
	switch aiConfig.Ai {
	case "openai":
		provider, err = openai.NewOpenAIProvider(aiConfig.APIKey, aiConfig.Proxy, a.config.OpenAI.URL)
		if err != nil {
			return fmt.Errorf("failed to initialize OpenAI provider: %w", err)
		}
	default:
		return fmt.Errorf("unsupported AI provider: %s", aiConfig.Ai)
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
	if ai := c.String("ai"); ai != "" {
		cfg.AI = ai
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
