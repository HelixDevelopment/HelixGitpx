// Package repo (event_outbox) writes hello.said events to hello.outbox_events
// inside the same pgx transaction as the counter UPSERT. Debezium's PostgreSQL
// connector streams the outbox table to Kafka via the EventRouter SMT.
package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EventOutbox implements domain.Emitter by inserting into hello.outbox_events.
// The actual Kafka emission is performed out-of-process by Debezium.
type EventOutbox struct {
	Pool  *pgxpool.Pool
	Topic string
}

type helloSaidPayload struct {
	Name     string `json:"name"`
	Greeting string `json:"greeting"`
	Count    int64  `json:"count"`
	At       string `json:"at"`
}

// Emit inserts a single outbox row in its own transaction. Use EmitInTx
// for callers that want to share an existing transaction with the counter UPSERT.
func (e *EventOutbox) Emit(ctx context.Context, name, greeting string, count int64) error {
	tx, err := e.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := e.EmitInTx(ctx, tx, name, greeting, count); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// EmitInTx is the transactional entrypoint for callers that own the tx.
func (e *EventOutbox) EmitInTx(ctx context.Context, tx pgx.Tx, name, greeting string, count int64) error {
	payload, err := json.Marshal(helloSaidPayload{
		Name:     name,
		Greeting: greeting,
		Count:    count,
		At:       time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}
	topic := e.Topic
	if topic == "" {
		topic = "hello.said"
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO hello.outbox_events(aggregate_id, topic, payload)
		VALUES ($1, $2, $3)`,
		name, topic, payload,
	)
	return err
}
