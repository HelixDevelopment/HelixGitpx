package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/handler"
	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/providers/github"
	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/providers/gitlab"
	"github.com/helixgitpx/platform/log"
)

func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("ADAPTER_POOL_HTTP_ADDR", ":8080")

	reg := handler.NewRegistry()
	reg.Register(adapter.GitHub, &github.Adapter{})
	reg.Register(adapter.GitLab, &gitlab.Adapter{})

	h := &handler.Handler{Registry: reg}

	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("adapter-pool listening", "addr", addr)
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
	lg.Info("adapter-pool stopped")
	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
