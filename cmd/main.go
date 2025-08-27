package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RakibulBh/AI-pr-reviewer/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Setup app variables
	config.Bootstrap(&config.BootstrapConfig{
		R:            r,
		GithubToken:  os.Getenv("GITHUB_REPO_WEBHOOK_SECRET"),
		GeminiApiKey: os.Getenv("GEMINI_API_KEY"),
	})

	// start the server
	log.Fatal(http.ListenAndServe(":3000", r))
}
