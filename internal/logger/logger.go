package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	LogFormatter string
	LogLevel     string
}

// NewLogger creates a new logger based on the provided configuration.
//
// Parameters:
// - cfg: a pointer to a Config struct that contains the configuration for the logger.
//
// Returns:
// - logger: a pointer to a slog.Logger struct that represents the new logger instance.
// - error: an error if there was an issue creating the logger.
func NewLogger(cfg *Config) (*slog.Logger, error) {
	var logLevel = new(slog.LevelVar)

	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		return nil, fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	logHandlerOpts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var logHandler slog.Handler

	if cfg.LogFormatter == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, logHandlerOpts)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOpts)
	}

	logger := slog.New(logHandler)

	return logger, nil
}
