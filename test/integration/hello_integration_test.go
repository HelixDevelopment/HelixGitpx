//go:build integration

// Package integration exercises real dependencies (Postgres, Kafka, Keycloak)
// via the compose stack. NEVER use mocks here — that violates Constitution §II.
package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func mustEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Fatalf("%s must be set — run `make compose-up` first", key)
	}
	return v
}

func getJSON(t *testing.T, ctx context.Context, url string) (*http.Response, map[string]interface{}) {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("request build: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("parse JSON: %v (body: %s)", err, body)
	}
	return resp, parsed
}

// TestHelloService_Greeting verifies the core business logic: the /v1/hello
// endpoint returns a personalized greeting with an incrementing counter.
// This is the primary integration test — a healthz-only check would be a
// bluff per CONST-035.
func TestHelloService_Greeting(t *testing.T) {
	base := mustEnv(t, "HELIXGITPX_HELLO_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, body := getJSON(t, ctx, base+"/v1/hello?name=IntegrationTest")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 OK, got %d (body: %v)", resp.StatusCode, body)
	}

	greeting, _ := body["greeting"].(string)
	if greeting == "" {
		t.Fatal("greeting must not be empty")
	}
	if !strings.Contains(greeting, "IntegrationTest") {
		t.Fatalf("greeting must contain the name 'IntegrationTest', got %q", greeting)
	}

	count, ok := body["count"].(float64)
	if !ok {
		t.Fatal("count field missing or not a number")
	}
	if count < 1 {
		t.Fatalf("count must be >= 1 after a greeting, got %v", count)
	}
}

// TestHelloService_CounterIncrements verifies the counter is monotonic
// across sequential requests — proving the Postgres-backed counter works.
func TestHelloService_CounterIncrements(t *testing.T) {
	base := mustEnv(t, "HELIXGITPX_HELLO_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	name := "MonotonicCheck_" + time.Now().Format("150405")
	url := base + "/v1/hello?name=" + name

	_, body1 := getJSON(t, ctx, url)
	count1 := body1["count"].(float64)

	_, body2 := getJSON(t, ctx, url)
	count2 := body2["count"].(float64)

	if count2 <= count1 {
		t.Fatalf("counter must be monotonically increasing: first=%v second=%v", count1, count2)
	}
}

// TestHelloService_EmptyNameRejected verifies input validation:
// an empty name must produce a 400 error, not a 200.
func TestHelloService_EmptyNameRejected(t *testing.T) {
	base := mustEnv(t, "HELIXGITPX_HELLO_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, _ := getJSON(t, ctx, base+"/v1/hello?name=")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400 for empty name, got %d", resp.StatusCode)
	}
}

// TestHelloService_HealthCheck verifies the service is alive.
// This is supplementary — the business logic tests above are the real proof.
func TestHelloService_HealthCheck(t *testing.T) {
	base := mustEnv(t, "HELIXGITPX_HELLO_URL")
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
