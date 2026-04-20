# ADR-0036 — MirrorMaker 2 over Confluent Replicator

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Cross-region Kafka replication options: MirrorMaker 2 (OSS, Strimzi-native),
Confluent Replicator (commercial, richer features), custom consumer→producer bridge.

## Decision

MirrorMaker 2. Strimzi ships it as a first-class CR (`KafkaMirrorMaker2`) with
topic/group whitelisting and exactly-once semantics.

## Consequences

- OSS-first stance preserved.
- Loss of Replicator's schema-aware routing; we accept this because topics are
  already namespaced by service.
- Monitoring is standard Kafka Connect metrics.

## Links

- Spec §LOCKED C-7
- impl/helixgitpx-platform/helm/mirrormaker2/
