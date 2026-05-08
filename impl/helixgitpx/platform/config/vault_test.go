package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/helixgitpx/platform/config"
)

type cfgWithVault struct {
	DSN string `env:"DSN" vault:"kv/hello#dsn" default:"postgres://default"`
}

type cfgWithTestSecret struct {
	Secret string `env:"SECRET" vault:"kv/test#secret_value"`
}

func TestLoad_VaultFallsBackToDefaultWhenAddrUnset(t *testing.T) {
	prevAddr, hadAddr := os.LookupEnv("VAULT_ADDR")
	prevToken, hadToken := os.LookupEnv("VAULT_TOKEN")
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")
	t.Cleanup(func() {
		if hadAddr {
			os.Setenv("VAULT_ADDR", prevAddr)
		}
		if hadToken {
			os.Setenv("VAULT_TOKEN", prevToken)
		}
	})
	var c cfgWithVault
	if err := config.Load(&c, config.Options{Prefix: "X"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.DSN != "postgres://default" {
		t.Errorf("DSN = %q, want default", c.DSN)
	}
}

func TestVaultResolver_ReadsSecretFromRealVault(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("SKIP-OK: #VAULT — set VAULT_ADDR and VAULT_TOKEN to run")
	}

	r := config.NewVaultResolver()
	if r == nil {
		t.Fatal("NewVaultResolver returned nil with VAULT_ADDR and VAULT_TOKEN set")
	}

	val, err := r.Read(context.Background(), "kv/test#secret_value")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if val != "test-secret-42" {
		t.Fatalf("want 'test-secret-42' got %q", val)
	}
}

func TestVaultResolver_InvalidPath(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("SKIP-OK: #VAULT — set VAULT_ADDR and VAULT_TOKEN to run")
	}

	r := config.NewVaultResolver()
	if r == nil {
		t.Skip("SKIP-OK: #VAULT — resolver nil")
	}

	_, err := r.Read(context.Background(), "kv/test#nonexistent_key")
	if err == nil {
		t.Fatal("expected error for nonexistent key")
	}
}

func TestVaultResolver_MissingHash(t *testing.T) {
	r := &config.VaultResolver{Addr: "http://localhost:8200", Token: "root"}
	_, err := r.Read(context.Background(), "kv/test-no-hash")
	if err == nil {
		t.Fatal("expected error for path without #key")
	}
}

func TestLoad_VaultTagResolvesSecret(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("SKIP-OK: #VAULT — set VAULT_ADDR and VAULT_TOKEN to run")
	}

	var c cfgWithTestSecret
	if err := config.Load(&c, config.Options{Prefix: "X"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.Secret != "test-secret-42" {
		t.Fatalf("want secret='test-secret-42' got %q", c.Secret)
	}
}
