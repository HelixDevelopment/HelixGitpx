# Kafka Topic Catalog

> Authoritative inventory of every Kafka topic in HelixGitpx: purpose, producers, consumers, partition key, retention, compaction, ACLs, and schema subject.

---

## Naming Convention

```
helixgitpx.<bounded-context>.<entity>.<action>[.<qualifier>]
```

- Bounded contexts: `auth`, `org`, `repo`, `upstream`, `sync`, `conflict`, `collab`, `ai`, `audit`, `billing`, `notify`, `policy`, `search`.
- Breaking schema change → new topic with `.vN` suffix; old remains until consumer migration is complete.
- Dead-letter topics: suffix `.dlq`.
- Compacted topics: suffix `.state` or documented per entry below.

---

## Cross-Cutting Defaults

- Replication factor: **3** in production.
- `min.insync.replicas`: **2**.
- Producer: idempotent + acks=all.
- Consumers: SASL/OAUTHBEARER via SPIFFE SVID.
- Schema: Avro in Karapace; envelope common across all topics.
- Tiered storage enabled in production for topics marked (T).

---

## Core Topics

### `helixgitpx.auth.events`
- **Producer**: auth-service
- **Consumers**: audit-projector, search-projector, notify-dispatcher, tenant-meter
- **Partitions**: 30
- **Key**: `user_id`
- **Retention**: 7 d (events flow to audit quickly)
- **Compaction**: no
- **Purpose**: login / logout / mfa / pat create / session revoke

### `helixgitpx.org.events`
- **Producer**: org-service
- **Consumers**: repo-service (invalidation), audit-projector, search-projector
- **Partitions**: 30
- **Key**: `org_id`
- **Retention**: 30 d (T)
- **Purpose**: org / team / membership create/update/delete

### `helixgitpx.repo.events`
- **Producer**: repo-service
- **Consumers**: sync-orchestrator, search-projector, audit-projector, ai-indexer, tenant-meter
- **Partitions**: 60
- **Key**: `repo_id`
- **Retention**: **compact** (event-sourced aggregate); 7 d before compaction
- **Compaction**: yes (`cleanup.policy=compact,delete`)
- **Purpose**: repo create / update / archive / transfer / settings / topics

### `helixgitpx.repo.ref.updated`
- **Producer**: git-ingress, adapter-pool (from upstream)
- **Consumers**: sync-orchestrator, conflict-detector, search-projector, live-events-fanout
- **Partitions**: 120
- **Key**: `repo_id`
- **Retention**: 30 d (T)
- **Purpose**: every branch/tag SHA change; includes `origin` (local|upstream|system|ai)

### `helixgitpx.repo.ref.watch.keepalive`
- **Producer**: git-ingress
- **Consumers**: live-events-service
- **Partitions**: 30
- **Key**: `repo_id`
- **Retention**: 1 d
- **Purpose**: heartbeat used to prove a subscription is alive in gRPC streams

### `helixgitpx.repo.commit.ingested`
- **Producer**: git-ingress
- **Consumers**: ai-indexer, search-projector, audit-projector
- **Partitions**: 60
- **Key**: `repo_id`
- **Retention**: 14 d
- **Purpose**: new commits visible to indexers

### `helixgitpx.repo.lfs.event`
- **Producer**: git-ingress
- **Consumers**: tenant-meter, audit-projector
- **Partitions**: 30
- **Key**: `repo_id`
- **Retention**: 30 d
- **Purpose**: LFS object uploaded / referenced / orphaned

---

## Upstream / Adapter / Webhook

### `helixgitpx.upstream.events`
- **Producer**: upstream-service
- **Consumers**: adapter-pool, audit-projector
- **Partitions**: 30
- **Key**: `upstream_id`
- **Retention**: 30 d
- **Purpose**: upstream CRUD, enabled/disabled, credential rotations

### `helixgitpx.upstream.webhook.received`
- **Producer**: webhook-gateway
- **Consumers**: conflict-detector, sync-orchestrator, audit-projector
- **Partitions**: 60
- **Key**: `repo_id` when derivable, else `upstream_id`
- **Retention**: 7 d (T)
- **Purpose**: canonicalised inbound webhook payloads (post-HMAC-verify, post-dedup)

### `helixgitpx.upstream.webhook.dlq`
- **Partitions**: 10
- **Retention**: 180 d
- **Purpose**: webhooks that failed normalisation / validation

