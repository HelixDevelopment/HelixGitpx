# Metrics Catalog

> Every metric HelixGitpx emits. Canonical prefix: **`helixgitpx_`**. Naming: `<noun>_<unit>`. Unit suffixes: `_total` (counter), `_seconds` / `_bytes` / `_ratio` / `_count` (gauge or histogram base unit).

Labels should be low-cardinality; never raw user IDs or IPs. Tenancy via `org_id` only when business-critical (and then at aggregate granularity).

---

## RED — Generic Per-Service

| Metric | Type | Labels | Notes |
|---|---|---|---|
| `helixgitpx_http_requests_total` | Counter | `service, method, path, code` | Path is normalised to route template |
| `helixgitpx_http_request_duration_seconds` | Histogram | `service, method, path, code` | Latency buckets |
| `helixgitpx_http_requests_in_flight` | Gauge | `service` | For HPA signal |
| `helixgitpx_grpc_server_started_total` | Counter | `service, grpc_service, grpc_method` | |
| `helixgitpx_grpc_server_handled_total` | Counter | `service, grpc_service, grpc_method, grpc_code` | |
| `helixgitpx_grpc_server_handling_seconds` | Histogram | … | |
| `helixgitpx_grpc_server_msg_received_total` | Counter | … | Streaming |
| `helixgitpx_grpc_server_msg_sent_total` | Counter | … | Streaming |
| `helixgitpx_service_up` | Gauge (1) | `service` | Always 1 when running |
| `helixgitpx_service_info` | Gauge (1) | `service, version, commit_sha, build_time` | |
| `helixgitpx_build_info` | Gauge (1) | `service, go_version, kotlin_version, node_version` | |

---

## Auth

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_auth_logins_total` | Counter | `result, method` |
| `helixgitpx_auth_mfa_challenges_total` | Counter | `result` |
| `helixgitpx_auth_tokens_minted_total` | Counter | `token_type` |
| `helixgitpx_auth_tokens_rejected_total` | Counter | `reason` |
| `helixgitpx_auth_refresh_reuse_detected_total` | Counter | — |
| `helixgitpx_auth_active_sessions` | Gauge | — |

---

## Git Ingress

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_git_push_events_total` | Counter | `repo_size_bucket, result` |
| `helixgitpx_git_push_bytes_total` | Counter | — |
| `helixgitpx_git_clone_bytes_total` | Counter | — |
| `helixgitpx_git_ref_updates_total` | Counter | `kind` |

---

## Sync

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_sync_jobs_total` | Counter | `status, trigger` |
| `helixgitpx_sync_duration_seconds` | Histogram | `status` |
| `helixgitpx_sync_fanout_targets` | Histogram | Number of upstreams per fan-out |
| `helixgitpx_sync_replication_lag_seconds` | Gauge | `upstream` |

---

## Upstream / Adapter

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_adapter_requests_total` | Counter | `provider, operation, status` |
| `helixgitpx_adapter_request_duration_seconds` | Histogram | `provider, operation` |
| `helixgitpx_adapter_rate_limit_remaining` | Gauge | `provider, upstream_id` |
| `helixgitpx_adapter_circuit_breaker_state` | Gauge | `provider` (0 closed 1 half 2 open) |
| `helixgitpx_adapter_auth_failures_total` | Counter | `provider, reason` |
| `helixgitpx_adapter_shadow_divergence_total` | Counter | `provider, kind` |
| `helixgitpx_webhook_events_total` | Counter | `provider, event, result` |
| `helixgitpx_webhook_dedup_hits_total` | Counter | `provider` |

---

## Conflict

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_conflicts_detected_total` | Counter | `kind` |
| `helixgitpx_conflicts_resolved_total` | Counter | `kind, strategy, decided_by` |
| `helixgitpx_conflicts_escalated_total` | Counter | `kind` |
| `helixgitpx_conflicts_auto_resolution_ratio` | Gauge | — |
| `helixgitpx_conflicts_time_to_resolve_seconds` | Histogram | `kind` |
| `helixgitpx_conflicts_open` | Gauge | `kind, priority` |
| `helixgitpx_conflict_undo_applied_total` | Counter | `kind` |

---

## AI

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_ai_requests_total` | Counter | `task, model, status` |
| `helixgitpx_ai_latency_seconds` | Histogram | `task, model` |
| `helixgitpx_ai_tokens_in_total` | Counter | `task, model` |
| `helixgitpx_ai_tokens_out_total` | Counter | `task, model` |
| `helixgitpx_ai_cost_usd` | Counter | `org_id, model` |
| `helixgitpx_ai_accept_rate` | Gauge | `task, model` |
| `helixgitpx_ai_confidence_bucket` | Histogram | `task` |
| `helixgitpx_ai_guardrail_reject_total` | Counter | `rail, reason` |
| `helixgitpx_ai_hallucination_score` | Gauge | `task, model` |
| `helixgitpx_ai_finetune_runs_total` | Counter | `status` |
| `helixgitpx_ai_model_active_version` | Gauge | `model, version` (1 active) |

