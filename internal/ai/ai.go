package ai

import (
	"context"
	"time"
)

// AIProvider defines the interface for AI providers
type AIProvider interface {
	GenerateCommitMessage(opts Options) (string, error)
}

// Config holds AI provider configuration
type Config struct {
	Ai        string
	APIKey    string
	Timeout   time.Duration
	MaxLength int
	Model     string
	Language  string
}

type Options struct {
	Context          context.Context
	Diff             string
	Model            string
	Language         string
	CommitType       string
	CustomConvention string
	MaxLength        int
}
