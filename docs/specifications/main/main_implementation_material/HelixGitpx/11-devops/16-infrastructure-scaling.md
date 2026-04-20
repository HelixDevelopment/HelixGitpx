# 16 — Infrastructure & Scaling

> **Document purpose**: Describe how HelixGitpx runs: **Kubernetes topology, multi-region deployment, autoscaling, storage, networking, cost envelopes, and failure domains**.

---

## 1. Deployment Environments

| Env | Purpose | Scale |
|---|---|---|
| `dev` | Per-engineer local cluster (Kind / minikube / k3d) | single node |
| `ci` | Ephemeral per-PR namespace on shared cluster | tiny |
| `staging` | Full parity with prod, reduced scale | 3 nodes per plane |
| `preprod` | Canary / load testing | prod-like |
| `prod-<region>` | Real customer traffic | see §3 |

GitOps repo: `helixgitpx-platform`. Every environment is a folder. Argo CD reconciles on commit.

---

## 2. Kubernetes Baseline

| Component | Choice |
|---|---|
| Distribution | Vanilla **Kubernetes 1.31+** on bare metal (MAAS/Talos) or EKS/GKE/AKS |
| CNI | **Cilium** (eBPF, L7 policies) |
| Service Mesh | **Istio Ambient** |
| Ingress | **Envoy Gateway** (Gateway API v1) + Cloudflare in front |
| DNS | CoreDNS + ExternalDNS |
| Cert | cert-manager + Let's Encrypt for public; internal mTLS via SPIRE |
| Storage | **Rook-Ceph** (bare metal) or **EBS/Premium SSD** (cloud) |
| Object storage | **MinIO** cluster (bare metal) or S3/R2 |
| Secrets | **HashiCorp Vault** + CSI driver |
| Admission | **Kyverno** (policy) + **Gatekeeper** (ABAC fallback) |
| Supply chain | **Connaisseur** or **Cosigned** verifying Cosign signatures |

### 2.1 Node Pools

| Pool | Purpose | Node type (example) |
|---|---|---|
| `system` | kube-system, observability | 4 vCPU / 16 GB × 3 |
| `app` | Stateless services | 16 vCPU / 64 GB × auto |
| `memory` | Redis, Dragonfly | 8 vCPU / 64 GB × 3 |
| `data` | Postgres, Kafka, OpenSearch, Qdrant, Meilisearch | 16 vCPU / 128 GB / NVMe × 3+ |
| `gpu` | LLM inference, training | 8 vCPU / 64 GB / 1× A100 or L4 |
| `build` | Ephemeral CI runners | spot / burstable |

---

## 3. Multi-Region Topology

- **Regions at GA**: `eu-west-1` (Frankfurt), `us-east-1` (N. Virginia). EU-only region for data-residency customers.
- **Active-Active** for stateless services; **Active-Passive** for Postgres primary.
- **Routing**: GeoDNS + Anycast via Cloudflare / Route53.
- **Data plane**:
  - Kafka MirrorMaker 2 mirrors per-topic with region prefix.
  - Postgres streaming replication + logical replication for tenant-scoped subsets.
  - Object storage bi-directional replication (S3 CRR / MinIO bucket replication).
- **Failover**: planned = minutes; unplanned = automatic with 30 s RPO budget.

### 3.1 Region Pinning

Each org has a `primary_region` column. Data-residency customers can pin strictly (no cross-region replication). Routing rules in API gateway ensure the correct region serves every request.

---

## 4. Service-Level Scaling

### 4.1 HorizontalPodAutoscaler

- Every `Deployment` has an HPA on the relevant signal.
- For request-driven services: custom metric `grpc_server_inflight_requests` or `http_requests_in_flight` (from the Prometheus adapter).
- For event-driven services: **KEDA** on Kafka consumer lag.
- For LLM inference: KEDA on queue depth + GPU utilisation.

### 4.2 VerticalPodAutoscaler (advisory)

- VPA in "Recommender" mode — does not evict, just advises limits during load tests.

### 4.3 Cluster Autoscaler / Karpenter

- **Karpenter** (cloud) for rapid node provisioning on appropriate instance types.
- On bare metal: maintain warm pools; scale via MAAS/Talos.

### 4.4 SLO-Backed Scaling

Scale targets tied to SLO budgets:
- If p99 latency is within budget, do not scale up.
- If error-budget burn > 2× hourly budget for 15 min, scale up and page on-call.

---

## 5. Capacity Model

### 5.1 GA Scale Targets

| Metric | Target |
|---|---|
| Organisations | 100 000 |
| Repos | 10 000 000 |
| Events/day | 1 000 000 000 |
| Live-event subscribers concurrent | 500 000 |
| Git push events/s peak | 10 000 |
| Webhook deliveries/s sustained | 30 000 |
| Search queries/s peak | 5 000 |
| LLM tokens/s | 1 000 000 |

### 5.2 Capacity per Service (sized for 10× headroom)

| Service | Replicas prod | CPU req / limit | Mem req / limit |
|---|---|---|---|
| api-gateway | 20 | 1 / 2 | 1 GiB / 2 GiB |
| auth-service | 10 | 500m / 1 | 512 / 1024 |
| repo-service | 20 | 1 / 2 | 1 / 2 |
| sync-orchestrator | 20 | 1 / 2 | 1 / 2 |
| conflict-resolver | 15 | 2 / 4 | 2 / 4 |
| adapter-pool | 40 | 1 / 2 | 1 / 2 |
| live-events-service | 30 | 1 / 2 | 2 / 4 |
| ai-service | 10 (+ GPU pool) | 2 / 4 | 4 / 8 |
| search-projectors × N | 5 each | 500m / 1 | 1 / 2 |
| billing | 3 | 500m / 1 | 512 / 1024 |
| audit-service | 5 | 500m / 1 | 512 / 1024 |

