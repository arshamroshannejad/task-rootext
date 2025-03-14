package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github/arshamroshannejad/task-rootext/config"
	"github/arshamroshannejad/task-rootext/internal/handler"
	"github/arshamroshannejad/task-rootext/internal/middleware"
	"github/arshamroshannejad/task-rootext/internal/repository"
	"github/arshamroshannejad/task-rootext/internal/service"
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
	r.MethodNotAllowed(handler.HttpMethodNotAllowedHandler)
	r.NotFound(handler.HttpRequestNotFound)
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, redisDB, zapLogger, cfg)
	userHandler := handler.NewUserHandler(userService)
	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository, zapLogger)
	postHandler := handler.NewPostHandler(postService)
	apiV1Router := chi.NewRouter()
	apiV1Router.Route("/auth", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterHandler)
		r.Post("/login", userHandler.LoginHandler)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JwtAuth(redisDB, zapLogger, cfg))
			r.Post("/logout", userHandler.LogoutHandler)
		})
	})
	apiV1Router.Route("/post", func(r chi.Router) {
		r.Get("/", postHandler.GetAllPostsHandler)
		r.Get("/{id}", postHandler.GetPostHandler)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JwtAuth(redisDB, zapLogger, cfg))
			r.Post("/", postHandler.CreatePostHandler)
			r.Put("/{id}", postHandler.UpdatePostHandler)
			r.Delete("/{id}", postHandler.DeletePostHandler)
			r.Post("/{id}/vote", postHandler.AddPostVoteHandler)
			r.Delete("/{id}/unvote", postHandler.RemovePostVoteHandler)
		})
	})
	r.Mount("/api/v1", apiV1Router)
	return r
}
