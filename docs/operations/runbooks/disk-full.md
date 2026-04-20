# Runbook — Disk full on PVC

**Alert:** `PersistentVolumeUsageHigh` (>85%) or `PersistentVolumeFillingUp` (>95%).

## Immediate checks

- Identify: `kubectl get pvc -A --sort-by=.status.capacity.storage`.
- Confirm growth rate: query Prometheus `kubelet_volume_stats_used_bytes` over 24h.

## Remediation

1. **Git blob storage PVC** — expand PVC via `kubectl patch pvc <name> -p '{"spec":{"resources":{"requests":{"storage":"NEWGi"}}}}'` (StorageClass must allow expansion).
2. **Postgres WAL PVC** — run `VACUUM FULL` off-hours + trigger manual WAL archive; CNPG supports online resize.
3. **Log PVC** — rotate / truncate; verify log retention policies.
4. **Kafka PVC** — reduce retention in `KafkaTopic` CRs (e.g., `retention.ms=604800000`) and await compaction.

## Escalation

If resize not possible and disk >95% → move affected workload to new node with larger PVC.
