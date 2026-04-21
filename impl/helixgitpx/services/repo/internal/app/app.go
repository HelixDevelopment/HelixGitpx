// Package app is the composition root for repo-service.
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/repo/internal/handler"
	"github.com/helixgitpx/helixgitpx/services/repo/internal/memstore"
	"github.com/helixgitpx/platform/log"
)

// Run wires the repo-service HTTP server and serves until ctx is done.
// Store backend defaults to the in-memory implementation; a Postgres
// backend lands when the repo.repos schema bootstrap is wired.
func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("REPO_HTTP_ADDR", ":8080")

	store := memstore.New()
	h := &handler.Handler{Store: store}

	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("repo-service listening", "addr", addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	lg.Info("repo-service stopped")
	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
