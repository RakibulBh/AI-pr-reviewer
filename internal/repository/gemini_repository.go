package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	// Create a context with a longer timeout for LLM processing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

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
					Description: "The commit ID, this can be found in the SHA of the diff",
				},
				"path": {
					Type:        genai.TypeString,
					Description: "The file path relative to repository root",
				},
				"line": {
					Type:        genai.TypeInteger,
					Description: "The line of the blob in the pull request diff that the comment applies to. For a multi-line comment, the last line of the range that your comment applies to.",
				},
			},
			Required: []string{"body", "commit_id", "path", "line"},
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
		ctx,
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
		if comment.SubjectType != "" && (comment.SubjectType != "file" && comment.SubjectType != "line") {
			return fmt.Errorf("comment %d: subject type must be file or line", i)
		}
	}
	return nil
}
