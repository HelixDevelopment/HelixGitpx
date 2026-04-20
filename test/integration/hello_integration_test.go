//go:build integration

// Package integration exercises real dependencies (Postgres, Kafka, Keycloak)
// via the compose stack. NEVER use mocks here — that violates Constitution §II.
package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestHelloService_HealthCheck hits the real hello service.
func TestHelloService_HealthCheck(t *testing.T) {
	base := os.Getenv("HELIXGITPX_HELLO_URL")
	if base == "" {
		t.Fatal("HELIXGITPX_HELLO_URL must be set — run `make compose-up` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/healthz", nil)
	if err != nil {
		t.Fatalf("request build: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 OK, got %d", resp.StatusCode)
	}
}
