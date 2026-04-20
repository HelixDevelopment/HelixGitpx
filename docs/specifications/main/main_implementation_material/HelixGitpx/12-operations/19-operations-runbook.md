# 19 — Operations Runbook

> **Document purpose**: A working **on-call handbook**. Every incident class has a runbook. Every runbook has: detection signal, severity, first response, diagnosis tree, mitigations, escalation, and postmortem template anchor.

---

## 1. On-Call Essentials

- **Primary / secondary** rotation, weekly shifts.
- Handover ritual: review open incidents, current deploys, recent changes.
- Pager goes to PagerDuty; auto-creates incident Slack channel with IC, Tech Lead, Communications roles.
- Commands:
  - `helixctl deploy status`
  - `helixctl incident open --severity=P1 --summary="…"`
  - `helixctl break-glass --user=<me> --reason="…" --duration=30m` (time-bound elevated access, recorded in audit).

---

## 2. Severity Definitions

| Sev | Example | Response time |
|---|---|---|
| **P1** | Customer-impacting outage, data loss risk, security breach | 15 min ack, 30 min update, fix ASAP |
| **P2** | Significant degradation; SLO burn fast | 30 min ack, 1 h update |
| **P3** | Non-critical degradation; workaround exists | Next business day |
| **P4** | Noise, technical debt | Sprint |

---

## 3. Runbook Index

### Platform

- `RB-001` — K8s node NotReady
- `RB-002` — Pod CrashLoopBackOff
- `RB-003` — Namespace resource quota exhausted
- `RB-004` — Certificate near expiry
- `RB-005` — Istio ambient data plane unhealthy
- `RB-006` — SPIRE server down

### Database / Messaging

- `RB-010` — Postgres primary down (HA failover)
- `RB-011` — Postgres replication lag high
- `RB-012` — Postgres disk full
- `RB-013` — Slow queries / lock storm
- `RB-020` — Kafka under-replicated partitions
- `RB-021` — Kafka consumer group lag critical
- `RB-022` — Kafka broker disk full
- `RB-023` — Schema registry down
- `RB-030` — Redis/Dragonfly memory exhaustion
- `RB-031` — Redis replication break
- `RB-040` — OpenSearch cluster red
- `RB-041` — Qdrant OOM
- `RB-042` — Meilisearch index corrupt

### Services

- `RB-100` — api-gateway 5xx spike
- `RB-101` — auth-service token minting failure
- `RB-110` — sync-orchestrator stuck workflows
- `RB-111` — webhook-gateway dedup cache down
- `RB-120` — conflict-resolver backlog
- `RB-121` — conflict-resolver AI confidence regression
- `RB-130` — adapter-pool provider rate-limited
- `RB-131` — adapter-pool auth failure storm
- `RB-140` — live-events-service stream pile-up
- `RB-150` — ai-service inference timeout
- `RB-151` — LLM model pulled from disk
- `RB-160` — billing meter lag

### Security

- `RB-200` — Suspected credential compromise
- `RB-201` — Secret leaked in a public commit (on our side)
- `RB-202` — Malicious plugin behaviour
- `RB-203` — DDoS at edge
- `RB-204` — Supply-chain alert (unsigned image detected)

### Data / Customer

- `RB-300` — Data integrity anomaly (event store vs. projection mismatch)
- `RB-301` — Accidental auto-apply (undo window)
- `RB-302` — GDPR erasure request processing

---

## 4. Sample Runbooks

### RB-010 — Postgres primary down

**Detect**: `PostgresPrimaryDown` alert; `up{job="postgres-primary"} == 0`.

**Severity**: P1.

**First actions**:

1. Acknowledge alert.
2. Confirm: `helixctl db status` and `kubectl -n data get pg-cluster -o yaml`.
3. If CNPG cluster, verify CNPG controller is alive: `kubectl -n cnpg-system logs deploy/cnpg-controller-manager`.

**Failover**:

- CNPG performs automatic failover when the primary is unreachable > 60 s.
- If stuck: `kubectl cnpg promote <cluster> <replica>`.

**Verify**:

- Check replication status: `kubectl cnpg status <cluster>`.
- Run `helixctl smoke db` (INSERT + SELECT).
- Watch app-side error budget burn.

**Comms**: update status page if customer-visible degradation observed.

**Post**: root-cause, PITR integrity check, follow-up actions.

---

### RB-020 — Kafka under-replicated partitions

**Detect**: `KafkaUnderReplicatedPartitions` alert.

**Severity**: P1 if > 1 partition for > 5 min.

**Check**:

```
kubectl -n kafka exec kafka-0 -- bin/kafka-topics.sh --bootstrap-server localhost:9092 --describe --under-replicated-partitions
```

**Common causes**:

- Broker disk full (→ RB-022).
- Broker pod crashed (→ check `kubectl logs`).
- Network partition.

