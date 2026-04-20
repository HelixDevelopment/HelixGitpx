# RB-011 — Postgres Replication Lag High

> **Severity default**: P1 (data freshness at risk; synchronous writes degraded)
> **Owner**: Database team
> **Last tested**: 2026-04-02

## 1. Detection

Alert: `PostgresReplicationLagHigh` — `pg_replication_lag_seconds > 30` for 5 min.

Supporting signals:
- Read replicas showing stale data in application logs (`db.replication_lag` errors).
- Synchronous replica not ack'ing: primary commit latency rising.
- `CNPG` status shows "Syncing" replica not progressing.

---

## 2. Context

HelixGitpx Postgres: 1 primary + 2 synchronous replicas + N async replicas.
- Sync replicas contribute to durability (quorum).
- Async replicas serve read-only queries (e.g. projections).

High sync-replica lag means writes stall; high async-replica lag means stale reads.

---

## 3. First 2 Minutes

1. Ack; identify **which** replica is lagging: `kubectl cnpg status helixgitpx`.
2. Check dashboard "Postgres — Replication" for per-replica lag history.
3. Is the primary healthy (CPU / IO)? Lag is often a symptom of primary-side pressure.

---

## 4. Diagnose

| Pattern | Likely cause | Check |
|---|---|---|
| All replicas lagging equally | Primary WAL generation exceeds network | `pg_stat_wal`, network saturation |
| One replica lagging | Replica-specific issue (CPU, disk, apply) | `kubectl top pod` + `pg_stat_replication` |
| Sudden spike after a deploy | Long transaction / backfill | `pg_stat_activity` long queries |
| Gradual creep over hours | Index rebuild / vacuum on replica | `pg_stat_progress_*` |
| Lag during off-peak | Backup / snapshot job taking IO | Backup schedule |

### Commands

```bash
# Top-level status
kubectl cnpg status helixgitpx

# Detailed replication state
kubectl -n data exec helixgitpx-1 -c postgres -- \
  psql -c 'SELECT * FROM pg_stat_replication;'

# WAL generation rate
kubectl -n data exec helixgitpx-1 -c postgres -- \
  psql -c "SELECT pg_wal_lsn_diff(pg_current_wal_lsn(), '0/0') AS wal_bytes;"

# Long-running transactions (can hold replica progress)
kubectl -n data exec helixgitpx-1 -c postgres -- \
  psql -c "SELECT pid, state, now()-xact_start AS xact_age, query FROM pg_stat_activity WHERE xact_start IS NOT NULL ORDER BY xact_start;"
```

---

## 5. Mitigations

### 5.1 Terminate blocking transactions

```bash
# Dry-run first
kubectl -n data exec helixgitpx-1 -c postgres -- \
  psql -c "SELECT pid, query FROM pg_stat_activity WHERE now() - xact_start > interval '15 minutes';"

# Terminate (coordinate with owner!)
kubectl -n data exec helixgitpx-1 -c postgres -- \
  psql -c "SELECT pg_terminate_backend(<pid>);"
```

### 5.2 Pause non-critical workloads

```bash
# Pause heavy projectors temporarily
helixctl kafka group pause --group search-projector --for=15m

# Pause AI indexing
helixctl flag set ai.indexer-enabled --off --for=15m
```

### 5.3 Accelerate replica apply

- If replica is bottlenecked on apply (single WAL apply process), consider enabling `synchronous_commit=remote_apply` adjustments during catch-up.
- Increase replica-side `max_standby_streaming_delay` if long queries on replica cause conflicts.

### 5.4 Rebuild the lagging replica

If a replica is irrecoverable:

```bash
kubectl cnpg reload helixgitpx
# or the drastic option, reclone the replica:
kubectl -n data delete pod helixgitpx-3  # CNPG will re-init
```

### 5.5 Throttle incoming writes (last resort)

```bash
helixctl flag set writes.read-only --on --duration=10m
```

Announces a maintenance window; writes reject with `common.unavailable`.

---

## 6. Verify Recovery

- Lag returns to < 1 s on all replicas.
- No `db.replication_lag` errors in service logs.
- `pg_stat_replication.write_lag / flush_lag / replay_lag` all < 1 s.
- CNPG status: all replicas "Ready".

---

## 7. Post-Incident

- Root-cause WAL generation spike vs. replica saturation.
- If long transaction: add a timeout in the originating service.
- If insufficient replica capacity: rightsize.
- If backup correlation: reschedule or run from a read replica instead.

---

## 8. Related

- RB-010 (Primary down)
- RB-012 (Disk full)
- RB-013 (Slow queries)
- [16-infrastructure-scaling.md §9] — RPO/RTO targets
