package main

import (
	"context"
	"fmt"
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

// CommitAction handles the generation of commit messages
func (a *App) CommitAction(c *cli.Context) error {
	// Load git diff
	diff, err := a.gitService.GetDiff(context.Background())
	if err != nil {
		return fmt.Errorf("❌ failed to get git diff: %v", err)
	}
	if diff == "" {
		return fmt.Errorf("❌ nothing staged for commit")
	}

	// Prepare AI configuration from CLI flags or config
	aiConfig := ai.Config{
		Ai:           c.String("ai"),
		Model:        c.String("model"),
		APIKey:       c.String("api-key"),
		Timeout:      time.Duration(c.Int("timeout")) * time.Second,
		MaxLength:    c.Int("maxLength"),
		Language:     c.String("lang"),
		MaxRedirects: c.Int("max-redirects"),
	}

	// Apply defaults from config if not provided via CLI
	if aiConfig.Ai == "" {
		aiConfig.Ai = a.config.AI
	}
	if aiConfig.Model == "" {
		aiConfig.Model = a.config.OpenAI.Model
	}
	if aiConfig.APIKey == "" {
		aiConfig.APIKey = a.config.OpenAI.APIKey
	}
	if aiConfig.Timeout == 0 {
		aiConfig.Timeout = time.Duration(a.config.Timeout) * time.Second
	}
	if aiConfig.MaxLength == 0 {
		aiConfig.MaxLength = a.config.MaxLength
	}
	if aiConfig.Language == "" {
		aiConfig.Language = a.config.Language
	}
	if aiConfig.MaxRedirects == 0 {
		aiConfig.MaxRedirects = a.config.MaxRedirects
	}

	if proxy := c.String("proxy"); proxy != "" {
		a.config.Proxy = proxy
	}

	// Validate required fields
	if aiConfig.Ai == "" {
		return fmt.Errorf("❌ AI provider must be specified (use --ai or set in config)")
	}
	if aiConfig.APIKey == "" {
		return fmt.Errorf("❌ API key must be specified (use --api-key or set in config)")
	}

	useGitmoji := c.Bool("emoji") || a.config.UseGitmoji

	// Initialize AI provider
	var provider ai.AIProvider
	switch aiConfig.Ai {
	case "openai":
		provider, err = openai.NewOpenAIProvider(aiConfig.APIKey, a.config.Proxy)
		if err != nil {
			return fmt.Errorf("❌ failed to initialize OpenAI provider: %v", err)
		}
	default:
		return fmt.Errorf("❌ unsupported AI provider: %s", aiConfig.Ai)
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
			CustomConvention: c.String("custom-convention"),
			MaxLength:        aiConfig.MaxLength,
			MaxRedirects:     aiConfig.MaxRedirects,
		})
	if err != nil {
		return fmt.Errorf("❌ failed to generate commit message: %v", err)
	}

	if useGitmoji {
		msg = utils.AddGitmojiToCommitMessage(msg)
	}

	fmt.Println("✅ Commit message generated. You can now run:")
	fmt.Printf("   git commit -m %q\n", msg)
	return nil
}

// ConfigAction handles configuration updates
func (a *App) ConfigAction(c *cli.Context) error {
	// Load existing config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("❌ failed to load config: %v", err)
	}

	// Update fields based on CLI arguments
	if c.String("ai") != "" {
		cfg.AI = c.String("ai")
	}
	if c.String("model") != "" {
		cfg.OpenAI.Model = c.String("model")
	}
	if c.String("api-key") != "" {
		cfg.OpenAI.APIKey = c.String("api-key")
	}
	if c.String("lang") != "" {
		cfg.Language = c.String("lang")
	}
	if c.Int("timeout") != 0 {
		cfg.Timeout = c.Int("timeout")
	}
	if c.Int("maxLength") != 0 {
		cfg.MaxLength = c.Int("maxLength")
	}
	if c.String("proxy") != "" {
		cfg.Proxy = c.String("proxy")
	}
	if c.String("commit-type") != "" {
		cfg.CommitType = c.String("commit-type")
	}
	if c.String("custom-convention") != "" {
		cfg.CustomConvention = c.String("custom-convention")
	}
	if c.IsSet("emoji") {
		cfg.UseGitmoji = c.Bool("emoji")
	}
	if c.Int("max-redirects") != 0 {
		cfg.MaxRedirects = c.Int("max-redirects")
	}

	// Save updated config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("❌ failed to save config: %v", err)
	}

	fmt.Println("✅ Configuration updated successfully")
	return nil
}
