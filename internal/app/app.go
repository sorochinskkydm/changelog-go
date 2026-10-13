package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sorochinskkydm/changelog-go/internal/config"
	"github.com/sorochinskkydm/changelog-go/internal/httpx"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	db 		Pinger
	server *http.Server
}

type HealthResponse struct {
	Status string `json:"status"`
	DB string `json:"db"`
}

func (app *App) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := app.db.Ping(ctx); err != nil {
		app.log.Error("Health check failed", "error", err)
		httpx.JSON(w, http.StatusServiceUnavailable, HealthResponse{Status: "error", DB: "down"})
		return 
	}

	httpx.JSON(w, http.StatusOK, HealthResponse{Status: "ok", DB: "up"})
}

func (app *App) Run(ctx context.Context) error {
	errorCh := make(chan error, 1)

	go func() {
		app.log.Info("server started", "addr", app.server.Addr)
		if err := app.server.ListenAndServe(); err != nil {
			errorCh <- err
		}
	}()

	select {
	case err := <-errorCh:
		return fmt.Errorf("listen :%w", err)
	case <-ctx.Done():
	}

	app.log.Info("Shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.cfg.HTTP.ShutdownTimeout)

	defer cancel()

	return app.server.Shutdown(shutdownCtx)

}

func New(cfg *config.Config, log *slog.Logger, db Pinger) *App {
	a := &App{cfg: cfg, log: log, db: db}
	a.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:           a.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return a
}

func (a *App) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(httpx.Logger(a.log))

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", a.health)
	})
	return router
}
