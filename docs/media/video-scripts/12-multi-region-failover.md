# Script — 12 Multi-region failover drill

**Track:** Operators · **Length:** 18 min · **Goal:** viewer can execute a controlled failover end-to-end.

## Body

1. **Architecture recap** — 0:30 – 2:00.
   Active region-a, standby region-b, MirrorMaker 2 + CNPG replica + GeoDNS.
2. **Pre-drill checklist** — 2:00 – 4:00.
   Comms staged, replication lag check, GeoDNS snapshot.
3. **Isolate region-a** — 4:00 – 6:00.
   NetworkPolicy blocks egress; alert fires.
4. **Promote region-b** — 6:00 – 9:00.
   CNPG patch, MirrorMaker stop, GeoDNS flip.
5. **Smoke tests** — 9:00 – 12:00.
   Login, org list, create PR, push.
6. **Data integrity checks** — 12:00 – 14:00.
   Row counts between regions; audit-log merkle root.
7. **Recovery** — 14:00 – 16:30.
   Re-enable region-a, re-seed replication reversed.
8. **Post-drill review** — 16:30 – 17:30.
   Timing log; what slowed the runbook; fix PRs.

## Wrap-up (17:30 – 18:00)
"RTO 5 min, RPO 30 s. Measured every quarter."

## Companion doc
`tools/dr/dr-drill-runbook.md`
