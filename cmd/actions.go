package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rezatg/gitc/internal/ai"
	"github.com/rezatg/gitc/internal/ai/generic"
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

// ConfigureAI constructs the AI configuration with proper validation
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
		URL:              c.String("url"),
	}

	// Apply config defaults
	a.applyConfigDefaults(cfg)

	// Validate required fields
	if err := a.validateConfig(a.config); err != nil {
		return nil, fmt.Errorf("invalid AI configuration: %w", err)
	}

	return cfg, nil
}

// generateCommitMessage generates a commit message based on git diff and AI configuration.
func (a *App) generateCommitMessage(ctx context.Context, diff string, cfg *ai.Config) (string, error) {
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
		MaxRedirects:     cfg.MaxRedirects,
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
	lines := strings.Split(msg, "\n")

	if len(lines) == 1 {
		fmt.Printf("   git commit -m \"%s\"\n", strings.TrimSpace(lines[0]))
	} else {
		fmt.Printf("   git commit -m \"%s\" \\\n", strings.TrimSpace(lines[0]))
		for i, line := range lines[1:] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if i < len(lines[1:])-1 {
				fmt.Printf("      -m \"%s\" \\\n", line)
			} else {
				fmt.Printf("      -m \"%s\"\n", line)
			}
		}
	}
	return nil
}

// ConfigAction handles configuration updates
func (a *App) ConfigAction(c *cli.Context) error {
	cfg := *a.config
	a.updateConfigFromFlags(&cfg, c)

	if err := a.validateConfig(&cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if err := config.Save(&cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

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
		switch cfg.Provider {
		case "openai":
			cfg.Model = "gpt-4o-mini"
		case "grok":
			cfg.Model = "grok-3"
		case "deepseek":
			cfg.Model = "deepseek-rag"
		default:
			cfg.Model = a.config.Model
		}
	}
	if cfg.APIKey == "" {
		cfg.APIKey = a.config.APIKey
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
	if cfg.URL == "" {
		switch cfg.Provider {
		case "openai":
			cfg.URL = "https://api.openai.com/v1/chat/completions"
		case "grok":
			cfg.URL = "https://api.x.ai/v1/chat/completions"
		case "deepseek":
			cfg.URL = "https://api.deepseek.com/v1/chat/completions" // فرضی
		default:
			cfg.URL = a.config.URL
		}
	}
}

// initAIProvider initializes the AI provider based on the configuration.
func (a *App) initAIProvider(cfg *ai.Config) (ai.AIProvider, error) {
	return generic.NewGenericProvider(cfg.APIKey, cfg.Proxy, cfg.URL, cfg.Provider)
}

// validateConfig checks if the AI configuration is valid.
func (a *App) validateConfig(cfg *config.Config) error {
	if cfg.Provider == "" {
		return fmt.Errorf("AI provider is required")
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if cfg.MaxLength <= 0 {
		return fmt.Errorf("max length must be positive")
	}
	return nil
}

// updateConfigFromFlags updates the configuration based on CLI flags.
func (a *App) updateConfigFromFlags(cfg *config.Config, c *cli.Context) {
	if provider := c.String("provider"); provider != "" {
		cfg.Provider = provider
	}
	if model := c.String("model"); model != "" {
		cfg.Model = model
	}
	if apiKey := c.String("api-key"); apiKey != "" {
		cfg.APIKey = apiKey
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
	if maxRedirects := c.Int("max-redirects"); maxRedirects != 0 {
		cfg.MaxRedirects = maxRedirects
	}
	if url := c.String("url"); url != "" {
		cfg.URL = url
	}
}