**Mitigation**:

- Rebalance by adding capacity or moving partitions (`kafka-reassign-partitions.sh`).
- If broker disk: free logs via `kafka-log-dirs.sh` + expand PVC.

**Verify**: under-replicated count returns to 0; consumer groups reading normally.

---

### RB-100 — API gateway 5xx spike

**Detect**: `APIGateway5xxHigh` — `sum(rate(http_requests_total{job="api-gateway",code=~"5.."}[5m])) / sum(rate(http_requests_total{job="api-gateway"}[5m])) > 0.02`.

**Severity**: P1.

**Diagnose**:

1. Check which backend RPCs are failing: drill-down dashboard by `grpc_server_handled_total{code!="OK"}`.
2. Look at Tempo for a sample 5xx trace — find the breaking dependency.
3. Check recent deploys: `helixctl deploy recent`.

**Mitigation**:

- If recent deploy: `helixctl rollback <service>`.
- If dependency down: engage that service's runbook.
- If rate-limit misconfigured: adjust via Envoy config flag.

---

### RB-130 — Adapter provider rate-limited

**Detect**: `UpstreamRateLimitLow` or `AdapterCircuitOpen`.

**Severity**: P2.

**Actions**:

1. Identify provider from alert labels.
2. Check `helixgitpx_adapter_rate_limit_remaining{provider=…}`.
3. Confirm credentials are healthy (rotation may not have applied).
4. Reduce adapter concurrency for that provider.
5. Switch to degraded mode (read-only mirror, pause non-essential fan-out pushes).
6. Communicate to affected customers if needed.

---

### RB-200 — Suspected credential compromise

**Severity**: P1 (security).

1. Force-revoke active sessions for suspected account.
2. Invalidate PATs and refresh tokens.
3. Rotate affected upstream credentials in Vault.
4. Run `helixctl audit trace --user=<id> --last=7d` and preserve evidence.
5. Engage security engineer on call.
6. Notify legal if applicable; GDPR 72 h clock if EU user data involved.

---

### RB-301 — Accidental auto-apply

**Detect**: customer complaint or `ConflictUndoRequest` metric spike.

**Severity**: P2 (can escalate).

1. Identify case: `helixctl conflict show <case_id>`.
2. If within 5-min window: trigger `helixctl conflict undo --case=<id>`.
3. If outside: open a reverse-merge PR on all affected upstreams.
4. Notify affected repo maintainers.
5. Post-incident: evaluate whether policy should be tightened.

---

## 5. Common Commands

```bash
# Health
helixctl platform status
helixctl service ls
helixctl deploy recent --service=<name> -n 10

# Logs / traces
helixctl logs <service> --since=15m --trace-id=<id>
helixctl trace <id>

# Data
helixctl db primary
helixctl db lag
helixctl kafka topics
helixctl kafka lag --group=<group>

# Operations
helixctl scale <service> --replicas=20
helixctl pod restart <service>/<pod>
helixctl rollback <service>
helixctl flag set <name> --off
helixctl conflict replay --case=<id>
helixctl dlq list --topic=<t>
helixctl dlq replay --topic=<t> --from=<offset> --to=<offset>
```

All commands require MFA; sensitive ones need a second operator's co-sign (via Slack slash-command).

---

## 6. Break-Glass Access

- Normal access is least-privilege.
- Emergency elevated access via `helixctl break-glass` — time-bound (default 30 m), reason required, automatically audited and announced in `#secops` channel.

---

## 7. Postmortem Template

```markdown
# Incident: <TITLE>
**Severity**: P?
**Date / Time (UTC)**: start → resolved
**Duration**: __ min
**Incident Commander**: @…
**Responders**: @…, @…

## Impact
- <who/what was affected, quantified>

## Timeline
- 11:02 UTC — alert X fired
- 11:04 — …

## Root Cause
<5 whys — ultimate cause>

## Triggers
<what made the risk realise now>

## Mitigation / Resolution
<what we did that worked>

## What went well
- …

## What didn't
- …

## Action items
| # | Item | Owner | Due | Status |

## Lessons Learned
- …
```

Publish to `postmortems/` folder; discussed at weekly review.

---

## 8. Communications

- Status page: auto-create incident on P1 / P2; update every 30 min minimum.
- Customer-facing language templates live in `runbooks/communications/`.
- Transparency over speed for final postmortem — 5-business-day target.

---

## 9. Escalation Tree

1. Primary on-call.
2. Secondary on-call (after 15 min without ack).
3. Service tech lead.
4. Platform lead / CTO.
5. External partners (provider support, cloud support, legal).

---

## 10. Drills

- **Quarterly**: full DR (region loss).
- **Monthly**: random pod / broker kill (chaos-mesh).
- **Weekly**: on-call simulation via tabletop.
- **Per release**: runbook review for touched services.

---

*— End of Operations Runbook —*
