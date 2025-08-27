package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RakibulBh/AI-pr-reviewer/internal/model"
)

type GithubRepository struct {
	Token string
}

func NewGithubRepository(token string) *GithubRepository {
	return &GithubRepository{
		Token: token,
	}
}

func (u *GithubRepository) CreateReviewComments(owner string, repo string, pullNumber int64, comment model.ReviewCommentRequest) error {
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

func (u *GithubRepository) ListPullRequestFiles(owner, repo string, pullNumber int64) ([]model.PRFile, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files", owner, repo, pullNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+u.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch PR files: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []model.PRFile
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}

	return files, nil
}
