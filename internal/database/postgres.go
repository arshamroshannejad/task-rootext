package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github/arshamroshannejad/task-rootext/config"
	"time"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", makeDsn(cfg))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Postgres.ConnMaxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func makeDsn(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database,
	)
}
