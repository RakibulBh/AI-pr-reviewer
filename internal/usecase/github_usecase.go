package usecase

import (
	"github.com/RakibulBh/AI-pr-reviewer/internal/model"
	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
	"github.com/RakibulBh/AI-pr-reviewer/internal/utils"
	"github.com/go-playground/webhooks/v6/github"
)

type GithubUsecase struct {
	repository *repository.GithubRepository
	gemini     *repository.GeminiRepository
}

func NewGithubUsecase(repository *repository.GithubRepository, gemini *repository.GeminiRepository) *GithubUsecase {
	return &GithubUsecase{
		repository: repository,
		gemini:     gemini,
	}
}

func (g *GithubUsecase) PullRequestReviewer(pullRequest github.PullRequestPayload) error {
	owner := pullRequest.Repository.Owner.Login
	repo := pullRequest.Repository.Name
	pullNumber := pullRequest.Number

	// Fetch file changes with code for context
	files, err := g.repository.ListPullRequestFiles(owner, repo, pullNumber)
	if err != nil {
		return err
	}

	// Get the reviews for the PR
	filesToJsonString, err := utils.StructToJSON(files)
	if err != nil {
		return err
	}
	reviews, err := g.gemini.GetCodeReviews(filesToJsonString)
	if err != nil {
		return err
	}

	// Create each revoew
	for _, review := range reviews {
		comment := model.ReviewCommentRequest{
			Body:      review.Body,
			CommitID:  review.CommitID,
			Path:      review.Path,
			StartLine: review.StartLine,
			StartSide: review.StartSide,
			Line:      review.Line,
			Side:      review.Side,
		}

		err = g.repository.CreateReviewComments(owner, repo, pullNumber, comment)
		if err != nil {
			return err
		}
	}

	return nil
}
