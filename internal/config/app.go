package config

import (
	"log"

	httpPackage "github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http"
	"github.com/RakibulBh/AI-pr-reviewer/internal/delivery/http/route"
	"github.com/RakibulBh/AI-pr-reviewer/internal/repository"
	"github.com/RakibulBh/AI-pr-reviewer/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type BootstrapConfig struct {
	R *chi.Mux
}

func Bootstrap(appConfig *BootstrapConfig) {
	// setup repositories
	githubRepository := repository.NewGithubRepository()

	// setup producer

	// configs
	hook, err := NewGithubHook()
	if err != nil {
		log.Fatalf("error creating new github webhook: %v\n", err)
	}

	// setup use cases
	githubUsecase := usecase.NewGithubUsecase(githubRepository)

	// setup controller
	githubController := httpPackage.NewGithubController(githubUsecase, hook)

	// setup middleware

	routeConfig := route.RouteConfig{
		R:                appConfig.R,
		GithubController: githubController,
	}
	routeConfig.Setup()
}
