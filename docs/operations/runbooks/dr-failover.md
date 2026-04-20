# Runbook — Production DR failover

**Alert:** `HelixGitpxRegionUnreachable` sustained 3 min, or declared incident.

## Go / No-go

Only trigger if:
- Region-a API unreachable from ≥2 independent monitoring vantage points.
- Region-b replication lag < 5 minutes.
- On-call SRE lead approves.

## Steps

Identical to `tools/dr/dr-drill-runbook.md` — but this is for real. Key differences from drill:

- **Comms:** Post to status page BEFORE step 4 (GeoDNS flip). Use template `comms/dr-failover-initial.md`.
- **Step 5:** Verify smoke tests pass against real customer data sample (5 orgs selected at random).
- **Step 8:** DO NOT re-enable region-a automatically — wait for post-incident review + root cause.

## Customer comms

Use `comms/dr-failover-initial.md` → `comms/dr-failover-restored.md` → post-mortem email.

## Escalation

Anything unexpected during execution → halt at current step, page SRE lead + CTO.
