// Package app is the composition root for orgteam-service.
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/handler"
	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/memstore"
	"github.com/helixgitpx/platform/log"
)

// Run wires the orgteam HTTP server and serves until ctx is done.
//
// Authentication is a cross-cutting concern normally provided by the
// auth middleware; for local dev we accept an `X-Actor` header. A real
// deployment wires the platform/auth UnaryInterceptor in front of this
// mux and rejects requests without a valid SVID + JWT.
func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("ORGTEAM_HTTP_ADDR", ":8080")
	store := memstore.New()

	actor := func(ctx context.Context) string {
		req, ok := ctx.Value(requestKey{}).(*http.Request)
		if !ok {
			return ""
		}
		return req.Header.Get("X-Actor")
	}

	mux := http.NewServeMux()
	(&handler.OrgHandler{Store: store, ActorFromCtx: actor}).RegisterOrgRoutes(mux)
	mux.Handle("POST /v1/orgs/{id}/residency", &handler.ResidencyHandler{
		Repo:         store,
		ActorFromCtx: actor,
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           withRequestInContext(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("orgteam-service listening", "addr", addr)
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
	lg.Info("orgteam-service stopped")
	return nil
}

type requestKey struct{}

func withRequestInContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestKey{}, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
