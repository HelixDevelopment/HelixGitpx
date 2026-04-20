# 34 — Cost Optimisation Guide

> **Document purpose**: Keep HelixGitpx's infrastructure spend aligned with its business value. Written for platform engineers, SREs, and FinOps partners. Covers how we measure cost, per-tenant attribution, and concrete optimisation levers with expected impact.

---

## 1. Principles

1. **Measure before optimising.** Kubecost + per-tenant usage events.
2. **Cost is a non-functional requirement.** Every major change declares a rough cost delta.
3. **Optimise for right-sized, not minimum.** Starving the system leads to outages, which are more expensive.
4. **Customer value first.** Never degrade customer SLOs to save money.
5. **Expose cost internally.** Dashboards per team; quarterly review.

---

## 2. Cost Breakdown (typical prod region, as of GA sizing)

| Component | Share | Notes |
|---|---|---|
| EKS + node compute | **55 %** | Stateless services, HPA-driven |
| Postgres (Aurora or CNPG on NVMe) | 12 % | HA + replicas |
| Kafka (MSK or Strimzi on NVMe) | 10 % | 3x replication + tiered storage |
| Object store (S3 / MinIO) | 8 % | Git packs, LFS, release assets, backups |
| GPU inference | 6 % | Scales with AI usage |
| OpenSearch / Meilisearch / Qdrant | 4 % | Indexes |
| Observability (Prometheus / Tempo / Loki) | 3 % | Retention-sensitive |
| Network egress | 2 % | Cross-AZ + out to upstreams |

Exact numbers vary by customer base + region pricing.

---

## 3. Measurement Stack

- **Kubecost** for K8s workloads; per-namespace + per-label cost allocation.
- **AWS Cost Explorer / GCP BigQuery billing export / Azure cost views** at provider level.
- **Custom usage events** (`helixgitpx.billing.usage`) feed tenant attribution.
- Unified "Cost & Capacity" Grafana folder merges these.
- Monthly FinOps review meeting.

Every service has a **"cost of this service"** dashboard with: compute, storage, network, per-tenant breakdown.

---

## 4. Per-Tenant Attribution

- Namespace labels (`org_id`) propagate from resource to Kubecost to billing-meter.
- Usage events are the authority for customer-facing billing; Kubecost is the authority for internal cost.
- Shared resources (e.g. Kafka cluster) allocated by **metered partition usage** (bytes produced + consumed).
- Object storage: prefix-based; bucket events tracked daily.

We compute a **per-org margin** monthly:

```
margin = revenue(org) − attributed_cost(org) − overhead_allocation
```

Low-margin / negative-margin orgs surface in a dashboard; sales+engineering decide whether to adjust plan, price, or optimisation.

---

## 5. Levers — By Category

### 5.1 Compute (biggest knob)

**Right-size requests/limits**
- VPA in Recommender mode surfaces over-provisioned pods.
- Quarterly tightening campaign; target < 20 % headroom on stable services.

**Spot / Preemptible**
- Stateless services on spot via Karpenter (expected interruption < 5 %).
- Kafka / Postgres / GPU NEVER on spot (durability concerns).
- Expected savings: **30-60 %** on affected pool; 15-25 % on total compute.

**ARM64 / Graviton**
- Go, Kotlin/JVM, Rust all multi-arch.
- Graviton saves ~20 % at same performance; most of our services moved.

**Scale-down windows**
- Non-prod environments scaled to zero on nights/weekends (`helixctl env sleep`).
- Saved ~$20k/mo on staging/dev.

**Right-size nodes**
- Larger instances pack better for memory-heavy services (Postgres, Kafka).
- Smaller for stateless bursty workloads.

### 5.2 Storage

**Postgres**
- Tier cold partitions (audit > 90 d) to cheaper storage class if supported (Aurora I/O-Optimized vs. Standard; Babelfish for heavy analytics).
- Compact / VACUUM FULL periodically on low-churn tables.
- Review index bloat quarterly.

**Kafka**
- Tiered storage (MSK Serverless alt. custom with S3 tier on Strimzi) — hot window on EBS, cold on S3.
- Evaluate retention per topic; many defaults can be shortened.
- Compression: `zstd` on high-throughput topics; ~40 % size reduction.

**OpenSearch / Meilisearch / Qdrant**
- ILM on OpenSearch: hot → warm (cheaper nodes) → cold (searchable snapshots).
- Qdrant: scalar quantisation if vector count > 10 M (4× compression).
- Meilisearch: routine purge of stale project indexes.

**Object store**
- S3 Intelligent-Tiering for LFS objects rarely accessed.
- Lifecycle rules: release-asset older versions → IA → Glacier; expire after customer retention + legal hold window.
- Versioning wisely used — not unlimited.

### 5.3 Network

**Cross-AZ traffic**
- Kafka producer co-location with consumer AZ when possible (`rack.id`).
- Istio Ambient minimises sidecar overhead.
- Pod topology-aware routing.

