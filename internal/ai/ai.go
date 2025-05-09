package ai

import (
	"context"
	"time"
)

// AIProvider defines the interface for AI providers
type AIProvider interface {
	GenerateCommitMessage(ctx context.Context, diff string, opts MessageOptions) (string, error)
}

// Config holds AI provider configuration
type Config struct {
	Ai               string
	APIKey           string
	Timeout          time.Duration
	MaxLength        int
	Model            string
	Language         string
	CustomConvention string
	MaxRedirects     int

	Proxy string
}

type MessageOptions struct {
	Model            string
	Language         string
	CommitType       string
	CustomConvention string
	MaxLength        int
	MaxRedirects     int
}
