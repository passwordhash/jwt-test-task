package config

import (
	"log/slog"
	"os"
)

func SetupLogger(env string) *slog.Logger {
	var handler slog.Handler
	handlerOpts := new(slog.HandlerOptions)
	w := os.Stdout

	switch env {
	case "dev":
		handlerOpts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(w, handlerOpts)
	case "prod":
		handlerOpts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(w, handlerOpts)
	default:
		handlerOpts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(w, handlerOpts)
	}

	return slog.New(handler)
}
