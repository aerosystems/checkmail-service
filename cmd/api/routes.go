package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(mux chi.Router) {
		// Public routes
		mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
		mux.Use(middleware.Heartbeat("/ping"))

		mux.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("doc.json"), // The url pointing to API definition
		))

		mux.Get("/v1/domain/{domainName}", app.BaseHandler.DomainRead)
		mux.Post("/v1/domain", app.BaseHandler.DomainCreate)
		mux.Patch("/v1/domain/{domainName}", app.BaseHandler.DomainUpdate)
		mux.Delete("/v1/domain/{domainName}", app.BaseHandler.DomainDelete)
		mux.Get("/v1/data/{data}", app.BaseHandler.Data)
		mux.Get("/v1/email/{email}", app.BaseHandler.Email)
	})

	return mux
}
