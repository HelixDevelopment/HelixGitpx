// Package pg wraps pgx/v5 with HelixGitpx-specific defaults and a thin
// migration runner (via pressly/goose). Callers get a *pgxpool.Pool ready
// for use; the package exposes typed sentinel errors so callers can classify
// failures without string matching.
package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrUnavailable signals the database is unreachable.
var ErrUnavailable = errors.New("pg: unavailable")

// Options configures Open.
type Options struct {
	DSN                 string
	MaxConns            int32
	MinConns            int32
	ConnectTimeout      time.Duration
	HealthCheckInterval time.Duration
}

// Open constructs a pool. Callers own Close().
func Open(ctx context.Context, opts Options) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("pg: parse DSN: %w", err)
	}
	if opts.MaxConns > 0 {
		cfg.MaxConns = opts.MaxConns
	}
	if opts.MinConns > 0 {
		cfg.MinConns = opts.MinConns
	}
	if opts.ConnectTimeout > 0 {
		cfg.ConnConfig.ConnectTimeout = opts.ConnectTimeout
	}
	if opts.HealthCheckInterval > 0 {
		cfg.HealthCheckPeriod = opts.HealthCheckInterval
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pg: new pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.Join(ErrUnavailable, err)
	}
	return pool, nil
}

// IsUnavailable reports whether err wraps ErrUnavailable.
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }

// Probe returns a function that pings the pool for health checks.
func Probe(pool *pgxpool.Pool) func(context.Context) error {
	return func(ctx context.Context) error {
		if pool == nil {
			return ErrUnavailable
		}
		return pool.Ping(ctx)
	}
}
