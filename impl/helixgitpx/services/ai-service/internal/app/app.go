package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/helixgitpx/helixgitpx/services/ai-service/internal/handler"
	"github.com/helixgitpx/helixgitpx/services/ai-service/internal/usecase"
	"github.com/helixgitpx/platform/log"
)

type echoLLM struct{}

func (e *echoLLM) Prompt(_ context.Context, _, prompt string) (string, error) {
	return "echo:" + prompt, nil
}

func Run(ctx context.Context, lg *log.Logger) error {
	addr := envOrDefault("AI_HTTP_ADDR", ":8080")

	llm := &echoLLM{}
	uc := &usecase.UseCases{LLM: llm}
	h := handler.NewHandler(uc)

	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		lg.Info("ai-service listening", "addr", addr)
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
	lg.Info("ai-service stopped")
	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
