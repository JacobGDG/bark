package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/JacobGDG/bark/cli"
)

func init() {
	env := os.Getenv("LOG_LEVEL")
	handlerOptions := &slog.HandlerOptions{}

	switch strings.ToLower(env) {
	case "debug":
		handlerOptions.Level = slog.LevelDebug
	case "info":
		handlerOptions.Level = slog.LevelInfo
	case "error":
		handlerOptions.Level = slog.LevelError
	default:
		handlerOptions.Level = slog.LevelWarn
	}

	var loggerHandler *slog.JSONHandler
	loggerHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	slog.SetDefault(slog.New(loggerHandler))
}

func main() {
	// slog.Debug("Debug")
	// slog.Info("Information")
	// slog.Warn("Warning")
	// slog.Error("Error")

	cmd, err := cli.NewRootCommand("v0.0.1")
	if err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}

	cmd.Execute()
}
