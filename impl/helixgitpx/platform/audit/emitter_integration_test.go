package audit_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/helixgitpx/platform/audit"
)

func auditPool(t *testing.T) *pgxpool.Pool {
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

func TestEmitter_Emit_InsertsIntoOutbox(t *testing.T) {
	pool := auditPool(t)
	ctx := context.Background()
	em := &audit.Emitter{Pool: pool, OutboxFQN: "hello.outbox_events"}

	ev := audit.Event{
		ActorUserID: "user-42",
		ActorIP:     "10.0.0.1",
		Action:      "org.create",
		Target:      "acme-corp",
		Details:     map[string]any{"plan": "pro"},
	}
	if err := em.Emit(ctx, ev); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	var (
		topic   string
		payload []byte
	)
	err := pool.QueryRow(ctx,
		`SELECT topic, payload FROM hello.outbox_events WHERE aggregate_id = $1 ORDER BY created_at DESC LIMIT 1`,
		"acme-corp",
	).Scan(&topic, &payload)
	if err != nil {
		t.Fatalf("read back outbox row: %v", err)
	}
	if topic != "audit.events" {
		t.Errorf("topic = %q, want audit.events", topic)
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("payload not JSON: %v", err)
	}
	if decoded["actor_user_id"] != "user-42" {
		t.Errorf("actor_user_id = %v, want user-42", decoded["actor_user_id"])
	}
	if decoded["action"] != "org.create" {
		t.Errorf("action = %v, want org.create", decoded["action"])
	}
	if decoded["target"] != "acme-corp" {
		t.Errorf("target = %v, want acme-corp", decoded["target"])
	}
	details, ok := decoded["details"].(map[string]any)
	if !ok || details["plan"] != "pro" {
		t.Errorf("details.plan = %v, want pro", decoded["details"])
	}

	_, _ = pool.Exec(ctx, "DELETE FROM hello.outbox_events WHERE aggregate_id = $1", "acme-corp")
}

func TestEmitter_EmitInTx_WithinTransaction(t *testing.T) {
	pool := auditPool(t)
	ctx := context.Background()
	em := &audit.Emitter{Pool: pool, OutboxFQN: "hello.outbox_events"}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		t.Fatalf("BeginTx: %v", err)
	}
	defer tx.Rollback(ctx)

	ev := audit.Event{
		ActorUserID: "user-99",
		Action:      "repo.delete",
		Target:      "repo-x",
	}
	if err := em.EmitInTx(ctx, tx, ev); err != nil {
		t.Fatalf("EmitInTx: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	var payload []byte
	err = pool.QueryRow(ctx,
		`SELECT payload FROM hello.outbox_events WHERE aggregate_id = $1`, "repo-x",
	).Scan(&payload)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	var decoded map[string]any
	_ = json.Unmarshal(payload, &decoded)
	if decoded["action"] != "repo.delete" {
		t.Errorf("action = %v, want repo.delete", decoded["action"])
	}
	if decoded["actor_user_id"] != "user-99" {
		t.Errorf("actor_user_id = %v, want user-99", decoded["actor_user_id"])
	}

	_, _ = pool.Exec(ctx, "DELETE FROM hello.outbox_events WHERE aggregate_id = $1", "repo-x")
}

func TestEmitter_NilPool_Panics(t *testing.T) {
	em := &audit.Emitter{Pool: nil, OutboxFQN: "hello.outbox_events"}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on nil Pool")
		}
	}()
	_ = em.Emit(context.Background(), audit.Event{Action: "test", Target: "t"})
}
