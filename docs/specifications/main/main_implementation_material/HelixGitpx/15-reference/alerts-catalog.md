# Alerts Catalog

> Every alert has: name, severity, signal, expression, rationale, runbook. Alerts fire on **symptoms**, not causes, unless a cause is both unambiguous and actionable.

---

## Severity Scale

| Sev | Trigger | Response SLA |
|---|---|---|
| **P1** | Customer-facing outage, data-loss risk, security breach | Page 24/7; ack ≤ 15 min |
| **P2** | Significant degradation; fast SLO burn | Page business hours, Slack 24/7 |
| **P3** | Noticeable but non-breaking | Slack |
| **P4** | Technical debt / deprecation | Ticket |

Every alert must have:
- `severity` label.
- `service` label.
- `runbook_url` annotation.
- `summary` (≤ 80 chars) and `description` annotations.
- `dashboard_url` annotation.

---

## Platform Health

### KubernetesNodeNotReady

```yaml
severity: P1
expr: kube_node_status_condition{condition="Ready",status!="true"} == 1
for: 5m
summary: "Node {{ $labels.node }} is NotReady"
runbook: RB-001
```

### PodCrashLoopBackOff

```yaml
severity: P2
expr: rate(kube_pod_container_status_restarts_total[10m]) > 0.2
for: 15m
runbook: RB-002
```

### CertificateNearExpiry

```yaml
severity: P2
expr: (certmanager_certificate_expiration_timestamp_seconds - time()) < 86400 * 14
runbook: RB-004
```

### SPIREServerDown

```yaml
severity: P1
expr: up{job="spire-server"} == 0
for: 3m
runbook: RB-006
```

---

## Postgres

### PostgresPrimaryDown

```yaml
severity: P1
expr: pg_up{role="primary"} == 0
for: 1m
runbook: RB-010
```

### PostgresReplicationLagHigh

```yaml
severity: P1
expr: pg_replication_lag_seconds > 30
for: 5m
runbook: RB-011
```

### PostgresDiskAlmostFull

```yaml
severity: P2
expr: (node_filesystem_avail_bytes{mountpoint="/var/lib/postgresql/data"} / node_filesystem_size_bytes) < 0.15
for: 15m
runbook: RB-012
```

### PostgresSlowQueries

```yaml
severity: P3
expr: pg_stat_statements_slow_queries_total > 0
for: 30m
runbook: RB-013
```

---

## Kafka

### KafkaUnderReplicatedPartitions

```yaml
severity: P1
expr: kafka_server_replicamanager_underreplicatedpartitions > 0
for: 5m
runbook: RB-020
```

### KafkaConsumerLagCritical

```yaml
severity: P2
expr: helixgitpx_kafka_consumer_lag > 10000
for: 5m
runbook: RB-021
```

### KafkaBrokerDiskFull

```yaml
severity: P1
expr: kafka_log_log_size_bytes / kafka_log_log_size_limit_bytes > 0.85
for: 10m
runbook: RB-022
```

### KafkaDLQFlood

```yaml
severity: P2
expr: rate(helixgitpx_kafka_dlq_events_total[5m]) > 100
for: 5m
runbook: RB-021
```

### SchemaRegistryDown

```yaml
severity: P1
expr: up{job="karapace"} == 0
for: 3m
runbook: RB-023
```

---

## Redis / Search / Vector

### RedisMemoryHigh

```yaml
severity: P2
expr: redis_memory_used_bytes / redis_memory_max_bytes > 0.85
for: 10m
runbook: RB-030
```

### OpenSearchClusterRed

```yaml
severity: P1
expr: opensearch_cluster_status{color="red"} == 1
for: 5m
runbook: RB-040
```

### QdrantOOMRisk

```yaml
severity: P2
expr: container_memory_working_set_bytes{container="qdrant"} / container_spec_memory_limit_bytes > 0.9
for: 10m
runbook: RB-041
```

### MeilisearchIndexUnhealthy

```yaml
severity: P2
expr: meilisearch_indexing_pending > 100000
for: 15m
runbook: RB-042
```

---

## Services

### APIGateway5xxHigh

```yaml
severity: P1
expr: |
  sum by (service) (rate(helixgitpx_http_requests_total{service="api-gateway",code=~"5.."}[5m]))
  / sum by (service) (rate(helixgitpx_http_requests_total{service="api-gateway"}[5m]))
  > 0.02
for: 5m
runbook: RB-100
```

