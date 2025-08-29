package main

import (
	"os"

	"github.com/RakibulBh/AI-pr-reviewer/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Setup app variables
	config.Bootstrap(&config.BootstrapConfig{
		R:                   r,
		GithubWebhookSecret: os.Getenv("GITHUB_REPO_WEBHOOK_SECRET"),
		GithubAccessToken:   os.Getenv("GITHUB_ACCESS_TOKEN"),
		GeminiApiKey:        os.Getenv("GEMINI_API_KEY"),
		Port:                "8080",
	})

}
