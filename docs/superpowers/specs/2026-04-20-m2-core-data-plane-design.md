# M2 Core Data Plane — Design Spec

| Field | Value |
|---|---|
| Status | APPROVED (pending user review) |
| Author | Милош Васић + Claude (brainstorming session 2026-04-20) |
| Milestone | M2 — Core Data Plane (Weeks 5–10 in `13-roadmap/17-milestones.md`) |
| Scope | Full 20-item roadmap §3 (items 19–38); nothing skipped |
| Sequencing | Approach 2 — observability-first; per-service Grafana gates |
| Supersedes | — |
| Implements | `docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md` §3 |

---

## 1. Context

M1 Foundation landed a monorepo, 14 shared Go packages, protobuf codegen, a scaffolding tool, a working hello service (code + integration test compiles clean; runtime verification deferred to M2 due to local podman socket latency), full CI catalog (`workflow_dispatch`-only), Kyverno/Checkov/ARC/Vault config-only artifacts, Nx/KMP/Docusaurus scaffolds, ADRs 0001–0005, runbook template, and a completion-matrix verifier (27/27 green). Tag `m1-foundation` at commit `70ab700`.

M2 builds the Core Data Plane: the Kubernetes cluster, the data-at-rest and data-in-motion infrastructure, service identity, and the full observability stack. Unlike M1 (file artifacts), M2 is infrastructure — its exit criterion requires a running cluster with 20 components deployed, hello reachable via ingress, and the full metrics/logs/traces/profiles pipeline visible in Grafana.

## 2. Goals

G1. Provision a local k3d cluster bootstrapped via Argo CD app-of-apps, with every M2 component managed as an `Application` CR (GitOps from the first install).

G2. Install the full Phase 2.1–2.5 set (20 numbered items) on the cluster, in observability-first order (Approach 2).

G3. Upgrade hello to the transactional outbox pattern. Verify Debezium captures the outbox table and emits to the same `hello.said` topic M1 used.

G4. Provide a staged-HA pattern: HA manifests are authoritative in Git (`staging`/`prod-eu` overlays); local overlay patches to single-replica for machine capacity.

G5. Ship full observability: Prometheus + Mimir + Loki + Tempo + Pyroscope + Grafana + Alertmanager, with Grafana dashboards provisioned from Git, PagerDuty integration placeholder-configured, and per-chart "Down" / "HighErrorRate" alert rules baseline.

G6. Prove the spine end-to-end on the cluster: curl → hello Ingress → gRPC → pg UPSERT + outbox write → Debezium → `hello.said` → Kafka consumer; metrics/logs/traces/profiles show in Grafana.

G7. Enable Istio Ambient mesh in application and data namespaces; verify zero-code mTLS between hello and its deps.

## 3. Non-goals

- Actual multi-region or DR (M8).
- Real PagerDuty, Let's Encrypt, external DNS provider accounts (manifests ship; activation is deferred).
- HA enforcement on the local machine — HA is expressed in `staging`/`prod-eu` overlays only; `local` is single-replica by design.
- Any service beyond hello. New services land in M3 (auth/org/team), M4 (repo/git-ingress/adapter-pool/webhook-gateway/upstream), M5 (sync-orchestrator/conflict-resolver/live-events).
- macOS/Windows developer support — M2 is Linux-first per `11-devops/26-on-prem-deployment.md`.

## 4. Locked constraints from brainstorming

| ID | Constraint | Source |
|---|---|---|
| C-1 (new) | K8s target: local k3d (`up.sh --m2`) deploying for real | Q1 |
| C-2 (new) | HA manifests authoritative in Git; `kustomize/overlays/local/` patches to single-replica | Q2 |
| C-3 (new) | Hello service upgraded to transactional outbox via Debezium; `hello.said` topic unchanged | Q3 |
| C-4 (new) | Istio Ambient included (zTunnel) in M2 | Q4 |
| C-5 (new) | Full 7-component observability (Prom + Mimir + Loki + Tempo + Pyroscope + Grafana + AM) | Q5 |
| C-6 (new) | Local-friendly externals: self-signed cert-manager, noop external-dns, MinIO S3, placeholder PagerDuty | Q6 |
| C-7 (new) | Observability-first sequencing (Approach 2); every data-service "done" gate requires Grafana visibility | Approaches |
| Inherited from M1 | GitHub Actions `workflow_dispatch`-only; portable compose wrapper; single git history; `mise.toml`; spine-first with completion matrix | M1 spec §4 |

## 5. Cluster bootstrap & GitOps root

