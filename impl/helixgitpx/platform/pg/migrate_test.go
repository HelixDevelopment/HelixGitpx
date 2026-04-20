package pg_test

import (
	"context"
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
