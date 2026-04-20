# RB-010 — Postgres Primary Down

> **Severity default**: P1
> **Owner**: Platform SRE
> **Last tested**: 2026-04-15 (monthly DR drill)

## 1. Detection

Alert: `PostgresPrimaryDown` — `up{job="postgres-primary"} == 0` for 1 min.

Other signals:
- Application errors with code `common.unavailable` or `db.deadlock` spiking.
- `pg_up{role="primary"}` == 0 in Grafana.
- CNPG controller warnings in `kubectl -n cnpg-system logs`.

---

## 2. Immediate Response (first 5 minutes)

1. **Acknowledge** the PagerDuty alert.
2. **Open incident channel** (`helixctl incident open --severity=P1`).
3. **Confirm** the alert is real — check synthetic probes (`/healthz` returns 503), `kubectl -n data get pg-cluster helixgitpx -o yaml`, and `kubectl cnpg status helixgitpx`.
4. **Announce** on status page: "Investigating connectivity issues."
5. **Silence** downstream symptom alerts for 15 min to reduce noise.

---

## 3. Diagnosis Tree

```
pg_up == 0 ?
├── Node NotReady?
│     → RB-001 (Node NotReady) — CNPG may reschedule automatically
├── PVC mount failure?
│     → Check storage class health; contact storage ops
├── Crash loop (OOM)?
│     → Check pod events; consider temporary memory bump
├── PG process stuck?
│     → SIGTERM stuck pod; let CNPG replace
└── Controller confused?
      → Restart cnpg-controller-manager deployment
```

Typical cause: node with the primary pod became unavailable.

---

## 4. Failover

### Automatic (default)

CNPG performs failover when primary unreachable > 60 s. Watch:

```bash
kubectl cnpg status helixgitpx
# Look for: "Current Primary" switch to a replica
```

### Manual (if stuck)

```bash
# Promote specific replica
kubectl cnpg promote helixgitpx helixgitpx-2

# Or trigger fresh elections
kubectl -n data delete pod helixgitpx-1   # the stuck primary pod
```

---

## 5. Verify Recovery

- `pg_up{role="primary"} == 1` for ≥ 2 min.
- Replica count ≥ 2: `kubectl cnpg status helixgitpx`.
- Writes succeed: `helixctl smoke db` (runs INSERT + SELECT in auth schema).
- Error-budget burn halts.
- `pg_replication_lag_seconds` returns to < 1 s.

Let auto-healing close any dependent symptom alerts; clear manual silences.

---

## 6. Post-Recovery

1. Update status page: resolved.
2. Run `pg_verify_checksums` (read-only) on the new primary.
3. Check for any writes that may have been accepted but not yet replicated before failover — compare Postgres LSN to Kafka outbox lag.
4. If any data loss suspected → escalate to DB team; consider PITR to just-before-failure.
5. Kick off RCA workflow (postmortem).

---

## 7. Communications Templates

**Public status page (during)**:
> We're investigating degraded database connectivity affecting writes. Reads remain available. Updates every 15 minutes.

**Public status page (resolved)**:
> Database connectivity has been restored. Writes are processing normally. A full incident report will be published within 5 business days.

---

## 8. Drill Instructions

Replay this scenario in staging:

```bash
# Confirm staging cluster
kubectl config current-context   # must contain "staging"

# Kill the primary
PRIMARY=$(kubectl -n data get pg-cluster helixgitpx \
  -o jsonpath='{.status.currentPrimary}')
kubectl -n data delete pod "${PRIMARY}"

# Start stopwatch and observe
watch -n 2 'kubectl cnpg status helixgitpx'
```

Expected: new primary promoted within 60 s.

---

## 9. Related

- ADR-0014 (Argo CD GitOps) — PR rollback path if the cause is config drift
- RB-011 (Postgres Replication Lag High)
- RB-012 (Postgres Disk Full)
- [16-infrastructure-scaling.md §9] — RPO/RTO contracts
