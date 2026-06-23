package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func init() {
	env := os.Getenv("LOG_LEVEL")
	handlerOptions := &slog.HandlerOptions{}

	fmt.Println(env)

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
	slog.Debug("Debug")
	slog.Info("Information")
	slog.Warn("Warning")
	slog.Error("Error")
}
