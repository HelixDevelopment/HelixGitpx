// Package config (vault.go) provides a Vault KV v2 resolver invoked by Load
// when a struct field carries a `vault:"path/to/key"` tag. The HTTP call
// targets $VAULT_ADDR with the token at $VAULT_TOKEN (populated by Vault
// Agent Injector in-cluster, or set manually for dev).
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// VaultResolver fetches secrets from Vault KV v2.
type VaultResolver struct {
	Addr   string
	Token  string
	Client *http.Client
}

// NewVaultResolver reads VAULT_ADDR and VAULT_TOKEN from env.
// Returns nil when either is unset (caller treats it as no-op).
func NewVaultResolver() *VaultResolver {
	addr := os.Getenv("VAULT_ADDR")
	tok := os.Getenv("VAULT_TOKEN")
	if addr == "" || tok == "" {
		return nil
	}
	return &VaultResolver{
		Addr:   strings.TrimRight(addr, "/"),
		Token:  tok,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Read fetches kv/data/<path>'s data[key]. Expects KV v2 layout.
// path must be in the form "<mount>/<secret-path>#<key>", e.g. "kv/hello#dsn".
func (r *VaultResolver) Read(ctx context.Context, path string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("vault: resolver is nil")
	}
	parts := strings.SplitN(path, "#", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("vault: expected mount/path#key, got %q", path)
	}
	kvPath := parts[0]
	key := parts[1]

	i := strings.Index(kvPath, "/")
	if i < 0 {
		return "", fmt.Errorf("vault: expected <mount>/<path>, got %q", kvPath)
	}
	url := fmt.Sprintf("%s/v1/%s/data/%s", r.Addr, kvPath[:i], kvPath[i+1:])

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", r.Token)
	resp, err := r.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault: %s", resp.Status)
	}

	var body struct {
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	val, ok := body.Data.Data[key]
	if !ok {
		return "", fmt.Errorf("vault: key %q not in %s", key, kvPath)
	}
	return val, nil
}
