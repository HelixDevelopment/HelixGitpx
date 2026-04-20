// Package pg (migrate.go) wraps github.com/pressly/goose/v3 to apply
// SQL migrations from a directory to a DSN. Used by services' migrate-job.
package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// MigrateOptions configures Migrate.
type MigrateOptions struct {
	DSN string
	Dir string // filesystem path containing *.sql migrations
}

// Migrate applies all Up migrations under opts.Dir to opts.DSN.
// Idempotent — re-applying after completion is a no-op.
func Migrate(ctx context.Context, opts MigrateOptions) error {
	db, err := sql.Open("pgx", opts.DSN)
	if err != nil {
		return fmt.Errorf("pg.Migrate: open: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("pg.Migrate: dialect: %w", err)
	}
	if err := goose.UpContext(ctx, db, opts.Dir); err != nil {
		return fmt.Errorf("pg.Migrate: up: %w", err)
	}
	return nil
}
