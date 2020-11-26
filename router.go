package main

import (
	"net/http"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) getRouter() *chi.Mux {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(
		cors.Handler,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.DefaultLogger,
		chiprometheus.NewMiddleware("whoami"),
	)

	r.Group(func(telemetry chi.Router) {
		telemetry.Get("/metrics", promhttp.Handler().ServeHTTP)
		telemetry.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			render.PlainText(w, r, `{"ping":"pong"}`)
		})
	})

	r.Route("/", func(api chi.Router) {
		api.Get("/*", s.whoamiHandler())
		api.Put("/*", s.whoamiHandler())
		api.Post("/*", s.whoamiHandler())
		api.Patch("/*", s.whoamiHandler())
		api.Delete("/*", s.whoamiHandler())
		api.Options("/*", s.whoamiHandler())
	})

	return r
}
