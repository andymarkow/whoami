package main

import (
	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		middleware.RealIP,
		middleware.Recoverer,
		middleware.DefaultLogger,
		chiprometheus.NewMiddleware("whoami"),
	)

	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/ping", pingHandler())
	r.Get("/env", envHandler())
	r.Get("/*", whoamiHandler())

	return r
}
