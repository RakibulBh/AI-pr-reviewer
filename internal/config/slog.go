package config

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	// Alternative: Text handler for development (human-readable)
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	slog.SetDefault(slog.New(textHandler))
}
