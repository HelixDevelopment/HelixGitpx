# 18 — Observability

> **Document purpose**: Specify the **metrics, logs, traces, profiling, SLOs, dashboards, and alerts** that let us run HelixGitpx confidently at scale.

---

## 1. The Three Pillars + One

| Pillar | Tool | Storage |
|---|---|---|
| Metrics | Prometheus → **Mimir** (HA / long-term) | Object store |
| Logs | Vector → **Loki** (+ OpenSearch mirror for compliance) | Object store |
| Traces | OTel Collector → **Tempo** | Object store |
| Profiles | Pyroscope/Parca (continuous profiling) | Object store |

Visualised in **Grafana**. Alerts through **Alertmanager** → PagerDuty + Slack + incident channels.

---

## 2. OpenTelemetry Everywhere

- Every service emits OTel traces, metrics, and logs via the **OTel SDK**.
- OTel Collector in-cluster fan-outs to the right backend.
- Resource attributes standardised (service.name, service.version, deployment.environment, k8s.namespace, k8s.pod.name, host.name, telemetry.sdk.*).
- **W3C traceparent** propagated: browsers → gateway → services → DB → upstream adapters.
- Every log line carries `trace_id` / `span_id`.

---

## 3. Metric Conventions

- Prefix: `helixgitpx_`.
- Naming: `<noun>_<unit>` (e.g. `events_delivered_total`, `sync_duration_seconds`).
- Units: base SI units; seconds not ms, bytes not MiB.
- Labels: low-cardinality; never raw user ids or IPs. Tenancy via `org_id` only when business-critical.
- Histogram buckets tuned per signal (latency: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]).

### 3.1 Required Metrics per Service

Every service exposes:

- `helixgitpx_service_up` (gauge, 1).
- `helixgitpx_service_info` (gauge, 1) with labels: version, commit_sha, build_time.
- `helixgitpx_build_info` (static, informational).
- **RED**: `http_requests_total`, `http_request_duration_seconds`, `http_requests_in_flight`.
- **gRPC RED**: `grpc_server_started_total`, `grpc_server_handled_total`, `grpc_server_handling_seconds`, `grpc_server_msg_received_total`, `grpc_server_msg_sent_total`.
- **USE** for nodes: CPU, memory, disk, network.
- **Saturation** signals: queue depths, worker busy ratios.

### 3.2 Business Metrics (examples)

- `helixgitpx_repos_created_total{org_id}`
- `helixgitpx_events_produced_total{topic}`
- `helixgitpx_sync_jobs_total{status}`
- `helixgitpx_conflicts_detected_total{kind}`
- `helixgitpx_ai_requests_total{task, model}`
- `helixgitpx_upstream_api_rate_limit_remaining{provider}`

---

## 4. Logging

- **Structured JSON** logs; schema in `helix-platform/logging`.
- Required fields: `ts, level, service, version, trace_id, span_id, msg`.
- Optional: `org_id, repo_id, user_id, correlation_id, error.type, error.stack`.
- Never log: secrets, tokens, full request bodies with PII, passwords, OTPs.
- **Log levels**: `DEBUG, INFO, WARN, ERROR`. Prod default `INFO`.
- **Sampling** for high-volume debug logs (head-based via OTel processor).
- **Retention**:
  - Application logs: 30 d hot in Loki, 90 d cold.
  - Audit mirror: 7 y in OpenSearch (legal).
  - Security-relevant logs: 1 y hot, 7 y cold.

---

## 5. Tracing

- **Sampling**: tail-based on the collector — keep 100 % of error/slow traces, 1 % of successful (configurable per org).
- Required spans: gRPC server / client, HTTP handlers, DB queries, Kafka produce / consume, Redis ops, upstream adapter calls, LLM invocations.
- Span attributes: stable keys from OTel semantic conventions.
- **Exemplars**: every p99 metric point links to an exemplar trace in Tempo.

---

## 6. Continuous Profiling

- **Pyroscope** agents (eBPF) on every node.
- Per-service CPU + heap profiles sampled continuously.
- Used for regression detection and hot-path discovery.
- Flame graphs accessible from Grafana.

---

## 7. SLIs & SLOs

### 7.1 User-Facing SLOs

| Service / Endpoint | SLI | SLO |
|---|---|---|
| API Gateway (reads) | Success rate | ≥ 99.95 % rolling 30 d |
| | Latency p99 | ≤ 300 ms |
| API Gateway (writes) | Success rate | ≥ 99.9 % |
| | Latency p99 | ≤ 500 ms |
| Git push ingress | Success rate | ≥ 99.9 % |
| | Time to replicate to first upstream | p99 ≤ 5 s |
| Live events | Delivery freshness (produce → client) | p99 ≤ 500 ms |
| | Subscription uptime | ≥ 99.95 % |
| Sync orchestrator | End-to-end sync success rate | ≥ 99.5 % |
| | Sync latency (small repo) p99 | ≤ 15 s |
| Conflict resolver auto-rate | % resolved without human | ≥ 75 % |
| Search | p99 latency | ≤ 150 ms |
| AI (conflict proposal) | p99 latency | ≤ 3 s |
| | Acceptance rate | ≥ 60 % |

