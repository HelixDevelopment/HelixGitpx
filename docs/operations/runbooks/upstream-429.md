# Runbook — Upstream 429 storm

**Alert:** `UpstreamRateLimitBurst` (adapter-pool 429 rate > 5/s sustained 2 min).

## Immediate checks

- Which provider: `rate(adapter_pool_upstream_requests_total{status="429"}[5m]) by (provider)`.
- Current concurrency per provider: `adapter_pool_inflight{provider="..."}`.
- Remaining quota window (from headers): `adapter_pool_quota_remaining{provider="..."}`.

## Remediation

1. **Single provider** — tighten per-provider token bucket via `sync-orchestrator` ConfigMap; flip `max_concurrent` down.
2. **GitHub secondary rate limit** — back off to 1 RPS for 5 min; the library retries with jitter automatically.
3. **If retrying DLQ spikes** — pause `fanout_push_wf` schedule; drain backlog under reduced concurrency.

## Escalation

Persistent 429s from paid provider → contact provider support with ticket number + traces.
