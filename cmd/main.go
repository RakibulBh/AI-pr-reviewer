package main

import (
	"log/slog"
	"os"
	"strconv"

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

	env := os.Getenv("ENV")

	// Extract the stored private key
	var keyData []byte
	var err error
	if env == "production" {
		keyDataStr := os.Getenv("GITHUB_BOT_PRIVATE_KEY")
		if keyDataStr == "" {
			slog.Error("GITHUB_BOT_PRIVATE_KEY environment variable is not set")
			os.Exit(1)
		}
		keyData = []byte(keyDataStr)
	} else {
		keyData, err = os.ReadFile("docs/credentials/bibi-the-monkey-code-reviewer.2025-08-29.private-key.pem")
		if err != nil {
			slog.Error("failed to read private key", "err", err)
			os.Exit(1)
		}
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		slog.Error("failed to parse private key", "err", err)
		os.Exit(1)
	}

	appIDStr := os.Getenv("APP_ID")
	if appIDStr == "" {
		slog.Error("APP_ID environment variable is not set")
	}

	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		slog.Error("Invalid APP_ID", "APP_ID", appID)
	}

	config.Bootstrap(&config.BootstrapConfig{
		R:            r,
		GeminiApiKey: os.Getenv("GEMINI_API_KEY"),
		Port:         os.Getenv("PORT"),
		Env:          env,

		// Github Repostored private key
		GithubWebhookSecret: os.Getenv("GITHUB_REPO_WEBHOOK_SECRET"),

		// Github Bot
		GithubBotPrivateKey: privateKey,
		AppID:               int64(appID),
	})

}