---

## Search

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_search_queries_total` | Counter | `backend, index, status` |
| `helixgitpx_search_query_latency_seconds` | Histogram | `backend, index` |
| `helixgitpx_search_projector_lag_seconds` | Gauge | `consumer` |
| `helixgitpx_search_index_docs` | Gauge | `index` |

---

## Live Events

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_events_subscriptions_active` | Gauge | `transport` |
| `helixgitpx_events_sent_total` | Counter | `transport, event_type` |
| `helixgitpx_events_dropped_total` | Counter | `reason` |
| `helixgitpx_events_delivery_seconds` | Histogram | `transport` |
| `helixgitpx_events_connections_total` | Counter | `transport` |
| `helixgitpx_events_disconnects_total` | Counter | `reason` |
| `helixgitpx_events_resume_stale_total` | Counter | — |

---

## Data Plane — Kafka / Postgres / Redis / Object Store

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_kafka_consumer_lag` | Gauge | `group, topic, partition` |
| `helixgitpx_kafka_producer_errors_total` | Counter | `topic` |
| `helixgitpx_kafka_dlq_events_total` | Counter | `topic` |
| `helixgitpx_kafka_schema_registry_cache_hit_ratio` | Gauge | — |
| `helixgitpx_pg_query_duration_seconds` | Histogram | `service, op` |
| `helixgitpx_pg_deadlocks_total` | Counter | `service` |
| `helixgitpx_pg_replication_lag_seconds` | Gauge | `replica` |
| `helixgitpx_redis_ops_total` | Counter | `service, op, status` |
| `helixgitpx_redis_op_latency_seconds` | Histogram | `service, op` |
| `helixgitpx_object_store_ops_total` | Counter | `service, op, status` |

---

## Security / Audit

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_policy_decisions_total` | Counter | `subject_kind, action, effect` |
| `helixgitpx_policy_decision_seconds` | Histogram | — |
| `helixgitpx_audit_events_total` | Counter | `action, outcome` |
| `helixgitpx_audit_anchor_success_total` | Counter | — |
| `helixgitpx_security_failed_logins_total` | Counter | `reason` |
| `helixgitpx_security_anomalies_total` | Counter | `kind` |

---

## Billing

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_billing_usage_units_total` | Counter | `org_id, meter` |
| `helixgitpx_billing_quota_exceeded_total` | Counter | `org_id, meter` |
| `helixgitpx_billing_spend_usd` | Counter | `org_id, category` |

---

## Platform / Meta

| Metric | Type | Labels |
|---|---|---|
| `helixgitpx_tenancy_active_orgs` | Gauge | — |
| `helixgitpx_tenancy_active_users_24h` | Gauge | — |
| `helixgitpx_tenancy_active_repos` | Gauge | — |

---

## SLO-Derived

| Metric | Type | Labels | Meaning |
|---|---|---|---|
| `helixgitpx_slo_error_budget_remaining` | Gauge | `slo` | 0–1 |
| `helixgitpx_slo_burn_rate` | Gauge | `slo, window` | ×1 means on track |

---

## Guidelines

- **Unit consistency**: seconds (not ms), bytes (not MiB).
- **Histogram buckets**: latency `[0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]`; bytes `[256, 1KiB, 4KiB, 16KiB, 64KiB, 256KiB, 1MiB, 4MiB, 16MiB, 64MiB]`.
- **Label cardinality**: target ≤ 10 k series per metric per cluster. Avoid unbounded labels.
- **Derived metrics**: prefer recording rules in Mimir over computing in dashboards.
- **Quantiles**: compute via histogram_quantile; never emit pre-aggregated quantiles for services that scale horizontally.