**Egress**
- NAT gateway per AZ cost add up; shared NAT only in non-prod.
- VPC endpoints for AWS APIs to avoid NAT traversal.
- Cloudflare caches static + OpenAPI (tiny but meaningful).

**Image pulls**
- Registry mirror in-cluster; images pulled once per node.

### 5.4 GPU / Inference

**Right-size models**
- 8B LoRA for most tasks; only conflict-resolution-complex uses 70B.
- Per-task model choice in LiteLLM router.

**Batching**
- vLLM continuous batching; target ~70 % GPU utilisation.
- Cold GPU pods pre-warmed during peak hours only.

**Spot / MIG**
- Inference pool tolerant to interruption (queue-replayed).
- MIG partitioning when full GPU not needed.

**Org-level budgets**
- Token budgets per org per plan prevent runaway spend.
- AI on HelixGitpx itself budgeted; overshoots investigated.

### 5.5 Observability

**Retention**
- Metrics: 13 months raw would be expensive — keep 30 d at 15 s, 1 y at 5 m (downsampled in Mimir).
- Traces: tail-sampling keeps 1 % of OK, 100 % of error/slow.
- Logs: 30 d hot in Loki, 90 d S3 cold.

**Cardinality**
- Cardinality explosion is the #1 cost driver.
- CI reviews metric PRs for high-cardinality labels.
- Regex relabelling in OTel collector to strip useless labels.

### 5.6 Data transfer to sub-processors

- Analytics pipelines exported as batched parquet, not streaming.
- Webhook delivery batched to external sinks.

---

## 6. Cost Budgets

Every quarter the FinOps partner sets budgets:

- Per-cluster per-environment monthly cap.
- Per-team chargeback targets.
- Shared-service fair-use thresholds.

Alerts:

- **Cost anomaly** (forecasted > budget by 10 %) → Slack notification.
- **Sudden spend spike** (> 20 % above rolling daily average) → page FinOps + SRE lead.

---

## 7. Chargeback / Showback

- **Internal teams**: showback model — teams see their cost but aren't billed.
- **Per-customer**: always full attribution; enterprise customers can see their own infra cost in a hidden dashboard for transparency.

Cost allocation label discipline is enforced at admission (Kyverno policy `helixgitpx-require-standard-labels`).

---

## 8. Capacity Planning

- Quarterly capacity review: headroom for 2× current peak.
- Long-lead items (Reserved Instances / Savings Plans / committed use contracts) renewed annually; break-even analysis documented.
- Strategic conversations about multi-cloud / bare-metal cost curve happen yearly.

---

## 9. Savings Plans / Reserved Capacity

- Compute Savings Plans (AWS) for baseline: covered hours 70 %.
- EC2 Reserved for Postgres / Kafka (durable workloads).
- GPU capacity commitments sized to p50 demand; burst on-demand.

---

## 10. Anti-Patterns to Avoid

- **Over-aggressive downscaling** during rolling deploys → SLO breach during scale-up.
- **Cold-starting GPU pods on user request** — ruins p99 inference latency.
- **Cheapest node class for PG/Kafka** — IOPS become the bottleneck; total cost higher due to retries.
- **Over-retention of observability data** — grows fast; rarely read.
- **No cost visibility per service owner** — optimisation never happens.

---

## 11. Quick Wins Checklist (repeat quarterly)

- [ ] Any pod with CPU < 10 % steady-state? Right-size.
- [ ] Any topic with retention > actual consumer lag by 10×? Shorten.
- [ ] Any table partition > 6 months untouched? Archive.
- [ ] Any Cloudwatch log group > retention without use? Shorten.
- [ ] Any S3 prefix without lifecycle rule? Add.
- [ ] Any GPU with < 40 % avg utilisation? Consolidate.
- [ ] Any dev/staging env running 24/7 without need? Schedule shutdowns.
- [ ] Any over-replicated service with capacity headroom > 3×? Reduce minReplicas.
- [ ] Any unused feature flag still emitting metrics? Retire.

---

## 12. Benchmarks (cost/performance)

Tracked in our FinOps dashboard:

| Unit | Target | Current (example) |
|---|---|---|
| Cost per 1000 API requests | $0.004 | $0.0037 |
| Cost per 1 GB pushed | $0.05 | $0.048 |
| Cost per 1000 events delivered | $0.002 | $0.0019 |
| Cost per 1000 AI tokens served | $0.08 | $0.062 |
| Cost per active repo / month | $0.35 | $0.34 |

Regressions trigger investigation.

---

## 13. Tooling

- `helixctl cost show --service=<n>` — quick view.
- `helixctl cost breakdown --org=<slug>` — per-tenant.
- `helixctl cost savings-opportunities` — machine-generated candidates (right-sizing, etc.).
- GitHub Action on PR: flags if expected cost of the change > $X/month without justification.

---

## 14. FinOps Governance

- Monthly review with engineering leads.
- Quarterly review with finance.
- Annual cost architecture review (multi-year capex / opex balance).
- Every major architectural proposal must include a cost model (ADR template has a section for this).

---

*— End of Cost Optimisation Guide —*
