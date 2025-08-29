package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/RakibulBh/AI-pr-reviewer/internal/usecase"
	"github.com/google/go-github/v74/github"
	"google.golang.org/genai"
)

type GithubController struct {
	usecase          *usecase.GithubUsecase
	webhookSecretKey string
	client           *genai.Client
}

func NewGithubController(usecase *usecase.GithubUsecase, webhookSecretKey string) *GithubController {
	return &GithubController{
		usecase:          usecase,
		webhookSecretKey: webhookSecretKey,
	}
}

func (c *GithubController) MainReciever(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(c.webhookSecretKey))
	if err != nil {
		slog.Error("error validating github webhook request payload", "error", err)
		return
	}

	// Parse the event
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		slog.Error("error parsing the webhook", "error", err)
		return
	}

	switch event := event.(type) {
	case *github.PullRequestEvent:
		pullRequest := event
		slog.Info("pull request event received")

		// Return 200 immediately to GitHub to prevent timeout
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Webhook received, processing in background"))

		// Process asynchronously with a longer timeout context
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()

			err := c.usecase.PullRequestReviewer(ctx, pullRequest)
			if err != nil {
				slog.Error("error reviewing pull request", "error", err)
				return
			}
		}()

	default:
		w.WriteHeader(http.StatusOK)
	}
}
