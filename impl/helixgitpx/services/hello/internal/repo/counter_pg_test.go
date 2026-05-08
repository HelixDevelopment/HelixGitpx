package repo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func dsn(t *testing.T) string {
	t.Helper()
	d := os.Getenv("TEST_PG_DSN")
	if d == "" {
		d = "postgres://helix:helix@localhost:15432/helixgitpx?sslmode=disable"
	}
	return d
}

func openPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn(t))
	if err != nil {
		t.Fatalf("open pool: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("SKIP-OK: #integration — Postgres not available: %v", err)
	}
	return pool
}

func TestCounterPG_Increment_FirstCall(t *testing.T) {
	pool := openPool(t)
	ctx := context.Background()

	name := "test_counter_first_" + time.Now().Format("20060102150405")
	c := &CounterPG{Pool: pool}

	n, err := c.Increment(ctx, name)
	if err != nil {
		t.Fatalf("Increment: %v", err)
	}
	if n != 1 {
		t.Errorf("first increment = %d, want 1", n)
	}

	// cleanup
	_, _ = pool.Exec(ctx, "DELETE FROM hello.greetings WHERE name = $1", name)
}

func TestCounterPG_Increment_MonotonicallyIncreasing(t *testing.T) {
	pool := openPool(t)
	ctx := context.Background()

	name := "test_counter_mono_" + time.Now().Format("20060102150405")
	c := &CounterPG{Pool: pool}

	n1, _ := c.Increment(ctx, name)
	n2, _ := c.Increment(ctx, name)
	n3, _ := c.Increment(ctx, name)

	if n1 != 1 {
		t.Errorf("n1 = %d, want 1", n1)
	}
	if n2 != 2 {
		t.Errorf("n2 = %d, want 2", n2)
	}
	if n3 != 3 {
		t.Errorf("n3 = %d, want 3", n3)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM hello.greetings WHERE name = $1", name)
}

func TestCounterPG_Increment_DifferentNames(t *testing.T) {
	pool := openPool(t)
	ctx := context.Background()

	ts := time.Now().Format("20060102150405")
	c := &CounterPG{Pool: pool}

	n1, _ := c.Increment(ctx, "test_counter_a_"+ts)
	n2, _ := c.Increment(ctx, "test_counter_b_"+ts)

	if n1 != 1 || n2 != 1 {
		t.Errorf("different names should both start at 1: n1=%d n2=%d", n1, n2)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM hello.greetings WHERE name LIKE $1", "test_counter_%_"+ts)
}
