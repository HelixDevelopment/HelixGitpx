// Package app is the composition root for webhook-gateway.
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/webhook-gateway/internal/handler"
	"github.com/helixgitpx/platform/log"
)

func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("WEBHOOK_HTTP_ADDR", ":8080")

	// Secret source: env var for now; a Vault-backed provider plugs in
	// via the same SecretProvider signature when platform/config/vault is
	// wired.
	secret := os.Getenv("WEBHOOK_SHARED_SECRET")
	secrets := func(_, _ string) (string, bool) {
		if secret == "" {
			return "", false
		}
		return secret, true
	}

	h := &handler.Handler{Secrets: secrets, Recorder: &handler.InMemoryRecorder{}}

	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("webhook-gateway listening", "addr", addr)
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
	lg.Info("webhook-gateway stopped")
	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
