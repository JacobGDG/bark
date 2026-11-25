package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/JacobGDG/bark/internal/cli"
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

	// Logs go to stderr so stdout stays clean for scripts reading bark's output.
	loggerHandler := slog.NewJSONHandler(os.Stderr, handlerOptions)
	slog.SetDefault(slog.New(loggerHandler))
}

func main() {
	cmd, err := cli.NewRootCommand("v0.0.1")
	if err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(cli.ExitError)
	}

	if err := cmd.Execute(); err != nil {
		var exitErr *cli.ExitCodeError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code)
		}
		os.Exit(cli.ExitError)
	}
}