### `helixgitpx.adapter.circuit`
- **Producer**: adapter-pool
- **Consumers**: observability-bridge (→ alerts), sync-orchestrator
- **Partitions**: 10
- **Key**: `provider`
- **Retention**: 7 d
- **Purpose**: circuit state changes per provider

---

## Sync & Conflict

### `helixgitpx.sync.job.events`
- **Producer**: sync-orchestrator (Temporal)
- **Consumers**: search-projector, audit-projector, tenant-meter
- **Partitions**: 60
- **Key**: `repo_id`
- **Retention**: 30 d
- **Purpose**: sync job lifecycle — queued, running, step, succeeded, failed, partial

### `helixgitpx.sync.completed`
- **Producer**: sync-orchestrator
- **Consumers**: live-events-fanout, tenant-meter
- **Partitions**: 60
- **Key**: `repo_id`
- **Retention**: 30 d
- **Purpose**: final result of a FanOutPush / InboundReconcile

### `helixgitpx.conflict.detected`
- **Producer**: conflict-detector
- **Consumers**: conflict-resolver, notify-dispatcher, live-events-fanout, audit-projector
- **Partitions**: 30
- **Key**: `repo_id`
- **Retention**: 180 d
- **Purpose**: every conflict case detection

### `helixgitpx.conflict.resolved`
- **Producer**: conflict-resolver
- **Consumers**: live-events-fanout, notify-dispatcher, audit-projector, ai-feedback-curator
- **Partitions**: 30
- **Key**: `repo_id`
- **Retention**: 180 d
- **Purpose**: applied / escalated / undone / cancelled

### `helixgitpx.conflict.ai.feedback`
- **Producer**: conflict-resolver, web/mobile clients
- **Consumers**: ai-feedback-curator
- **Partitions**: 30
- **Key**: `org_id`
- **Retention**: 18 months (per privacy policy)
- **Purpose**: accept / reject / edit feedback on proposals

---

## Collab (PRs / Issues / Comments)

### `helixgitpx.pr.events`
- **Producer**: pr-service
- **Consumers**: search-projector, ai-indexer, notify-dispatcher, audit-projector
- **Partitions**: 60
- **Key**: `repo_id`
- **Retention**: compact + 30 d
- **Compaction**: yes
- **Purpose**: PR create/update/close/merge/reopen/review

### `helixgitpx.issue.events`
- **Producer**: issue-service
- **Consumers**: search-projector, ai-indexer, notify-dispatcher
- **Partitions**: 60
- **Key**: `repo_id`
- **Retention**: compact + 30 d
- **Purpose**: issue lifecycle

### `helixgitpx.collab.crdt.ops`
- **Producer**: crdt-service
- **Consumers**: collab.issue / pr projectors; live-events-fanout
- **Partitions**: 60
- **Key**: `doc_id`
- **Retention**: compact + 30 d
- **Purpose**: metadata CRDT ops (labels, milestones, issue body)

### `helixgitpx.release.events`
- **Producer**: repo-service
- **Consumers**: search-projector, notify-dispatcher
- **Partitions**: 30
- **Key**: `repo_id`
- **Retention**: 180 d
- **Purpose**: release published / draft / asset uploaded

---

## AI

### `helixgitpx.ai.prompt.run`
- **Producer**: ai-service
- **Consumers**: ai-feedback-curator, tenant-meter, audit-projector
- **Partitions**: 30
- **Key**: `org_id`
- **Retention**: 90 d
- **Purpose**: every prompt execution (metadata only; no content)

### `helixgitpx.ai.finetune.requested`
- **Producer**: ai-feedback-curator
- **Consumers**: ai-finetune-worker
- **Partitions**: 5
- **Key**: `model`
- **Retention**: 30 d
- **Purpose**: trigger a training job

### `helixgitpx.ai.model.promoted`
- **Producer**: ai-service
- **Consumers**: live-events-fanout (admin users), audit-projector
- **Partitions**: 5
- **Retention**: 365 d
- **Purpose**: new active model version

---

## Audit / Billing / Policy

### `helixgitpx.audit.events`
- **Producer**: every service via outbox
- **Consumers**: audit-projector, compliance-exporter, siem-bridge
- **Partitions**: 60
- **Key**: `org_id` (fallback `system`)
- **Retention**: 365 d on Kafka + long-term in OpenSearch / PG
- **Purpose**: authoritative audit stream

