# Runbook — Kafka broker failure

**Alert:** `KafkaBrokerDown` (Strimzi metric: `strimzi_kafka_broker_up == 0`)

## Immediate checks (1 min)

- `kubectl -n kafka get pods | grep kafka-`
- `kubectl -n kafka logs <broker-pod> --tail=200`
- Confirm quorum: `kubectl -n kafka exec helixgitpx-kafka-0 -- bin/kafka-metadata-quorum.sh`

## Diagnosis

If pod pending/evicted, check node pressure: `kubectl top nodes`.
If pod crashlooping, check logs for `OutOfMemoryError`, `Corrupt record`, or `Too many open files`.

## Remediation

1. **Single broker down** — Strimzi auto-recovers; monitor for 5 min.
2. **Multiple brokers down** — drain affected node, cordon, let Strimzi reschedule.
3. **Corrupt log segment** — delete segment per Strimzi docs and let ISR catch up.

## Escalation

Persistent ISR shrink beyond 10 min → page SRE lead.
