package config_test

import (
	"os"
	"testing"

	"github.com/helixgitpx/platform/config"
)

type cfgWithVault struct {
	DSN string `env:"DSN" vault:"kv/hello#dsn" default:"postgres://default"`
}

func TestLoad_VaultFallsBackToDefaultWhenAddrUnset(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")
	var c cfgWithVault
	if err := config.Load(&c, config.Options{Prefix: "X"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.DSN != "postgres://default" {
		t.Errorf("DSN = %q, want default", c.DSN)
	}
}
