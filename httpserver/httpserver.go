package httpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	echoprom "github.com/andyglass/echo-prometheus"
	"github.com/andyglass/whoami/config"
	"github.com/andyglass/whoami/httpserver/handlers"
	"github.com/andyglass/whoami/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	logURLExcludes []string
	healthzStatus  = 200
)

type WebServer struct {
	cfg    *config.Config
	router *echo.Echo
}

func New(cfg *config.Config) *WebServer {
	logURLExcludes = cfg.LogURLExcludes
	return &WebServer{
		cfg:    cfg,
		router: newRouter(cfg.Version),
	}
}

// Middleware with delayed response
func responceDelayMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		delayQuery := c.Request().URL.Query().Get("delay")
		if delayQuery == "" {
			return next(c)
		}
		delay, err := time.ParseDuration(delayQuery)
		if err != nil {
			logger.App.Error("Delay duration provided but could not be parsed from query param")
			return next(c)
		}
		time.Sleep(delay)
		return next(c)
	}
}

func newRouter(version string) *echo.Echo {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	configMetrics := echoprom.NewConfig()
	configMetrics.Buckets = []float64{
		0.05, // 50ms
		0.1,  // 100ms
		0.5,  // 500ms
		1,    // 1s
		2.5,  // 2.5s
		5,    // 5s
		10,   // 10s
	}

	e.Use(echoprom.MiddlewareWithConfig(configMetrics))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			for _, exclude := range logURLExcludes {
				if strings.HasPrefix(c.Path(), exclude) {
					return true
				}
			}
			return false
		},
		CustomTimeFormat: "2006-01-02T15:04:05.000Z",
		Format: `{"time":"${time_custom}","request_id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","status":${status},"user_agent":"${user_agent}",` +
			`"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf(`{"version":"%s"}`, version))
	}, responceDelayMiddleware)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(healthzStatus, fmt.Sprintf(`{"status":"%d"}`, healthzStatus))
	}, responceDelayMiddleware)

	e.POST("/healthz", func(c echo.Context) error {
		byteContent, err := io.ReadAll(c.Request().Body)
		if err != nil || len(byteContent) == 0 {
			return c.String(http.StatusBadRequest, "Wrong POST request content or empty body")
		}
		defer c.Request().Body.Close()
		status, err := strconv.Atoi(string(byteContent))
		if err != nil || status < 200 || status > 599 {
			return c.String(http.StatusBadRequest, "Wrong http status code")
		}
		healthzStatus = status
		return c.NoContent(http.StatusNoContent)
	})

	h := handlers.New()

	api := e.Group("/api", responceDelayMiddleware)
	api.Any("*", h.WhoamiJSON)

	plain := e.Group("/", responceDelayMiddleware)
	plain.Any("*", h.WhoamiPlain)

	return e
}

func (w *WebServer) Start() {
	web := w.router
	web.Server.Addr = w.cfg.ServerAddr
	logger.App.Infof("Starting http server on %s", w.cfg.ServerAddr)
	if err := web.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.App.Fatal(errors.WithStack(err))
	}
}

func (w *WebServer) Shutdown() {
	logger.App.Info("Http server graceful shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := w.router.Server.Shutdown(ctx); err != nil {
		logger.App.Fatalf("Http server graceful shutdown failed: %+v", err)
	}
}
