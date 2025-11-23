package config

import (
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func configureLogging(level string, t string) {
	var handler slog.Handler
	logConfig := &slog.HandlerOptions{
		Level: logLevel(level),
	}

	if t == "console" {
		handler = slog.NewTextHandler(os.Stdout, logConfig)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, logConfig)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func logLevel(l string) slog.Level {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
