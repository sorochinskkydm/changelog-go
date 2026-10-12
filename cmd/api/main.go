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

	if err := app.New(cfg, log).Run(ctx); err != nil {
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