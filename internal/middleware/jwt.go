package middleware

import (
	"context"
	"errors"
	"github.com/arshamroshannejad/task-rootext/config"
	"github.com/arshamroshannejad/task-rootext/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func JwtAuth(redisDB *redis.Client, zapLogger *zap.Logger, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				helpers.WriteJson(w, http.StatusUnauthorized, helpers.M{"error": "authorization header not provided"})
				return
			}
			authToken := strings.Split(authHeader, " ")
			if len(authToken) != 2 || strings.ToLower(authToken[0]) != "bearer" {
				helpers.WriteJson(w, http.StatusUnauthorized, helpers.M{"error": "invalid Authorization header"})
				return
			}
			token, err := helpers.IsTokenValid(authToken[1], cfg.App.Secret)
			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					helpers.WriteJson(w, http.StatusUnauthorized, helpers.M{"error": "token is expired"})
				default:
					zapLogger.Error("invalid token: failed to parse token", zap.Any("error", err))
					helpers.WriteJson(w, http.StatusUnauthorized, helpers.M{"error": "token is invalid"})
				}
				return
			}
			if _, err := redisDB.Get(context.Background(), authHeader).Result(); err == nil {
				helpers.WriteJson(w, http.StatusUnauthorized, helpers.M{"error": "token is expired"})
				return
			}
			claims, err := helpers.GetClaims(token)
			if err != nil {
				zapLogger.Error("invalid claims: failed to parse token", zap.Any("error", err))
				helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "user_id", claims["user_id"]))
			r = r.WithContext(context.WithValue(r.Context(), "email", claims["email"]))
			r = r.WithContext(context.WithValue(r.Context(), "exp", claims["exp"]))
			next.ServeHTTP(w, r)
		})
	}
}
