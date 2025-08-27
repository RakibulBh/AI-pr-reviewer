package usecase

import (
	"github.com/RakibulBh/AI-pr-reviewer/internal/model"
	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
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
	_, err := g.repository.ListPullRequestFiles(owner, repo, pullNumber)
	if err != nil {
		return err
	}

	comment := model.ReviewCommentRequest{
		Body:      "Great stuff!",
		CommitID:  "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		Path:      "file1.txt",
		StartLine: 1,
		StartSide: "RIGHT",
		Line:      2,
		Side:      "RIGHT",
	}

	err = g.repository.CreateReviewComments(owner, repo, pullNumber, comment)
	if err != nil {
		return err
	}

	return nil
}
