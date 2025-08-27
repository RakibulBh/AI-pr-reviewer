package http

import (
	"fmt"
	"net/http"

	"github.com/RakibulBh/AI-pr-reviewer/internal/usecase"
	"github.com/go-playground/webhooks/v6/github"
)

type GithubController struct {
	usecase *usecase.GithubUsecase
	hook    *github.Webhook
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
			// ok event wasn't one of the ones asked to be parsed
		}
	}
	switch payload.(type) {
	case github.ReleasePayload:
		release := payload.(github.ReleasePayload)
		// Do whatever you want from here...
		fmt.Printf("%+v", release)

	case github.PingPayload:
		ping := payload.(github.PingPayload)
		fmt.Printf("ping request babe: %v", ping.Sender)

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		// Do whatever you want from here...
		fmt.Print(pullRequest.Sender)
		fmt.Printf("%+v", pullRequest)
	}
}
