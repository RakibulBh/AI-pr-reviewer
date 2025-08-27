package usecase

import (
	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
)

type GithubUsecase struct {
	repository *repository.GithubRepository
}

type ReviewCommentRequest struct {
	Body      string `json:"body"`
	CommitID  string `json:"commit_id"`
	Path      string `json:"path"`
	StartLine int    `json:"start_line,omitempty"`
	StartSide string `json:"start_side,omitempty"`
	Line      int    `json:"line"`
	Side      string `json:"side"`
}

func NewGithubUsecase(repository *repository.GithubRepository) *GithubUsecase {
	return &GithubUsecase{
		repository: repository,
	}
}

func (c *GithubUsecase) PullRequestReviwer() error {
	comment := ""
	owner := ""
	repo := ""
	pullNumber := 1
	err := c.repository.CreateReviewComments(comment, owner, repo, pullNumber)
	if err != nil {
		return err
	}

	return nil
}
