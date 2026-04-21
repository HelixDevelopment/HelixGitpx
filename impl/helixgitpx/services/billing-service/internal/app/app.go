// Package app is the composition root for billing-service.
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/handler"
	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/provider"
	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/usecase"
	"github.com/helixgitpx/platform/log"
)

func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("BILLING_HTTP_ADDR", ":8080")

	prov := &provider.Stripe{APIKey: os.Getenv("STRIPE_API_KEY")}
	uc := &usecase.UseCases{Prov: prov}
	h := &handler.Handler{UseCases: uc}

	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("billing-service listening", "addr", addr)
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
	lg.Info("billing-service stopped")
	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
