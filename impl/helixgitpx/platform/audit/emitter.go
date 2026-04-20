package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Event is the canonical audit payload. Published to topic audit.events.
type Event struct {
	At          time.Time      `json:"at"`
	ActorUserID string         `json:"actor_user_id"`
	ActorIP     string         `json:"actor_ip,omitempty"`
	Action      string         `json:"action"`
	Target      string         `json:"target"`
	Details     map[string]any `json:"details,omitempty"`
}

// MarshalJSON defaults At to now() when unset.
func (e Event) MarshalJSON() ([]byte, error) {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	type alias Event
	return json.Marshal(alias(e))
}

// Emitter writes audit events to the emitting service's local outbox table.
// Debezium captures the outbox and routes to topic audit.events via the
// EventRouter SMT.
type Emitter struct {
	Pool      *pgxpool.Pool
	OutboxFQN string // e.g. "auth.outbox_events" or "org.outbox_events"
}

// Emit inserts one audit event into the outbox in its own tx.
func (e *Emitter) Emit(ctx context.Context, ev Event) error {
	tx, err := e.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := e.EmitInTx(ctx, tx, ev); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// EmitInTx is the transactional entrypoint for callers that own the tx.
func (e *Emitter) EmitInTx(ctx context.Context, tx pgx.Tx, ev Event) error {
	payload, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx,
		fmt.Sprintf(`INSERT INTO %s(aggregate_id, topic, payload) VALUES ($1, $2, $3)`, e.OutboxFQN),
		ev.Target, "audit.events", payload)
	return err
}
