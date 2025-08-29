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
					Description: "The commit ID, this can be found in the SHA of the diff",
				},
				"path": {
					Type:        genai.TypeString,
					Description: "The file path relative to repository root",
				},
				"position": {
					Type:        genai.TypeInteger,
					Description: "The position in the diff where you want to add a review comment. Note this value is not the same as the line number in the file. The position value equals the number of lines down from the first @@ hunk header in the file you want to add a comment. The line just below the @@ line is position 1, the next line is position 2, and so on. The position in the diff continues to increase through lines of whitespace and additional hunks until the beginning of a new file.",
				},
				"subject_type": {
					Type:        genai.TypeString,
					Enum:        []string{"line", "file"},
					Description: "The level at which the comment is targeted.",
				},
			},
			Required: []string{"body", "commit_id", "path", "subject_type"},
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
		if comment.Position <= 0 {
			return fmt.Errorf("comment %d: position must be greater than 0", i)
		}
		if comment.SubjectType != "" && (comment.SubjectType != "file" && comment.SubjectType != "line") {
			return fmt.Errorf("comment %d: subject type must be file or line", i)
		}
	}
	return nil
}
