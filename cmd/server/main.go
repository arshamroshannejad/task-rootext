package main

import (
	"fmt"
	"github.com/arshamroshannejad/task-rootext/config"
	"github.com/arshamroshannejad/task-rootext/internal/database"
	"github.com/arshamroshannejad/task-rootext/internal/logger"
	"github.com/arshamroshannejad/task-rootext/internal/router"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//	@title						task-rootext
//	@version					0.1.0
//	@host						localhost:8000
//	@description				API like reddit application
//	@termsOfService				http://swagger.io/terms/
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				Arsham Roshannejad
//	@contact.url				arshamroshannejad.ir
//	@contact.email				arshamdev2001@gmail.com
//	@license.name				MIT
//	@license.url				https://www.mit.edu/~amini/LICENSE.md
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func main() {
	cfg, err := config.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize config variables: %v", err))
	}
	zapLog, err := logger.New(cfg.App.Debug)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}
	defer zapLog.Sync()
	db, err := database.OpenDB(cfg)
	if err != nil {
		zapLog.Fatal("Failed to connect postgres", zap.Error(err))
	}
	defer db.Close()
	zapLog.Info("Postgres connected", zap.String("Host", cfg.Postgres.Host), zap.Int("Port", cfg.Postgres.Port))
	redisDB, err := database.OpenRedis(cfg)
	if err != nil {
		zapLog.Fatal("Failed to connect redis", zap.Error(err))
	}
	defer redisDB.Close()
	zapLog.Info("Redis connected", zap.String("Host", cfg.Redis.Host), zap.Int("Port", cfg.Redis.Port))
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port),
		Handler:      router.SetupRoutes(db, redisDB, zapLog, cfg),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	zapLog.Info("Starting server", zap.String("Host", cfg.App.Host), zap.Int("Port", cfg.App.Port))
	if err := srv.ListenAndServe(); err != nil {
		zapLog.Fatal("Failed to start server", zap.Error(err))
	}
}
