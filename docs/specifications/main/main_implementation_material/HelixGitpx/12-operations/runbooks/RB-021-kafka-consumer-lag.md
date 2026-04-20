# RB-021 — Kafka Consumer Lag Critical

> **Severity default**: P2 (P1 if lag growing > 10× baseline)
> **Owner**: Platform SRE
> **Last tested**: 2026-04-01

## 1. Detection

Alert: `KafkaConsumerLagCritical` — `helixgitpx_kafka_consumer_lag > 10000` for 5 min.

Other signs:
- Live events noticeably delayed from customer reports.
- Downstream services (search indexer, projections) falling behind.
- `helixgitpx_events_delivery_seconds` p99 degrading.

---

## 2. First Moves

1. Ack alert; open incident if growth rate suggests escalation.
2. Identify affected group: `{{ $labels.group }}` from alert.
3. Check Grafana dashboard "Kafka Consumer Lag" — which partitions, how fast growing?
4. Check whether a **deploy** happened recently on the consuming service.

---

## 3. Diagnosis

Most common causes, in order of frequency:

| Cause | Symptom | Action |
|---|---|---|
| Slow consumer (single-partition hotspot) | One partition lag spikes while others fine | Redistribute (rebalance or repartition) |
| Consumer OOMing / crashing | Pod restarts spike, lag grows per restart | Fix memory; increase limit |
| Downstream DB slow | Consumer latency p99 degraded | Check DB; throttle if needed |
| Schema registry flaky | Deserialisation errors | Check Karapace health |
| Too few replicas | Throughput capped | Scale via KEDA override |
| Quota / ACL change | Consumer `NotAuthorizedException` | Check ACLs |
| Poison message | Specific offset repeatedly failing | Move to DLQ with documented offset |

### Commands

```bash
# Show lag per partition for the group
kubectl -n kafka exec kafka-0 -- bin/kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --describe --group <GROUP>

# Service logs
helixctl logs <service> --since=5m | grep -E "error|lag|retry" | head -50

# Consumer pod metrics
kubectl -n helixgitpx top pods -l app.kubernetes.io/name=<service>

# Most recent consumed offset per partition
helixctl kafka consumer offsets --group <GROUP>
```

---

## 4. Mitigations

### 4.1 Scale out

```bash
# Temporary replica bump (before KEDA takes over)
helixctl scale <service> --replicas=<current * 2>

# Or trigger KEDA override
kubectl -n helixgitpx annotate scaledobject <name> \
  autoscaling.keda.sh/paused-replicas="<N>"
```

Note: partitions cap parallelism. If lag is on one partition, more replicas won't help until you repartition.

### 4.2 Throttle downstream dependency

If DB is the bottleneck, reduce producer rate on the ingestion side OR pause consumer briefly:

```bash
helixctl kafka group pause --group <GROUP> --for=10m
```

### 4.3 Quarantine a poison message

Find the offset that repeatedly fails (logs show same offset in retries):

```bash
helixctl kafka message show --topic <t> --partition <p> --offset <o>
helixctl kafka move --topic <t> --partition <p> --offset <o> --to dlq
```

Then skip:

```bash
helixctl kafka offset set --group <GROUP> --topic <t> --partition <p> --offset <o+1>
```

Always log the action and the offset into the incident channel.

### 4.4 Repartition

If a partition is chronically hot (tenant skew), repartition the topic:

```bash
kubectl -n kafka exec kafka-0 -- bin/kafka-reassign-partitions.sh \
  --bootstrap-server localhost:9092 \
  --reassignment-json-file /tmp/plan.json \
  --execute
```

This is a platform-level operation — coordinate with the Kafka owner.

### 4.5 Tiered storage

If disk full was the cause (linked to RB-022), tiered storage may already be offloading; verify with metrics.

---

## 5. Verify Recovery

- `helixgitpx_kafka_consumer_lag{group=<GROUP>}` trending toward 0.
- `helixgitpx_events_delivery_seconds` p99 back within SLO.
- No new `kafka.dlq` entries.

---

## 6. Post-Incident

- Document the root cause: hotspot / deploy regression / dependency / poison / capacity.
- Open follow-ups:
  - If hotspot: evaluate partition key.
  - If poison: add validation at producer side.
  - If dependency: capacity plan.
- Update this runbook with lessons learned.

---

## 7. Related

- RB-020 (Under-replicated partitions)
- RB-022 (Broker disk full)
- RB-023 (Schema registry down)
- [18-observability.md §8] — alerting philosophy

---

## 8. Drill

Replay in staging:

```bash
# Pause consumer group for 5 min to build lag
helixctl kafka group pause --group <GROUP> --for=5m --env=staging

# Observe alert fires on expected cadence
# Resume and observe recovery
helixctl kafka group resume --group <GROUP>
```
