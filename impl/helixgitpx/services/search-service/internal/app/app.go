// Package app is the composition root for search-service.
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/search-service/internal/engines"
	"github.com/helixgitpx/helixgitpx/services/search-service/internal/handler"
	"github.com/helixgitpx/platform/log"
)

// Run wires the search-service HTTP server and serves until ctx is done.
//
// Engines (Meilisearch, Qdrant, Zoekt) are registered via env vars at
// boot time. When no endpoint is configured, search returns empty results.
func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("SEARCH_HTTP_ADDR", ":8080")

	var es []engines.Engine
	// Future: register Meilisearch, Qdrant, Zoekt clients here based on env.

	h := &handler.Handler{Engines: es, Timeout: 2 * time.Second}
	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("search-service listening", "addr", addr, "engines", len(es))
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
	lg.Info("search-service stopped")
	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
