package middleware

import (
	"github.com/arshamroshannejad/task-rootext/config"
	"github.com/go-chi/cors"
	"net/http"
)

func CorsMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: cfg.App.CorsOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
