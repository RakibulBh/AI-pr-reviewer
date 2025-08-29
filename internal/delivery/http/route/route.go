package route

import (
	"log/slog"
	"net/http"

	httpPackage "github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http"
	"github.com/go-chi/chi/v5"
)

type RouteConfig struct {
	R                *chi.Mux
	GithubController *httpPackage.GithubController
	HealthController *httpPackage.HealthController
}

func (c *RouteConfig) Setup() {
	c.SetupWebhookRoute()
	c.SetupMetricRoutes()
	c.SetupCatchAllRoute()
}

func (c *RouteConfig) SetupMetricRoutes() {
	c.R.Get("/health", c.HealthController.Health)
}

func (c *RouteConfig) SetupWebhookRoute() {
	c.R.Post("/webhook", c.GithubController.MainReciever)
}

func (c *RouteConfig) SetupCatchAllRoute() {
	c.R.NotFound(func(w http.ResponseWriter, r *http.Request) {
		slog.Warn("404 - Route not found", "method", r.Method, "path", r.URL.Path, "user_agent", r.UserAgent())
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	})
}
