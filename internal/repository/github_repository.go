package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/RakibulBh/AI-pr-reviewer/internal/model"
)

type GithubRepository struct {
	WebhookSecret string
	AccessToken   string
}

func NewGithubRepository(webhookSecret string, accessToken string) *GithubRepository {
	return &GithubRepository{
		WebhookSecret: webhookSecret,
		AccessToken:   accessToken,
	}
}

func (u *GithubRepository) CreateReviewComments(owner string, repo string, pullNumber int64, comment *model.ReviewCommentRequest) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", owner, repo, pullNumber)

	// marshal the request body
	payload, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	slog.Debug("the payload being sent to the URL", "payload", string(payload), "url", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// set headers
	req.Header.Set("Authorization", "Bearer "+u.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("failed to create review comment", "status", resp.Status, "response", string(body))
		return fmt.Errorf("failed to create review comment: %s - %s", resp.Status, string(body))
	}

	return nil
}

func (u *GithubRepository) ListPullRequestFiles(owner, repo string, pullNumber int64, pageNumber int32) ([]model.PRFile, error) {
	slog.Debug("trying to fetch diffs", "owner", owner, "repo", repo, "pullNumber", pullNumber)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files?per_page=10&page=%d", owner, repo, pullNumber, pageNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch PR files: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []model.PRFile
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}

	return files, nil
}
