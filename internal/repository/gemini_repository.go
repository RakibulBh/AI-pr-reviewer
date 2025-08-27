package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RakibulBh/AI-pr-reviewer/internal/model"
	"github.com/RakibulBh/AI-pr-reviewer/internal/utils"
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

func (g *GeminiRepository) GetCodeReviews(code string) ([]model.ReviewCommentRequest, error) {
	parts := []*genai.Part{
		{Text: code},
	}

	content := []*genai.Content{
		{Parts: parts},
	}

	// Get technical requirements
	fileContent, err := utils.ReadRepositoryRuleFile("main.md")
	if err != nil {
		return nil, err
	}

	// Generate system prompt
	systemPrompt := utils.GenerateCodeReviewPrompt(fileContent)

	// Setup response schema for structured output
	responseSchema := &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"body": {
					Type:        genai.TypeString,
					Description: "The review comment text",
				},
				"commit_id": {
					Type:        genai.TypeString,
					Description: "The commit ID (will be set by system)",
				},
				"path": {
					Type:        genai.TypeString,
					Description: "The file path relative to repository root",
				},
				"start_line": {
					Type:        genai.TypeInteger,
					Description: "Optional starting line number for multi-line comments",
				},
				"start_side": {
					Type:        genai.TypeString,
					Description: "Optional starting side for multi-line comments",
				},
				"line": {
					Type:        genai.TypeInteger,
					Description: "The line number where the issue occurs",
				},
				"side": {
					Type:        genai.TypeString,
					Description: "The side of the diff (RIGHT for new code)",
				},
			},
			Required: []string{"body", "path", "line", "side"},
		},
	}

	// Setup the configuration
	cfg := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Role: "system",
			Parts: []*genai.Part{
				{
					Text: systemPrompt,
				},
			},
		},
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	}

	result, err := g.client.Models.GenerateContent(
		g.ctx,
		"gemini-2.0-flash",
		content,
		cfg,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract and validate the response
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// Parse the JSON response
	responseText := result.Candidates[0].Content.Parts[0].Text
	var reviewComments []model.ReviewCommentRequest

	if err := json.Unmarshal([]byte(responseText), &reviewComments); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Validate the response
	if err := validateReviewComments(reviewComments); err != nil {
		return nil, fmt.Errorf("invalid response format: %w", err)
	}

	return reviewComments, nil
}

// validateReviewComments validates the structure and content of review comments
func validateReviewComments(comments []model.ReviewCommentRequest) error {
	for i, comment := range comments {
		if comment.Body == "" {
			return fmt.Errorf("comment %d: body cannot be empty", i)
		}
		if comment.Path == "" {
			return fmt.Errorf("comment %d: path cannot be empty", i)
		}
		if comment.Line <= 0 {
			return fmt.Errorf("comment %d: line must be greater than 0", i)
		}
		if comment.Side == "" {
			return fmt.Errorf("comment %d: side cannot be empty", i)
		}
		// Validate side values
		if comment.Side != "LEFT" && comment.Side != "RIGHT" {
			return fmt.Errorf("comment %d: side must be either 'LEFT' or 'RIGHT'", i)
		}
		// Validate start_line if provided
		if comment.StartLine > 0 && comment.StartLine > comment.Line {
			return fmt.Errorf("comment %d: start_line cannot be greater than line", i)
		}
	}
	return nil
}
