package pg_test

import (
	"context"
	"os"
	"testing"

	"github.com/helixgitpx/platform/pg"
)

func TestMigrate_InvalidDSN(t *testing.T) {
	err := pg.Migrate(context.Background(), pg.MigrateOptions{
		DSN: "postgres://invalid-host-name-that-does-not-exist:5432/db?sslmode=disable",
		Dir: "/tmp/does-not-matter-test-will-fail-earlier",
	})
	if err == nil {
		t.Fatal("expected error for invalid DSN")
	}
}

func TestOpen_RealPostgres_PoolsAndPings(t *testing.T) {
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		t.Skip("SKIP-OK: #PG — set TEST_PG_DSN to run")
	}

	ctx := context.Background()
	pool, err := pg.Open(ctx, pg.Options{DSN: dsn})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Ping after Open: %v", err)
	}
}

func TestOpen_InvalidDSN_ReturnsUnavailable(t *testing.T) {
	_, err := pg.Open(context.Background(), pg.Options{
		DSN: "postgres://invalid-host:5432/db?sslmode=disable",
	})
	if err == nil {
		t.Fatal("expected error for invalid DSN")
	}
	if !pg.IsUnavailable(err) {
		t.Fatalf("expected ErrUnavailable, got: %v", err)
	}
}

func TestProbe_ReturnsErrorOnNilPool(t *testing.T) {
	probe := pg.Probe(nil)
	if err := probe(context.Background()); err == nil {
		t.Fatal("expected error from nil pool probe")
	}
}
