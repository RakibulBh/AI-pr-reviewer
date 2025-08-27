package config

import (
	"context"
	"log"

	httpPackage "github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http"
	"github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http/route"
	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
	"github.com/RakibulBh/AI-pr-reviewer/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type BootstrapConfig struct {
	R            *chi.Mux
	GithubToken  string
	GeminiApiKey string
}

func Bootstrap(appConfig *BootstrapConfig) {
	// configs
	hook, err := NewGithubHook(appConfig.GithubToken)
	if err != nil {
		log.Fatalf("error creating new github webhook: %v\n", err)
	}
	client, err := NewGeminiClient(context.Background(), appConfig.GeminiApiKey)
	if err != nil {
		log.Fatalf("error creating gemini client: %v\n", err)
	}

	// setup repositories
	githubRepository := repository.NewGithubRepository(appConfig.GithubToken)
	geminiRepository := repository.NewGeminiRepository(client, context.Background())

	// setup use cases
	githubUsecase := usecase.NewGithubUsecase(githubRepository, geminiRepository)

	// setup controller
	githubController := httpPackage.NewGithubController(githubUsecase, hook)

	// setup middleware

	routeConfig := route.RouteConfig{
		R:                appConfig.R,
		GithubController: githubController,
	}
	routeConfig.Setup()
}
