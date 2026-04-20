# 22 — Billing & Metering

> **Document purpose**: Define how HelixGitpx **meters usage, enforces plan limits, calculates invoices, and handles dunning** — end to end, from event → meter → invoice → payment → accounting.

---

## 1. Principles

1. **Events are the source of truth.** Every billable action produces a Kafka event. Meters aggregate events. Invoices are projections of meters. If accounting ever disagrees, we can replay from events.
2. **Deterministic.** Given the same event stream, we always compute the same invoice. Idempotent consumers; ordered partition keys.
3. **Tenant-isolated.** `org_id` on every event; meters keyed by `(org_id, meter, period)`.
4. **No surprise charges.** Soft limit → notice. Hard limit → throttle. Overages require an opt-in upgrade or explicit on-demand consent.
5. **Transparent.** Customers can see live usage, projected spend, and historical invoices via API and UI.

---

## 2. Meters

| Meter | Unit | Trigger | Notes |
|---|---|---|---|
| `repos.active`                | repos | monthly snapshot | repos with any push/PR/issue in period |
| `storage.git.bytes`           | byte-months | nightly snapshot | Git objects + packs |
| `storage.lfs.bytes`           | byte-months | nightly | LFS objects |
| `storage.assets.bytes`        | byte-months | nightly | Release assets, attachments |
| `egress.git.bytes`            | bytes | per clone/fetch | counted at git-ingress |
| `egress.api.bytes`            | bytes | per API response | api-gateway metric |
| `webhook.deliveries`          | count | per delivery | outbound webhooks |
| `upstream.connections`        | count | snapshot | enabled upstreams |
| `sync.jobs`                   | count | per job completion | successful + partial |
| `ai.tokens.input`             | tokens | per prompt run | see [07-ai/10-llm-self-learning.md] |
| `ai.tokens.output`            | tokens | per prompt run | |
| `ai.inference.seconds.gpu`    | seconds | per GPU second | for "AI Pro" tier |
| `search.queries`              | count | per query | user-initiated only |
| `events.delivered.live`       | count | per delivery | includes WebSocket/gRPC stream |
| `cicd.minutes`                | minutes | per CI job | if CI runners hosted |
| `users.seats`                 | users | monthly snapshot | seat-based plans |

Each meter emits to `helixgitpx.billing.usage` with `{org_id, meter, units, unit_type, period_month, recorded_at}`.

---

## 3. Plans

| Plan | Target | Pricing | Notes |
|---|---|---|---|
| **Free** | Hobby, OSS | $0 | 25 repos, 2 upstreams, 100 k AI tokens/month |
| **Team** | Small teams | $15/user/mo | 500 repos, 10 upstreams, 10 M AI tokens, basic SLO |
| **Business** | Growing SaaS | $49/user/mo | unlimited repos, unlimited upstreams, 100 M AI tokens, advanced SLO, SSO, audit exports |
| **Enterprise** | Regulated / large | Custom | on-prem / dedicated / residency pinning, 24/7 SLA, legal addenda |
| **OSS Program** | Public OSS projects | $0 Business-equivalent | application-based |
| **Education / Non-Profit** | Schools, charities | 50 % discount | application-based |

Plan limits are JSON on `billing.plans.limits`:

```json
{
  "repos": 500,
  "upstreams_per_org": 10,
  "ai_tokens_monthly": 10000000,
  "storage_bytes": 107374182400,
  "egress_bytes_monthly": 536870912000,
  "users_seats": -1,
  "sla_tier": "business",
  "audit_retention_days": 365,
  "residency_pinning": false
}
```

`-1` means unlimited.

---

## 4. Enforcement

### 4.1 Soft & Hard Limits

- **Soft** (80 %) — email + in-app notice + dashboard badge.
- **Hard** (100 %) — action blocked with `common.quota_exceeded`. Read-only access preserved.

Hard-limit predicates are evaluated in the **gatekeeper middleware** at the API gateway + per-service hot paths (e.g. push accept hook).

Enforcement is **eventually consistent** — we accept up to a small over-run before block engages, but rate-limit recovery is instant.

### 4.2 Overage Options

- Team: no overages; users must upgrade.
- Business: opt-in "Flex" mode with on-demand token / byte pricing disclosed up-front.
- Enterprise: contractual overage rates.

### 4.3 Grace Periods

- Payment past due: 14-day grace. Day 7: email. Day 14: downgrade to Free (read-only on overages).
- Trial-to-paid conversion: 30-day grace.

---

## 5. Metering Pipeline

```
producers ──► Kafka topic  helixgitpx.billing.usage
                     │
                     ▼
           billing-meter consumer (stateful, sharded by org_id)
                     │
                     ├─► Postgres: billing.usage_records (append-only)
                     ├─► Aggregated hourly roll-ups → billing.usage_hourly
                     └─► Invoice generator (monthly, Temporal workflow)
                                 │
                                 ▼
                     Postgres: billing.invoices
                                 │
                                 ▼
                     External billing provider (Stripe / Paddle)
```

