# RB-100 — API Gateway 5xx Spike

> **Severity default**: P1
> **Owner**: API platform team
> **Last tested**: 2026-03-28

## 1. Detection

Alert: `APIGateway5xxHigh` — 5xx error rate on `api-gateway` > 2 % for 5 min.

Supporting signals:
- Synthetic probes failing from ≥ 2 regions.
- Customer reports in `#support-escalations`.
- `helixgitpx_http_requests_in_flight` climbing (queuing).
- Grafana "API Gateway — Overview" in the red.

---

## 2. First 2 Minutes

1. Acknowledge alert; **open incident** (`helixctl incident open --severity=P1`).
2. Quick triage: is it all customers (global) or a subset (tenant-scoped)?
   - Dashboard filter by `org_id` top-10.
3. Check the **Deploys** feed: was anything released in the last 2 h?
4. Check **Status Page** internal feed for active incidents in upstream dependencies.

---

## 3. Diagnose

Typical causes, in order of frequency:

| Suspicion | Quick check | Confirm |
|---|---|---|
| Recent deploy regression | `helixctl deploy recent -n 10` | Correlate PR → start of spike |
| Backend service outage | `helixctl service ls --unhealthy` | Drill into failing RPC |
| Database degraded | Postgres dashboard | `pg_replication_lag`, slow queries |
| Cache eviction storm | Redis dashboard | hit-ratio drop |
| Auth provider down | Auth dashboard | `auth.tokens_rejected_total` spike |
| Rate-limit misconfig | Envoy stats | 429s paired with 5xxs |
| Thundering herd on expired cache | Cache miss spike | coordinate with `cache.stampede` alert |
| DDoS / burst | Cloudflare analytics | WAF rule |

### Useful commands

```bash
# Error rate by backend
helixctl logs api-gateway --since=5m --filter='severity=ERROR' --limit=200 \
  | jq -r '.upstream_service' | sort | uniq -c | sort -rn

# Find a sample failing trace (Tempo)
helixctl trace search --service=api-gateway --code=5xx --since=5m --limit=1

# Error breakdown
promql "sum by (upstream) (rate(helixgitpx_http_requests_total{service=\"api-gateway\",code=~\"5..\"}[2m]))"
```

---

## 4. Mitigations

### 4.1 Rollback the last deploy (most common)

```bash
helixctl rollback api-gateway
# or the specific problematic backend service
helixctl rollback repo-service
```

Argo Rollouts undoes the canary. Observe error rate drops within 1-3 min.

### 4.2 Shed load

If the cause is a single expensive endpoint:

```bash
# Temporarily rate-limit the hot path at the gateway
helixctl gateway set-limit --route=/api/v1/repos/:id/tree --rps=100
```

Or flip a feature flag off:

```bash
helixctl flag set ui.conflict-resolver-three-pane --off
```

### 4.3 Scale out (if capacity-bound)

```bash
helixctl scale api-gateway --replicas=60
```

HPA may already be scaling; verify it isn't capped.

### 4.4 Read-only mode (last resort)

If writes are the crisis and reads are fine:

```bash
helixctl flag set writes.read-only --on --reason="INCIDENT-xxx"
```

This blocks writes cluster-wide; reads + exports continue. Announce on status page immediately.

### 4.5 Engage the backend owner

If it's localised to one backend service, page that team directly.

---

## 5. Verify Recovery

- `helixgitpx_http_requests_total{code=~"5.."}` rate drops below 0.5 %.
- Synthetic probes green from all regions.
- `helixgitpx_http_requests_in_flight` returns to baseline.
- Error-budget burn rate back below 1×.

---

## 6. Communication

**During**:
> "We're investigating elevated error rates on the API. No data loss observed. Updates every 15 minutes."

**Resolved**:
> "The elevated error rate has been resolved. Services are operating normally. A detailed incident report will be published within 5 business days."

---

## 7. Post-Incident

- Write the postmortem using the template in [19-operations-runbook.md §7].
- If rollback was the mitigation: root-cause the regression; add a test to prevent recurrence.
- If alert fired too late / too early: tune threshold + evaluation window.
- Review queuing signals; consider earlier HPA scale-up triggers.

---

## 8. Related

- RB-101 (Auth token minting failure)
- RB-011 (Postgres replication lag)
- RB-140 (Live events delivery lag)
- RB-203 (DDoS)
