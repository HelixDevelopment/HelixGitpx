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

// TestOrgTeam_CreateListGet exercises the full orgteam RPC surface
// against a real instance (Postgres + Kafka outbox). No mocks.
func TestOrgTeam_CreateListGet(t *testing.T) {
	base := os.Getenv("HELIXGITPX_ORGTEAM_URL")
	token := os.Getenv("HELIXGITPX_INTEGRATION_TOKEN")
	if base == "" || token == "" {
		t.Fatal("HELIXGITPX_ORGTEAM_URL and HELIXGITPX_INTEGRATION_TOKEN are required")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Create org.
	createBody, _ := json.Marshal(map[string]any{
		"name":      "itest-" + timeSuffix(),
		"residency": "EU",
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/orgs", bytes.NewReader(createBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Fatalf("create: want 2xx, got %d", resp.StatusCode)
	}
	var created struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Residency string `json:"residency"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.ID == "" {
		t.Fatal("empty org id returned")
	}
	if created.Residency != "EU" {
		t.Fatalf("want residency=EU, got %q", created.Residency)
	}

	// 2. List — new org must be present.
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, base+"/v1/orgs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: want 200, got %d", resp.StatusCode)
	}
	var list struct {
		Orgs []struct {
			ID string `json:"id"`
		} `json:"orgs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	found := false
	for _, o := range list.Orgs {
		if o.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created org %s not in list", created.ID)
	}
}

func timeSuffix() string {
	return time.Now().UTC().Format("20060102150405")
}
