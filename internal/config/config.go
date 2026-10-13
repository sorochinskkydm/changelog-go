package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
	Log  LogConfig
}

type HTTPConfig struct {
	Port            int           `env:"SERVER_PORT" envDefault:"3000"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	DBName   string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxConns int32    `env:"DB_MAX_CONNS" envDefault:"10"`
}

type LogConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

func (config DBConfig) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(config.User, config.Password),
		Host: fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path: config.DBName,
		RawQuery: url.Values{"sslmode": {config.SSLMode}}.Encode(),
	}

	return u.String()
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
