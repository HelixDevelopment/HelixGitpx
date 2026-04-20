# DR Drill — Region Loss Simulation

**Status:** Required quarterly. Most recent: <update-after-each-drill>.

## Pre-drill checklist

- [ ] All personnel in #helixgitpx-dr-drill Slack channel.
- [ ] Customer comms template staged (no outbound notifications for drill).
- [ ] Region-b replication lag verified < 30s via `SELECT * FROM pg_stat_replication;`.
- [ ] MirrorMaker 2 lag metrics green.
- [ ] GeoDNS failover policy snapshot saved to `/tmp/geodns-pre-drill.json`.

## Drill steps

1. **T+0** — Isolate region-a. Apply network policy blocking egress from region-a
   to region-b: `kubectl --context region-a apply -f tools/dr/isolate-region-a.yaml`.
2. **T+30s** — Confirm alert `HelixGitpxRegionUnreachable` fires in Alertmanager.
3. **T+60s** — Promote region-b CNPG cluster to primary:
   `kubectl --context region-b patch cluster helixgitpx-pg-b -p '{"spec":{"replica":{"enabled":false}}}' --type merge`.
4. **T+2m** — Flip GeoDNS policy to route 100% to eu-west-2:
   `kubectl --context region-b apply -f tools/dr/geodns-failover-to-b.yaml`.
5. **T+3m** — Stop MirrorMaker 2 (prevent split-brain writes):
   `kubectl --context region-b scale deploy/mirrormaker2 --replicas=0`.
6. **T+5m** — Smoke-test API from external location; validate:
   - `/healthz` returns 200.
   - Login via Keycloak succeeds.
   - `/api/v1/orgs` returns expected orgs.
   - A new PR can be created end-to-end.
7. **T+15m** — Data integrity check:
   `psql -c "SELECT count(*) FROM org.organizations" against both regions.`
8. **T+30m** — Begin recovery. Re-enable region-a, re-establish replication
   *reversed* (a is now follower). Document any manual resync.
9. **T+60m** — Post-drill review. Capture failed steps, RTO (target: 5m),
   RPO (target: 30s) into `docs/operations/dr-drill-<date>.md`.

## Success criteria

- Recovery Time Objective (RTO) ≤ 5 minutes.
- Recovery Point Objective (RPO) ≤ 30 seconds.
- No data loss for completed transactions.
- No more than 10 failed user requests during failover window.
- Runbook clarity: operator could execute without any SME present.

## Rollback

If drill goes wrong and causes actual customer impact:
1. Re-activate region-a immediately (remove network isolation).
2. Flip GeoDNS back to region-a preferred.
3. Disable the region-b primary promotion if not already propagated.
4. Escalate to #helixgitpx-incident.
