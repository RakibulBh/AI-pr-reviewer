package route

import (
	"github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http"
	"github.com/go-chi/chi/v5"
)

type RouteConfig struct {
	R                *chi.Mux
	GithubController *http.GithubController
}

func (c *RouteConfig) Setup() {
	c.SetupWebhookRoute()
}

func (c *RouteConfig) SetupWebhookRoute() {
	c.R.Post("/webhook", c.GithubController.MainReciever)
}
