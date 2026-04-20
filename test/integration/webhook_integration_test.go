//go:build integration

package integration

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestWebhook_AcceptsValidSignature POSTs a synthetic GitHub push event
// signed with the shared HMAC secret. Real webhook-gateway must return 2xx.
func TestWebhook_AcceptsValidSignature(t *testing.T) {
	base := os.Getenv("HELIXGITPX_WEBHOOK_URL")
	secret := os.Getenv("HELIXGITPX_WEBHOOK_SECRET")
	if base == "" || secret == "" {
		t.Fatal("HELIXGITPX_WEBHOOK_URL and HELIXGITPX_WEBHOOK_SECRET are required")
	}

	payload := []byte(`{"ref":"refs/heads/main","repository":{"full_name":"o/r"}}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/webhooks/github", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", sig)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		t.Fatalf("want 2xx, got %d", resp.StatusCode)
	}
}

// TestWebhook_RejectsTamperedPayload ensures signature verification catches
// a modified body. This is the attack we defend against — constitution §II.
func TestWebhook_RejectsTamperedPayload(t *testing.T) {
	base := os.Getenv("HELIXGITPX_WEBHOOK_URL")
	secret := os.Getenv("HELIXGITPX_WEBHOOK_SECRET")
	if base == "" || secret == "" {
		t.Fatal("HELIXGITPX_WEBHOOK_URL and HELIXGITPX_WEBHOOK_SECRET are required")
	}

	original := []byte(`{"ref":"refs/heads/main","repository":{"full_name":"o/r"}}`)
	tampered := []byte(`{"ref":"refs/heads/evil","repository":{"full_name":"o/r"}}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(original)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/webhooks/github", bytes.NewReader(tampered))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", sig)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 401/400 on tamper, got %d", resp.StatusCode)
	}
}
