package repository

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiRepository struct {
	client *genai.Client
	ctx    context.Context
}

func NewGeminiRepository(client *genai.Client, ctx context.Context) *GeminiRepository {
	return &GeminiRepository{
		client: client,
		ctx:    ctx,
	}
}

func (g *GeminiRepository) GenerateText(prompt string) (*genai.GenerateContentResponse, error) {
	parts := []*genai.Part{
		{Text: prompt},
	}

	content := []*genai.Content{
		{Parts: parts},
	}

	result, err := g.client.Models.GenerateContent(
		g.ctx,
		"gemini-2.0-flash",
		content,
		nil, // Use default generation config
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	return result, nil
}
