package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github/arshamroshannejad/task-rootext/config"
	"github/arshamroshannejad/task-rootext/internal/helpers"
	"github/arshamroshannejad/task-rootext/internal/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func SetupRoutes(db *sql.DB, redisDB *redis.Client, zapLogger *zap.Logger, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RedirectSlashes)
	r.Use(chiMiddleware.CleanPath)
	r.Use(chiMiddleware.Heartbeat("/heartbeat"))
	r.Use(chiMiddleware.Timeout(time.Second * 30))
	r.Use(middleware.CorsMiddleware(cfg))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		body := "The requested resource could not be found. Please check the URL and try again."
		helpers.WriteJson(w, http.StatusNotFound, helpers.M{"error": body})
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		body := "The method specified in the request is not allowed for the resource." +
			"Please check the allowed methods and try again."
		helpers.WriteJson(w, http.StatusMethodNotAllowed, helpers.M{"error": body})
	})
	return r
}
