package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"merch-shop/internal/api"
	"merch-shop/internal/handler"
	mdlwr "merch-shop/internal/middleware"
	"merch-shop/internal/service"
)

const (
	timeout = time.Second * 60
)

func New(secret string, service *service.Services) *chi.Mux {
	r := chi.NewRouter()

	corsMiddleware := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(
		middleware.Timeout(timeout),
		corsMiddleware,
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
	)

	h := handler.New(secret, service)

	r.Group(func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Mount("/api/auth", api.Handler(h))
		})

		r.Group(func(r chi.Router) {
			r.Use(mdlwr.AuthMiddleware(secret))
			r.Mount("/api", api.Handler(h))
		})
	})

	return r
}