### `helixgitpx.billing.usage`
- **Producer**: tenant-meter (fan-in from many sources)
- **Consumers**: billing-meter, tenant-meter
- **Partitions**: 30
- **Key**: `org_id`
- **Retention**: 90 d
- **Purpose**: metered usage events

### `helixgitpx.billing.invoice.generated`
- **Producer**: billing-service
- **Consumers**: notify-dispatcher, payments-bridge
- **Partitions**: 5
- **Retention**: 365 d
- **Purpose**: invoice produced

### `helixgitpx.policy.bundle.deployed`
- **Producer**: policy-service
- **Consumers**: every service (in-process OPA refresh listener)
- **Partitions**: 5
- **Retention**: 90 d
- **Purpose**: rollout of a new OPA bundle version

---

## Notify / Search / Live

### `helixgitpx.notify.requested`
- **Producer**: notify-subscribers (projections from domain events)
- **Consumers**: notify-dispatcher
- **Partitions**: 30
- **Key**: `user_id`
- **Retention**: 7 d
- **Purpose**: notification candidates waiting to be delivered

### `helixgitpx.notify.delivery.status`
- **Producer**: notify-dispatcher
- **Consumers**: audit-projector, live-events-fanout
- **Partitions**: 30
- **Retention**: 30 d
- **Purpose**: sent / failed / dropped per channel

### `helixgitpx.search.index.changed`
- **Producer**: search-projector
- **Consumers**: observability-bridge
- **Partitions**: 10
- **Retention**: 7 d
- **Purpose**: telemetry on index throughput (not a hot path)

### `helixgitpx.live.fanout.firehose`
- **Producer**: live-events-fanout
- **Consumers**: live-events-service instances
- **Partitions**: 60
- **Key**: `scope` (derived)
- **Retention**: 1 h
- **Purpose**: in-cluster fan-out of live events to subscribers

---

## Dead Letter Topics

| Source | DLQ topic | Retention |
|---|---|---|
| `helixgitpx.repo.ref.updated` | `helixgitpx.repo.ref.updated.dlq` | 90 d |
| `helixgitpx.sync.completed`   | `helixgitpx.sync.completed.dlq`   | 90 d |
| `helixgitpx.upstream.webhook.received` | `helixgitpx.upstream.webhook.received.dlq` | 180 d |
| `helixgitpx.conflict.detected` | `helixgitpx.conflict.detected.dlq` | 180 d |
| `helixgitpx.audit.events`      | `helixgitpx.audit.events.dlq`      | 365 d |
| `helixgitpx.ai.prompt.run`     | `helixgitpx.ai.prompt.run.dlq`     | 90 d |

All DLQs consumed by `dlq-inspector` service; operators review via `helixctl dlq list/replay`.

---

## Access Control

Topic ACLs enforced via SASL/OAUTHBEARER + Kafka ACLs:

- Producers: per-service principal (`spiffe://.../ns/helixgitpx/sa/<svc>`).
- Consumers: consumer group = `<service>-<purpose>`; read ACL only for topics they need.
- No `*` wildcards in production ACLs.
- Admin topics (`__consumer_offsets`, `__transaction_state`) inaccessible to services.

---

## Schema Registry Subjects

Every topic `T` has a subject `T-value` in Karapace. Key subjects (`T-key`) exist only for entities whose key is complex (rare).

Compatibility: **BACKWARD** default; **BACKWARD_TRANSITIVE** for critical topics (`helixgitpx.audit.events`, `helixgitpx.billing.usage`).

---

## Retention Overrides

A topic's retention may be temporarily extended by an operator for investigation:

```bash
helixctl kafka topic alter \
  --topic helixgitpx.repo.ref.updated \
  --retention 30d → 90d \
  --reason "HGX-4501 investigation" \
  --ttl 14d
```

Auto-reverts on TTL expiry; audit-logged.

---

## Adding a Topic

1. Decide name per convention.
2. Register Avro schema in Karapace.
3. Add entry in this catalog.
4. Update infra-as-code (Strimzi `KafkaTopic` CR).
5. Add metrics / dashboards / alerts.
6. Document producers + consumers in service READMEs.
7. Link to / update ADR if this introduces a new event pattern.

---

*— End of Topic Catalog —*
