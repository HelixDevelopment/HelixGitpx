# Script — 15 Incident response in production

**Track:** Operators · **Length:** 14 min · **Goal:** viewer walks a realistic incident from page to resolution.

## Cold open
Pager buzzes at 03:12 UTC. "HelixGitpxServiceDown: webhook-gateway."

## Body

1. **Acknowledge + assess** — 0:30 – 2:00.
   Open status page, check scope, declare incident.
2. **Establish a room** — 2:00 – 3:00.
   Slack channel, incident commander, comms lead.
3. **Follow the runbook** — 3:00 – 6:00.
   Jump to the runbook link in the alert; execute step-by-step.
4. **Identify root cause** — 6:00 – 9:00.
   Logs in Loki, traces in Tempo, profile in Pyroscope.
5. **Mitigate** — 9:00 – 11:00.
   Feature flag; partial rollback via Argo CD.
6. **Comms + status page** — 11:00 – 12:30.
   Update every 30 min; customer-facing summary.
7. **Post-incident review** — 12:30 – 13:45.
   Blameless review template; follow-up issues filed.

## Wrap-up (13:45 – 14:00)
"Every runbook starts as a post-mortem. Write good ones."

## Companion doc
`docs/operations/runbooks/`
