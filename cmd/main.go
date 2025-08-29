package main

import (
	"log/slog"
	"os"

	"github.com/RakibulBh/AI-pr-reviewer/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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

	// Extract the stored private key
	keyData, err := os.ReadFile("docs/credentials/bibi-the-monkey-code-reviewer.2025-08-29.private-key.pem")
	if err != nil {
		slog.Error("failed to read private key", "err", err)
		os.Exit(1)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		slog.Error("failed to parse private key", "err", err)
		os.Exit(1)
	}

	config.Bootstrap(&config.BootstrapConfig{
		R:            r,
		GeminiApiKey: os.Getenv("GEMINI_API_KEY"),
		Port:         "8080",

		// Github Repostored private key
		GithubWebhookSecret: os.Getenv("GITHUB_REPO_WEBHOOK_SECRET"),

		// Github Bot
		GithubBotPrivateKey: privateKey,
		AppID:               1864066,
	})

}
