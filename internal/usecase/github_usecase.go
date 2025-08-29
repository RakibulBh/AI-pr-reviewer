package usecase

import (
	"log/slog"
	"time"

	"github.com/RakibulBh/AI-pr-reviewer/internal/model"
	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
	"github.com/RakibulBh/AI-pr-reviewer/internal/utils/json"
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

	// Loop until there are no more pages of files to review
	var pageCount int32 = 0
	for {
		files, err := g.repository.ListPullRequestFiles(owner, repo, pullNumber, pageCount)
		if err != nil {
			slog.Error("error fetching diffs", "error", err, "owner", owner, "repo", repo, "pullNumber", pullNumber)
			return err
		}

		// If there are no files break the loop
		if len(files) <= 0 {
			break
		}

		// Get the reviews for the PR
		filesToJsonString, err := json.StructToJSON(files)
		if err != nil {
			return err
		}
		reviews, err := g.gemini.GetCodeReviews(filesToJsonString)
		if err != nil {
			slog.Error("error getting code reviews from LLM", "error", err)
			return err
		}
		slog.Info("reviews have been created by the LLM", "number_of_reviews", len(reviews))

		// Create each review
		for _, review := range reviews {
			commitID := pullRequest.PullRequest.Head.Sha

			comment := &model.ReviewCommentRequest{
				Body:        review.Body,
				CommitID:    commitID,
				Path:        review.Path,
				Position:    review.Position,
				SubjectType: review.SubjectType,
			}

			time.Sleep(time.Second * 5)

			err = g.repository.CreateReviewComments(owner, repo, pullNumber, comment)
			if err != nil {
				slog.Error("error creating review comment", "error", err, "comment", comment)
				return err
			}
		}

		pageCount++
		time.Sleep(time.Second * 15)
	}

	return nil
}