**Cluster:** k3d cluster `helix` (1 server + 2 agents), Cilium CNI (replacing k3s's default flannel), ingress-nginx, 8 namespaces pre-created. Spin-up via `impl/helixgitpx-platform/k8s-local/up.sh --m2`:

1. Preflight: ensure ≥ 48 GiB free RAM (abort if not).
2. Pre-pull heavy images (postgres, kafka, opensearch, vault, mimir, loki, tempo).
3. `k3d cluster create --k3s-arg '--flannel-backend=none@server:*' --k3s-arg '--disable-network-policy@server:*'`.
4. Install Cilium via Helm (replaces flannel).
5. Install ingress-nginx via Helm.
6. Create namespaces: `helix-system`, `helix-identity`, `istio-system`, `helix-data`, `helix-cache`, `helix-secrets`, `helix-observability`, `helix`.
7. Install Argo CD from `argocd/bootstrap/kustomization.yaml`.
8. Apply `argocd/applicationset/app-of-apps.yaml` — root Application that reconciles everything under `argocd/applications/`.
9. Write `/etc/hosts` entries for `hello.helix.local`, `grafana.helix.local`, `vault.helix.local`, `argocd.helix.local`, `kafka.helix.local`.

**GitOps tree:**

```
impl/helixgitpx-platform/
├── argocd/
│   ├── bootstrap/          Argo CD install (kustomization)
│   ├── applicationset/
│   │   └── app-of-apps.yaml
│   └── applications/       25 Application CRs (one per helm/ child)
├── helm/
│   ├── cilium/
│   ├── ingress-nginx/
│   ├── cert-manager/
│   ├── external-dns/
│   ├── minio/
│   ├── prometheus-stack/   Prom + Grafana + AM (+dashboards, +datasources)
│   ├── mimir/
│   ├── loki/
│   ├── tempo/
│   ├── pyroscope/
│   ├── spire/
│   ├── istio-base/
│   ├── istio-ambient/      istiod + ztunnel
│   ├── cnpg-operator/
│   ├── cnpg-cluster/       the actual Postgres HA Cluster CR
│   ├── strimzi-operator/
│   ├── kafka-cluster/      Kafka + KafkaTopic + KafkaConnect + KafkaConnector
│   ├── karapace/
│   ├── debezium/
│   ├── dragonfly/
│   ├── meilisearch/
│   ├── opensearch/
│   ├── qdrant/
│   └── vault/
├── kustomize/
│   └── overlays/
│       ├── local/          replicas:1, resource requests halved, self-signed issuer
│       ├── staging/        full HA, Let's Encrypt DNS-01, real external-dns
│       └── prod-eu/        staging + stricter limits + data-residency
└── sql/
    └── schemas.sql         per-service schema + RLS baseline
```

**Sync waves** via `argocd.argoproj.io/sync-wave` annotation:

| Wave | Components |
|---|---|
| -10 | cilium |
| -5 | cert-manager, ingress-nginx, external-dns, minio |
| -3 | prometheus-stack, mimir, loki, tempo, pyroscope |
| 0 | spire, istio-base, istio-ambient |
| 5 | cnpg-operator, strimzi-operator |
| 7 | cnpg-cluster, kafka-cluster, dragonfly, meilisearch, opensearch, qdrant, vault |
| 9 | karapace, debezium |
| 10 | hello |

Observability installs before any data service so the data services' rollouts are visible immediately.

## 6. Observability stack

- **kube-prometheus-stack**: Prometheus, Alertmanager, Grafana, node-exporter, kube-state-metrics. Grafana dashboards and datasources provisioned from `helm/prometheus-stack/grafana/{dashboards,datasources}/` (ConfigMaps).
- **Mimir** (`grafana/mimir-distributed`): remote-write target for Prometheus. Local = single-binary, 7-day retention. Staging = microservices mode, 90-day retention.
- **Loki** (`grafana/loki`): single-binary locally; microservices in staging. Promtail DaemonSet ships pod logs.
- **Tempo** (`grafana/tempo`): OTLP ingest on `:4317`. `OTEL_EXPORTER_OTLP_ENDPOINT` env var for every service points here.
- **Pyroscope** (`grafana/pyroscope`): scrapes Go `/debug/pprof/*`. `platform/telemetry` grows a `RegisterPprof(mux *http.ServeMux)` helper.
- **Alertmanager**: routes to PagerDuty integration key fetched from `kv/alerting/pagerduty/integration-key` (Vault in staging; placeholder Secret locally). Local routes to a null receiver if the secret is absent.

Every Helm chart produced by M2 ships two baseline alert rules in `helm/<chart>/templates/alerts.yaml`: `<Chart>Down` (replicas = 0 ≥ 2 min) and `<Chart>HighErrorRate` (5xx > 5 % over 5 min).

## 7. Cert-manager, external-dns, ingress-nginx, MinIO

- **cert-manager**: two `ClusterIssuer` resources — `selfsigned-ca` (local) and `letsencrypt-dns01` (staging/prod-eu). Overlay selects the active one.
- **external-dns**: `noop` webhook provider locally; `cloudflare`/`route53` in staging (provider creds in Vault).
- **ingress-nginx**: serves L7. Hostnames (local): `hello.helix.local`, `grafana.helix.local`, `argocd.helix.local`, `vault.helix.local`, `kafka.helix.local`.
- **MinIO**: single-instance locally; 4-node distributed in staging. Bucket provisioning via a post-install Job: `cnpg-backups`, `opensearch-snapshots`, `loki-chunks`, `tempo-traces`, `pyroscope-profiles`, `mimir-blocks`. Root credentials in Vault; each consumer has its own IAM-equivalent policy.

## 8. SPIRE + Istio Ambient

**SPIRE** (`spiffe/spire` chart, server + agent DaemonSet):
- Trust domain `helixgitpx.local` / `helixgitpx.dev`.
- K8s PSAT attestor; workload entries auto-created via `spire-k8s-workload-registrar` (one SPIFFE ID per ServiceAccount).
- `platform/spire.NewFetcher` lifts the M1 `TODO(M2)` and wires `go-spiffe/v2/workloadapi`.

**Istio Ambient** (3 Applications: `istio-base` CRDs + operator, `istiod` control plane, `ztunnel` data plane + Istio CNI):
- Namespace labels: `helix`, `helix-data`, `helix-cache`, `helix-secrets` are labeled `istio.io/dataplane-mode: ambient`. Observability and system namespaces opt out.
- Istio consumes SPIRE SVIDs via SPIFFE CSRA integration (`meshConfig.caProvider: spiffe://helixgitpx.local`).
- Verification: `tcpdump` on a zTunnel pod shows mTLS for hello→pg connections (no plaintext).

## 9. Postgres via CNPG + per-service schemas + RLS + migrations

- **CNPG operator**, namespace-scoped to `helix-data`.
- **`cnpg-cluster`** Helm chart produces `Cluster` CR `helix-pg`:
  - Local: 1 instance, 10 Gi PVC, Barman backups → MinIO `cnpg-backups/`.
  - Staging: 3 instances + streaming replication, 100 Gi, same backup target.
  - `postgresql.parameters.wal_level=logical` (for Debezium).
- **Per-service schemas:** `impl/helixgitpx-platform/sql/schemas.sql` declares `hello`, `auth`, `repo`, `sync`, `conflict`, `upstream`, `collab`, `events`, `platform` and one `<name>_svc` role per schema. Applied via post-install Job.
- **RLS baseline:** every application table gets `ENABLE ROW LEVEL SECURITY` with a `USING (true)` policy initially. Later milestones tighten.
- **Goose migrations:** `platform/pg/migrate.go` wraps `goose/v3`; hello's chart gets a `migrate-job.yaml` pre-upgrade Hook. Services get a `migrate` subcommand on their main binary.
- **PITR drill:** `cnpg-cluster-restore-drill` CronJob in staging only (monthly). Local is disabled.

## 10. Kafka + Karapace + Debezium + outbox

- **Strimzi operator** (CRDs `Kafka`, `KafkaTopic`, `KafkaConnect`, `KafkaConnector`).
- **`kafka-cluster`** chart: 1-broker KRaft local; 3-broker RF=3 min-ISR=2 staging. Topics: `hello.said` (7d), `hello.outbox` (30d), `helixgitpx._internal.*` (reserved).
- **Karapace**: BSR-compatible REST on `:8081`. `platform/kafka.ResolveFn` hook wires to Karapace (lifts M1 `TODO(M2)`).
- **Debezium Connect cluster**: `KafkaConnect` CR + a `KafkaConnector` named `hello-outbox` — source `helix.pg:5432`, schema `hello`, table `outbox_events`, sink topic `hello.said`.
- **Outbox in hello:** replace `internal/repo/event_kafka.go` with `internal/repo/event_outbox.go` — INSERTs to `hello.outbox_events(id uuid pk, name, greeting, count, at, sent_at)` inside the same transaction as the counter UPSERT. Remove the direct Kafka producer. External behavior unchanged.

## 11. Caches + search: Dragonfly, Meilisearch, OpenSearch, Qdrant

- **Dragonfly** replaces Redis. Wire-compatible with go-redis v9; `platform/redis` unchanged.
- **Meilisearch**: 1 instance local, 3 instances + master in staging. Primary key from Vault.
- **OpenSearch**: 1 node local, 3 nodes staging. MinIO snapshot repo `opensearch-snapshots/`. ILM policy `logs-7d` ships as ConfigMap.
- **Qdrant**: 1 node local, distributed in staging.

Each chart ships a `ServiceMonitor` (Prom CRD) and Grafana dashboards from Git.

## 12. Vault (Raft + Shamir)

- Vault StatefulSet, Raft storage.
- Local: 1 replica, Shamir 3-key/2-threshold; unseal keys in `vault-unseal-keys` Secret (unsafe, local only).
- Staging: 3 replicas, auto-unseal via cloud KMS (provider-dependent; GCP KMS default).
- `vault-bootstrap` Job initialises on first install.
- Policies added: `pg-backup-writer`, `opensearch-backup-writer`, `loki-chunk-writer`, `tempo-trace-writer`, `mimir-block-writer`, `pyroscope-profile-writer`, `hello`.
- K8s auth method, ServiceAccount-to-policy bindings. Vault Agent Injector for sidecar-based secret delivery.

## 13. Hello service updates

- **Outbox pattern** (see §10) — single code change with the largest semantic impact.
- **SVID-terminated gRPC**: `platform/grpc.NewServer` accepts `*spire.Fetcher`; when present, terminates mTLS with the workload SVID. Tests add a fake SVID source.
- **OTLP export wired**: `OTEL_EXPORTER_OTLP_ENDPOINT` → Tempo's OTLP gRPC receiver.
- **Pprof handlers**: `platform/telemetry.RegisterPprof(mux)` registers `/debug/pprof/*` on the health mux. Hello's health server calls it.
- **ServiceMonitor** + PodMonitor in the helm chart so Prom picks metrics up automatically.
- **Config from Vault**: `platform/config` grows a `vault:"path/to/key"` tag; hello's DSNs now come from Vault KV (falling back to env vars for local convenience).
- **Argo CD-managed deploy**: `argocd/applications/hello.yaml` points at `impl/helixgitpx/services/hello/deploy/helm/`. Chart adds migrate-job hook + Debezium KafkaConnector CR.

## 14. Error handling, testing, completion matrix

**Cluster-level error handling:**
- PodDisruptionBudgets on every chart (`minAvailable: replicas-1`).
- Argo CD automated self-heal + prune on every Application.
- Health CRD scripts (Lua) for Strimzi `Kafka`, CNPG `Cluster`, Vault `StatefulSet`, `KafkaConnector`.
- `SIGTERM` → context cancel in every Go main.

**Testing:**

| Layer | Tool | M2 ships |
|---|---|---|
| Helm unit | `helm unittest` | `helm/<chart>/tests/` per chart |
| Kustomize | `kubectl kustomize` + diff | overlay-vs-base comparisons |
| Policy | Kyverno + Checkov (unchanged from M1) | coverage now spans all M2 manifests |
| Cluster integration | `scripts/verify-m2-cluster.sh` | polls Argo CD for 25 Applications Synced+Healthy |
| Spine integration | `scripts/verify-m2-spine.sh` | curl via Ingress + grpcurl + outbox→topic in < 2s + trace in Tempo + logs in Loki + metrics in Prom |
| SPIRE | `kubectl exec` hello pod + `spire-agent api fetch` | returns SVID |
| Istio Ambient | tcpdump via ztunnel | shows mTLS on hello→pg path |

**Completion matrix** — 20 roadmap items, each with artifact + gate. See §15 Exit criteria for the explicit list; also embedded verbatim in the M2 plan document as the verify-m2-cluster.sh script contents.

## 15. Exit criteria (explicit, all 20 items + end-to-end)

| # | Item | Gate |
|---|---|---|
| 19 | Staging cluster (k3d local) | `kubectl get nodes` shows 3 Ready; Cilium DaemonSet Running |
| 20 | Argo CD reconciling | `argocd app list` shows 25 Applications; all Synced + Healthy |
| 21 | cert-manager + external-dns | TLS cert issued for `hello.helix.local`; external-dns log shows noop webhook calls |
| 22 | SPIRE + SVID fetch | `kubectl exec hello-<pod> -- spire-agent api fetch` returns an SVID |
| 23 | CNPG HA | `kubectl get cluster -n helix-data helix-pg` Ready; 1 replica local / 3 staging |
| 24 | Schemas + RLS | `psql -c '\dn'` shows 9 schemas; RLS enabled on every `hello.*` table |
| 25 | Goose migrations | Hello's migrate-job pod completes with exit 0; `hello.greetings` table exists |
| 26 | PITR to object store | `mc ls minio/cnpg-backups/` shows WAL files within 5 min of Cluster Ready |
| 27 | Strimzi Kafka | `kubectl get kafka -n helix-data` Ready; 3 topics exist |
| 28 | Karapace | `curl http://karapace:8081/subjects` returns 200 + JSON |
| 29 | Debezium | `KafkaConnector hello-outbox` Running; ≥ 1 event delivered |
| 30 | Outbox pattern | INSERT into `hello.outbox_events` → event on `hello.said` within 2 s |
| 31 | Dragonfly | `dragonfly-cli ping` = PONG; hello cache hit works |
| 32 | Meilisearch | `curl :7700/health` returns 200 |
| 33 | OpenSearch ILM + snapshots | `_cat/health` green; `_snapshot/_status` shows completed snapshot |
| 34 | Qdrant | `/readyz` returns 200 |
| 35 | Vault | `vault status` Initialized + Unsealed |
| 36 | Prom/Mimir/Loki/Tempo/Pyroscope | all five Applications Synced + Healthy; Grafana datasources connect |
| 37 | Grafana provisioning from Git | `grafana-cli` lists ≥ 10 dashboards, all marked "from disk" |
| 38 | Alertmanager → PagerDuty (placeholder) | AM config renders; `curl -s http://alertmanager:9093/api/v2/status` shows one valid route to PagerDuty receiver |

**End-to-end spine:**

- `curl https://hello.helix.local/v1/hello?name=world` (via ingress-nginx + self-signed cert) → `{"greeting":"hello, world","count":1}`.
- `grpcurl hello.helix.local:443 helixgitpx.hello.v1.HelloService/SayHello` → same payload.
- Kafka consumer on `hello.said` sees the event emitted via outbox within 2 s.
- Grafana dashboard for hello shows: latency histogram, log line from Loki, sampled trace with pg + redis + kafka spans, CPU flamegraph in Pyroscope.

## 16. Risks & mitigations

| Risk | Mitigation |
|---|---|
| k3d single-host OOM with full stack | `local` overlay halves requests/limits; `up.sh --m2` refuses to start below 48 GiB free RAM |
| Heavy image pulls time out | `up.sh --m2` pre-pulls postgres, kafka, opensearch, vault, mimir, loki, tempo before cluster creation |
| Cilium + k3d flannel coexistence | `up.sh --m2` disables flannel and network-policy in k3d; Cilium installed before any workload schedules |
| Debezium needs `wal_level=logical` | CNPG cluster values set it explicitly |
| Istio Ambient + Cilium CNI | Native-pod-CNI mode for both; documented in `argocd/applications/istio-ambient.yaml` |
| SPIRE datastore ephemeral locally | SQLite on SPIRE server pod; acceptable — re-registration on restart is fast |
| PagerDuty key absent | Alertmanager routes to null receiver when the integration-key Secret doesn't exist |
| macOS dev host | Explicitly out of scope; documented in `up.sh --m2`'s README |
| Vault unseal key leak (local Secret) | Only local unseal keys are in-cluster; staging/prod use cloud KMS auto-unseal |
| Argo CD loop on CRD-before-instance ordering | Sync waves enforce operators (wave 5) land before their CR instances (wave 7) |

## 17. Open questions

None. All blocking decisions are locked in §4.

## 18. References

- Roadmap: `docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md` §3 (items 19–38)
- On-prem-deployment baseline: `docs/specifications/.../11-devops/26-on-prem-deployment.md`
- Infra scaling expectations: `docs/specifications/.../11-devops/16-infrastructure-scaling.md`
- Observability spec: `docs/specifications/.../09-observability/18-observability.md`
- Data model + schemas: `docs/specifications/.../03-data/`, `docs/specifications/.../16-schemas/*.sql`
- Existing manifests (reference): `docs/specifications/.../18-manifests/`
- M1 foundation: `docs/superpowers/specs/2026-04-20-m1-foundation-design.md`, tag `m1-foundation`

— End of M2 Core Data Plane design —
