package config

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	httpPackage "github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http"
	"github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http/route"
	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
	"github.com/RakibulBh/AI-pr-reviewer/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type BootstrapConfig struct {
	R                   *chi.Mux
	Port                string
	GithubWebhookSecret string
	GithubAccessToken   string
	GeminiApiKey        string
}

func Bootstrap(appConfig *BootstrapConfig) {
	// configs
	hook, err := NewGithubHook(appConfig.GithubWebhookSecret)
	if err != nil {
		log.Fatalf("error creating new github webhook: %v\n", err)
	}
	slog.Info("Github Hook Listener has been succsessfully connected")
	client, err := NewGeminiClient(context.Background(), appConfig.GeminiApiKey)
	if err != nil {
		log.Fatalf("error creating gemini client: %v\n", err)
	}
	slog.Info("LLM client has been created and connected")

	// Setup logger
	SetupLogger()

	// setup repositories
	githubRepository := repository.NewGithubRepository(appConfig.GithubWebhookSecret, appConfig.GithubAccessToken)
	geminiRepository := repository.NewGeminiRepository(client, context.Background())

	// setup use cases
	githubUsecase := usecase.NewGithubUsecase(githubRepository, geminiRepository)

	// setup controller
	githubController := httpPackage.NewGithubController(githubUsecase, hook)
	healthController := httpPackage.NewHealthController()

	// setup middleware

	routeConfig := route.RouteConfig{
		R:                appConfig.R,
		GithubController: githubController,
		HealthController: healthController,
	}

	// Setup routes and start server
	routeConfig.Setup()
	err = http.ListenAndServe(":"+appConfig.Port, appConfig.R)
	if err != nil {
		slog.Warn("major error starting server", "error", err)
		return
	}
	slog.Info("app is now running", "port", appConfig.Port)
}
