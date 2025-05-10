package openai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/rezatg/gitc/internal/ai"
	"github.com/rezatg/gitc/pkg/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

const (
	defaultOpenAIURL = "https://api.openai.com/v1/chat/completions"
	systemPrompt     = "You are an AI assistant that generates Git commit messages."
)

// OpenAIProvider implements the AIProvider interface for OpenAI
type OpenAIProvider struct {
	apiKey string
	client *fasthttp.Client
	url    string
}

type Request struct {
	Model string `json:"model"`
	// Store       bool      `json:"store"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey, proxy, url string) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}
	if url == "" {
		url = defaultOpenAIURL
	}

	client := &fasthttp.Client{
		MaxConnsPerHost: 10,
	}

	// Configure proxy if provided
	if proxy != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
	}

	return &OpenAIProvider{
		apiKey: apiKey,
		client: client,
		url:    url,
	}, nil
}

// GenerateCommitMessage generates a commit message using OpenAI API
func (p *OpenAIProvider) GenerateCommitMessage(ctx context.Context, diff string, opts ai.MessageOptions) (string, error) {
	prompt := utils.GetPromptForSingleCommit(diff, opts.CommitType, opts.CustomConvention, opts.Language)
	reqBody := Request{
		Model: opts.Model,
		// Store: false,
		Messages: []Message{
			{"system", systemPrompt},
			{"user", prompt},
		},
		MaxTokens:   max(512, opts.MaxLength), // More tokens for complete messages
		Temperature: 0.7,                      // Slightly creative but controlled
	}

	jsonData, err := sonic.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode JSON: %v", err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(p.url)
	req.Header.SetMethod("POST")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(jsonData)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := p.client.DoRedirects(req, resp, opts.MaxRedirects); err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	var res Response
	if err = sonic.Unmarshal(resp.Body(), &res); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if statusCode := resp.StatusCode(); statusCode != fasthttp.StatusOK {
		if res.Error.Message != "" {
			return "", fmt.Errorf("API error [%d]: %s", statusCode, res.Error.Message)
		}

		return "", fmt.Errorf("API returned status %d: %s", statusCode, resp.Body())
	}

	if res.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", res.Error.Message)
	} else if len(res.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	commitMessage := strings.TrimSpace(res.Choices[0].Message.Content)

	// Ensure message is not empty
	if commitMessage == "" {
		return "", errors.New("empty commit message generated")
	}

	// if len(commitMessage) > opts.MaxLength {
	// 	commitMessage = commitMessage[:opts.MaxLength]
	// }

	return commitMessage, nil
}