### 7.2 Platform SLOs

- Postgres write latency p99 ≤ 20 ms.
- Kafka produce latency p99 ≤ 30 ms.
- Redis GET p99 ≤ 2 ms.
- Node availability ≥ 99.99 %.

### 7.3 Error Budgets

- Each SLO yields an error budget.
- Release freeze trigger: budget < 20 % with 50 % of window left; no feature launches until recovered.
- Dashboards show remaining budget and burn rate.

---

## 8. Alerting

Alerts are tiered. Fire on **symptom**, not cause, where possible.

| Severity | Response | Example |
|---|---|---|
| **P1** | Page (24/7) | API 5xx > 2 % for 5 min; data loss suspected |
| **P2** | Page (business hours), Slack always | Consumer lag high; upstream rate-limit near zero |
| **P3** | Slack only | Disk > 75 %; dependency CVE > High |
| **P4** | Ticket | Flaky test; deprecation warning |

### 8.1 Critical Alerts (examples)

- **ServiceDown** — `up == 0` for 2 min (P1).
- **ErrorBudgetFastBurn** — 2 % budget spent in 1 h (P1).
- **ErrorBudgetSlowBurn** — 10 % spent in 6 h (P2).
- **KafkaUnderReplicatedPartitions** — > 0 for 5 min (P1).
- **PostgresReplicationLag** — > 30 s (P1).
- **RedisMemoryHigh** — > 85 % (P2).
- **AdapterCircuitOpen** — open > 1 min (P2).
- **ConflictBacklogGrowing** — p99 age > 6 h (P2).
- **SLAHealthIssue** — any managed customer's tenant-scoped SLO below threshold.

### 8.2 Alert Hygiene

- Every alert has a runbook link.
- No "unknown" pages — acknowledged and routed within 15 min.
- Alert-review ritual: weekly; tune or retire noisy alerts.
- Quarterly "silent day": ensure no non-actionable alerts in 7 d.

---

## 9. Dashboards

### 9.1 Standard dashboards

- **Platform overview**: health by region, SLO burn, open incidents, deploy activity.
- **Per-service RED**: rate, errors, duration with drill-down.
- **Kafka cluster**: broker health, topic throughput, consumer lag.
- **Postgres**: QPS, latency, replication, locks, cache hit.
- **LLM**: tokens/s, accept rate, cost.
- **Conflicts**: detected, auto-resolved, escalated, backlog.
- **Tenant health**: per-org view for customer success.
- **Cost**: per-service, per-tenant (Kubecost + our meter).

### 9.2 Authoring & Sharing

- Dashboards JSON in Git; reviewed like code.
- Mixin pattern: reusable dashboard fragments (`helixgitpx-mixins`).
- Grafana folder per team.

---

## 10. Synthetic Monitoring

- **Checkly / Blackbox exporter / custom k6** runs probes from multiple regions every minute:
  - Login + create repo + push + pull.
  - API health endpoints.
  - WebSocket subscribe.
  - Upstream connectivity.
- Latency + availability dashboards per probe.
- Alerts if any probe fails repeatedly.

---

## 11. Audit & Compliance Visibility

- Audit queries exposed via [05-apis] `/api/v1/audit/events`.
- Dashboards for:
  - Admin actions per day.
  - Failed logins.
  - Role changes.
  - Data export requests.
- Weekly compliance report auto-generated.

---

## 12. Observability of Observability

- Meta-metrics: collector queue length, dropped logs/metrics/spans, ingester lag.
- Alerts on the observability platform itself.
- If Prometheus is down, we still want to know; Blackbox + mobile paging fallback.

---

## 13. Cost & Retention

- Aggressive sampling on low-value signals.
- Ship only what you will look at.
- Long-term retention:
  - Metrics: 1 y at 5 m resolution, 13 m at raw via downsampling.
  - Logs: 30 d hot, 90 d cold.
  - Traces: 7 d (critical — error/slow — 90 d).
  - Profiles: 14 d.
- Tenant cost attribution via labels → billing meter.

---

## 14. Data Classification in Observability

- Telemetry redaction rules in the Collector.
- Redact emails / tokens / secrets via regex + allowlists.
- Tenant-scoped RBAC in Grafana (Enterprise or OSS + proxy) so customers (enterprise tier) can see only their own dashboards.

---

## 15. Incident Management

- PagerDuty / Opsgenie as the on-call system.
- Incident channel auto-created in Slack with responders, status, roles.
- Status page (`status.helixgitpx.example.com`) updated via **statuspage-bot** tied to incident creation.
- Post-incident: blameless retro template; 5-whys; action items tracked.

---

*— End of Observability —*
