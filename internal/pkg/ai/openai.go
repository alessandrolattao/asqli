package ai

import (
	"context"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClientInterface defines the interface for interacting with OpenAI's API
type OpenAIClientInterface interface {
	CreateCompletion(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// OpenAIClient implements the OpenAIClientInterface
type OpenAIClient struct {
	client *openai.Client
}

// NewClient creates a new AI client
func NewClient() (*Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is not set")
	}

	openaiClient := &OpenAIClient{
		client: openai.NewClient(apiKey),
	}

	return &Client{
		openaiClient: openaiClient,
	}, nil
}

// CreateCompletion sends a chat completion request to OpenAI and returns the response
func (c *OpenAIClient) CreateCompletion(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature: 0,
		TopP:        0.1,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

