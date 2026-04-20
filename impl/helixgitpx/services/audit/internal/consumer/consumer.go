// Package consumer reads audit.events from Kafka and appends to audit.events
// via the SECURITY DEFINER function audit.append_event (which bypasses the
// no-update/no-delete rules on the append-only table).
package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Consumer ties a kgo Client + pg pool together.
type Consumer struct {
	Client *kgo.Client
	Pool   *pgxpool.Pool
}

type rawEvent struct {
	At          time.Time      `json:"at"`
	ActorUserID string         `json:"actor_user_id"`
	ActorIP     string         `json:"actor_ip"`
	Action      string         `json:"action"`
	Target      string         `json:"target"`
	Details     map[string]any `json:"details"`
}

// Run loops fetching records and inserting via audit.append_event. Exits on ctx.Done().
func (c *Consumer) Run(ctx context.Context) error {
	for {
		fetches := c.Client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return nil
		}
		if errs := fetches.Errors(); len(errs) > 0 {
			return fmt.Errorf("consumer: fetch errors: %v", errs)
		}
		fetches.EachRecord(func(r *kgo.Record) {
			_ = c.handle(ctx, r)
		})
		if err := c.Client.CommitUncommittedOffsets(ctx); err != nil {
			return err
		}
	}
}

func (c *Consumer) handle(ctx context.Context, r *kgo.Record) error {
	var ev rawEvent
	if err := json.Unmarshal(r.Value, &ev); err != nil {
		return err
	}
	details, _ := json.Marshal(ev.Details)
	_, err := c.Pool.Exec(ctx,
		`SELECT audit.append_event($1, $2, $3, $4, $5, $6::jsonb)`,
		ev.At, ev.ActorUserID, ev.ActorIP, ev.Action, ev.Target, string(details))
	return err
}
