package config

import (
	"log/slog"
	"os"
)

func SetupLogger(env string) {
	if env == "production" {
		// Production: JSON handler writing to file
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fallback to stdout if file creation fails
			logFile = os.Stdout
		}

		jsonHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
		slog.SetDefault(slog.New(jsonHandler))
	} else {
		// Development: Text handler for human-readable output
		textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		slog.SetDefault(slog.New(textHandler))
	}
}
