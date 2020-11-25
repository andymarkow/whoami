package main

import (
	"net/http"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) getRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		middleware.RealIP,
		middleware.Recoverer,
		middleware.DefaultLogger,
		chiprometheus.NewMiddleware("whoami"),
	)

	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, `{"ping":"pong"}`)
	})
	r.Get("/*", s.whoamiHandler())

	return r
}
