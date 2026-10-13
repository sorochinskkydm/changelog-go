package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sorochinskkydm/changelog-go/internal/config"
)

// fakeDB — подделка базы для теста: возвращает ту ошибку, которую мы задали.
type fakeDB struct {
	err error
}

func (f fakeDB) Ping(ctx context.Context) error {
	return f.err
}

func TestHealth(t *testing.T) {
	tests := []struct {
		name       string
		dbErr      error
		wantStatus int
		wantBody   string
	}{
		{"db up", nil, http.StatusOK, `{"status":"ok","db":"up"}` + "\n"},
		{"db down", errors.New("connection refused"), http.StatusServiceUnavailable, `{"status":"error","db":"down"}` + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(io.Discard, nil))
			a := New(&config.Config{}, log, fakeDB{err: tt.dbErr})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
			rec := httptest.NewRecorder()
			a.routes().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if rec.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
