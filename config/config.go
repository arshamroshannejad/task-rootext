package config

import (
	"bytes"
	_ "embed"
	"github.com/spf13/viper"
	"time"
)

//go:embed config.yaml
var Configurations []byte

type Postgres struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

type Redis struct {
	Host     string
	Port     int
	Password string
}

type App struct {
	Host           string
	Port           int
	Debug          bool
	BaseAPI        string
	Secret         string
	AccessHourTTL  time.Duration
	RefreshHourTTL time.Duration
	CorsOrigins    []string
	CorsMaxAge     int
}

type Config struct {
	Postgres *Postgres
	Redis    *Redis
	App      *App
}

func New() (*Config, error) {
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(Configurations)); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
