// Package app is the composition root for opa-bundle-server.
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/handler"
	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/store"
	"github.com/helixgitpx/platform/log"
)

// Run wires the opa-bundle-server HTTP server and serves until ctx is done.
func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("OPA_BUNDLE_HTTP_ADDR", ":8080")

	s := store.New()
	if path := os.Getenv("OPA_BUNDLE_SEED"); path != "" {
		if err := seedFromFile(s, path); err != nil {
			lg.Error("seed failed", "err", err.Error())
		} else {
			lg.Info("seeded bundle from disk", "path", path)
		}
	}

	h := &handler.Handler{Store: s}
	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("opa-bundle-server listening", "addr", addr)
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
	lg.Info("opa-bundle-server stopped")
	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func seedFromFile(s *store.Store, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	meta, err := domain.NewBundle("seed", "1.0.0", "local", content, false, time.Now())
	if err != nil {
		return err
	}
	meta = s.Put(meta, content)
	meta.ActivatedAt = time.Now()
	return s.Activate(meta.ID, meta)
}