### AuthTokenMintFailures

```yaml
severity: P1
expr: rate(helixgitpx_auth_tokens_rejected_total[5m]) > 50
for: 5m
runbook: RB-101
```

### SyncOrchestratorBacklog

```yaml
severity: P2
expr: helixgitpx_kafka_consumer_lag{group="sync-orchestrator-main"} > 5000
for: 10m
runbook: RB-110
```

### WebhookDedupCacheDown

```yaml
severity: P2
expr: up{job="webhook-gateway",instance=~".*redis.*"} == 0
for: 2m
runbook: RB-111
```

### ConflictQueueGrowing

```yaml
severity: P2
expr: histogram_quantile(0.99, rate(helixgitpx_conflicts_time_to_resolve_seconds_bucket[1h])) > 21600
for: 30m
runbook: RB-120
```

### ConflictAutoResolveDrop

```yaml
severity: P2
expr: helixgitpx_conflicts_auto_resolution_ratio < 0.6
for: 1h
runbook: RB-121
```

### AdapterCircuitOpen

```yaml
severity: P2
expr: helixgitpx_adapter_circuit_breaker_state == 2
for: 1m
runbook: RB-130
```

### UpstreamAuthFailureSpike

```yaml
severity: P2
expr: rate(helixgitpx_adapter_auth_failures_total[5m]) > 5
for: 5m
runbook: RB-131
```

### UpstreamRateLimitLow

```yaml
severity: P2
expr: helixgitpx_adapter_rate_limit_remaining < 100
for: 5m
runbook: RB-130
```

### LiveEventsDeliveryLagHigh

```yaml
severity: P2
expr: histogram_quantile(0.99, rate(helixgitpx_events_delivery_seconds_bucket[5m])) > 0.5
for: 5m
runbook: RB-140
```

### LiveEventsDroppedSpike

```yaml
severity: P2
expr: rate(helixgitpx_events_dropped_total[5m]) > 100
for: 5m
runbook: RB-140
```

### AIInferenceTimeouts

```yaml
severity: P2
expr: |
  sum(rate(helixgitpx_ai_requests_total{status="timeout"}[5m]))
  / sum(rate(helixgitpx_ai_requests_total[5m]))
  > 0.05
for: 10m
runbook: RB-150
```

### AIModelRegression

```yaml
severity: P2
expr: (helixgitpx_ai_accept_rate - helixgitpx_ai_accept_rate offset 6h) < -0.05
for: 6h
runbook: RB-151
```

### BillingMeterLag

```yaml
severity: P3
expr: helixgitpx_kafka_consumer_lag{group="billing-meter"} > 10000
for: 30m
runbook: RB-160
```

---

## Security

### SuspectedCredentialCompromise

```yaml
severity: P1
expr: increase(helixgitpx_security_anomalies_total{kind="geo_jump"}[10m]) > 0
runbook: RB-200
```

### SecretLeakDetected

```yaml
severity: P1
expr: increase(helixgitpx_gitleaks_findings_total{severity="critical"}[1h]) > 0
runbook: RB-201
```

### DDoSDetected

```yaml
severity: P1
expr: rate(cloudflare_requests_blocked_total[1m]) > 10000
for: 2m
runbook: RB-203
```

### UnsignedImageAdmitted

```yaml
severity: P1
expr: increase(kyverno_policy_results_total{result="fail",reason="unsigned"}[5m]) > 0
runbook: RB-204
```

---

## SLO Burn Rate

### FastBurn — 2% budget in 1h

```yaml
severity: P1
expr: helixgitpx_slo_burn_rate{window="1h"} > 14.4
for: 2m
```

### SlowBurn — 10% in 6h

```yaml
severity: P2
expr: helixgitpx_slo_burn_rate{window="6h"} > 6
for: 15m
```

Burn rate math follows Google SRE workbook: 2 % / 1 h ≈ 14.4× baseline.

---

## Meta — Observability Health

### PrometheusDown

```yaml
severity: P1
expr: up{job="prometheus"} == 0
for: 3m
```

### ObservabilityCollectorDropping

```yaml
severity: P2
expr: rate(otelcol_processor_dropped_spans[5m]) > 100
for: 10m
```

---

## Alert Hygiene

- Any alert firing > 3×/week without action → review; tighten or retire.
- No orphan alerts — every alert maps to a runbook ID.
- Silence expiries never > 7 days without justification.
- Weekly alert-review ritual owned by the SRE rotation.