### 5.1 Idempotency

- Every usage event carries a `usage_event_id` (UUIDv7).
- `billing.usage_records.id` is that same UUIDv7; DB upsert on conflict → no-op.
- Meter workers can safely replay.

### 5.2 Late-Arriving Events

- Events within 48 h of close of billing period included in current invoice.
- Older → reversal credit on next period's invoice.

### 5.3 Shifting Time Zones

- Billing period boundaries are **UTC month-end** regardless of customer's local time.
- All stored timestamps are UTC; display-time conversion only.

---

## 6. Invoice Generation (Temporal Workflow)

Workflow `GenerateMonthlyInvoice(org_id, period_month)`:

1. Snapshot meter totals from `billing.usage_hourly` where `period = period_month`.
2. Compute **plan-inclusive** usage (subtract plan quotas).
3. Compute **overage** tiered pricing.
4. Apply discounts (OSS, education, promo codes).
5. Compute taxes (via tax-provider integration — Avalara / Stripe Tax).
6. Compute totals; store in `billing.invoices`.
7. Publish to `helixgitpx.billing.invoice.generated`.
8. Submit to external billing provider (creates charge).
9. On success → `paid_at`; on failure → dunning workflow.
10. Send email invoice + webhook notification.

---

## 7. Payments

### 7.1 Processor

- **Stripe** primary; **Paddle** alternative for EU/VAT-heavy customers.
- PCI handled by processor; HelixGitpx never touches card data.
- Webhooks from processor feed `billing.payments` table.

### 7.2 Dunning

- Day 0: charge attempt.
- Day 1, 3, 5, 7: retry (auto).
- Day 7: email to admin.
- Day 14: downgrade + persistent UI banner.
- Day 30: suspension (read-only + export).
- Day 60: deletion notice (30-day warning).
- Day 90: data deletion (irreversible) unless legal hold.

### 7.3 Refunds

- Processed manually by support for abuse / billing errors.
- Audit-logged.

---

## 8. Cost Attribution (Internal)

Beyond customer billing, HelixGitpx attributes **internal** cost to each org for observability:

- **Compute** (Kubecost per namespace × per-tenant allocation).
- **Storage** (PVC / S3 tag).
- **LLM inference** (per-prompt token usage → GPU-hour allocation).
- **Egress**.

Result: per-org margin dashboard for the business team.

---

## 9. API

- `GET /api/v1/orgs/{id}/billing/plan` — current plan + limits.
- `POST /api/v1/orgs/{id}/billing/upgrade` — change plan.
- `GET /api/v1/orgs/{id}/billing/usage?period=YYYY-MM` — usage rollup.
- `GET /api/v1/orgs/{id}/billing/invoices` — list.
- `GET /api/v1/orgs/{id}/billing/invoices/{id}` — detail + PDF.
- `POST /api/v1/orgs/{id}/billing/payment-methods` — add card (redirect to processor).
- `DELETE /api/v1/orgs/{id}/billing/subscription` — cancel (downgrade on next cycle).

Enterprise-only:
- `POST /api/v1/admin/billing/credits` — apply credit.
- `POST /api/v1/admin/billing/adjustments` — adjustment.

---

## 10. Observability

Metrics (see [15-reference/metrics-catalog.md]):
- `helixgitpx_billing_usage_units_total`
- `helixgitpx_billing_quota_exceeded_total`
- `helixgitpx_billing_spend_usd`

Dashboards:
- Usage vs. plan quota per org.
- Projected monthly spend per org.
- Invoice pipeline health (generation latency, failures).
- Dunning funnel.

Alerts:
- Billing meter lag (`RB-160`).
- Invoice generation failed.
- Payment processor webhook backlog.

---

## 11. Audit

Every billing action (plan change, limit override, credit grant, subscription cancel) writes an `audit.events` entry. Retained per customer retention policy (min 7 y for paid accounts).

---

## 12. Compliance

- **SOX**: segregation of duties — engineers cannot issue credits; only billing-admin role.
- **GDPR**: invoices are legal records; retention mandated independent of erasure requests.
- **Tax**: EU VAT + US state sales tax computed via provider; nexus tracked.

---

## 13. Testing

- Deterministic replay test: golden dataset of usage events → expected invoice bytes.
- Fuzz: malformed usage payloads.
- Chaos: processor outage (simulate Stripe 500) — retries + alert.
- Contract tests against processor webhooks.

---

## 14. Customer-Facing Transparency

- Live usage page with drill-down.
- Projected end-of-month invoice estimate.
- Explanation of every line item (link to meter docs).
- Downloadable CSV of raw metered events.

---

*— End of Billing & Metering —*
