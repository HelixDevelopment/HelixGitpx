// Command audit-merkle walks audit.events for the prior hour, builds a
// SHA-256 Merkle tree, and writes the root into audit.anchors. Runs as a
// Kubernetes CronJob.
package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/helixgitpx/helixgitpx/services/audit/internal/merkle"
	"github.com/helixgitpx/platform/log"
	"github.com/helixgitpx/platform/pg"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	lg := log.New(log.Options{Level: "info", Service: "audit-merkle"})

	dsn := os.Getenv("AUDIT_POSTGRES_DSN")
	if dsn == "" {
		lg.Error("AUDIT_POSTGRES_DSN required")
		os.Exit(1)
	}
	pool, err := pg.Open(ctx, pg.Options{DSN: dsn})
	if err != nil {
		lg.Error("pg.Open", "err", err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	now := time.Now().UTC().Truncate(time.Hour)
	from := now.Add(-time.Hour)
	to := now

	rows, err := pool.Query(ctx, `
		SELECT id, details FROM audit.events
		 WHERE at >= $1 AND at < $2 ORDER BY at, id`, from, to)
	if err != nil {
		lg.Error("query", "err", err.Error())
		os.Exit(1)
	}
	defer rows.Close()

	var leaves [][]byte
	for rows.Next() {
		var id string
		var det json.RawMessage
		if err := rows.Scan(&id, &det); err != nil {
			lg.Error("scan", "err", err.Error())
			os.Exit(1)
		}
		leaves = append(leaves, append([]byte(id), det...))
	}
	if len(leaves) == 0 {
		lg.Info("no events in window; skipping anchor", "from", from, "to", to)
		return
	}

	root := merkle.Root(leaves)
	if _, err := pool.Exec(ctx,
		`INSERT INTO audit.anchors(period_start, period_end, merkle_root)
		 VALUES ($1, $2, $3)`, from, to, root); err != nil {
		lg.Error("insert anchor", "err", err.Error())
		os.Exit(1)
	}
	lg.Info("anchor written", "from", from, "to", to, "leaves", len(leaves))
}
