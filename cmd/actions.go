package cmd

import (
	"context"
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
func (a *App) ConfigureAI(c *cli.Context) (*ai.Config, error) {
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
		UseGitmoji:       !c.Bool("no-emoji") && c.Bool("emoji"),
	}

	// Apply config defaults
	a.applyConfigDefaults(cfg)

	// Validate required fields
	if err := a.validateConfig(a.config); err != nil {
		return nil, fmt.Errorf("invalid AI configuration: %w", err)
	}

	return cfg, nil
}

// GenerateCommitMessage generates a commit message based on git diff and AI configuration.
func (a *App) generateCommitMessage(ctx context.Context, diff string, cfg *ai.Config) (string, error) {
	// Initialize AI provider
	provider, err := a.initAIProvider(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	opts := ai.MessageOptions{
		Model:            cfg.Model,
		Language:         cfg.Language,
		CommitType:       cfg.CommitType,
		CustomConvention: cfg.CustomConvention,
		MaxLength:        cfg.MaxLength,
	}

	msg, err := provider.GenerateCommitMessage(ctx, diff, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Apply Gitmoji if enabled
	if cfg.UseGitmoji {
		msg = utils.AddGitmojiToCommitMessage(msg)
	}

	return msg, nil
}

// CommitAction handles the generation of commit messages
func (a *App) CommitAction(c *cli.Context) error {
	// Stage all changes if --all (-a) flag is set
	if c.Bool("all") {
		if err := a.gitService.StageAll(c.Context); err != nil {
			return fmt.Errorf("❌ failed to stage changes: %v", err)
		}

		fmt.Println("✅ All changes staged successfully")
	}
	// Fetch git diff for staged changes
	diff, err := a.gitService.GetDiff(c.Context)
	if err != nil {
		return fmt.Errorf("❌ failed to get git diff: %v", err)
	} else if diff == "" {
		return fmt.Errorf("❌ nothing staged for commit")
	}

	// Configure AI settings
	cfg, err := a.ConfigureAI(c)
	if err != nil {
		return fmt.Errorf("❌ failed to build AI config: %w", err)
	}

	// Generate commit message
	msg, err := a.generateCommitMessage(c.Context, diff, cfg)
	if err != nil {
		return fmt.Errorf("❌ failed to generate commit message: %w", err)
	}

	fmt.Println("✅ Commit message generated. You can now run:")
	fmt.Printf("   git commit -m \"%s\"\n", strings.ReplaceAll(msg, "\n", ""))
	return nil
}

// ConfigAction handles configuration updates
func (a *App) ConfigAction(c *cli.Context) error {
	// Create a copy of the current config
	cfg := *a.config

	// Update config fields from CLI flags
	a.updateConfigFromFlags(&cfg, c)

	// Validate updated config
	if err := a.validateConfig(&cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Save updated config
	if err := config.Save(&cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Update app's config
	a.config = &cfg
	fmt.Println("✅ Configuration updated successfully")
	return nil
}

// applyConfigDefaults sets default values for unset AI configuration fields.
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
}

// validateConfig checks if the application configuration is valid.
func (a *App) validateConfig(cfg *config.Config) error {
	switch {
	case cfg.Provider == "":
		return fmt.Errorf("AI provider is required")
	case cfg.OpenAI.APIKey == "":
		return fmt.Errorf("API key is required")
	case cfg.Timeout <= 0:
		return fmt.Errorf("timeout must be positive")
	case cfg.MaxLength <= 0:
		return fmt.Errorf("max length must be positive")
	}
	return nil
}

// initAIProvider initializes the AI provider based on the configuration.
func (a *App) initAIProvider(cfg *ai.Config) (ai.AIProvider, error) {
	switch cfg.Provider {
	case "openai":
		return openai.NewOpenAIProvider(cfg.APIKey, cfg.Proxy, a.config.OpenAI.URL)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.Provider)
	}
}

// updateConfigFromFlags updates the configuration based on CLI flags..
func (a *App) updateConfigFromFlags(cfg *config.Config, c *cli.Context) {
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
	if c.IsSet("no-emoji") {
		cfg.UseGitmoji = !c.Bool("no-emoji")
	} else if c.IsSet("emoji") {
		cfg.UseGitmoji = c.Bool("emoji")
	}
	if maxRedirects := c.Int("max-redirects"); c.Int("max-redirects") != 0 {
		cfg.MaxRedirects = maxRedirects
	}
}
