package config

import (
	"context"
	"log/slog"
	"os"
)

func SetupLogger(env string) {
	if env == "production" {
		// Production: Dual output - JSON to file, Text to terminal
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("Failed to open log file, falling back to stdout only", "error", err)
			// Fallback to text handler on stdout only
			textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: false,
			})
			slog.SetDefault(slog.New(textHandler))
			return
		}

		// Create a multi-handler that writes JSON to file and text to terminal
		multiHandler := NewMultiHandler(
			slog.NewJSONHandler(logFile, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: false,
			}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: false,
			}),
		)
		slog.SetDefault(slog.New(multiHandler))
	} else {
		// Development: Text handler for human-readable output
		textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		slog.SetDefault(slog.New(textHandler))
	}
}

// MultiHandler writes to multiple handlers simultaneously
type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Enable if any handler is enabled for this level
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	// Write to all handlers
	for _, h := range m.handlers {
		if h.Enabled(ctx, record.Level) {
			if err := h.Handle(ctx, record); err != nil {
				// Continue with other handlers even if one fails
				continue
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return NewMultiHandler(newHandlers...)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return NewMultiHandler(newHandlers...)
}
