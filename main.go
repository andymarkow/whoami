package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/andymarkow/whoami/internal/config"
	"github.com/andymarkow/whoami/internal/httpserver"
	"github.com/andymarkow/whoami/internal/logger"
	"github.com/andymarkow/whoami/internal/telemetry"
)

var (
	Version = "0.0.0-dev"
)

func main() {
	telemetry.Init(Version)

	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("config.NewConfig: %w", err))
	}

	l, err := logger.NewLogger(&logger.Config{
		LogFormatter: cfg.LogFormatter,
		LogLevel:     cfg.LogLevel,
	})
	if err != nil {
		panic(fmt.Errorf("NewLogger: %w", err))
	}
	slog.SetDefault(l)

	srv := httpserver.NewHTTPServer(cfg)

	go func() {
		slog.Info(fmt.Sprintf("Starting http server on address %s:%s", cfg.ServerHost, cfg.ServerPort))
		if err := srv.Start(); err != nil {
			slog.Error("srv.Start: %w", err)
			os.Exit(1)
		}
	}()

	// Gracefully shutdown the web server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Http server graceful shutdown initiated")
	if err := srv.Shutdown(); err != nil {
		slog.Error("srv.Shutdown: %w", err)
		os.Exit(1)
	}
}
