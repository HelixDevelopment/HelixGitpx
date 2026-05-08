package pg_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/helixgitpx/platform/pg"
)

func TestOpen_InvalidDSNFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := pg.Open(ctx, pg.Options{DSN: "not-a-valid-dsn"})
	if err == nil {
		t.Fatalf("expected error for invalid DSN")
	}
}

func TestIsUnavailable(t *testing.T) {
	if !pg.IsUnavailable(pg.ErrUnavailable) {
		t.Errorf("ErrUnavailable not classified as unavailable")
	}
	if pg.IsUnavailable(errors.New("other")) {
		t.Errorf("arbitrary error classified as unavailable")
	}
}

func TestProbe_NilPool(t *testing.T) {
	probe := pg.Probe(nil)
	err := probe(context.Background())
	if err == nil {
		t.Fatal("expected error for nil pool probe")
	}
	if !pg.IsUnavailable(err) {
		t.Errorf("probe(nil) err = %v, want ErrUnavailable", err)
	}
}

func TestOpen_UnreachableHost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := pg.Open(ctx, pg.Options{DSN: "postgres://nobody:nopass@localhost:19999/nonexistent?sslmode=disable"})
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
	if !pg.IsUnavailable(err) {
		t.Errorf("Open(unreachable) err = %v, want ErrUnavailable wrapping", err)
	}
}

func TestOpen_RealPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		dsn = "postgres://helix:helix@localhost:15432/helixgitpx?sslmode=disable"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pg.Open(ctx, pg.Options{DSN: dsn, MaxConns: 4, MinConns: 1})
	if err != nil {
		t.Skipf("SKIP-OK: #integration — Postgres not available at %s: %v", dsn, err)
	}
	defer pool.Close()

	if pool == nil {
		t.Fatal("pool is nil after successful Open")
	}

	var result int
	if err := pool.QueryRow(ctx, "SELECT 1").Scan(&result); err != nil {
		t.Fatalf("SELECT 1 failed: %v", err)
	}
	if result != 1 {
		t.Errorf("SELECT 1 = %d, want 1", result)
	}
}

func TestProbe_RealPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		dsn = "postgres://helix:helix@localhost:15432/helixgitpx?sslmode=disable"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pg.Open(ctx, pg.Options{DSN: dsn})
	if err != nil {
		t.Skipf("SKIP-OK: #integration — Postgres not available at %s: %v", dsn, err)
	}
	defer pool.Close()

	probe := pg.Probe(pool)
	if probeErr := probe(ctx); probeErr != nil {
		t.Errorf("probe(healthy pool) = %v, want nil", probeErr)
	}
}

func TestOpen_OptionsApplied(t *testing.T) {
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		dsn = "postgres://helix:helix@localhost:15432/helixgitpx?sslmode=disable"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pg.Open(ctx, pg.Options{
		DSN:                 dsn,
		MaxConns:            2,
		MinConns:            1,
		ConnectTimeout:      3 * time.Second,
		HealthCheckInterval: 10 * time.Second,
	})
	if err != nil {
		t.Skipf("SKIP-OK: #integration — Postgres not available at %s: %v", dsn, err)
	}
	defer pool.Close()

	stat := pool.Stat()
	if stat.MaxConns() != 2 {
		t.Errorf("MaxConns = %d, want 2", stat.MaxConns())
	}
}
