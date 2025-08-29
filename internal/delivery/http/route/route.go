package route

import (
	"github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http"
	"github.com/go-chi/chi/v5"
)

type RouteConfig struct {
	R                *chi.Mux
	GithubController *http.GithubController
	HealthController *http.HealthController
}

func (c *RouteConfig) Setup() {
	c.SetupWebhookRoute()
	c.SetupMetricRoutes()
}

func (c *RouteConfig) SetupMetricRoutes() {
	c.R.Get("/health", c.HealthController.Health)
}

func (c *RouteConfig) SetupWebhookRoute() {
	c.R.Post("/webhook", c.GithubController.MainReciever)
}
