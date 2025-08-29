package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/RakibulBh/AI-pr-reviewer/internal/usecase"
	"github.com/go-playground/webhooks/v6/github"
	"google.golang.org/genai"
)

type GithubController struct {
	usecase *usecase.GithubUsecase
	hook    *github.Webhook
	client  *genai.Client
}

func NewGithubController(usecase *usecase.GithubUsecase, hook *github.Webhook) *GithubController {
	return &GithubController{
		usecase: usecase,
		hook:    hook,
	}
}

func (c *GithubController) MainReciever(w http.ResponseWriter, r *http.Request) {
	payload, err := c.hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			slog.Error("github event not found", "event", payload)
			// ok event wasn't one of the ones asked to be parsed
		}
		http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
		return
	}

	switch payload.(type) {
	case github.ReleasePayload:
		release := payload.(github.ReleasePayload)
		// Do whatever you want from here...
		fmt.Printf("%+v", release)
		w.WriteHeader(http.StatusOK)

	case github.PingPayload:
		ping := payload.(github.PingPayload)
		fmt.Printf("ping request babe: %v", ping.Sender)
		w.WriteHeader(http.StatusOK)

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
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
			slog.Info("pull request review completed successfully")
		}()

	default:
		w.WriteHeader(http.StatusOK)
	}
}
