package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GithubRepository struct {
	Token string
}

func NewGithubRepository(token string) *GithubRepository {
	return &GithubRepository{
		Token: token,
	}
}

func (u *GithubRepository) CreateReviewComments(comment string, owner string, repo string, pullNumber int) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", owner, repo, pullNumber)

	// marshal the request body
	payload, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// set headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+u.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to create review comment: %s", resp.Status)
	}

	return nil
}
