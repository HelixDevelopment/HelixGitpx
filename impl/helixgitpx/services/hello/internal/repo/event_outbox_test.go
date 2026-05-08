package repo

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func openPoolOutbox(t *testing.T) *pgxpool.Pool {
	t.Helper()
	d := os.Getenv("TEST_PG_DSN")
	if d == "" {
		d = "postgres://helix:helix@localhost:15432/helixgitpx?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, d)
	if err != nil {
		t.Fatalf("open pool: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("SKIP-OK: #integration — Postgres not available: %v", err)
	}
	return pool
}

func TestEventOutbox_Emit_InsertsRow(t *testing.T) {
	pool := openPoolOutbox(t)
	ctx := context.Background()

	e := &EventOutbox{Pool: pool, Topic: "test.hello.said"}
	if err := e.Emit(ctx, "world", "hello, world", int64(1)); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	var aggregateID, topic string
	var payload []byte
	err := pool.QueryRow(ctx, `
		SELECT aggregate_id, topic, payload FROM hello.outbox_events
		WHERE aggregate_id = $1 ORDER BY created_at DESC LIMIT 1`, "world",
	).Scan(&aggregateID, &topic, &payload)
	if err != nil {
		t.Fatalf("query outbox: %v", err)
	}
	if aggregateID != "world" {
		t.Errorf("aggregate_id = %q, want world", aggregateID)
	}
	if topic != "test.hello.said" {
		t.Errorf("topic = %q, want test.hello.said", topic)
	}

	var p struct {
		Name     string `json:"name"`
		Greeting string `json:"greeting"`
		Count    int64  `json:"count"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if p.Name != "world" {
		t.Errorf("payload.Name = %q, want world", p.Name)
	}
	if p.Greeting != "hello, world" {
		t.Errorf("payload.Greeting = %q", p.Greeting)
	}
	if p.Count != 1 {
		t.Errorf("payload.Count = %d, want 1", p.Count)
	}

	// cleanup
	_, _ = pool.Exec(ctx, "DELETE FROM hello.outbox_events WHERE aggregate_id = $1", "world")
}

func TestEventOutbox_Emit_DefaultTopic(t *testing.T) {
	pool := openPoolOutbox(t)
	ctx := context.Background()

	e := &EventOutbox{Pool: pool} // Topic is empty
	if err := e.Emit(ctx, "testuser", "hello, testuser", int64(5)); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	var topic string
	err := pool.QueryRow(ctx, `
		SELECT topic FROM hello.outbox_events
		WHERE aggregate_id = $1 LIMIT 1`, "testuser",
	).Scan(&topic)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if topic != "hello.said" {
		t.Errorf("default topic = %q, want hello.said", topic)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM hello.outbox_events WHERE aggregate_id = $1", "testuser")
}

func TestEventOutbox_Emit_CountPreserved(t *testing.T) {
	pool := openPoolOutbox(t)
	ctx := context.Background()

	e := &EventOutbox{Pool: pool, Topic: "test.count"}
	if err := e.Emit(ctx, "countuser", "hi", int64(42)); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	var payload []byte
	_ = pool.QueryRow(ctx, `
		SELECT payload FROM hello.outbox_events WHERE aggregate_id = $1`, "countuser",
	).Scan(&payload)

	var p struct {
		Count int64 `json:"count"`
	}
	_ = json.Unmarshal(payload, &p)
	if p.Count != 42 {
		t.Errorf("count = %d, want 42", p.Count)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM hello.outbox_events WHERE aggregate_id = $1", "countuser")
}
