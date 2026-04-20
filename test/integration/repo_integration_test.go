//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRepo_CreateAndBindUpstream(t *testing.T) {
	base := os.Getenv("HELIXGITPX_REPO_URL")
	orgID := os.Getenv("HELIXGITPX_INTEGRATION_ORG_ID")
	token := os.Getenv("HELIXGITPX_INTEGRATION_TOKEN")
	if base == "" || orgID == "" || token == "" {
		t.Fatal("HELIXGITPX_REPO_URL, HELIXGITPX_INTEGRATION_ORG_ID, HELIXGITPX_INTEGRATION_TOKEN required")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create repo.
	body, _ := json.Marshal(map[string]any{
		"name":   "itest-repo-" + timeSuffix(),
		"org_id": orgID,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/repos", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Fatalf("create repo: want 2xx, got %d", resp.StatusCode)
	}
	var repo struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if repo.ID == "" {
		t.Fatal("no repo id")
	}

	// Bind upstream.
	bind, _ := json.Marshal(map[string]any{
		"provider":  "generic",
		"url":       "https://example.invalid/owner/repo.git",
		"direction": "write",
	})
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/repos/"+repo.ID+"/bindings", bytes.NewReader(bind))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Fatalf("bind: want 2xx, got %d", resp.StatusCode)
	}
}
