package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CounterPG implements domain.Counter using Postgres UPSERT.
type CounterPG struct {
	Pool *pgxpool.Pool
}

func (c *CounterPG) Increment(ctx context.Context, name string) (int64, error) {
	var n int64
	err := c.Pool.QueryRow(ctx, `
        INSERT INTO hello.greetings(name, count, last_said_at)
        VALUES ($1, 1, NOW())
        ON CONFLICT (name) DO UPDATE
          SET count = hello.greetings.count + 1,
              last_said_at = NOW()
        RETURNING count`, name).Scan(&n)
	return n, err
}