---

## 6. Storage Sizing

| Store | Initial | 1-year estimate | Replication |
|---|---|---|---|
| Postgres | 2 TiB | 20 TiB | Primary + 2 replicas, streaming + cross-region |
| Kafka | 5 TiB local + S3 tier | 25 TiB hot / 500 TiB cold | RF=3, tiered |
| OpenSearch | 2 TiB hot / 10 TiB warm / cold searchable | 50 TiB | RF=2 + snapshot |
| Meilisearch | 200 GiB | 1 TiB | read replicas |
| Qdrant | 500 GiB | 5 TiB | RF=2 |
| Object store (Git packs, LFS, assets) | 10 TiB | 200 TiB | RF=3 / erasure 8+3 |
| Redis | 50 GiB | 500 GiB | RF=2 + AOF |

---

## 7. Network

### 7.1 Segmentation

- Public zone: Cloudflare edge + Envoy Gateway (TLS termination).
- Service zone: Istio mesh, mTLS-only, default-deny NetworkPolicies.
- Data zone: Postgres/Kafka/OpenSearch only reachable from service zone; no external IPs.
- Egress zone: explicit NAT gateway; all outbound Git upstream traffic labelled and audited.

### 7.2 Rate Limits / DDoS

- L4 at Cloudflare; L7 via Envoy Gateway with token-bucket extensions.
- Per-token / per-IP limits (see [07-rest-api.md §10](../04-apis/07-rest-api.md)).
- Selective tarpit for abusive IPs.

---

## 8. Failure Domains

| Failure | Blast radius | Mitigation |
|---|---|---|
| Pod crash | 1 replica | K8s restart; others serve |
| Node failure | pods on node | PDB ≥ 2; topology spread |
| AZ failure | replicas in AZ | multi-AZ deployment, quorum preserved |
| Region failure | region | GeoDNS cutover; secondary promoted |
| Postgres primary | DB writes | PITR replica promoted; app reconnects |
| Kafka broker | partition | RF=3, min-ISR=2 |
| Service mesh control plane | configuration updates | ambient data plane keeps serving |

PodDisruptionBudgets: `maxUnavailable: 1` for replicas ≥ 3.

---

## 9. Backups & DR

- **Postgres**: PITR base backup nightly + continuous WAL to object storage; restore tested monthly.
- **Kafka**: topic tiered storage + nightly metadata snapshot to Git.
- **OpenSearch**: snapshot repository on S3, daily.
- **Qdrant / Meilisearch**: volume snapshot + logical dump.
- **Object storage**: versioned + Object Lock for compliance buckets.
- **Vault**: Raft backup daily, keys split per Shamir.
- **DR drills**: full region failover rehearsed quarterly.

### 9.1 RPO / RTO Targets

| Component | RPO | RTO |
|---|---|---|
| Postgres transactional | ≤ 30 s | ≤ 15 min |
| Kafka | ≤ 30 s | ≤ 5 min |
| Object store | ≤ 5 min | ≤ 15 min |
| OpenSearch (rebuildable from Kafka) | tolerant | ≤ 60 min |
| Qdrant / Meilisearch (rebuildable) | tolerant | ≤ 60 min |

---

## 10. Observability Infrastructure

- **Metrics**: Prometheus (Mimir or Thanos for long-term + HA).
- **Logs**: Vector → Loki; mirror selected to OpenSearch for long-term compliance.
- **Traces**: OTel Collector → Tempo (or Jaeger).
- **Alerting**: Alertmanager → PagerDuty/Opsgenie; runbooks linked per alert.
- **Dashboards**: Grafana (SLO, capacity, per-service RED dashboards).

See [18-observability.md](../09-observability/18-observability.md).

---

## 11. Cost Controls

- **Resource quotas** per namespace.
- **Per-tenant cost allocation** via Kubecost + OpenTelemetry usage events → billing service.
- **Scheduled scale-down** on non-production environments overnight.
- **Spot/preemptible** for stateless workloads where safe (not for Kafka/Postgres).
- **GPU partitioning** (MIG) for LLM inference efficiency.
- Monthly cost review against budget per component.

---

## 12. Edge & CDN

- Cloudflare (or Fastly) in front for static assets, WebSocket proxying, DDoS, WAF.
- API cached only for truly public read-only endpoints (OpenAPI spec, public repo listings, favicons).
- Service-worker offline for web.

---

## 13. Hardware Assumptions (for bare-metal deployment)

- Each data-plane node: 2× AMD EPYC 9554 or Intel Sapphire Rapids, 1 TB RAM, 8× NVMe, 2× 25 GbE, IPMI.
- GPU node: 2× A100 80GB or L40S ×4 per box, 2 TB NVMe cache.
- Rack-level: dual power, dual ToR switches, BGP-ECMP.
- 3 racks minimum per region for AZ analogues.

---

## 14. Runbooks

- Add / remove node.
- Add / remove broker / partition.
- Scale Postgres replica.
- Rolling K8s upgrade (Talos node-by-node).
- Certificate rotation (automated; manual trigger documented).

See [19-operations-runbook.md](../12-operations/19-operations-runbook.md).

---

*— End of Infrastructure & Scaling —*
