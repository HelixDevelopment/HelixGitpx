//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// TestAuth_WhoAmI_ForbiddenWithoutToken verifies the auth-service rejects
// an unauthenticated request to /whoami with 401/403.
func TestAuth_WhoAmI_ForbiddenWithoutToken(t *testing.T) {
	base := os.Getenv("HELIXGITPX_AUTH_URL")
	if base == "" {
		t.Fatal("HELIXGITPX_AUTH_URL must be set — run `make compose-up` first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/v1/whoami", nil)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 401/403 unauthenticated, got %d", resp.StatusCode)
	}
}

// TestAuth_WhoAmI_AcceptsValidToken requires HELIXGITPX_INTEGRATION_TOKEN
// (minted out-of-band against the ephemeral Keycloak). The test fails if the
// token is missing — NO mocks, per Constitution §II §2.
func TestAuth_WhoAmI_AcceptsValidToken(t *testing.T) {
	base := os.Getenv("HELIXGITPX_AUTH_URL")
	token := os.Getenv("HELIXGITPX_INTEGRATION_TOKEN")
	if base == "" || token == "" {
		t.Fatal("HELIXGITPX_AUTH_URL and HELIXGITPX_INTEGRATION_TOKEN are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/v1/whoami", nil)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !strings.Contains(body.Email, "@") {
		t.Fatalf("expected email in response, got %q", body.Email)
	}
}
