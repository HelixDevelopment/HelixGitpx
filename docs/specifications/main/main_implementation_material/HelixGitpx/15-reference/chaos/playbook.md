# Chaos Engineering Playbook

> We deliberately break HelixGitpx in controlled ways to prove our redundancy works and to find surprises. This document is the **curated list of chaos experiments**, their hypotheses, and how to run them safely.

Tooling: **Chaos Mesh** (K8s-native) + **Litmus** (scenario framework) + custom **chaos-bot** for orchestration. Experiments are YAML-committed in `chaos/experiments/` and run in staging continuously + prod on Game Days.

---

## Principles

1. **Hypothesis first.** Every experiment states what we expect to happen. If reality differs, it's either a bug or a learning.
2. **Blast radius capped.** Experiments target defined namespaces / pods / traffic percentages. Abort switch is one command.
3. **Observability required.** No experiment runs without dashboards watching.
4. **Communication.** Game Days announced ≥ 24 h; continuous staging chaos advertised on status-internal.
5. **Learn and ship.** Every failed hypothesis → issue → fix → regression test.

---

## Tier 1 — Continuous (run in staging 24/7)

Small, bounded, low-impact. Catches simple regressions.

### C1.1 — Random pod kill (any service)
**Hypothesis**: service pods restart within 30 s without client impact beyond retries.
**Chaos**: `PodChaos` every 30 min, one random pod per experiment.
**Safeguards**: exclude pods in stateful sets holding leader lease; skip during deploys.
**Success**: synthetic probes remain green.

### C1.2 — Network delay injection (5 % of traffic)
**Hypothesis**: added 100 ms p99 latency on inter-service calls stays within error budget.
**Chaos**: `NetworkChaos` `delay` for 10 min every 4 h.
**Success**: SLO burn-rate stays < 2×.

### C1.3 — CPU pressure on app pods
**Hypothesis**: HPA scales up when CPU is pressured, and services remain responsive.
**Chaos**: `StressChaos` CPU 80 % for 15 min on random app pod.

### C1.4 — DNS flake (5 % of lookups)
**Hypothesis**: clients retry DNS and recover.
**Chaos**: `DNSChaos` `random` on 5 % of requests.

---

## Tier 2 — Weekly Scheduled

More impactful — run in staging with alerting muted for expected symptoms.

### C2.1 — Kill one Postgres primary
**Hypothesis**: CNPG promotes a replica in ≤ 60 s; writes pause then resume; no data loss.
**Chaos**: delete primary pod.
**Success**: `pg_up{role="primary"}` recovers ≤ 90 s; `sync_replication_lag` returns to baseline; client-side errors only in that window, all retryable.

### C2.2 — Kafka broker eviction
**Hypothesis**: under-replicated partitions appear briefly then recover; no consumer lag build-up > 1 min.
**Chaos**: cordon + drain one broker; let Strimzi reschedule.

### C2.3 — OpenSearch node down
**Hypothesis**: cluster turns yellow, queries served from other nodes; recovery on node return.

### C2.4 — Vault sealed (soft)
**Hypothesis**: services holding cached SVIDs remain functional; new token issuance paused; alert fires; unseal restores normal operation in ≤ 5 min.

### C2.5 — Upstream adapter 429 storm
**Hypothesis**: circuit opens per provider; work pauses for that upstream; degraded mode notification surfaces; other providers unaffected.
**Chaos**: `HTTPChaos` on adapter-pool egress to specific provider returns 429 for 10 min.

---

## Tier 3 — Monthly Game Days

Planned, cross-team. Write-up published.

### G3.1 — Full AZ failure
**Hypothesis**: multi-AZ topology survives; SLO unaffected after ≤ 3 min.
**Chaos**: cordon all nodes in one zone; drain.

### G3.2 — Full region failure (DR rehearsal)
**Hypothesis**: traffic cuts over to secondary region; writes queue on client, replay on recovery; RPO ≤ 30 s.
**Chaos**: block egress to primary region at the edge; promote secondary.

### G3.3 — LLM inference pool down
**Hypothesis**: conflict auto-apply rate temporarily drops; CRDT / policy paths unaffected; users get "AI unavailable" messaging on AI-only features.

### G3.4 — Object storage degraded
**Hypothesis**: LFS uploads/downloads retry; Git push accept pauses; alerts fire; customers see honest error messaging.

### G3.5 — Kafka mirror partition (cross-region)
**Hypothesis**: each region operates independently; divergence event logged; manual reconciliation tool converges state after link heals.

### G3.6 — Certificate expiry simulation
**Hypothesis**: rotation path triggers well before expiry; if it fails, page fires at T-14 days.
**Chaos**: set a cert to expire in 1 day in a staging env.

### G3.7 — Secret leak response
**Hypothesis**: tabletop — a fake leaked PAT detection triggers auto-revoke within 10 min and notifies the owner.

### G3.8 — Supply-chain compromise
**Hypothesis**: an unsigned image fails admission; audit fires; nothing reaches prod.

---

## Authoring an Experiment

1. Write a YAML under `chaos/experiments/<tier>/<id>-<short-name>.yaml`:

```yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: random-kill-repo-service
  namespace: helixgitpx
  labels:
    chaos.helixgitpx.io/tier: "1"
    chaos.helixgitpx.io/owner: "repo-service"
spec:
  action: pod-kill
  mode: one
  selector:
    namespaces: [helixgitpx]
    labelSelectors:
      app.kubernetes.io/name: repo-service
  scheduler:
    cron: "*/30 * * * *"
  duration: "10s"
```

2. Add hypothesis + success criteria in `chaos/experiments/<id>.md` (same basename).
3. PR review: chaos captain + service owner + SRE.
4. Merge → Argo CD deploys to staging.
5. Observe for one week; promote to "approved recurring".

---

## Safety

- **Kill-switch**: `helixctl chaos stop-all` disables all active experiments instantly.
- **Blocked periods**: `chaos-bot` refuses to start experiments during active P1 incidents, during change freezes, or when the cluster has degraded capacity.
- **Production experiments** require on-call acknowledgement AND two-person approval.
- **Customer-visible experiments**: only on opt-in infrastructure (our own dogfood org first).

---

## Reporting

Every experiment produces a report:

- Who / When / Target
- Hypothesis
- Observed metrics (traces, latencies, error rates)
- Hypothesis held? Variance?
- Action items (if any)
- Linked incidents (if any occurred)

Stored as PRs in `chaos/reports/YYYY/` with dashboards snapshotted.

---

## Metrics of the Chaos Program

- Number of experiments run / week.
- Number of action items opened.
- % of experiments where hypothesis held.
- MTTR improvements tracked over time.

Targets: ≥ 5 new experiments quarterly; ≥ 2 Game Days quarterly; all issues found triaged within 7 days.
