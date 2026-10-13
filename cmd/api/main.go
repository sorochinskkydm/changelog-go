package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sorochinskkydm/changelog-go/internal/app"
	"github.com/sorochinskkydm/changelog-go/internal/config"
	"github.com/sorochinskkydm/changelog-go/internal/platform/postgres"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()

	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	log := newLogger(cfg.Log.Level)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Error("connect to db", "error", err)
		os.Exit(1)
	}

	defer db.Close()
	log.Info("connected to db", "host", cfg.DB.Host, "db", cfg.DB.DBName)

	if err := app.New(cfg, log, app.Pinger(db)).Run(ctx); err != nil {
		log.Error("App stopped with error", "error", err)
		os.Exit(1)
	}
	log.Info("App stopped")

}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}
