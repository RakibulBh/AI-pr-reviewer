package repository

import (
	"context"
	"crypto/rsa"
	"log/slog"

	"github.com/google/go-github/v74/github"
)

type GithubRepository struct {
	WebhookSecret string
	PrivateKey    *rsa.PrivateKey
}

func NewGithubRepository(webhookSecret string, privateKey *rsa.PrivateKey) *GithubRepository {
	return &GithubRepository{
		WebhookSecret: webhookSecret,
		PrivateKey:    privateKey,
	}
}

func (u *GithubRepository) CreateReviewComments(ctx context.Context, client *github.Client, owner string, repo string, pullNumber int, comment *github.PullRequestComment) error {

	_, _, err := client.PullRequests.CreateComment(ctx, owner, repo, pullNumber, comment)
	if err != nil {
		return err
	}

	return nil
}

func (u *GithubRepository) ListPullRequestFiles(ctx context.Context, client *github.Client, owner, repo string, pullNumber int, pageNumber int) ([]*github.CommitFile, error) {
	slog.Debug("trying to fetch diffs", "owner", owner, "repo", repo, "pullNumber", pullNumber)

	opts := &github.ListOptions{
		Page: int(pageNumber),
	}

	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pullNumber, opts)
	if err != nil {
		return nil, err
	}

	return files, nil
}
