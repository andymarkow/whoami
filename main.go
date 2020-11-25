package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

type server struct {
	cfg *config
	web *http.Server
}

func main() {
	// Server context initialization
	srv := &server{
		cfg: getConfig(),
	}
	srv.web = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", srv.cfg.ServerHost, srv.cfg.ServerPort),
		Handler:      srv.getRouter(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// OS Signal handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("Starting web server on port %s", srv.cfg.ServerPort)
	go func() {
		if err := srv.web.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	log.Info("Web server graceful shutdown in progress")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.web.Shutdown(ctx); err != nil {
		log.Errorf("Web server graceful shutdown has failed: %+v", err)
	}
	log.Info("Web server gracefully stopped")
}
