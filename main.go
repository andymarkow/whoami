package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/andyglass/whoami/config"
	"github.com/andyglass/whoami/httpserver"
	"github.com/andyglass/whoami/logger"
	"github.com/andyglass/whoami/telemetry"
)

var (
	Version = "0.0.1"
)

func main() {
	cfg := config.New(Version)
	logger.App = logger.New(cfg.LogLevel)
	telemetry.Init(Version)

	webserver := httpserver.New(cfg)
	go webserver.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	webserver.Shutdown()
}
