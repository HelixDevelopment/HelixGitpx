# RB-130 — Adapter Provider Rate-Limited

> **Severity default**: P2 (P1 if affects > 20 % of customers)
> **Owner**: Integrations team
> **Last tested**: 2026-04-12

## 1. Detection

Alert: `UpstreamRateLimitLow` — `helixgitpx_adapter_rate_limit_remaining{provider=*}` < 100 for > 5 min.
Related: `AdapterCircuitOpen` — adapter circuit breaker tripped for provider.

Supporting signals:
- Customers report "push completed locally but didn't reach GitHub".
- `helixgitpx_adapter_requests_total{status="rate_limited"}` spike.
- `helixgitpx_conflicts_detected_total{kind="ref_divergence"}` rising (drift because we can't push).

---

## 2. First 2 Minutes

1. Ack alert.
2. Identify **which provider**: look at alert labels (`provider=github`).
3. Determine **scope**: one org's token exhausted vs. global provider-wide throttle.
4. Check provider status page (some 429s are rate limits *they* have changed).

---

## 3. Diagnose

### 3.1 One org is over budget

Typical pattern: a specific PAT or App installation hitting its hourly ceiling.

```bash
# Show rate limit per upstream id
promql "helixgitpx_adapter_rate_limit_remaining{provider=\"github\"}"

# Top consumers
helixctl logs adapter-pool --since=30m \
  --filter='provider=github AND status=rate_limited' \
  | jq -r '.upstream_id' | sort | uniq -c | sort -rn
```

### 3.2 Global provider-wide

- Provider has tightened limits.
- Global GitHub Actions rate abuse from unrelated traffic (rare on our egress IP).
- Our fan-out concurrency too high.

### 3.3 Credential problem

Sometimes 429 masks a 401 (old tokens returning 429 rather than 401). Verify with:

```bash
helixctl adapter test-credentials --upstream <id>
```

---

## 4. Mitigations

### 4.1 Throttle adapter concurrency for provider

```bash
# Drop global concurrency for this provider
helixctl adapter set-concurrency --provider=github --value=20

# Or pause fan-out temporarily
helixctl flag set sync.fanout-paused --on --provider=github --reason="rate limits"
```

Reads + local operations unaffected; only outbound writes pause.

### 4.2 Offer customer alternatives

For the affected org:

- Suggest rotating to a **GitHub App** (higher rate limit than PAT).
- If Enterprise, enable **dedicated token pool**.
- Reduce upstream concurrency in their org settings.

### 4.3 Switch to shadow mode (for new customer in import)

If the issue is a brand-new import overwhelming their quota:

```bash
helixctl upstream set-shadow --upstream <id> --on
```

Their other upstreams unaffected; new upstream will catch up once throttle clears.

### 4.4 Open circuit manually

If many 429s are wasted:

```bash
helixctl adapter circuit open --provider=github --for=15m --reason="rate-limited"
```

---

## 5. Rate Limit Budget Policy

- Adapter pool maintains a **token bucket per (provider, upstream_id)** mirroring upstream's rate limit header.
- When remaining < 10 % → dec concurrency.
- When remaining < 3 % → pause non-critical operations (read-only mirrors, metadata refresh).
- When 0 → pause all, schedule resume at `reset_at`.

This should happen automatically; if it isn't, investigate:

```bash
helixctl adapter metrics --provider=github
```

---

## 6. Verify Recovery

- `helixgitpx_adapter_rate_limit_remaining` trending up.
- `helixgitpx_adapter_circuit_breaker_state == 0` (closed).
- Backlog of queued fan-out jobs draining (Kafka lag).
- No new customer complaints in 30 min.

---

## 7. Post-Incident

- If systemic throttling: review our default concurrency defaults with the provider team.
- If abuse pattern detected: flag the customer for a sales/support conversation.
- If provider silently tightened: update our runbook + docs; notify customers if broadly impacting.

Communication templates (for customer-visible slow fan-out):

> "Sync to [Provider X] is currently delayed due to upstream rate limits. Your data is safe; operations are queued and will replay automatically. Updates as they develop."

---

## 8. Related

- RB-131 (Adapter auth failure storm)
- RB-110 (Sync orchestrator stuck workflows)
- RB-301 (Accidental auto-apply) — rate limits can contribute to divergence
- [10-git-provider-integrations.md §Rate Limits]
