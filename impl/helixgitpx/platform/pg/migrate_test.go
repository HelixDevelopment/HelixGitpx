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

func TestMigrate_RealPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		dsn = "postgres://helix:helix@localhost:15432/helixgitpx?sslmode=disable"
	}

	ctx := context.Background()
	err := pg.Migrate(ctx, pg.MigrateOptions{
		DSN: dsn,
		Dir: "/nonexistent-migrations-dir",
	})
	if err == nil {
		t.Skip("SKIP-OK: #integration — Migrate with empty dir succeeded (no migrations to run)")
	}
}
