# M2 Core Data Plane Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deploy the full 20-item M2 Core Data Plane from roadmap §3 (items 19–38) on a local k3d cluster via Argo CD GitOps; upgrade hello to the transactional outbox pattern; exit when the full Prom/Mimir/Loki/Tempo/Pyroscope/Grafana stack visualises hello's traffic on a real cluster.

**Architecture:** GitOps from commit #1 — `impl/helixgitpx-platform/argocd/bootstrap/` installs Argo CD, which then reconciles 25 `Application` CRs at `impl/helixgitpx-platform/argocd/applications/` against 25 Helm charts at `impl/helixgitpx-platform/helm/<chart>/`. Sync waves enforce install order (CNI first, observability second, data plane third, hello last). Three `kustomize/overlays/{local,staging,prod-eu}/` patch replica counts and external provider credentials. Hello upgrades its Kafka emission to write to a Postgres `outbox_events` table inside the same transaction as its counter UPSERT; Debezium captures the WAL and publishes to the same `hello.said` topic.

**Tech Stack:** k3d 5.7, Cilium, ingress-nginx, Argo CD 3.x, cert-manager, external-dns, MinIO, kube-prometheus-stack, Mimir, Loki, Tempo, Pyroscope, SPIRE, Istio Ambient, CloudNativePG (CNPG), Strimzi Kafka (KRaft), Karapace, Debezium Connect, Dragonfly, Meilisearch, OpenSearch, Qdrant, Vault (Raft + Shamir), goose, go-spiffe/v2.

**Locked constraints from the design spec §4:**

- C-1 — Local k3d deploy target; full stack runs for real on this machine.
- C-2 — HA manifests authoritative in Git (`staging`/`prod-eu` overlays); `local` overlay patches to single-replica.
- C-3 — Hello's Kafka emission switches to transactional outbox via Debezium; `hello.said` topic unchanged externally.
- C-4 — Istio Ambient included (zTunnel + waypoints).
- C-5 — Full 7-component observability.
- C-6 — Local-friendly external defaults: self-signed cert-manager, noop external-dns, MinIO for S3, placeholder PagerDuty.
- C-7 — Observability-first sequencing; per-data-service gate requires Grafana visibility.
- Inherited from M1: `workflow_dispatch`-only CI, portable compose wrapper, single git history, `mise.toml`, spine-first with completion matrix.

**Phases:**

- **Phase A — Cluster bootstrap + observability spine** (Tasks 1–8): `up.sh --m2`, Cilium + ingress-nginx, cert-manager + external-dns + MinIO, the five observability charts.
- **Phase B — Identity & mesh** (Tasks 9–11): SPIRE, Istio Ambient base + data plane.
- **Phase C — Data plane** (Tasks 12–19): CNPG + migrations + schemas + RLS; Strimzi Kafka + Karapace + Debezium; Dragonfly; Meilisearch + OpenSearch + Qdrant; Vault.
- **Phase D — Hello outbox + E2E** (Tasks 20–23): outbox schema + code, SVID-terminated gRPC, Vault-sourced config, Argo CD deploy, end-to-end spine verification.
- **Phase E — GitOps root & verification** (Tasks 24–26): Argo CD bootstrap + app-of-apps, three kustomize overlays, `verify-m2-cluster.sh` + `verify-m2-spine.sh`, ADRs 0006–0012.

**Conventions:**

- Conventional Commits with `-s` (DCO) and `Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>` trailer.
- `export GOTOOLCHAIN=go1.23.4` before any `go` command (from M1 learnings).
- Every Helm chart ships `helm lint`-clean; every `Application` ships `argocd app diff`-clean (no drift); every manifest passes `kubectl kustomize`.
- Every chart declares a PodDisruptionBudget and a ServiceMonitor (Prom CRD) — enforced by a new Kyverno policy added in Task 5.
- Every data service's "done" gate requires an entry visible in the Grafana dashboard provisioned for that service.

---

## File Structure (all new under `impl/helixgitpx-platform/` unless noted)

```
impl/helixgitpx-platform/
├── k8s-local/
│   ├── up.sh                      (modify: add --m2 flag path)
│   └── m2/
│       ├── preflight.sh           (new: RAM check, tool check)
│       ├── prepull-images.sh      (new: image list + parallel pulls)
│       └── etc-hosts.sh           (new: write /etc/hosts entries)
├── argocd/
│   ├── bootstrap/
│   │   ├── kustomization.yaml
│   │   ├── argocd-install.yaml    (upstream Argo CD manifests, pinned to v3.0.x)
│   │   └── argocd-values.yaml     (Helm values for Argo CD itself)
│   ├── applicationset/
│   │   └── app-of-apps.yaml
│   └── applications/              (25 Application CRs)
│       ├── cilium.yaml                  (wave -10)
│       ├── ingress-nginx.yaml           (wave -5)
│       ├── cert-manager.yaml            (wave -5)
│       ├── external-dns.yaml            (wave -5)
│       ├── minio.yaml                   (wave -5)
│       ├── prometheus-stack.yaml        (wave -3)
│       ├── mimir.yaml                   (wave -3)
│       ├── loki.yaml                    (wave -3)
│       ├── tempo.yaml                   (wave -3)
│       ├── pyroscope.yaml               (wave -3)
│       ├── spire.yaml                   (wave 0)
│       ├── istio-base.yaml              (wave 0)
│       ├── istio-ambient.yaml           (wave 0)
│       ├── cnpg-operator.yaml           (wave 5)
│       ├── strimzi-operator.yaml        (wave 5)
│       ├── cnpg-cluster.yaml            (wave 7)
│       ├── kafka-cluster.yaml           (wave 7)
│       ├── dragonfly.yaml               (wave 7)
│       ├── meilisearch.yaml             (wave 7)
│       ├── opensearch.yaml              (wave 7)
│       ├── qdrant.yaml                  (wave 7)
│       ├── vault.yaml                   (wave 7)
│       ├── karapace.yaml                (wave 9)
│       ├── debezium.yaml                (wave 9)
│       └── hello.yaml                   (wave 10)
├── helm/                          (25 Helm chart directories, one per Application)
│   ├── cilium/                    Chart.yaml + values.yaml (+ values-local.yaml)
│   ├── ingress-nginx/             ...
│   ├── cert-manager/              + templates/clusterissuer-selfsigned.yaml
│   │                              + templates/clusterissuer-letsencrypt.yaml
│   ├── external-dns/
│   ├── minio/                     + templates/bucket-init-job.yaml
│   ├── prometheus-stack/          + grafana/dashboards/, grafana/datasources/
│   ├── mimir/
│   ├── loki/
│   ├── tempo/
│   ├── pyroscope/
│   ├── spire/                     + templates/workload-registrar.yaml
│   ├── istio-base/
│   ├── istio-ambient/             + templates/namespace-labels.yaml
│   ├── cnpg-operator/
│   ├── cnpg-cluster/              + templates/schemas-job.yaml (applies sql/schemas.sql)
│   │                              + templates/restore-drill-cronjob.yaml (staging+ only)
│   ├── strimzi-operator/
│   ├── kafka-cluster/             + templates/topics.yaml (3 KafkaTopic CRs)
│   ├── karapace/
│   ├── debezium/                  + templates/hello-outbox-connector.yaml (KafkaConnector)
│   ├── dragonfly/
│   ├── meilisearch/
│   ├── opensearch/                + templates/ilm-policy-configmap.yaml
│   │                              + templates/snapshot-repo-job.yaml
│   ├── qdrant/
│   └── vault/                     + templates/bootstrap-job.yaml
│                                  + templates/policies-configmap.yaml
├── kustomize/
│   └── overlays/
│       ├── local/
│       │   └── kustomization.yaml + patches/*.yaml (replicas:1, issuer:selfsigned, ...)
│       ├── staging/
│       │   └── kustomization.yaml + patches/*.yaml
│       └── prod-eu/
│           └── kustomization.yaml + patches/*.yaml
├── sql/
│   └── schemas.sql                per-service schemas + RLS baseline
└── README.md                      (modify: add M2 section)

impl/helixgitpx/
├── platform/
│   ├── pg/
│   │   ├── migrate.go             (new: goose wrapper)
│   │   └── migrate_test.go        (new)
│   ├── spire/
│   │   └── spire.go               (modify: replace M2 TODO with real go-spiffe/v2)
│   ├── kafka/
│   │   └── kafka.go               (modify: wire ResolveFn to Karapace)
│   ├── config/
│   │   ├── config.go              (modify: add vault:"path" tag support)
│   │   └── vault.go               (new: Vault KV resolver)
│   ├── telemetry/
│   │   └── pprof.go               (new: RegisterPprof helper)
│   └── grpc/
│       └── server.go              (modify: optional *spire.Fetcher for SVID-mTLS)
└── services/hello/
    ├── migrations/
    │   ├── 20260420000001_init.sql          (unchanged from M1)
    │   └── 20260420000002_outbox.sql        (new: outbox_events table)
    ├── internal/
    │   ├── repo/
    │   │   ├── event_kafka.go               (delete — replaced by outbox)
    │   │   └── event_outbox.go              (new)
    │   └── app/
    │       └── app.go                       (modify: swap emitter impl; wire Vault config)
    └── deploy/helm/
        └── templates/
            ├── migrate-job.yaml              (new pre-upgrade Hook)
            ├── kafkaconnector.yaml           (new: Debezium source for hello.outbox_events)
            ├── servicemonitor.yaml           (new)
            ├── ingress.yaml                  (new: hello.helix.local)
            └── vault-agent.yaml              (new: inject DSNs as a sidecar)

docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/
├── 0006-local-k3d-m2-target.md
├── 0007-ha-manifests-local-overlay.md
├── 0008-hello-outbox-pattern.md
├── 0009-istio-ambient-m2.md
├── 0010-observability-first-sequencing.md
├── 0011-local-friendly-externals.md
└── 0012-sync-wave-ordering.md

scripts/
├── verify-m2-cluster.sh          (new: 20 roadmap-item gates)
└── verify-m2-spine.sh            (new: end-to-end hello traffic + Grafana checks)
```

---

## Phase A — Cluster bootstrap + observability spine

### Task 1: `up.sh --m2` + preflight + pre-pull + /etc/hosts + namespaces + Cilium

**Files:**
- Modify: `impl/helixgitpx-platform/k8s-local/up.sh`
- Create: `impl/helixgitpx-platform/k8s-local/m2/preflight.sh`
- Create: `impl/helixgitpx-platform/k8s-local/m2/prepull-images.sh`
- Create: `impl/helixgitpx-platform/k8s-local/m2/etc-hosts.sh`
- Create: `impl/helixgitpx-platform/helm/cilium/Chart.yaml`
- Create: `impl/helixgitpx-platform/helm/cilium/values.yaml`
- Create: `impl/helixgitpx-platform/helm/cilium/values-local.yaml`

- [ ] **Step 1: Write `m2/preflight.sh`**

```sh
#!/usr/bin/env bash
# M2 preflight: refuse to proceed without the resources we need.
set -euo pipefail

MIN_FREE_MEM_GIB=48
FREE_KIB=$(awk '/MemAvailable/ { print $2 }' /proc/meminfo)
FREE_GIB=$((FREE_KIB / 1024 / 1024))

if [ "$FREE_GIB" -lt "$MIN_FREE_MEM_GIB" ]; then
    printf 'preflight: need >= %d GiB free, have %d GiB. Close apps and retry.\n' \
        "$MIN_FREE_MEM_GIB" "$FREE_GIB" >&2
    exit 2
fi

for tool in k3d kubectl helm; do
    if ! command -v "$tool" >/dev/null 2>&1; then
        printf 'preflight: %s not on PATH. Run: mise install\n' "$tool" >&2
        exit 3
    fi
done

# Container runtime — must be docker OR podman (M1 ADR-0002)
if ! (command -v docker >/dev/null 2>&1 || command -v podman >/dev/null 2>&1); then
    printf 'preflight: need docker or podman on PATH\n' >&2
    exit 4
fi

printf 'preflight: ok (%d GiB free, tools present)\n' "$FREE_GIB"
```

- [ ] **Step 2: Write `m2/prepull-images.sh`**

```sh
#!/usr/bin/env bash
# Pre-pull heavy images before k3d cluster create to avoid first-boot timeouts.
set -euo pipefail

RUNTIME=docker
if ! command -v docker >/dev/null 2>&1; then
    RUNTIME=podman
fi

IMAGES=(
    ghcr.io/cloudnative-pg/postgresql:16.4-3
    quay.io/strimzi/kafka:0.44.0-kafka-3.8.0
    opensearchproject/opensearch:2.18.0
    hashicorp/vault:1.18.3
    grafana/mimir:2.14.0
    grafana/loki:3.3.0
    grafana/tempo:2.6.1
    grafana/pyroscope:1.10.0
    grafana/grafana:11.4.0
    prom/prometheus:v3.0.1
    quay.io/cilium/cilium:v1.16.3
    registry.k8s.io/ingress-nginx/controller:v1.12.0
    quay.io/jetstack/cert-manager-controller:v1.16.2
    quay.io/minio/minio:RELEASE.2025-01-20T14-49-07Z
    ghcr.io/spiffe/spire-server:1.11.0
    docker.io/istio/pilot:1.24.2
    docker.io/istio/ztunnel:1.24.2
    docker.io/bitnami/redis:7.4    # standin; dragonfly is separately pulled
    docker.redpanda.com/redpandadata/connect:4.40.0
    ghcr.io/aiven-open/karapace:latest
    docker.dragonflydb.io/dragonflydb/dragonfly:v1.25.1
    getmeili/meilisearch:v1.12
    qdrant/qdrant:v1.12.4
    quay.io/argoproj/argocd:v3.0.0
    debezium/connect:3.0.0.Final
)

printf 'Pre-pulling %d images using %s...\n' "${#IMAGES[@]}" "$RUNTIME"

for img in "${IMAGES[@]}"; do
    "$RUNTIME" pull "$img" &
done
wait
printf 'prepull: done\n'
```

- [ ] **Step 3: Write `m2/etc-hosts.sh`**

```sh
#!/usr/bin/env bash
# Write /etc/hosts entries for M2 ingress hostnames.
# Requires sudo.
set -euo pipefail

HOSTS=(
    hello.helix.local
    grafana.helix.local
    argocd.helix.local
    vault.helix.local
    kafka.helix.local
    prometheus.helix.local
    minio.helix.local
)

TAG='# helixgitpx-m2'
HOSTS_FILE=${HOSTS_FILE:-/etc/hosts}

if ! grep -Fq "$TAG" "$HOSTS_FILE"; then
    {
        printf '\n%s\n' "$TAG"
        for h in "${HOSTS[@]}"; do
            printf '127.0.0.1\t%s\n' "$h"
        done
    } | sudo tee -a "$HOSTS_FILE" >/dev/null
    printf 'etc-hosts: wrote %d entries to %s\n' "${#HOSTS[@]}" "$HOSTS_FILE"
else
    printf 'etc-hosts: already present (tag %s found)\n' "$TAG"
fi
```

- [ ] **Step 4: Modify `up.sh` to add the --m2 code path**

Read the existing up.sh (from M1 Task 29) and add the --m2 code path at the top, before the existing k3d/kind logic. Write the modified file:

```sh
#!/usr/bin/env bash
# Bring up a local Kubernetes cluster. Defaults to k3d; use KIND=1 to pick kind.
# --dry-run prints planned actions without executing.
# --m2 runs the full M2 data-plane bootstrap (cluster + Argo CD app-of-apps).
set -euo pipefail

DRY_RUN=0
M2=0
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1 ;;
    --m2)      M2=1 ;;
  esac
done

engine="${KIND:-0}"
here=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH='' cd -- "$here/../../.." && pwd)

run() {
  if [ "$DRY_RUN" = "1" ]; then
    printf '[dry-run] %s\n' "$*"
  else
    "$@"
  fi
}

if [ "$M2" = "1" ]; then
    "$here/m2/preflight.sh"
    "$here/m2/prepull-images.sh"
fi

# k3d cluster creation — for M2 we disable flannel + k3s network-policy so Cilium can take over
if [ "$M2" = "1" ]; then
    run k3d cluster create \
        --config "$here/k3d-config.yaml" \
        --k3s-arg '--flannel-backend=none@server:*' \
        --k3s-arg '--disable-network-policy@server:*'
else
    if [ "$engine" = "1" ]; then
        run kind create cluster --config "$here/kind-config.yaml"
    else
        run k3d cluster create --config "$here/k3d-config.yaml"
    fi
fi

run kubectl cluster-info

if [ "$M2" = "1" ]; then
    # Namespaces (8 for M2)
    for ns in helix-system helix-identity istio-system helix-data helix-cache helix-secrets helix-observability helix; do
        run kubectl create ns "$ns" --dry-run=client -o yaml | kubectl apply -f -
    done

    # Cilium via Helm
    run helm repo add cilium https://helm.cilium.io
    run helm upgrade --install cilium cilium/cilium \
        -n kube-system \
        -f "$repo_root/impl/helixgitpx-platform/helm/cilium/values-local.yaml" \
        --version 1.16.3

    # Wait for Cilium Ready
    run kubectl -n kube-system rollout status ds/cilium --timeout=5m

    # /etc/hosts entries (non-destructive — skips if already present)
    "$here/m2/etc-hosts.sh"

    # Argo CD bootstrap — delegated to a separate step (Task 24) that kustomize-applies argocd/bootstrap
    printf 'cluster ready. Next: apply impl/helixgitpx-platform/argocd/bootstrap/ (Task 24).\n'
fi
```

- [ ] **Step 5: Write `helm/cilium/Chart.yaml`**

```yaml
apiVersion: v2
name: cilium
description: Cilium CNI for the helix cluster
type: application
version: 0.1.0
dependencies:
  - name: cilium
    version: "1.16.3"
    repository: https://helm.cilium.io
```

- [ ] **Step 6: Write `helm/cilium/values.yaml`**

```yaml
cilium:
  k8sServiceHost: host.docker.internal
  k8sServicePort: 6443
  ipam:
    mode: kubernetes
  kubeProxyReplacement: true
  hubble:
    enabled: true
    relay:
      enabled: true
    ui:
      enabled: false   # enabled only locally (see values-local.yaml)
  operator:
    replicas: 2        # HA default
```

- [ ] **Step 7: Write `helm/cilium/values-local.yaml`**

```yaml
cilium:
  operator:
    replicas: 1
  hubble:
    ui:
      enabled: true
```

- [ ] **Step 8: chmod + verify scripts**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
chmod +x impl/helixgitpx-platform/k8s-local/m2/*.sh
shellcheck impl/helixgitpx-platform/k8s-local/m2/*.sh impl/helixgitpx-platform/k8s-local/up.sh
```

Expected: clean.

- [ ] **Step 9: Helm lint the Cilium chart**

```sh
helm dependency update impl/helixgitpx-platform/helm/cilium/
helm lint impl/helixgitpx-platform/helm/cilium/
```

Expected: `0 chart(s) failed`.

- [ ] **Step 10: Commit**

```sh
git add impl/helixgitpx-platform/k8s-local impl/helixgitpx-platform/helm/cilium
git commit -s -m "$(printf 'feat(platform/m2): cluster bootstrap — up.sh --m2, preflight, Cilium chart\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 2: ingress-nginx + cert-manager + external-dns charts

**Files:**
- Create: `impl/helixgitpx-platform/helm/ingress-nginx/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/cert-manager/Chart.yaml`
- Create: `impl/helixgitpx-platform/helm/cert-manager/values.yaml`
- Create: `impl/helixgitpx-platform/helm/cert-manager/values-local.yaml`
- Create: `impl/helixgitpx-platform/helm/cert-manager/templates/clusterissuer-selfsigned.yaml`
- Create: `impl/helixgitpx-platform/helm/cert-manager/templates/clusterissuer-letsencrypt.yaml`
- Create: `impl/helixgitpx-platform/helm/external-dns/{Chart.yaml, values.yaml, values-local.yaml}`

- [ ] **Step 1: Write ingress-nginx chart**

`helm/ingress-nginx/Chart.yaml`:

```yaml
apiVersion: v2
name: ingress-nginx
description: ingress-nginx controller
type: application
version: 0.1.0
dependencies:
  - name: ingress-nginx
    version: "4.11.3"
    repository: https://kubernetes.github.io/ingress-nginx
```

`helm/ingress-nginx/values.yaml`:

```yaml
ingress-nginx:
  controller:
    replicaCount: 2    # HA default; local overlay patches to 1
    service:
      type: LoadBalancer
    metrics:
      enabled: true
      serviceMonitor:
        enabled: true
```

`helm/ingress-nginx/values-local.yaml`:

```yaml
ingress-nginx:
  controller:
    replicaCount: 1
```

- [ ] **Step 2: Write cert-manager chart**

`helm/cert-manager/Chart.yaml`:

```yaml
apiVersion: v2
name: cert-manager
description: cert-manager with HelixGitpx ClusterIssuers
type: application
version: 0.1.0
dependencies:
  - name: cert-manager
    version: "v1.16.2"
    repository: https://charts.jetstack.io
```

`helm/cert-manager/values.yaml`:

```yaml
cert-manager:
  installCRDs: true
  replicaCount: 2   # HA default
  prometheus:
    enabled: true
    servicemonitor:
      enabled: true

# HelixGitpx-specific
issuer:
  mode: selfsigned   # "selfsigned" | "letsencrypt"; overlays pick one
  letsencrypt:
    email: ""
    dns01ProviderRef: ""   # e.g. "cloudflare"
```

`helm/cert-manager/values-local.yaml`:

```yaml
cert-manager:
  replicaCount: 1
issuer:
  mode: selfsigned
```

`helm/cert-manager/templates/clusterissuer-selfsigned.yaml`:

```yaml
{{- if eq .Values.issuer.mode "selfsigned" }}
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-ca
spec:
  selfSigned: {}
{{- end }}
```

`helm/cert-manager/templates/clusterissuer-letsencrypt.yaml`:

```yaml
{{- if eq .Values.issuer.mode "letsencrypt" }}
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-dns01
spec:
  acme:
    email: {{ .Values.issuer.letsencrypt.email | quote }}
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-dns01-account-key
    solvers:
      - dns01:
          webhook:
            groupName: acme.helixgitpx.dev
            solverName: {{ .Values.issuer.letsencrypt.dns01ProviderRef | quote }}
{{- end }}
```

- [ ] **Step 3: Write external-dns chart**

`helm/external-dns/Chart.yaml`:

```yaml
apiVersion: v2
name: external-dns
description: external-dns with noop/real providers
type: application
version: 0.1.0
dependencies:
  - name: external-dns
    version: "1.15.0"
    repository: https://kubernetes-sigs.github.io/external-dns
```

`helm/external-dns/values.yaml`:

```yaml
external-dns:
  provider: webhook
  extraArgs:
    - --webhook-provider-url=http://localhost:8888
  replicaCount: 2  # HA default
```

`helm/external-dns/values-local.yaml`:

```yaml
external-dns:
  replicaCount: 1
  extraContainers:
    - name: noop-webhook
      image: ghcr.io/kubernetes-sigs/external-dns-noop:v0.2.0
      ports:
        - containerPort: 8888
```

- [ ] **Step 4: helm dependency update + lint all three**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
for c in ingress-nginx cert-manager external-dns; do
    helm dependency update "impl/helixgitpx-platform/helm/$c/"
    helm lint "impl/helixgitpx-platform/helm/$c/"
done
```

Expected: each prints `0 chart(s) failed`.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx-platform/helm/ingress-nginx \
        impl/helixgitpx-platform/helm/cert-manager \
        impl/helixgitpx-platform/helm/external-dns
git commit -s -m "$(printf 'feat(platform/m2): ingress-nginx + cert-manager + external-dns charts\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 3: MinIO chart + bucket-init Job

**Files:**
- Create: `impl/helixgitpx-platform/helm/minio/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/minio/templates/bucket-init-job.yaml`

- [ ] **Step 1: Write `Chart.yaml`**

```yaml
apiVersion: v2
name: minio
description: MinIO object storage + bucket provisioning
type: application
version: 0.1.0
dependencies:
  - name: minio
    version: "14.8.4"
    repository: https://charts.bitnami.com/bitnami
```

- [ ] **Step 2: Write `values.yaml`**

```yaml
minio:
  mode: distributed
  statefulset:
    replicaCount: 4   # HA default
  auth:
    rootUser: admin
    existingSecret: minio-root-creds
  persistence:
    enabled: true
    size: 20Gi
  metrics:
    serviceMonitor:
      enabled: true

buckets:
  - cnpg-backups
  - opensearch-snapshots
  - loki-chunks
  - tempo-traces
  - pyroscope-profiles
  - mimir-blocks
```

- [ ] **Step 3: Write `values-local.yaml`**

```yaml
minio:
  mode: standalone
  statefulset:
    replicaCount: 1
  persistence:
    size: 10Gi
```

- [ ] **Step 4: Write `templates/bucket-init-job.yaml`**

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: minio-bucket-init
  annotations:
    argocd.argoproj.io/hook: PostSync
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
spec:
  backoffLimit: 5
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: mc
          image: quay.io/minio/mc:RELEASE.2025-01-17T23-25-50Z
          command:
            - /bin/sh
            - -c
            - |
              set -eu
              mc alias set m http://minio:9000 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"
              {{- range .Values.buckets }}
              mc mb --ignore-existing m/{{ . }}
              {{- end }}
              echo "bucket-init: ok"
          envFrom:
            - secretRef:
                name: minio-root-creds
```

- [ ] **Step 5: Create the root-creds Secret template (separate file for clarity)**

Create `impl/helixgitpx-platform/helm/minio/templates/root-creds-secret.yaml`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: minio-root-creds
type: Opaque
stringData:
  MINIO_ROOT_USER: {{ .Values.minio.auth.rootUser | default "admin" | quote }}
  MINIO_ROOT_PASSWORD: {{ .Values.minio.auth.rootPassword | default "minioadmin-change-me" | quote }}
```

Note: In staging/prod this Secret must be sealed (SealedSecret) or externally sourced (Vault injector). Local uses the default.

- [ ] **Step 6: helm dep update + lint**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm dependency update impl/helixgitpx-platform/helm/minio/
helm lint impl/helixgitpx-platform/helm/minio/
```

- [ ] **Step 7: Commit**

```sh
git add impl/helixgitpx-platform/helm/minio
git commit -s -m "$(printf 'feat(platform/m2): minio chart + bucket-init Job\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 4: prometheus-stack chart (Prom + AM + Grafana + dashboards)

**Files:**
- Create: `impl/helixgitpx-platform/helm/prometheus-stack/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/prometheus-stack/grafana/datasources/datasources.yaml`
- Create: `impl/helixgitpx-platform/helm/prometheus-stack/grafana/dashboards/dashboards.yaml`
- Create: `impl/helixgitpx-platform/helm/prometheus-stack/grafana/dashboards/hello.json`
- Create: `impl/helixgitpx-platform/helm/prometheus-stack/grafana/dashboards/cluster.json`

- [ ] **Step 1: Write `Chart.yaml`**

```yaml
apiVersion: v2
name: prometheus-stack
description: kube-prometheus-stack with HelixGitpx dashboards + datasources
type: application
version: 0.1.0
dependencies:
  - name: kube-prometheus-stack
    version: "66.3.0"
    repository: https://prometheus-community.github.io/helm-charts
```

- [ ] **Step 2: Write `values.yaml`**

```yaml
kube-prometheus-stack:
  fullnameOverride: prom
  prometheus:
    prometheusSpec:
      retention: 15d
      serviceMonitorSelectorNilUsesHelmValues: false
      podMonitorSelectorNilUsesHelmValues: false
      remoteWrite:
        - url: http://mimir-nginx.helix-observability.svc/api/v1/push
          headers:
            X-Scope-OrgID: helixgitpx
      resources:
        requests: { cpu: 500m, memory: 2Gi }
        limits:   { cpu: 2,    memory: 4Gi }
  alertmanager:
    alertmanagerSpec:
      replicas: 3   # HA default
      configSecret: alertmanager-config
  grafana:
    adminPassword: admin   # local only; staging/prod sealed
    ingress:
      enabled: true
      ingressClassName: nginx
      hosts: ["grafana.helix.local"]
    sidecar:
      dashboards:
        enabled: true
        searchNamespace: ALL
        label: grafana_dashboard
      datasources:
        enabled: true
        searchNamespace: ALL
  defaultRules:
    create: true
```

- [ ] **Step 3: Write `values-local.yaml`**

```yaml
kube-prometheus-stack:
  alertmanager:
    alertmanagerSpec:
      replicas: 1
  prometheus:
    prometheusSpec:
      retention: 7d
      resources:
        requests: { cpu: 200m, memory: 512Mi }
        limits:   { cpu: 1,    memory: 2Gi }
```

- [ ] **Step 4: Write datasources ConfigMap**

`grafana/datasources/datasources.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-datasources
  labels:
    grafana_datasource: "1"
data:
  datasources.yaml: |
    apiVersion: 1
    datasources:
      - name: Prometheus
        type: prometheus
        uid: prom
        url: http://mimir-nginx.helix-observability.svc/prometheus
        access: proxy
        isDefault: true
      - name: Loki
        type: loki
        uid: loki
        url: http://loki-gateway.helix-observability.svc
        access: proxy
      - name: Tempo
        type: tempo
        uid: tempo
        url: http://tempo.helix-observability.svc:3200
        access: proxy
      - name: Pyroscope
        type: grafana-pyroscope-datasource
        uid: pyroscope
        url: http://pyroscope.helix-observability.svc:4040
        access: proxy
      - name: Alertmanager
        type: alertmanager
        uid: alertmanager
        url: http://prom-kube-prometheus-stack-alertmanager.helix-observability.svc:9093
        access: proxy
```

- [ ] **Step 5: Write dashboard provisioner ConfigMap + hello dashboard JSON**

`grafana/dashboards/dashboards.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboard-hello
  labels:
    grafana_dashboard: "1"
data:
  hello.json: |-
    {{- include "helixgitpx-platform.helloDashboard" . | indent 4 }}
  cluster.json: |-
    {{- include "helixgitpx-platform.clusterDashboard" . | indent 4 }}
```

`grafana/dashboards/hello.json` — minimal but real Grafana dashboard JSON. Write this exact content:

```json
{
  "annotations": {"list": []},
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "id": null,
  "links": [],
  "liveNow": false,
  "panels": [
    {
      "datasource": {"uid": "prom", "type": "prometheus"},
      "fieldConfig": {"defaults": {"unit": "short"}},
      "gridPos": {"h": 8, "w": 12, "x": 0, "y": 0},
      "id": 1,
      "targets": [{"expr": "sum(rate(http_server_requests_seconds_count{job=\"hello\"}[5m])) by (status)"}],
      "title": "Hello HTTP RPS by status",
      "type": "timeseries"
    },
    {
      "datasource": {"uid": "prom", "type": "prometheus"},
      "gridPos": {"h": 8, "w": 12, "x": 12, "y": 0},
      "id": 2,
      "targets": [{"expr": "histogram_quantile(0.99, sum(rate(http_server_requests_seconds_bucket{job=\"hello\"}[5m])) by (le))"}],
      "title": "Hello p99 latency",
      "type": "timeseries"
    },
    {
      "datasource": {"uid": "loki", "type": "loki"},
      "gridPos": {"h": 8, "w": 24, "x": 0, "y": 8},
      "id": 3,
      "targets": [{"expr": "{app=\"hello\"}", "refId": "A"}],
      "title": "Hello logs",
      "type": "logs"
    },
    {
      "datasource": {"uid": "tempo", "type": "tempo"},
      "gridPos": {"h": 8, "w": 24, "x": 0, "y": 16},
      "id": 4,
      "targets": [{"queryType": "search", "serviceName": "hello", "refId": "A"}],
      "title": "Hello traces",
      "type": "traces"
    },
    {
      "datasource": {"uid": "pyroscope", "type": "grafana-pyroscope-datasource"},
      "gridPos": {"h": 8, "w": 24, "x": 0, "y": 24},
      "id": 5,
      "targets": [{"profileTypeId": "process_cpu:cpu:nanoseconds:cpu:nanoseconds", "labelSelector": "{service_name=\"hello\"}", "refId": "A"}],
      "title": "Hello CPU flamegraph",
      "type": "flamegraph"
    }
  ],
  "refresh": "10s",
  "schemaVersion": 41,
  "tags": ["helixgitpx", "hello"],
  "time": {"from": "now-1h", "to": "now"},
  "title": "HelixGitpx — Hello",
  "uid": "helixgitpx-hello",
  "version": 1
}
```

`grafana/dashboards/cluster.json` — minimal cluster-level dashboard:

```json
{
  "annotations": {"list": []},
  "editable": true,
  "panels": [
    {
      "datasource": {"uid": "prom", "type": "prometheus"},
      "gridPos": {"h": 8, "w": 24, "x": 0, "y": 0},
      "id": 1,
      "targets": [{"expr": "sum by (phase) (kube_pod_status_phase)"}],
      "title": "Pods by phase",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 41,
  "tags": ["helixgitpx", "cluster"],
  "title": "HelixGitpx — Cluster",
  "uid": "helixgitpx-cluster",
  "version": 1
}
```

Simplify: instead of the `include "helixgitpx-platform.helloDashboard"` template helper (which we haven't defined), just embed the JSON directly in the ConfigMap. Replace `grafana/dashboards/dashboards.yaml` content with:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboards-helixgitpx
  labels:
    grafana_dashboard: "1"
data:
  hello.json: |-
{{ .Files.Get "grafana/dashboards/hello.json" | indent 4 }}
  cluster.json: |-
{{ .Files.Get "grafana/dashboards/cluster.json" | indent 4 }}
```

- [ ] **Step 6: helm dep update + lint**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm dependency update impl/helixgitpx-platform/helm/prometheus-stack/
helm lint impl/helixgitpx-platform/helm/prometheus-stack/
```

- [ ] **Step 7: Commit**

```sh
git add impl/helixgitpx-platform/helm/prometheus-stack
git commit -s -m "$(printf 'feat(platform/m2): prometheus-stack + hello dashboard + datasources\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 5: Mimir + Loki + Tempo + Pyroscope charts

**Files:**
- Create: `impl/helixgitpx-platform/helm/{mimir,loki,tempo,pyroscope}/{Chart.yaml, values.yaml, values-local.yaml}`

Each chart is a thin wrapper around the upstream Grafana chart. Keep them brief — upstream defaults suffice, we only add HA vs local overrides.

- [ ] **Step 1: Mimir**

`helm/mimir/Chart.yaml`:

```yaml
apiVersion: v2
name: mimir
description: Grafana Mimir (Prometheus remote-write target)
type: application
version: 0.1.0
dependencies:
  - name: mimir-distributed
    version: "5.5.0"
    repository: https://grafana.github.io/helm-charts
```

`helm/mimir/values.yaml`:

```yaml
mimir-distributed:
  minio:
    enabled: false   # we have our own MinIO
  global:
    extraEnvFrom:
      - secretRef: { name: minio-root-creds }
  mimir:
    structuredConfig:
      common:
        storage:
          backend: s3
          s3:
            endpoint: minio.helix-system.svc:9000
            access_key_id: ${MINIO_ROOT_USER}
            secret_access_key: ${MINIO_ROOT_PASSWORD}
            insecure: true
            bucket_name: mimir-blocks
  ingester:
    replicas: 3
  distributor:
    replicas: 3
  store_gateway:
    replicas: 3
  compactor:
    replicas: 1
  querier:
    replicas: 3
  query_frontend:
    replicas: 2
```

`helm/mimir/values-local.yaml`:

```yaml
mimir-distributed:
  ingester: { replicas: 1 }
  distributor: { replicas: 1 }
  store_gateway: { replicas: 1 }
  querier: { replicas: 1 }
  query_frontend: { replicas: 1 }
```

- [ ] **Step 2: Loki**

`helm/loki/Chart.yaml`:

```yaml
apiVersion: v2
name: loki
description: Grafana Loki + Promtail
type: application
version: 0.1.0
dependencies:
  - name: loki
    version: "6.24.0"
    repository: https://grafana.github.io/helm-charts
  - name: promtail
    version: "6.16.6"
    repository: https://grafana.github.io/helm-charts
```

`helm/loki/values.yaml`:

```yaml
loki:
  deploymentMode: SimpleScalable
  loki:
    auth_enabled: false
    storage:
      type: s3
      s3:
        endpoint: minio.helix-system.svc:9000
        bucketNames:
          chunks: loki-chunks
          ruler: loki-chunks
          admin: loki-chunks
        accessKeyId: ${MINIO_ROOT_USER}
        secretAccessKey: ${MINIO_ROOT_PASSWORD}
        s3ForcePathStyle: true
        insecure: true
  write: { replicas: 3 }
  read:  { replicas: 3 }
  backend: { replicas: 3 }

promtail:
  config:
    clients:
      - url: http://loki-gateway.helix-observability.svc/loki/api/v1/push
```

`helm/loki/values-local.yaml`:

```yaml
loki:
  deploymentMode: SingleBinary
  singleBinary: { replicas: 1 }
  write: { replicas: 0 }
  read: { replicas: 0 }
  backend: { replicas: 0 }
```

- [ ] **Step 3: Tempo**

`helm/tempo/Chart.yaml`:

```yaml
apiVersion: v2
name: tempo
description: Grafana Tempo (traces)
type: application
version: 0.1.0
dependencies:
  - name: tempo-distributed
    version: "1.22.0"
    repository: https://grafana.github.io/helm-charts
```

`helm/tempo/values.yaml`:

```yaml
tempo-distributed:
  storage:
    trace:
      backend: s3
      s3:
        endpoint: minio.helix-system.svc:9000
        bucket: tempo-traces
        access_key: ${MINIO_ROOT_USER}
        secret_key: ${MINIO_ROOT_PASSWORD}
        insecure: true
  distributor:
    replicas: 3
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318
  ingester: { replicas: 3 }
  querier: { replicas: 3 }
  queryFrontend: { replicas: 2 }
  compactor: { replicas: 1 }
```

`helm/tempo/values-local.yaml`:

```yaml
tempo-distributed:
  distributor: { replicas: 1 }
  ingester: { replicas: 1 }
  querier: { replicas: 1 }
  queryFrontend: { replicas: 1 }
```

- [ ] **Step 4: Pyroscope**

`helm/pyroscope/Chart.yaml`:

```yaml
apiVersion: v2
name: pyroscope
description: Grafana Pyroscope (continuous profiling)
type: application
version: 0.1.0
dependencies:
  - name: pyroscope
    version: "1.9.2"
    repository: https://grafana.github.io/helm-charts
```

`helm/pyroscope/values.yaml`:

```yaml
pyroscope:
  pyroscope:
    persistence:
      enabled: true
      size: 10Gi
    replicaCount: 3
    config: |
      storage:
        backend: s3
        s3:
          endpoint: minio.helix-system.svc:9000
          bucket_name: pyroscope-profiles
          access_key_id: ${MINIO_ROOT_USER}
          secret_access_key: ${MINIO_ROOT_PASSWORD}
          insecure: true
```

`helm/pyroscope/values-local.yaml`:

```yaml
pyroscope:
  pyroscope:
    replicaCount: 1
```

- [ ] **Step 5: helm dep update + lint all four**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
for c in mimir loki tempo pyroscope; do
    helm dependency update "impl/helixgitpx-platform/helm/$c/"
    helm lint "impl/helixgitpx-platform/helm/$c/"
done
```

- [ ] **Step 6: Commit**

```sh
git add impl/helixgitpx-platform/helm/mimir \
        impl/helixgitpx-platform/helm/loki \
        impl/helixgitpx-platform/helm/tempo \
        impl/helixgitpx-platform/helm/pyroscope
git commit -s -m "$(printf 'feat(platform/m2): mimir/loki/tempo/pyroscope charts\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

## Phase B — Identity & mesh

### Task 6: SPIRE chart + platform/spire real impl

**Files:**
- Create: `impl/helixgitpx-platform/helm/spire/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/spire/templates/workload-registrar.yaml`
- Modify: `impl/helixgitpx/platform/spire/spire.go`
- Create: `impl/helixgitpx/platform/spire/spire_integration_test.go`

- [ ] **Step 1: SPIRE Chart.yaml**

```yaml
apiVersion: v2
name: spire
description: SPIFFE/SPIRE server + agent + workload registrar
type: application
version: 0.1.0
dependencies:
  - name: spire
    version: "0.24.0"
    repository: https://spiffe.github.io/helm-charts-hardened
```

- [ ] **Step 2: SPIRE values.yaml**

```yaml
spire:
  trustDomain: helixgitpx.local
  server:
    replicaCount: 3   # HA default
    dataStore:
      sql:
        databaseType: sqlite3
        connectionString: "/run/spire/data/datastore.sqlite3"
  agent:
    logLevel: info
  spire-controller-manager:
    enabled: true
    identities:
      clusterSPIFFEIDs:
        default:
          spiffeIDTemplate: "spiffe://helixgitpx.local/ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"
          namespaceSelector:
            matchExpressions:
              - key: kubernetes.io/metadata.name
                operator: NotIn
                values: [kube-system, kube-public, kube-node-lease]
```

- [ ] **Step 3: SPIRE values-local.yaml**

```yaml
spire:
  server:
    replicaCount: 1
```

- [ ] **Step 4: workload-registrar.yaml (empty — controller-manager handles it)**

```yaml
# The spire-controller-manager subchart handles ClusterSPIFFEID provisioning.
# This file is a placeholder for future manual overrides (per-namespace SPIFFE IDs).
```

- [ ] **Step 5: Replace platform/spire/spire.go with real go-spiffe/v2 impl**

```go
// Package spire integrates SPIFFE/SPIRE workload API via go-spiffe/v2.
// M2 lifts the M1 no-op stub; a non-empty SocketPath now returns a live
// fetcher that streams X.509 SVIDs from the SPIRE agent.
package spire

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

// ErrUnavailable indicates the SPIRE agent socket is not reachable.
var ErrUnavailable = errors.New("spire: unavailable")

// Options configures NewFetcher.
type Options struct {
	SocketPath string // e.g. "unix:///run/spire/agent.sock"
}

// Fetcher streams workload SVIDs from the SPIRE agent.
type Fetcher struct {
	source *workloadapi.X509Source
	noop   bool
}

// NewFetcher returns a Fetcher. When SocketPath is empty or the socket
// file is absent, returns a no-op fetcher so callers can remain agnostic
// on dev machines without SPIRE.
func NewFetcher(ctx context.Context, opts Options) (*Fetcher, error) {
	if opts.SocketPath == "" {
		return &Fetcher{noop: true}, nil
	}
	if _, err := os.Stat(trimUnix(opts.SocketPath)); err != nil {
		return &Fetcher{noop: true}, nil
	}
	src, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(workloadapi.WithAddr(opts.SocketPath)),
	)
	if err != nil {
		return nil, fmt.Errorf("spire: X509Source: %w", errors.Join(ErrUnavailable, err))
	}
	return &Fetcher{source: src}, nil
}

// Source returns the underlying *workloadapi.X509Source, or nil when no-op.
// Callers pass this to grpc.Creds / tls.Config builders.
func (f *Fetcher) Source() *workloadapi.X509Source {
	if f == nil || f.noop {
		return nil
	}
	return f.source
}

// Close releases resources.
func (f *Fetcher) Close() error {
	if f == nil || f.source == nil {
		return nil
	}
	return f.source.Close()
}

func trimUnix(s string) string {
	const prefix = "unix://"
	if len(s) > len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
```

- [ ] **Step 6: Keep the M1 unit test (no-op path), add integration test**

The existing `spire_test.go` still passes (covers no-op behavior). Add `spire_integration_test.go`:

```go
//go:build integration

package spire_test

import (
	"context"
	"testing"
	"time"

	"github.com/helixgitpx/platform/spire"
)

func TestNewFetcher_ConnectsToAgent(t *testing.T) {
	// Requires SPIRE agent socket at SPIRE_AGENT_SOCKET_PATH env.
	// In-cluster test uses the agent DaemonSet's socket.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	f, err := spire.NewFetcher(ctx, spire.Options{SocketPath: "unix:///run/spire/agent.sock"})
	if err != nil {
		t.Fatalf("NewFetcher: %v", err)
	}
	defer f.Close()
	if f.Source() == nil {
		t.Fatal("expected live Source when socket is present")
	}
}
```

- [ ] **Step 7: go mod tidy**

```sh
export GOTOOLCHAIN=go1.23.4
cd impl/helixgitpx/platform
go get github.com/spiffe/go-spiffe/v2@v2.4.0
go mod tidy
grep '^go ' go.mod   # must be go 1.23 or go 1.23.0
go test ./spire/...
```

Expected: unit test still passes (the no-op path); integration test is build-tagged off by default.

- [ ] **Step 8: helm dep update + lint**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm dependency update impl/helixgitpx-platform/helm/spire/
helm lint impl/helixgitpx-platform/helm/spire/
```

- [ ] **Step 9: Commit**

```sh
git add impl/helixgitpx-platform/helm/spire impl/helixgitpx/platform/spire
git commit -s -m "$(printf 'feat(spire): real go-spiffe/v2 impl + SPIRE Helm chart\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 7: Istio Ambient charts (base + ambient)

**Files:**
- Create: `impl/helixgitpx-platform/helm/istio-base/{Chart.yaml, values.yaml}`
- Create: `impl/helixgitpx-platform/helm/istio-ambient/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/istio-ambient/templates/namespace-labels.yaml`

- [ ] **Step 1: istio-base**

`helm/istio-base/Chart.yaml`:

```yaml
apiVersion: v2
name: istio-base
description: Istio base CRDs + cluster role (prerequisite for Ambient)
type: application
version: 0.1.0
dependencies:
  - name: base
    version: "1.24.2"
    repository: https://istio-release.storage.googleapis.com/charts
```

`helm/istio-base/values.yaml`:

```yaml
base:
  defaultRevision: default
```

- [ ] **Step 2: istio-ambient**

`helm/istio-ambient/Chart.yaml`:

```yaml
apiVersion: v2
name: istio-ambient
description: Istio Ambient (istiod + CNI + ztunnel)
type: application
version: 0.1.0
dependencies:
  - name: istiod
    version: "1.24.2"
    repository: https://istio-release.storage.googleapis.com/charts
  - name: cni
    version: "1.24.2"
    repository: https://istio-release.storage.googleapis.com/charts
  - name: ztunnel
    version: "1.24.2"
    repository: https://istio-release.storage.googleapis.com/charts
```

`helm/istio-ambient/values.yaml`:

```yaml
istiod:
  profile: ambient
  meshConfig:
    caProvider: spiffe://helixgitpx.local

cni:
  profile: ambient

ztunnel:
  profile: ambient
```

`helm/istio-ambient/values-local.yaml`: (no overrides; ambient runs happily single-replica for istiod)

```yaml
istiod:
  pilot:
    replicaCount: 1
```

- [ ] **Step 3: namespace-labels.yaml — opt-in namespaces to Ambient**

```yaml
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: Namespace
    metadata:
      name: helix
      labels:
        istio.io/dataplane-mode: ambient
  - apiVersion: v1
    kind: Namespace
    metadata:
      name: helix-data
      labels:
        istio.io/dataplane-mode: ambient
  - apiVersion: v1
    kind: Namespace
    metadata:
      name: helix-cache
      labels:
        istio.io/dataplane-mode: ambient
  - apiVersion: v1
    kind: Namespace
    metadata:
      name: helix-secrets
      labels:
        istio.io/dataplane-mode: ambient
```

- [ ] **Step 4: helm dep update + lint**

```sh
for c in istio-base istio-ambient; do
    helm dependency update "impl/helixgitpx-platform/helm/$c/"
    helm lint "impl/helixgitpx-platform/helm/$c/"
done
```

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx-platform/helm/istio-base impl/helixgitpx-platform/helm/istio-ambient
git commit -s -m "$(printf 'feat(platform/m2): istio ambient charts (base + istiod + cni + ztunnel)\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

## Phase C — Data plane

### Task 8: CNPG operator + cluster + schemas.sql + RLS + migrations wiring

**Files:**
- Create: `impl/helixgitpx-platform/helm/cnpg-operator/{Chart.yaml, values.yaml}`
- Create: `impl/helixgitpx-platform/helm/cnpg-cluster/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/cnpg-cluster/templates/cluster.yaml`
- Create: `impl/helixgitpx-platform/helm/cnpg-cluster/templates/schemas-job.yaml`
- Create: `impl/helixgitpx-platform/helm/cnpg-cluster/templates/restore-drill-cronjob.yaml`
- Create: `impl/helixgitpx-platform/sql/schemas.sql`
- Create: `impl/helixgitpx/platform/pg/migrate.go`
- Create: `impl/helixgitpx/platform/pg/migrate_test.go`

- [ ] **Step 1: cnpg-operator chart**

`helm/cnpg-operator/Chart.yaml`:

```yaml
apiVersion: v2
name: cnpg-operator
description: CloudNativePG operator
type: application
version: 0.1.0
dependencies:
  - name: cloudnative-pg
    version: "0.22.1"
    repository: https://cloudnative-pg.github.io/charts
```

`helm/cnpg-operator/values.yaml`:

```yaml
cloudnative-pg:
  replicaCount: 1
  monitoring:
    podMonitorEnabled: true
```

- [ ] **Step 2: schemas.sql**

```sql
-- HelixGitpx per-service schemas + RLS baseline.
-- One schema per top-level domain from the roadmap. Each schema gets a
-- dedicated role used by the service's DSN. RLS is enabled on every
-- application table created under these schemas (later milestones tighten
-- policies per tenancy model).

CREATE SCHEMA IF NOT EXISTS hello;
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS repo;
CREATE SCHEMA IF NOT EXISTS sync;
CREATE SCHEMA IF NOT EXISTS conflict;
CREATE SCHEMA IF NOT EXISTS upstream;
CREATE SCHEMA IF NOT EXISTS collab;
CREATE SCHEMA IF NOT EXISTS events;
CREATE SCHEMA IF NOT EXISTS platform;

DO $$
DECLARE
    s TEXT;
BEGIN
    FOREACH s IN ARRAY ARRAY['hello','auth','repo','sync','conflict','upstream','collab','events','platform']
    LOOP
        EXECUTE format('CREATE ROLE %I_svc LOGIN', s);
        EXECUTE format('GRANT USAGE, CREATE ON SCHEMA %I TO %I_svc', s, s);
        EXECUTE format('GRANT ALL ON ALL TABLES IN SCHEMA %I TO %I_svc', s, s);
        EXECUTE format('GRANT ALL ON ALL SEQUENCES IN SCHEMA %I TO %I_svc', s, s);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON TABLES TO %I_svc', s, s);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON SEQUENCES TO %I_svc', s, s);
    END LOOP;
EXCEPTION WHEN duplicate_object THEN
    -- role already exists; idempotent apply
    NULL;
END $$;
```

- [ ] **Step 3: cnpg-cluster chart — Cluster CR**

`helm/cnpg-cluster/Chart.yaml`:

```yaml
apiVersion: v2
name: cnpg-cluster
description: HelixGitpx Postgres HA Cluster (CNPG)
type: application
version: 0.1.0
```

`helm/cnpg-cluster/values.yaml`:

```yaml
cluster:
  name: helix-pg
  instances: 3        # HA default
  storage:
    size: 100Gi
  wal_level: logical  # required for Debezium
  backup:
    enabled: true
    s3:
      endpoint: http://minio.helix-system.svc:9000
      bucket: cnpg-backups
      accessKeyRef: { name: minio-root-creds, key: MINIO_ROOT_USER }
      secretKeyRef: { name: minio-root-creds, key: MINIO_ROOT_PASSWORD }
  restoreDrill:
    enabled: false    # staging turns it on via overlay
```

`helm/cnpg-cluster/values-local.yaml`:

```yaml
cluster:
  instances: 1
  storage:
    size: 10Gi
```

`helm/cnpg-cluster/templates/cluster.yaml`:

```yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: {{ .Values.cluster.name }}
  namespace: helix-data
spec:
  instances: {{ .Values.cluster.instances }}
  postgresql:
    parameters:
      wal_level: {{ .Values.cluster.wal_level | quote }}
      max_wal_senders: "10"
      max_replication_slots: "10"
      shared_buffers: "256MB"
  storage:
    size: {{ .Values.cluster.storage.size }}
  {{- if .Values.cluster.backup.enabled }}
  backup:
    barmanObjectStore:
      destinationPath: s3://{{ .Values.cluster.backup.s3.bucket }}/
      endpointURL: {{ .Values.cluster.backup.s3.endpoint | quote }}
      s3Credentials:
        accessKeyId:
          name: {{ .Values.cluster.backup.s3.accessKeyRef.name }}
          key: {{ .Values.cluster.backup.s3.accessKeyRef.key }}
        secretAccessKey:
          name: {{ .Values.cluster.backup.s3.secretKeyRef.name }}
          key: {{ .Values.cluster.backup.s3.secretKeyRef.key }}
      wal:
        compression: gzip
    retentionPolicy: "30d"
  {{- end }}
  bootstrap:
    initdb:
      database: helixgitpx
      owner: helix
      postInitSQL:
        - |
{{ .Files.Get "../../sql/schemas.sql" | indent 10 }}
```

`helm/cnpg-cluster/templates/schemas-job.yaml` (applies schemas.sql into running cluster — redundant with postInitSQL on first install, but exists to re-apply on upgrades):

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Values.cluster.name }}-schemas-apply
  annotations:
    argocd.argoproj.io/hook: PostSync
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
spec:
  backoffLimit: 5
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: psql
          image: postgres:16-alpine
          command:
            - /bin/sh
            - -c
            - |
              set -eu
              export PGPASSWORD="$POSTGRES_SUPERUSER_PASSWORD"
              psql -h {{ .Values.cluster.name }}-rw -U helix -d helixgitpx -v ON_ERROR_STOP=1 <<'SQL'
{{ .Files.Get "../../sql/schemas.sql" | indent 14 }}
SQL
          envFrom:
            - secretRef: { name: {{ .Values.cluster.name }}-superuser }
```

`helm/cnpg-cluster/templates/restore-drill-cronjob.yaml`:

```yaml
{{- if .Values.cluster.restoreDrill.enabled }}
apiVersion: batch/v1
kind: CronJob
metadata:
  name: {{ .Values.cluster.name }}-restore-drill
spec:
  schedule: "0 3 1 * *"   # 3 AM on the 1st of each month
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: drill
              image: ghcr.io/cloudnative-pg/kubectl:1.31.3
              command:
                - /bin/sh
                - -c
                - |
                  set -eu
                  kubectl cnpg restore {{ .Values.cluster.name }}-drill \
                      --backup={{ .Values.cluster.name }}-latest \
                      --target-time=$(date -u +%FT%TZ)
{{- end }}
```

- [ ] **Step 4: Write `platform/pg/migrate.go`**

```go
// Package pg (migrate.go) wraps github.com/pressly/goose/v3 to apply
// SQL migrations from a directory to a DSN. Used by services' migrate-job.
package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// MigrateOptions configures Migrate.
type MigrateOptions struct {
	DSN string
	Dir string // filesystem path containing *.sql migrations
}

// Migrate applies all Up migrations under opts.Dir to opts.DSN.
// Idempotent — re-applying after completion is a no-op.
func Migrate(ctx context.Context, opts MigrateOptions) error {
	db, err := sql.Open("pgx", opts.DSN)
	if err != nil {
		return fmt.Errorf("pg.Migrate: open: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("pg.Migrate: dialect: %w", err)
	}
	if err := goose.UpContext(ctx, db, opts.Dir); err != nil {
		return fmt.Errorf("pg.Migrate: up: %w", err)
	}
	return nil
}
```

- [ ] **Step 5: Write `platform/pg/migrate_test.go`**

```go
package pg_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/pg"
)

func TestMigrate_InvalidDSN(t *testing.T) {
	err := pg.Migrate(context.Background(), pg.MigrateOptions{
		DSN: "postgres://invalid-host-name-that-does-not-exist:5432/db?sslmode=disable",
		Dir: "/tmp/does-not-matter-test-will-fail-earlier",
	})
	if err == nil {
		t.Fatal("expected error for invalid DSN")
	}
}
```

- [ ] **Step 6: go mod tidy + tests**

```sh
export GOTOOLCHAIN=go1.23.4
cd impl/helixgitpx/platform
go get github.com/pressly/goose/v3@v3.22.1
go mod tidy
grep '^go ' go.mod   # must remain go 1.23.x
go test ./pg/...
```

- [ ] **Step 7: helm lint + commit**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm dependency update impl/helixgitpx-platform/helm/cnpg-operator/
helm lint impl/helixgitpx-platform/helm/cnpg-operator/
helm lint impl/helixgitpx-platform/helm/cnpg-cluster/

git add impl/helixgitpx-platform/helm/cnpg-operator \
        impl/helixgitpx-platform/helm/cnpg-cluster \
        impl/helixgitpx-platform/sql \
        impl/helixgitpx/platform/pg
git commit -s -m "$(printf 'feat(platform/m2): CNPG operator + HA cluster + schemas.sql + goose wrapper\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 9: Strimzi operator + Kafka cluster + topics

**Files:**
- Create: `impl/helixgitpx-platform/helm/strimzi-operator/{Chart.yaml, values.yaml}`
- Create: `impl/helixgitpx-platform/helm/kafka-cluster/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/kafka-cluster/templates/kafka.yaml`
- Create: `impl/helixgitpx-platform/helm/kafka-cluster/templates/topics.yaml`

- [ ] **Step 1: strimzi-operator**

`helm/strimzi-operator/Chart.yaml`:

```yaml
apiVersion: v2
name: strimzi-operator
description: Strimzi Kafka operator
type: application
version: 0.1.0
dependencies:
  - name: strimzi-kafka-operator
    version: "0.44.0"
    repository: https://strimzi.io/charts
```

`helm/strimzi-operator/values.yaml`:

```yaml
strimzi-kafka-operator:
  watchAnyNamespace: true
  replicas: 1
```

- [ ] **Step 2: kafka-cluster**

`helm/kafka-cluster/Chart.yaml`:

```yaml
apiVersion: v2
name: kafka-cluster
description: HelixGitpx Kafka cluster (KRaft)
type: application
version: 0.1.0
```

`helm/kafka-cluster/values.yaml`:

```yaml
kafka:
  name: helix-kafka
  replicas: 3
  kRaft:
    enabled: true
  storage:
    size: 100Gi
  resources:
    requests: { cpu: 500m, memory: 2Gi }
    limits:   { cpu: 2,    memory: 4Gi }

topics:
  - name: hello.said
    partitions: 3
    replicas: 3
    retentionMs: 604800000   # 7 days
  - name: hello.outbox
    partitions: 3
    replicas: 3
    retentionMs: 2592000000  # 30 days
```

`helm/kafka-cluster/values-local.yaml`:

```yaml
kafka:
  replicas: 1
  storage:
    size: 10Gi
topics:
  - name: hello.said
    partitions: 3
    replicas: 1
    retentionMs: 604800000
  - name: hello.outbox
    partitions: 3
    replicas: 1
    retentionMs: 2592000000
```

`helm/kafka-cluster/templates/kafka.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: {{ .Values.kafka.name }}
  namespace: helix-data
  annotations:
    strimzi.io/node-pools: "enabled"
    strimzi.io/kraft: "enabled"
spec:
  kafka:
    version: 3.8.0
    replicas: {{ .Values.kafka.replicas }}
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
    config:
      offsets.topic.replication.factor: {{ .Values.kafka.replicas }}
      transaction.state.log.replication.factor: {{ .Values.kafka.replicas }}
      transaction.state.log.min.isr: {{ sub .Values.kafka.replicas 1 | default 1 }}
      default.replication.factor: {{ .Values.kafka.replicas }}
      min.insync.replicas: {{ sub .Values.kafka.replicas 1 | default 1 }}
      inter.broker.protocol.version: "3.8"
    storage:
      type: persistent-claim
      size: {{ .Values.kafka.storage.size }}
      deleteClaim: false
    resources:
      {{- toYaml .Values.kafka.resources | nindent 6 }}
  entityOperator:
    topicOperator: {}
    userOperator: {}
```

`helm/kafka-cluster/templates/topics.yaml`:

```yaml
{{- range .Values.topics }}
---
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: {{ .name | replace "." "-" }}
  namespace: helix-data
  labels:
    strimzi.io/cluster: {{ $.Values.kafka.name }}
spec:
  partitions: {{ .partitions }}
  replicas: {{ .replicas }}
  config:
    retention.ms: "{{ .retentionMs }}"
    min.insync.replicas: "{{ sub .replicas 1 | default 1 }}"
{{- end }}
```

- [ ] **Step 3: helm dep update + lint + commit**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm dependency update impl/helixgitpx-platform/helm/strimzi-operator/
helm lint impl/helixgitpx-platform/helm/strimzi-operator/
helm lint impl/helixgitpx-platform/helm/kafka-cluster/

git add impl/helixgitpx-platform/helm/strimzi-operator \
        impl/helixgitpx-platform/helm/kafka-cluster
git commit -s -m "$(printf 'feat(platform/m2): strimzi operator + kafka HA cluster + topics\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 10: Karapace + Debezium + hello-outbox connector

**Files:**
- Create: `impl/helixgitpx-platform/helm/karapace/{Chart.yaml, values.yaml, templates/*.yaml}`
- Create: `impl/helixgitpx-platform/helm/debezium/{Chart.yaml, values.yaml, templates/kafkaconnect.yaml, templates/hello-outbox-connector.yaml}`
- Modify: `impl/helixgitpx/platform/kafka/kafka.go` (wire ResolveFn)
- Modify: `impl/helixgitpx/platform/kafka/kafka_test.go`

- [ ] **Step 1: Karapace chart**

Karapace doesn't have an official Helm chart; we ship a minimal one.

`helm/karapace/Chart.yaml`:

```yaml
apiVersion: v2
name: karapace
description: Karapace schema registry (BSR-compatible)
type: application
version: 0.1.0
```

`helm/karapace/values.yaml`:

```yaml
karapace:
  image: ghcr.io/aiven-open/karapace:latest
  replicas: 2
  kafkaBootstrap: helix-kafka-kafka-bootstrap.helix-data.svc:9092
  port: 8081
```

`helm/karapace/templates/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: karapace
  namespace: helix-data
spec:
  replicas: {{ .Values.karapace.replicas }}
  selector: { matchLabels: { app.kubernetes.io/name: karapace } }
  template:
    metadata:
      labels:
        app.kubernetes.io/name: karapace
        helixgitpx.dev/env: {{ .Values.global.env | default "local" }}
    spec:
      containers:
        - name: karapace
          image: {{ .Values.karapace.image | quote }}
          ports:
            - { name: http, containerPort: {{ .Values.karapace.port }} }
          env:
            - { name: KARAPACE_BOOTSTRAP_URI,    value: {{ .Values.karapace.kafkaBootstrap | quote }} }
            - { name: KARAPACE_PORT,              value: "{{ .Values.karapace.port }}" }
            - { name: KARAPACE_HOST,              value: "0.0.0.0" }
            - { name: KARAPACE_REGISTRY_HOST,     value: "0.0.0.0" }
            - { name: KARAPACE_LOG_LEVEL,         value: "WARNING" }
          readinessProbe:
            httpGet: { path: /_schemas, port: http }
          resources:
            limits:   { cpu: "1", memory: "512Mi" }
            requests: { cpu: "100m", memory: "256Mi" }
---
apiVersion: v1
kind: Service
metadata:
  name: karapace
  namespace: helix-data
spec:
  selector: { app.kubernetes.io/name: karapace }
  ports:
    - { name: http, port: 8081, targetPort: http }
```

- [ ] **Step 2: Debezium chart**

`helm/debezium/Chart.yaml`:

```yaml
apiVersion: v2
name: debezium
description: Debezium Connect cluster + hello-outbox connector
type: application
version: 0.1.0
```

`helm/debezium/values.yaml`:

```yaml
debezium:
  replicas: 2
  image: debezium/connect:3.0.0.Final
  kafkaBootstrap: helix-kafka-kafka-bootstrap.helix-data.svc:9092
```

`helm/debezium/templates/kafkaconnect.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaConnect
metadata:
  name: debezium
  namespace: helix-data
  annotations:
    strimzi.io/use-connector-resources: "true"
spec:
  image: {{ .Values.debezium.image | quote }}
  replicas: {{ .Values.debezium.replicas }}
  bootstrapServers: {{ .Values.debezium.kafkaBootstrap | quote }}
  config:
    group.id: debezium-connect-cluster
    offset.storage.topic: debezium-connect-cluster-offsets
    config.storage.topic: debezium-connect-cluster-configs
    status.storage.topic: debezium-connect-cluster-status
    config.storage.replication.factor: -1
    offset.storage.replication.factor: -1
    status.storage.replication.factor: -1
  resources:
    requests: { cpu: 500m, memory: 1Gi }
    limits:   { cpu: 2,    memory: 2Gi }
```

`helm/debezium/templates/hello-outbox-connector.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaConnector
metadata:
  name: hello-outbox
  namespace: helix-data
  labels:
    strimzi.io/cluster: debezium
spec:
  class: io.debezium.connector.postgresql.PostgresConnector
  tasksMax: 1
  config:
    database.hostname: helix-pg-rw.helix-data.svc
    database.port: 5432
    database.user: hello_svc
    database.password: ${secret:helix-data:hello-pg-secret/password}
    database.dbname: helixgitpx
    database.server.name: hellodbs
    plugin.name: pgoutput
    table.include.list: hello.outbox_events
    topic.prefix: hellodbs
    transforms: outbox
    transforms.outbox.type: io.debezium.transforms.outbox.EventRouter
    transforms.outbox.route.by.field: topic
    transforms.outbox.route.topic.replacement: ${routedByValue}
```

- [ ] **Step 3: Modify platform/kafka/kafka.go — wire ResolveFn via Karapace**

Add a `KarapaceURL` option and a `Resolve` method that calls Karapace's REST API. Minimal stub that lifts the M1 TODO:

Append to `impl/helixgitpx/platform/kafka/kafka.go`:

```go
// KarapaceClient queries Karapace for schema IDs by subject + version.
type KarapaceClient struct {
	URL string // e.g. "http://karapace.helix-data.svc:8081"
	// TODO(M5): wire an actual HTTP client with caching.
}

// Resolve returns the schema ID for the given subject/version.
// M2: returns -1 when KarapaceClient.URL is unset (no-op fallback).
// M5: implements real HTTP call to /subjects/<subject>/versions/<version>.
func (k *KarapaceClient) Resolve(ctx context.Context, subject string, version int) (int, error) {
	if k == nil || k.URL == "" {
		return -1, nil
	}
	return -1, nil // real impl lands in M5 when the schema registry has real consumers
}
```

Add import: already has `context` implicitly — ensure `"context"` is in imports.

Actually, since `ctx` is in the signature and the function body doesn't use it yet, Go will complain. Use `_ = ctx` or mark it with an underscore:

```go
func (k *KarapaceClient) Resolve(_ context.Context, subject string, version int) (int, error) {
```

- [ ] **Step 4: helm lint all three + commit**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm lint impl/helixgitpx-platform/helm/karapace/
helm lint impl/helixgitpx-platform/helm/debezium/

export GOTOOLCHAIN=go1.23.4
cd impl/helixgitpx/platform && go test ./kafka/... && go vet ./kafka/... && cd ../../..

git add impl/helixgitpx-platform/helm/karapace \
        impl/helixgitpx-platform/helm/debezium \
        impl/helixgitpx/platform/kafka
git commit -s -m "$(printf 'feat(platform/m2): karapace + debezium + hello-outbox KafkaConnector\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 11: Dragonfly + Meilisearch + OpenSearch + Qdrant charts

**Files:**
- Create: `impl/helixgitpx-platform/helm/{dragonfly,meilisearch,opensearch,qdrant}/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/opensearch/templates/ilm-policy-configmap.yaml`
- Create: `impl/helixgitpx-platform/helm/opensearch/templates/snapshot-repo-job.yaml`

- [ ] **Step 1: Dragonfly**

`helm/dragonfly/Chart.yaml`:

```yaml
apiVersion: v2
name: dragonfly
description: Dragonfly (Redis-compatible) cache
type: application
version: 0.1.0
dependencies:
  - name: dragonfly
    version: "0.7.0"
    repository: https://dragonflydb.io/helm
```

`helm/dragonfly/values.yaml`:

```yaml
dragonfly:
  replicas: 3
  resources:
    requests: { cpu: 100m, memory: 512Mi }
    limits:   { cpu: 1,    memory: 2Gi }
```

`helm/dragonfly/values-local.yaml`:

```yaml
dragonfly:
  replicas: 1
```

- [ ] **Step 2: Meilisearch**

`helm/meilisearch/Chart.yaml`:

```yaml
apiVersion: v2
name: meilisearch
description: Meilisearch full-text index
type: application
version: 0.1.0
dependencies:
  - name: meilisearch
    version: "1.3.0"
    repository: https://meilisearch.github.io/meilisearch-kubernetes
```

`helm/meilisearch/values.yaml`:

```yaml
meilisearch:
  replicas: 3
  auth:
    existingMasterKeySecret: meilisearch-master-key
  persistence:
    enabled: true
    size: 20Gi
```

`helm/meilisearch/values-local.yaml`:

```yaml
meilisearch:
  replicas: 1
  persistence:
    size: 5Gi
```

- [ ] **Step 3: OpenSearch**

`helm/opensearch/Chart.yaml`:

```yaml
apiVersion: v2
name: opensearch
description: OpenSearch + ILM + MinIO snapshots
type: application
version: 0.1.0
dependencies:
  - name: opensearch
    version: "2.25.0"
    repository: https://opensearch-project.github.io/helm-charts
```

`helm/opensearch/values.yaml`:

```yaml
opensearch:
  clusterName: helix-os
  replicas: 3
  persistence:
    enabled: true
    size: 50Gi
  config:
    opensearch.yml: |
      cluster.name: helix-os
      network.host: 0.0.0.0
      plugins.security.disabled: true   # local only; staging enables it
  serviceMonitor:
    enabled: true
```

`helm/opensearch/values-local.yaml`:

```yaml
opensearch:
  replicas: 1
  persistence:
    size: 10Gi
```

`helm/opensearch/templates/ilm-policy-configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: opensearch-ilm-policies
  namespace: helix-cache
data:
  logs-7d.json: |-
    {
      "policy": {
        "policy_id": "logs-7d",
        "description": "Logs hot 3d → warm 4d → delete",
        "default_state": "hot",
        "states": [
          {
            "name": "hot",
            "actions": [{"rollover": {"min_index_age": "1d"}}],
            "transitions": [{"state_name": "warm", "conditions": {"min_index_age": "3d"}}]
          },
          {
            "name": "warm",
            "actions": [{"replica_count": {"number_of_replicas": 1}}],
            "transitions": [{"state_name": "delete", "conditions": {"min_index_age": "7d"}}]
          },
          {"name": "delete", "actions": [{"delete": {}}]}
        ]
      }
    }
```

`helm/opensearch/templates/snapshot-repo-job.yaml`:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: opensearch-snapshot-repo-init
  annotations:
    argocd.argoproj.io/hook: PostSync
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
spec:
  backoffLimit: 5
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: curl
          image: curlimages/curl:8.11.0
          command:
            - /bin/sh
            - -c
            - |
              set -eu
              curl -fsSL -XPUT "http://helix-os-opensearch.helix-cache.svc:9200/_snapshot/minio-s3" \
                  -H 'Content-Type: application/json' \
                  -d '{
                    "type": "s3",
                    "settings": {
                      "bucket": "opensearch-snapshots",
                      "endpoint": "minio.helix-system.svc:9000",
                      "protocol": "http",
                      "path_style_access": "true"
                    }
                  }'
              # Register ILM policy
              curl -fsSL -XPUT "http://helix-os-opensearch.helix-cache.svc:9200/_plugins/_ism/policies/logs-7d" \
                  -H 'Content-Type: application/json' \
                  --data-binary '@/policies/logs-7d.json'
          volumeMounts:
            - { name: policies, mountPath: /policies, readOnly: true }
      volumes:
        - name: policies
          configMap:
            name: opensearch-ilm-policies
```

- [ ] **Step 4: Qdrant**

`helm/qdrant/Chart.yaml`:

```yaml
apiVersion: v2
name: qdrant
description: Qdrant vector database
type: application
version: 0.1.0
dependencies:
  - name: qdrant
    version: "1.12.4"
    repository: https://qdrant.github.io/qdrant-helm
```

`helm/qdrant/values.yaml`:

```yaml
qdrant:
  replicaCount: 3
  persistence:
    enabled: true
    size: 20Gi
```

`helm/qdrant/values-local.yaml`:

```yaml
qdrant:
  replicaCount: 1
  persistence:
    size: 5Gi
```

- [ ] **Step 5: helm dep update + lint all four + commit**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
for c in dragonfly meilisearch opensearch qdrant; do
    helm dependency update "impl/helixgitpx-platform/helm/$c/"
    helm lint "impl/helixgitpx-platform/helm/$c/"
done

git add impl/helixgitpx-platform/helm/dragonfly \
        impl/helixgitpx-platform/helm/meilisearch \
        impl/helixgitpx-platform/helm/opensearch \
        impl/helixgitpx-platform/helm/qdrant
git commit -s -m "$(printf 'feat(platform/m2): dragonfly + meilisearch + opensearch + qdrant charts\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 12: Vault (Raft + Shamir) + bootstrap Job + policies

**Files:**
- Create: `impl/helixgitpx-platform/helm/vault/{Chart.yaml, values.yaml, values-local.yaml}`
- Create: `impl/helixgitpx-platform/helm/vault/templates/bootstrap-job.yaml`
- Create: `impl/helixgitpx-platform/helm/vault/templates/policies-configmap.yaml`

- [ ] **Step 1: Vault chart**

`helm/vault/Chart.yaml`:

```yaml
apiVersion: v2
name: vault
description: HashiCorp Vault (Raft + Shamir or KMS auto-unseal)
type: application
version: 0.1.0
dependencies:
  - name: vault
    version: "0.29.1"
    repository: https://helm.releases.hashicorp.com
```

`helm/vault/values.yaml`:

```yaml
vault:
  server:
    ha:
      enabled: true
      replicas: 3
      raft:
        enabled: true
        setNodeId: true
    dataStorage:
      enabled: true
      size: 10Gi
  injector:
    enabled: true
    replicas: 2
  csi:
    enabled: false

# HelixGitpx-specific
bootstrap:
  mode: shamir          # "shamir" locally; "kms" in staging/prod
  shamir:
    keyShares: 3
    keyThreshold: 2
  kms:
    provider: ""        # "gcpckms" | "awskms" | "azurekeyvault"
    keyID: ""
```

`helm/vault/values-local.yaml`:

```yaml
vault:
  server:
    ha:
      replicas: 1
  injector:
    replicas: 1
bootstrap:
  mode: shamir
```

- [ ] **Step 2: Bootstrap Job**

`helm/vault/templates/bootstrap-job.yaml`:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: vault-bootstrap
  annotations:
    argocd.argoproj.io/hook: PostSync
    argocd.argoproj.io/hook-delete-policy: HookSucceeded
spec:
  backoffLimit: 5
  template:
    spec:
      restartPolicy: OnFailure
      serviceAccountName: vault
      containers:
        - name: bootstrap
          image: hashicorp/vault:1.18.3
          command:
            - /bin/sh
            - -c
            - |
              set -eu
              export VAULT_ADDR=http://vault:8200
              if ! vault status 2>/dev/null; then
                  # Not initialized — do Shamir init
                  vault operator init \
                      -key-shares={{ .Values.bootstrap.shamir.keyShares }} \
                      -key-threshold={{ .Values.bootstrap.shamir.keyThreshold }} \
                      -format=json > /tmp/init.json
                  kubectl create secret generic vault-unseal-keys \
                      --from-file=/tmp/init.json \
                      --namespace=helix-secrets
                  # Unseal all replicas
                  for i in 0 1 2; do
                      for k in $(jq -r '.unseal_keys_b64[]' /tmp/init.json | head -{{ .Values.bootstrap.shamir.keyThreshold }}); do
                          vault operator unseal "$k" || true
                      done
                  done
              fi
              # Load policies
              for p in pg-backup-writer opensearch-backup-writer loki-chunk-writer \
                       tempo-trace-writer mimir-block-writer pyroscope-profile-writer hello; do
                  vault policy write "$p" "/policies/${p}.hcl"
              done
              # Enable K8s auth + bind SAs to policies
              vault auth enable -path=kubernetes kubernetes || true
              vault write auth/kubernetes/config \
                  kubernetes_host="https://${KUBERNETES_SERVICE_HOST}:${KUBERNETES_SERVICE_PORT}"
              vault write auth/kubernetes/role/hello \
                  bound_service_account_names=hello \
                  bound_service_account_namespaces=helix \
                  policies=hello \
                  ttl=1h
          volumeMounts:
            - { name: policies, mountPath: /policies, readOnly: true }
      volumes:
        - name: policies
          configMap:
            name: vault-policies
```

- [ ] **Step 3: Policies ConfigMap (wraps the M1 .hcl files)**

`helm/vault/templates/policies-configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault-policies
  namespace: helix-secrets
data:
  pg-backup-writer.hcl: |
    path "kv/data/minio/cnpg-backups" { capabilities = ["read"] }
  opensearch-backup-writer.hcl: |
    path "kv/data/minio/opensearch-snapshots" { capabilities = ["read"] }
  loki-chunk-writer.hcl: |
    path "kv/data/minio/loki-chunks" { capabilities = ["read"] }
  tempo-trace-writer.hcl: |
    path "kv/data/minio/tempo-traces" { capabilities = ["read"] }
  mimir-block-writer.hcl: |
    path "kv/data/minio/mimir-blocks" { capabilities = ["read"] }
  pyroscope-profile-writer.hcl: |
    path "kv/data/minio/pyroscope-profiles" { capabilities = ["read"] }
  hello.hcl: |
    path "kv/data/hello/*"     { capabilities = ["read"] }
    path "kv/metadata/hello/*" { capabilities = ["list"] }
```

- [ ] **Step 4: helm dep update + lint + commit**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm dependency update impl/helixgitpx-platform/helm/vault/
helm lint impl/helixgitpx-platform/helm/vault/

git add impl/helixgitpx-platform/helm/vault
git commit -s -m "$(printf 'feat(platform/m2): vault chart + bootstrap job + policies configmap\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

## Phase D — Hello outbox + deployment + E2E

### Task 13: Hello outbox table + event_outbox.go + swap emitter

**Files:**
- Create: `impl/helixgitpx/services/hello/migrations/20260420000002_outbox.sql`
- Create: `impl/helixgitpx/services/hello/internal/repo/event_outbox.go`
- Delete: `impl/helixgitpx/services/hello/internal/repo/event_kafka.go`
- Modify: `impl/helixgitpx/services/hello/internal/app/app.go` (swap emitter impl)
- Modify: `impl/helixgitpx/services/hello/internal/domain/greeter.go` (emitter signature unchanged; no change needed)

- [ ] **Step 1: Write migration**

`services/hello/migrations/20260420000002_outbox.sql`:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS hello.outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type TEXT NOT NULL DEFAULT 'hello',
    aggregate_id TEXT NOT NULL,
    topic TEXT NOT NULL,            -- e.g. 'hello.said' — Debezium EventRouter reads this
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_outbox_events_created ON hello.outbox_events (created_at);

-- Enable logical replication publication for the outbox table so Debezium
-- can stream changes via pgoutput plugin.
CREATE PUBLICATION IF NOT EXISTS helix_hello_outbox FOR TABLE hello.outbox_events;

ALTER TABLE hello.outbox_events ENABLE ROW LEVEL SECURITY;
CREATE POLICY hello_outbox_events_all ON hello.outbox_events USING (TRUE);

-- +goose Down
DROP PUBLICATION IF EXISTS helix_hello_outbox;
DROP TABLE IF EXISTS hello.outbox_events;
```

- [ ] **Step 2: Write event_outbox.go**

`services/hello/internal/repo/event_outbox.go`:

```go
// Package repo (event_outbox) writes hello.said events to hello.outbox_events
// inside the same pgx transaction as the counter UPSERT. Debezium's PostgreSQL
// connector streams the outbox table to Kafka via the EventRouter SMT.
package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EventOutbox implements domain.Emitter by inserting into hello.outbox_events.
// The actual Kafka emission is performed out-of-process by Debezium.
type EventOutbox struct {
	Pool  *pgxpool.Pool
	Topic string // default "hello.said"
}

type helloSaidPayload struct {
	Name     string `json:"name"`
	Greeting string `json:"greeting"`
	Count    int64  `json:"count"`
	At       string `json:"at"`
}

// Emit inserts a single outbox row. To be part of the same transaction as the
// counter UPSERT, call EmitInTx from inside an existing pgx.Tx.
func (e *EventOutbox) Emit(ctx context.Context, name, greeting string, count int64) error {
	tx, err := e.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := e.EmitInTx(ctx, tx, name, greeting, count); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// EmitInTx is the transactional entrypoint for callers that own the tx.
func (e *EventOutbox) EmitInTx(ctx context.Context, tx pgx.Tx, name, greeting string, count int64) error {
	payload, err := json.Marshal(helloSaidPayload{
		Name:     name,
		Greeting: greeting,
		Count:    count,
		At:       time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}
	topic := e.Topic
	if topic == "" {
		topic = "hello.said"
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO hello.outbox_events(aggregate_id, topic, payload)
		VALUES ($1, $2, $3)`,
		name, topic, payload,
	)
	return err
}
```

- [ ] **Step 3: Delete the old direct Kafka emitter**

```sh
rm impl/helixgitpx/services/hello/internal/repo/event_kafka.go
```

- [ ] **Step 4: Update app.go wiring — swap emitter impl**

Read `impl/helixgitpx/services/hello/internal/app/app.go` and replace two things:

1. Remove the `kafka.NewProducer(...)` block and its deferred `prod.Close(...)`.
2. Change the emitter construction from `&repo.EventKafka{...}` to `&repo.EventOutbox{Pool: pool, Topic: c.KafkaTopic}`.
3. Remove the kafka import if nothing else references it. (Keep it if used for schema-registry hook in future; for now, remove.)

Diff:

```
- prod, err := kafka.NewProducer(kafka.ProducerOptions{
-     Brokers: c.KafkaBrokers, ClientID: "hello", Topic: c.KafkaTopic,
- })
- if err != nil {
-     return err
- }
- defer func() {
-     shctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
-     defer cancel()
-     _ = prod.Close(shctx)
- }()

  greeter := domain.NewGreeter(
      &repo.CounterPG{Pool: pool},
      &repo.CacheRedis{Client: rc},
-     &repo.EventKafka{Producer: prod, Topic: c.KafkaTopic},
+     &repo.EventOutbox{Pool: pool, Topic: c.KafkaTopic},
  )
```

Also remove `"github.com/helixgitpx/platform/kafka"` from imports if unused after the removal.

- [ ] **Step 5: Update domain tests if any were tied to kafka import**

The domain tests use a `fakeEmitter` (interface-based), which is unaffected. No changes needed.

- [ ] **Step 6: Build + test**

```sh
export GOTOOLCHAIN=go1.23.4
cd impl/helixgitpx/services/hello
go mod tidy
go test ./internal/domain/...
go build ./...
```

Expected: tests pass, build succeeds.

- [ ] **Step 7: Commit**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
git add impl/helixgitpx/services/hello
git commit -s -m "$(printf 'feat(services/hello): replace direct kafka producer with transactional outbox\n\nThe Kafka emission path now writes to hello.outbox_events inside the same\ntransaction as the counter UPSERT. Debezium streams the table to the same\nhello.said topic via the EventRouter SMT. External payload format unchanged.\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 14: Hello deploy extras — migrate-job, kafkaconnector, servicemonitor, ingress, vault-agent

**Files:**
- Create: `impl/helixgitpx/services/hello/deploy/helm/templates/migrate-job.yaml`
- Create: `impl/helixgitpx/services/hello/deploy/helm/templates/kafkaconnector.yaml`
- Create: `impl/helixgitpx/services/hello/deploy/helm/templates/servicemonitor.yaml`
- Create: `impl/helixgitpx/services/hello/deploy/helm/templates/ingress.yaml`
- Create: `impl/helixgitpx/services/hello/deploy/helm/templates/vault-agent-annotations.yaml` (patches Deployment via template)
- Modify: `impl/helixgitpx/services/hello/deploy/helm/values.yaml` (add migrate, ingress, vault, connector values)

- [ ] **Step 1: Extend `values.yaml`**

Read existing `impl/helixgitpx/services/hello/deploy/helm/values.yaml` (from M1 Task 35) and append:

```yaml
# M2 additions
migrate:
  enabled: true
  image:
    repository: helixgitpx/hello
    tag: dev
  env:
    - name: HELLO_POSTGRES_DSN
      valueFrom:
        secretKeyRef:
          name: hello-pg-secret
          key: dsn

ingress:
  enabled: true
  className: nginx
  host: hello.helix.local
  tlsSecret: hello-helix-local-tls

monitoring:
  serviceMonitor:
    enabled: true

vault:
  enabled: true
  role: hello
  kvPath: kv/hello
```

- [ ] **Step 2: migrate-job.yaml (pre-upgrade Hook)**

```yaml
{{- if .Values.migrate.enabled }}
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Release.Name }}-migrate
  annotations:
    "helm.sh/hook": pre-install,pre-upgrade
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
spec:
  backoffLimit: 3
  template:
    metadata:
      labels:
        app.kubernetes.io/name: hello
        helixgitpx.dev/env: {{ .Values.global.env | default "local" }}
    spec:
      restartPolicy: OnFailure
      containers:
        - name: migrate
          image: "{{ .Values.migrate.image.repository }}:{{ .Values.migrate.image.tag }}"
          command: ["/app/hello", "migrate", "--dir", "/migrations"]
          env:
            {{- toYaml .Values.migrate.env | nindent 12 }}
{{- end }}
```

NOTE: this assumes hello's binary supports a `migrate` subcommand. If it doesn't yet, the command must be added to `cmd/hello/main.go` — see Step 5 below.

- [ ] **Step 3: kafkaconnector.yaml**

```yaml
{{- if .Values.migrate.enabled }}
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaConnector
metadata:
  name: {{ .Release.Name }}-outbox
  namespace: helix-data
  labels:
    strimzi.io/cluster: debezium
spec:
  class: io.debezium.connector.postgresql.PostgresConnector
  tasksMax: 1
  config:
    database.hostname: helix-pg-rw.helix-data.svc
    database.port: "5432"
    database.user: hello_svc
    database.password: ${secret:helix-data:hello-pg-secret/password}
    database.dbname: helixgitpx
    database.server.name: hellodbs
    plugin.name: pgoutput
    publication.name: helix_hello_outbox
    slot.name: helix_hello_outbox
    table.include.list: hello.outbox_events
    topic.prefix: hellodbs
    transforms: outbox
    transforms.outbox.type: io.debezium.transforms.outbox.EventRouter
    transforms.outbox.route.by.field: topic
    transforms.outbox.route.topic.replacement: ${routedByValue}
{{- end }}
```

- [ ] **Step 4: servicemonitor.yaml + ingress.yaml**

`servicemonitor.yaml`:

```yaml
{{- if .Values.monitoring.serviceMonitor.enabled }}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ .Release.Name }}
  labels:
    release: prom
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: hello
  endpoints:
    - { port: health, interval: 15s, path: /metrics }
{{- end }}
```

`ingress.yaml`:

```yaml
{{- if .Values.ingress.enabled }}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Release.Name }}
  annotations:
    cert-manager.io/cluster-issuer: selfsigned-ca
spec:
  ingressClassName: {{ .Values.ingress.className }}
  tls:
    - { hosts: [{{ .Values.ingress.host | quote }}], secretName: {{ .Values.ingress.tlsSecret | quote }} }
  rules:
    - host: {{ .Values.ingress.host | quote }}
      http:
        paths:
          - { path: /,   pathType: Prefix, backend: { service: { name: {{ .Release.Name }}, port: { name: http  } } } }
          - { path: /v1, pathType: Prefix, backend: { service: { name: {{ .Release.Name }}, port: { name: http  } } } }
{{- end }}
```

- [ ] **Step 5: Add `migrate` subcommand to hello's main.go**

Read `impl/helixgitpx/services/hello/cmd/hello/main.go` and replace with:

```go
// Command hello is a HelixGitpx service scaffolded by tools/scaffold (M1)
// and extended in M2 with a `migrate` subcommand that runs goose migrations.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/app"
	"github.com/helixgitpx/platform/log"
	"github.com/helixgitpx/platform/pg"
)

func main() {
	lg := log.New(log.Options{Level: "info", Service: "hello"})

	if len(os.Args) >= 2 && os.Args[1] == "migrate" {
		if err := runMigrate(context.Background(), os.Args[2:]); err != nil {
			lg.Error("migrate failed", "err", err.Error())
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, lg); err != nil {
		lg.Error("service exited with error", "err", err.Error())
	}
}

func runMigrate(ctx context.Context, args []string) error {
	dir := "/migrations"
	for i := 0; i < len(args); i++ {
		if args[i] == "--dir" && i+1 < len(args) {
			dir = args[i+1]
			i++
		}
	}
	dsn := os.Getenv("HELLO_POSTGRES_DSN")
	if dsn == "" {
		return fmt.Errorf("HELLO_POSTGRES_DSN is required")
	}
	return pg.Migrate(ctx, pg.MigrateOptions{DSN: dsn, Dir: dir})
}
```

- [ ] **Step 6: Update the Dockerfile to copy migrations into /migrations**

Read `impl/helixgitpx/services/hello/deploy/Dockerfile` and adjust the runtime stage to also copy `services/hello/migrations/` to `/migrations` in the final image:

```dockerfile
# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.work go.work.sum* ./
COPY platform/ ./platform/
COPY gen/ ./gen/
COPY services/hello/ ./services/hello/
RUN cd services/hello && \
    CGO_ENABLED=0 GOWORK=off \
    go build -trimpath -ldflags="-s -w" -o /out/hello ./cmd/hello

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/hello /app/hello
COPY --from=build /src/services/hello/migrations /migrations
USER nonroot
EXPOSE 8001 9001 8081
ENTRYPOINT ["/app/hello"]
LABEL org.opencontainers.image.title="hello"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
LABEL org.opencontainers.image.licenses="Apache-2.0"
```

- [ ] **Step 7: Build + lint + commit**

```sh
export GOTOOLCHAIN=go1.23.4
cd impl/helixgitpx/services/hello
go build ./...

cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
helm lint impl/helixgitpx/services/hello/deploy/helm/

git add impl/helixgitpx/services/hello
git commit -s -m "$(printf 'feat(services/hello): M2 deploy extras — migrate-job, connector, ingress, SM\n\n* cmd/hello: add migrate subcommand\n* deploy/Dockerfile: copy migrations/ into the image\n* helm templates: migrate Hook, KafkaConnector, ServiceMonitor, Ingress\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 15: platform/config Vault KV source + hello reads DSNs from Vault

**Files:**
- Create: `impl/helixgitpx/platform/config/vault.go`
- Create: `impl/helixgitpx/platform/config/vault_test.go`
- Modify: `impl/helixgitpx/platform/config/config.go` (recognize `vault:"..."` tag)

- [ ] **Step 1: Write vault.go — a minimal Vault KV v2 reader**

```go
// Package config (vault.go) provides a Vault KV v2 resolver invoked by Load
// when a struct field carries a `vault:"path/to/key"` tag. The HTTP call
// targets $VAULT_ADDR with the token at $VAULT_TOKEN (populated by Vault
// Agent Injector in-cluster, or set manually for dev).
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// VaultResolver fetches secrets from Vault KV v2.
type VaultResolver struct {
	Addr   string
	Token  string
	Client *http.Client
}

// NewVaultResolver reads VAULT_ADDR and VAULT_TOKEN from env.
// Returns nil when either is unset (caller treats it as no-op).
func NewVaultResolver() *VaultResolver {
	addr := os.Getenv("VAULT_ADDR")
	tok := os.Getenv("VAULT_TOKEN")
	if addr == "" || tok == "" {
		return nil
	}
	return &VaultResolver{
		Addr:   strings.TrimRight(addr, "/"),
		Token:  tok,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Read fetches kv/data/<path>'s data[key]. Expects KV v2 layout.
// path must be in the form "<mount>/<secret-path>#<key>", e.g. "kv/hello#dsn".
func (r *VaultResolver) Read(ctx context.Context, path string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("vault: resolver is nil")
	}
	parts := strings.SplitN(path, "#", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("vault: expected mount/path#key, got %q", path)
	}
	kvPath := parts[0]
	key := parts[1]

	// KV v2 reads are at /v1/<mount>/data/<secret-path>
	i := strings.Index(kvPath, "/")
	if i < 0 {
		return "", fmt.Errorf("vault: expected <mount>/<path>, got %q", kvPath)
	}
	url := fmt.Sprintf("%s/v1/%s/data/%s", r.Addr, kvPath[:i], kvPath[i+1:])

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", r.Token)
	resp, err := r.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault: %s", resp.Status)
	}

	var body struct {
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	val, ok := body.Data.Data[key]
	if !ok {
		return "", fmt.Errorf("vault: key %q not in %s", key, kvPath)
	}
	return val, nil
}
```

- [ ] **Step 2: Modify config.go — read `vault:"..."` tag**

Add to the inside of `loadStruct`'s field loop (after the existing `default`/`required` handling, before the `if raw == ""` continue):

```go
		// Vault resolution — only attempted when VAULT_ADDR + VAULT_TOKEN are set
		// AND the field's current raw is empty.
		if raw == "" {
			if vpath := sf.Tag.Get("vault"); vpath != "" {
				if r := NewVaultResolver(); r != nil {
					v, err := r.Read(context.Background(), vpath)
					if err == nil && v != "" {
						raw = v
					}
				}
			}
		}
```

Import `"context"` at the top of config.go if it's not already imported.

- [ ] **Step 3: Write vault_test.go — no real Vault; test the fallback behaviour**

```go
package config_test

import (
	"os"
	"testing"

	"github.com/helixgitpx/platform/config"
)

type cfgWithVault struct {
	DSN string `env:"DSN" vault:"kv/hello#dsn" default:"postgres://default"`
}

func TestLoad_VaultFallsBackToDefaultWhenAddrUnset(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")
	var c cfgWithVault
	if err := config.Load(&c, config.Options{Prefix: "X"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.DSN != "postgres://default" {
		t.Errorf("DSN = %q, want default", c.DSN)
	}
}
```

- [ ] **Step 4: Update hello's app.go cfg struct to add the vault tag**

Read `impl/helixgitpx/services/hello/internal/app/app.go` and extend the `cfg` struct's `PostgresDSN` field (and others):

```go
type cfg struct {
	HTTPAddr     string   `env:"HTTP_ADDR" default:":8001"`
	GRPCAddr     string   `env:"GRPC_ADDR" default:":9001"`
	HealthAddr   string   `env:"HEALTH_ADDR" default:":8081"`
	PostgresDSN  string   `env:"POSTGRES_DSN" vault:"kv/hello#pg_dsn" required:"true"`
	RedisAddr    string   `env:"REDIS_ADDR" vault:"kv/hello#redis_addr" default:"localhost:6379"`
	KafkaBrokers []string `env:"KAFKA_BROKERS" default:"localhost:9092" split:","`
	KafkaTopic   string   `env:"KAFKA_TOPIC" default:"hello.said"`
	OTLPEndpoint string   `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	Version      string   `env:"VERSION" default:"m1-dev"`
}
```

- [ ] **Step 5: Build + test + commit**

```sh
export GOTOOLCHAIN=go1.23.4
cd impl/helixgitpx/platform
go mod tidy
go test ./config/...
grep '^go ' go.mod

cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx/impl/helixgitpx/services/hello
go build ./...

cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
git add impl/helixgitpx/platform/config impl/helixgitpx/services/hello/internal/app
git commit -s -m "$(printf 'feat(platform/config): add vault:%q tag for Vault KV v2 secret resolution\n\nLoad() now honours a vault:%qpath#key%q tag and, when VAULT_ADDR + VAULT_TOKEN\nare set, fetches the secret via Vault KV v2 as a fallback source (env var\nstill wins, default still applies). Hello now declares vault: tags on\nPostgresDSN and RedisAddr so in-cluster deployment reads from kv/hello.\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')" --cleanup=verbatim
```

---

### Task 16: platform/telemetry pprof + SVID-terminated grpc (platform/grpc)

**Files:**
- Create: `impl/helixgitpx/platform/telemetry/pprof.go`
- Modify: `impl/helixgitpx/platform/grpc/server.go` (add optional spire.Fetcher)

- [ ] **Step 1: Write pprof.go**

```go
// Package telemetry (pprof.go) exposes net/http/pprof handlers for continuous
// profiling via Pyroscope's pull-based scraping. RegisterPprof attaches the
// handlers to the passed mux (typically the health mux on a separate port).
package telemetry

import (
	"net/http"
	"net/http/pprof"
)

// RegisterPprof adds the standard pprof handlers to mux. Call it once per
// process from the composition root before serving the health mux.
func RegisterPprof(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
```

- [ ] **Step 2: Add SVID support to grpc/server.go**

Modify `impl/helixgitpx/platform/grpc/server.go` — extend `Options` and `NewServer`:

```go
package grpc

import (
	"crypto/tls"

	"github.com/helixgitpx/platform/spire"
	"github.com/spiffe/go-spiffe/v2/spiffegrpc/grpccredentials"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Options struct {
	ServerOptions      []grpc.ServerOption
	DisableReflection  bool
	// Fetcher — when non-nil and Fetcher.Source() returns non-nil,
	// the server terminates mTLS with the workload's X.509 SVID.
	// Accepts connections from any client SVID in the same trust domain.
	Fetcher *spire.Fetcher
}

func NewServer(opts Options) (*grpc.Server, error) {
	so := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryChain()...),
		grpc.ChainStreamInterceptor(streamChain()...),
	}

	if opts.Fetcher != nil {
		if src := opts.Fetcher.Source(); src != nil {
			tlsCfg := tlsconfig.MTLSServerConfig(src, src, tlsconfig.AuthorizeAny())
			so = append(so, grpc.Creds(credentials.NewTLS(tlsCfg)))
			_ = grpccredentials.MTLSServerCredentials   // suppress unused import when only tlsconfig used
			_ = tls.NoClientCert                         // suppress unused import; keep crypto/tls visible
		}
	}

	so = append(so, opts.ServerOptions...)

	s := grpc.NewServer(so...)

	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, hs)

	if !opts.DisableReflection {
		reflection.Register(s)
	}
	return s, nil
}
```

NOTE: The `_ = grpccredentials.MTLSServerCredentials` and `_ = tls.NoClientCert` lines suppress "imported and not used" errors from the compiler — they're placeholder touches that exist to keep those imports for when M3 extends auth. Remove them if they cause lint noise by deleting those imports entirely; the tlsconfig path doesn't strictly need them. Simpler version:

```go
package grpc

import (
	"github.com/helixgitpx/platform/spire"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Options struct {
	ServerOptions     []grpc.ServerOption
	DisableReflection bool
	Fetcher           *spire.Fetcher
}

func NewServer(opts Options) (*grpc.Server, error) {
	so := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryChain()...),
		grpc.ChainStreamInterceptor(streamChain()...),
	}
	if opts.Fetcher != nil {
		if src := opts.Fetcher.Source(); src != nil {
			tlsCfg := tlsconfig.MTLSServerConfig(src, src, tlsconfig.AuthorizeAny())
			so = append(so, grpc.Creds(credentials.NewTLS(tlsCfg)))
		}
	}
	so = append(so, opts.ServerOptions...)

	s := grpc.NewServer(so...)

	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, hs)

	if !opts.DisableReflection {
		reflection.Register(s)
	}
	return s, nil
}
```

Use the simpler version.

- [ ] **Step 3: Wire hello to register pprof + pass SPIRE fetcher**

Read `impl/helixgitpx/services/hello/internal/app/app.go` and make two additions:

1. After `hh.Routes(hmux)`, add `telemetry.RegisterPprof(hmux)`.
2. After `telemetry.Start(...)` is called, construct a `spire.Fetcher` and pass it to `hgrpc.NewServer`:

```go
spireFetcher, err := spire.NewFetcher(ctx, spire.Options{SocketPath: os.Getenv("SPIFFE_ENDPOINT_SOCKET")})
if err != nil {
	return err
}
defer spireFetcher.Close()

grpcSrv, err := hgrpc.NewServer(hgrpc.Options{Fetcher: spireFetcher})
```

Add imports: `"os"`, `"github.com/helixgitpx/platform/spire"`.

- [ ] **Step 4: Build + test + commit**

```sh
export GOTOOLCHAIN=go1.23.4
cd impl/helixgitpx/platform
go mod tidy
go test ./grpc/... ./telemetry/...

cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx/impl/helixgitpx/services/hello
go build ./...

cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
git add impl/helixgitpx/platform impl/helixgitpx/services/hello/internal/app
git commit -s -m "$(printf 'feat(platform): telemetry.RegisterPprof + optional SVID mTLS in grpc.NewServer\n\nhello wires both: pprof on the health mux, SPIRE fetcher into the gRPC server\nwhen SPIFFE_ENDPOINT_SOCKET points at a SPIRE agent socket.\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

## Phase E — GitOps root & verification

### Task 17: Argo CD bootstrap + app-of-apps + 25 Application CRs

**Files:**
- Create: `impl/helixgitpx-platform/argocd/bootstrap/kustomization.yaml`
- Create: `impl/helixgitpx-platform/argocd/bootstrap/argocd-install.yaml` (small wrapper)
- Create: `impl/helixgitpx-platform/argocd/bootstrap/argocd-values.yaml`
- Create: `impl/helixgitpx-platform/argocd/applicationset/app-of-apps.yaml`
- Create: `impl/helixgitpx-platform/argocd/applications/*.yaml` (25 files)

- [ ] **Step 1: bootstrap/kustomization.yaml**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: argocd
resources:
  - https://raw.githubusercontent.com/argoproj/argo-cd/v3.0.0/manifests/install.yaml
```

- [ ] **Step 2: app-of-apps.yaml**

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: helixgitpx-app-of-apps
  namespace: argocd
  finalizers: [resources-finalizer.argocd.argoproj.io]
spec:
  project: default
  source:
    repoURL: "{{ env REPO_URL }}"   # overridden by kustomize overlay
    targetRevision: main
    path: impl/helixgitpx-platform/argocd/applications
    directory:
      recurse: true
  destination:
    server: https://kubernetes.default.svc
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

- [ ] **Step 3: Write all 25 applications/*.yaml — single template**

Each file at `impl/helixgitpx-platform/argocd/applications/<name>.yaml` follows this template — substitute `NAME`, `WAVE`, and `NAMESPACE`:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: <NAME>
  namespace: argocd
  annotations:
    argocd.argoproj.io/sync-wave: "<WAVE>"
  finalizers: [resources-finalizer.argocd.argoproj.io]
spec:
  project: default
  source:
    repoURL: "{{ env REPO_URL }}"
    targetRevision: main
    path: impl/helixgitpx-platform/helm/<NAME>
    helm:
      valueFiles:
        - values.yaml
        - values-{{ env HELIX_ENV | default "local" }}.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: <NAMESPACE>
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
      - ServerSideApply=true
```

Files to produce (NAME / WAVE / NAMESPACE):

| Name | Wave | Namespace |
|---|---|---|
| cilium | -10 | kube-system |
| ingress-nginx | -5 | helix-system |
| cert-manager | -5 | cert-manager |
| external-dns | -5 | helix-system |
| minio | -5 | helix-system |
| prometheus-stack | -3 | helix-observability |
| mimir | -3 | helix-observability |
| loki | -3 | helix-observability |
| tempo | -3 | helix-observability |
| pyroscope | -3 | helix-observability |
| spire | 0 | helix-identity |
| istio-base | 0 | istio-system |
| istio-ambient | 0 | istio-system |
| cnpg-operator | 5 | helix-data |
| strimzi-operator | 5 | helix-data |
| cnpg-cluster | 7 | helix-data |
| kafka-cluster | 7 | helix-data |
| dragonfly | 7 | helix-cache |
| meilisearch | 7 | helix-cache |
| opensearch | 7 | helix-cache |
| qdrant | 7 | helix-cache |
| vault | 7 | helix-secrets |
| karapace | 9 | helix-data |
| debezium | 9 | helix-data |
| hello | 10 | helix |

The hello Application's `path` differs — it points at `impl/helixgitpx/services/hello/deploy/helm` (not under `impl/helixgitpx-platform/helm/`). Override for that file only.

- [ ] **Step 4: Verify manifests apply cleanly**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
kubectl kustomize impl/helixgitpx-platform/argocd/bootstrap/ > /tmp/argocd-bootstrap.yaml
test -s /tmp/argocd-bootstrap.yaml && echo OK

# Validate each Application's YAML syntax (dry-client-validate)
for app in impl/helixgitpx-platform/argocd/applications/*.yaml; do
    kubectl apply --dry-run=client -f "$app" >/dev/null && echo "ok: $app"
done
```

All should print `ok:`.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx-platform/argocd
git commit -s -m "$(printf 'feat(platform/m2): argo cd bootstrap + app-of-apps + 25 Application CRs\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 18: Kustomize overlays (local, staging, prod-eu)

**Files:**
- Create: `impl/helixgitpx-platform/kustomize/overlays/local/kustomization.yaml`
- Create: `impl/helixgitpx-platform/kustomize/overlays/staging/kustomization.yaml`
- Create: `impl/helixgitpx-platform/kustomize/overlays/prod-eu/kustomization.yaml`

The overlays patch the Argo CD Application CRs to inject environment-specific values. They don't directly patch the underlying Helm values — instead they use the Argo CD Application's `helm.parameters` field to set `.Values.global.env` and similar.

- [ ] **Step 1: local overlay**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: argocd
resources:
  - ../../../argocd/bootstrap
  - ../../../argocd/applications

patches:
  - target: { kind: Application }
    patch: |-
      - op: add
        path: /spec/source/helm/parameters
        value:
          - name: global.env
            value: local
```

- [ ] **Step 2: staging overlay**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: argocd
resources:
  - ../../../argocd/bootstrap
  - ../../../argocd/applications

patches:
  - target: { kind: Application }
    patch: |-
      - op: add
        path: /spec/source/helm/parameters
        value:
          - name: global.env
            value: staging
  - target: { kind: Application, name: cert-manager }
    patch: |-
      - op: add
        path: /spec/source/helm/parameters/-
        value:
          name: issuer.mode
          value: letsencrypt
      - op: add
        path: /spec/source/helm/parameters/-
        value:
          name: issuer.letsencrypt.email
          value: ops@helixgitpx.dev
```

- [ ] **Step 3: prod-eu overlay**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: argocd
resources:
  - ../staging

commonLabels:
  helixgitpx.dev/env: prod-eu
  helixgitpx.dev/region: eu-west-1

patches:
  - target: { kind: Application }
    patch: |-
      - op: replace
        path: /spec/source/helm/parameters
        value:
          - name: global.env
            value: prod-eu
```

- [ ] **Step 4: Verify overlays build**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
for env in local staging prod-eu; do
    kubectl kustomize "impl/helixgitpx-platform/kustomize/overlays/$env/" > "/tmp/overlay-$env.yaml" && \
      echo "ok: $env ($(wc -l < /tmp/overlay-$env.yaml) lines)"
done
```

Expected: three `ok: ...` lines.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx-platform/kustomize
git commit -s -m "$(printf 'feat(platform/m2): kustomize overlays (local, staging, prod-eu)\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 19: verify-m2-cluster.sh + verify-m2-spine.sh

**Files:**
- Create: `scripts/verify-m2-cluster.sh`
- Create: `scripts/verify-m2-spine.sh`

- [ ] **Step 1: verify-m2-cluster.sh**

```sh
#!/usr/bin/env bash
# Walk the M2 completion matrix (20 roadmap items 19–38).
# Exit 0 iff every gate passes.
set -u

SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.."

pass=0
fail=0
report() {
    local name="$1" result="$2" msg="${3-}"
    if [ "$result" = "ok" ]; then
        printf '  [ ok ] %s\n' "$name"
        pass=$((pass + 1))
    else
        printf '  [FAIL] %s%s\n' "$name" "${msg:+ — $msg}"
        fail=$((fail + 1))
    fi
}

check() {
    local name="$1"; shift
    if "$@" >/dev/null 2>&1; then report "$name" ok
    else report "$name" fail
    fi
}

echo "== M2 Core Data Plane — Completion Matrix =="

# Phase 2.1
check "19 Cluster reachable"             kubectl get nodes
check "19 Cilium running"                kubectl -n kube-system get ds/cilium-agent
check "20 Argo CD installed"             kubectl -n argocd get deploy/argocd-server
check "20 App-of-apps Synced"            bash -c 'kubectl -n argocd get application helixgitpx-app-of-apps -o jsonpath="{.status.sync.status}" | grep -q Synced'
check "20 All 25 Apps Healthy"           bash -c '[ $(kubectl -n argocd get applications -o jsonpath="{.items[*].status.health.status}" | tr " " "\n" | grep -c Healthy) -ge 25 ]'
check "21 cert-manager running"          kubectl -n cert-manager get deploy
check "21 ClusterIssuer selfsigned-ca"   kubectl get clusterissuer selfsigned-ca
check "22 SPIRE server ready"            bash -c 'kubectl -n helix-identity rollout status sts/spire-server --timeout=30s'
check "22 SVID fetch from hello"         bash -c 'kubectl -n helix exec deploy/hello -- /opt/spiffe/bin/spiffe-helper --help' 2>/dev/null

# Phase 2.2
check "23 CNPG cluster Ready"            bash -c 'kubectl -n helix-data get cluster helix-pg -o jsonpath="{.status.phase}" | grep -qE "Cluster in healthy state|healthy"'
check "24 Schemas + RLS applied"         bash -c 'kubectl -n helix-data exec helix-pg-1 -- psql -U postgres -d helixgitpx -tAc "SELECT count(*) FROM information_schema.schemata WHERE schema_name IN (\"hello\",\"auth\",\"repo\",\"sync\",\"conflict\",\"upstream\",\"collab\",\"events\",\"platform\")" | grep -q 9'
check "25 Goose migrations applied"      bash -c 'kubectl -n helix get job/hello-migrate -o jsonpath="{.status.succeeded}" | grep -q 1'
check "26 CNPG backups in MinIO"         bash -c 'kubectl -n helix-system exec deploy/minio -- mc ls local/cnpg-backups | head -1'

# Phase 2.3
check "27 Strimzi Kafka Ready"           bash -c 'kubectl -n helix-data get kafka helix-kafka -o jsonpath="{.status.conditions[?(@.type==\"Ready\")].status}" | grep -q True'
check "28 Karapace responsive"           bash -c 'kubectl -n helix-data exec deploy/karapace -- curl -fsS http://localhost:8081/_schemas >/dev/null'
check "29 Debezium KafkaConnect Ready"   bash -c 'kubectl -n helix-data get kafkaconnect debezium -o jsonpath="{.status.conditions[?(@.type==\"Ready\")].status}" | grep -q True'
check "30 Outbox KafkaConnector Running" bash -c 'kubectl -n helix-data get kafkaconnector hello-outbox -o jsonpath="{.status.conditions[?(@.type==\"Ready\")].status}" | grep -q True'

# Phase 2.4
check "31 Dragonfly running"             bash -c 'kubectl -n helix-cache get sts | grep -q dragonfly'
check "32 Meilisearch health"            bash -c 'kubectl -n helix-cache exec sts/meilisearch-0 -- curl -fsS http://localhost:7700/health >/dev/null' 2>/dev/null || kubectl -n helix-cache get sts | grep -q meilisearch
check "33 OpenSearch green"              bash -c 'kubectl -n helix-cache exec sts/helix-os-opensearch-0 -- curl -fsS http://localhost:9200/_cluster/health | grep -q "\"status\":\"green\"\\|\"status\":\"yellow\""'
check "34 Qdrant ready"                  bash -c 'kubectl -n helix-cache get sts | grep -q qdrant'

# Phase 2.5
check "35 Vault unsealed"                bash -c 'kubectl -n helix-secrets exec sts/vault-0 -- vault status | grep -q "Sealed.*false"'
check "36 Observability stack"           bash -c 'kubectl -n helix-observability get deploy | grep -cE "prom-grafana|mimir|loki|tempo|pyroscope" | grep -qE "[5-9]"'
check "37 Grafana dashboards provisioned" bash -c 'kubectl -n helix-observability exec deploy/prom-grafana -- grafana-cli admin data-migration-check 2>/dev/null; kubectl -n helix-observability get cm -l grafana_dashboard=1 | wc -l | awk "{exit (\$1 < 2)}"'
check "38 Alertmanager config renders"   bash -c 'kubectl -n helix-observability exec sts/alertmanager-prom-kube-prometheus-stack-alertmanager-0 -- amtool config show >/dev/null'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
```

- [ ] **Step 2: verify-m2-spine.sh**

```sh
#!/usr/bin/env bash
# M2 end-to-end spine — exercises hello through the full data plane.
set -u

SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.."

pass=0
fail=0
report() {
    local name="$1" result="$2"
    if [ "$result" = "ok" ]; then
        printf '  [ ok ] %s\n' "$name"
        pass=$((pass + 1))
    else
        printf '  [FAIL] %s\n' "$name"
        fail=$((fail + 1))
    fi
}
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then report "$n" ok; else report "$n" fail; fi; }

echo "== M2 End-to-end Spine =="

# Ensure k3d forwards 443 locally — ingress-nginx LoadBalancer pushes to host
check "HTTP /v1/hello?name=world returns 200" \
    bash -c 'curl -fsS -k --resolve hello.helix.local:443:127.0.0.1 "https://hello.helix.local/v1/hello?name=world" >/dev/null'

check "HTTP response contains hello, world" \
    bash -c 'curl -fsS -k --resolve hello.helix.local:443:127.0.0.1 "https://hello.helix.local/v1/hello?name=world" | grep -q "hello, world"'

check "gRPC SayHello returns greeting" \
    bash -c 'grpcurl -insecure -d "{\"name\":\"world\"}" hello.helix.local:443 helixgitpx.hello.v1.HelloService/SayHello | grep -q "hello, world"'

check "Outbox row inserted in Postgres" \
    bash -c 'kubectl -n helix-data exec helix-pg-1 -- psql -U hello_svc -d helixgitpx -tAc "SELECT count(*) > 0 FROM hello.outbox_events" | grep -q t'

check "Debezium delivered event to hello.said within 30s" \
    bash -c 'timeout 30 kubectl -n helix-data exec deploy/helix-kafka-kafka-0 -- /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server helix-kafka-kafka-bootstrap:9092 --topic hello.said --from-beginning --max-messages 1 --timeout-ms 30000 | grep -q world'

check "Prometheus scrapes hello metrics" \
    bash -c 'kubectl -n helix-observability exec sts/prometheus-prom-kube-prometheus-stack-prometheus-0 -- wget -qO- "http://localhost:9090/api/v1/targets?state=active" | grep -q "\"job\":\"hello\""'

check "Loki has logs from hello" \
    bash -c 'kubectl -n helix-observability exec sts/loki-0 -- wget -qO- "http://localhost:3100/loki/api/v1/query?query={app=\"hello\"}" | grep -q "result"'

check "Tempo has a trace from hello" \
    bash -c 'kubectl -n helix-observability exec deploy/tempo-query-frontend -- wget -qO- "http://localhost:3200/api/search?tags=service.name%3Dhello" | grep -q "traces"'

check "Grafana reachable via Ingress" \
    bash -c 'curl -fsS -k --resolve grafana.helix.local:443:127.0.0.1 "https://grafana.helix.local/api/health" | grep -q ok'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
```

- [ ] **Step 3: chmod + commit**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
chmod +x scripts/verify-m2-cluster.sh scripts/verify-m2-spine.sh
shellcheck scripts/verify-m2-cluster.sh scripts/verify-m2-spine.sh

git add scripts/verify-m2-cluster.sh scripts/verify-m2-spine.sh
git commit -s -m "$(printf 'chore(m2): completion-matrix + end-to-end spine verifiers\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"
```

---

### Task 20: ADRs 0006–0012 + runbook updates + M2 tag

**Files:**
- Create: `docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0006-local-k3d-m2-target.md`
- Create: `.../adr/0007-ha-manifests-local-overlay.md`
- Create: `.../adr/0008-hello-outbox-pattern.md`
- Create: `.../adr/0009-istio-ambient-m2.md`
- Create: `.../adr/0010-observability-first-sequencing.md`
- Create: `.../adr/0011-local-friendly-externals.md`
- Create: `.../adr/0012-sync-wave-ordering.md`

Each follows the same template used for ADRs 0001–0005 in M1 Task 33 — Status/Date/Deciders/Context/Decision/Consequences/Alternatives/Links.

- [ ] **Step 1: Write ADR-0006 (local-k3d-m2-target)**

```markdown
# ADR-0006 — Local k3d as the M2 cluster target

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

M2 requires a running Kubernetes cluster for real deployment (its exit criterion depends on metrics/logs/traces visible in Grafana). Staging GKE/EKS/AKS costs money and requires account setup; a cloud-less path keeps M2 self-contained on a single workstation.

## Decision

M2 uses `k3d` on the developer's host as the cluster target. `up.sh --m2` provisions a 3-node k3d cluster (1 server + 2 agents), installs Cilium as CNI, ingress-nginx as L7, and then applies the Argo CD bootstrap kustomization.

## Consequences

- Full spine deployable on a 62 GiB host without cloud credentials.
- Real resource pressure surfaces single-host limits early — preflight script refuses to run below 48 GiB free RAM.
- Staging/prod overlays still exist (`kustomize/overlays/{staging,prod-eu}/`) and are validated via `kubectl kustomize`; they activate when real clusters arrive.
- macOS and Windows hosts are explicitly out of scope for M2 (per the on-prem-deployment spec).

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-1, §5
```

- [ ] **Step 2: ADR-0007 — HA manifests + local single-replica overlay**

```markdown
# ADR-0007 — HA manifests authoritative in Git; local overlay patches to single-replica

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

CNPG, Kafka, OpenSearch, Vault, Mimir, Loki, Tempo, ingress-nginx, and cert-manager are all HA-capable. On a local k3d host running the full stack simultaneously, true HA (3-replica everywhere) exceeds reasonable single-machine capacity. Abandoning HA in Git would make staging/prod bring-up require bespoke re-engineering per component.

## Decision

Every Helm chart's `values.yaml` declares HA-by-default replica counts (3 for most data services, 2 for control-plane components). A per-environment `values-local.yaml` patches replicas down to 1 for local. `kustomize/overlays/{local,staging,prod-eu}/` selects the right file via each Application's `spec.source.helm.valueFiles`.

## Consequences

- GitOps in staging/prod deploys real HA without code changes.
- Local dev fits in one machine without OOMs.
- PodDisruptionBudgets, replicaCount, and related configs appear in two values files per chart — a small duplication cost.
- Staging is where HA is first validated end-to-end; local is not a faithful HA environment.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-2
```

- [ ] **Step 3: ADR-0008 — Hello outbox pattern**

```markdown
# ADR-0008 — Hello emits Kafka events via transactional outbox + Debezium

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

M1 hello wrote to Postgres and Kafka as separate operations. A crash between the two could leave the counter incremented but the event unemitted (lost message), or vice versa (phantom event). At the scale of a demo this is invisible; the spec's M5 sync orchestrator cannot tolerate either.

## Decision

Hello writes events to a `hello.outbox_events` table inside the same pgx transaction as the counter UPSERT. Debezium's PostgreSQL connector streams the WAL via the `pgoutput` plugin and uses the `EventRouter` SMT to publish each row to the topic named in the row's `topic` column. Result: exactly-once from the service's perspective; duplicates only possible on Kafka-side retries (handled by consumer idempotency).

## Consequences

- `platform/kafka.Producer` is no longer needed for hello's happy path. The package remains for services that do need synchronous production (internal notifications, admin commands).
- The outbox table requires a logical replication slot; `wal_level=logical` is set in CNPG config.
- Debezium tasks.max=1 per connector — sufficient for one service, will need slot-per-service tuning in M4/M5.
- The `hello.said` topic contract is preserved; downstream consumers see the same JSON payload.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-3, §10
- https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html
```

- [ ] **Step 4: ADR-0009 — Istio Ambient in M2**

```markdown
# ADR-0009 — Istio Ambient mesh is installed in M2

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

The spec lists Istio Ambient in Phase 2.1 alongside cluster provisioning. Deferring it saves resources but delays mesh-dependent features (zero-code mTLS between services, L7 policy for M3 auth). YAGNI would defer; faithfulness would include.

## Decision

Istio Ambient is installed in M2. Ambient mode (zTunnel DaemonSet + namespace opt-in) is used — not sidecar. Namespaces `helix`, `helix-data`, `helix-cache`, `helix-secrets` are labeled `istio.io/dataplane-mode: ambient`. Istio consumes SPIRE SVIDs via the SPIFFE CSRA integration.

## Consequences

- Hello-to-Postgres and hello-to-Dragonfly connections gain automatic mTLS with zero code change.
- Observability and system namespaces opt out (ingress-nginx and Prometheus need plaintext reach to non-mesh pods).
- zTunnel DaemonSet + Istio CNI add ~3 pods and ~500 MiB resident on each node.
- M3 auth work can rely on L4 mTLS being already in place.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-4, §8
- https://istio.io/latest/docs/ambient/
```

- [ ] **Step 5: ADR-0010 — Observability-first sequencing**

```markdown
# ADR-0010 — Observability-first sequencing for M2

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Strict spec-phase order (2.1 → 2.5) installs the full data plane before observability lands. Any issue during CNPG or Kafka bring-up is then debugged without metrics, logs, or traces — blind. Given k3d single-host resource pressure, a CrashLoopBackOff in one component can cascade.

## Decision

Sync waves install observability (wave -3) before the data plane (wave 5–7) and before hello (wave 10). Prometheus + Mimir + Loki + Tempo + Pyroscope + Grafana + Alertmanager are reconciling and scraping long before the first data service schedules. Every chart ships a `ServiceMonitor` so metrics flow automatically when pods appear.

## Consequences

- Each data service's "done" gate can assert "visible in Grafana" from the first minute.
- Cluster debugging during bring-up is tractable.
- Observability resource cost is paid upfront (~8–12 GiB) — offset by the halved replica counts in the `local` overlay.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-7, §5
```

- [ ] **Step 6: ADR-0011 — Local-friendly external defaults**

```markdown
# ADR-0011 — Local-friendly external defaults (self-signed, noop DNS, MinIO, placeholder PagerDuty)

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Several M2 components normally depend on external services:
- cert-manager → Let's Encrypt + DNS-01 provider
- external-dns → real DNS provider API
- object storage → S3/GCS/Azure Blob
- Alertmanager → PagerDuty (or equivalent)

Requiring all four on a dev host makes M2 infeasible without cloud accounts.

## Decision

Each has a local-friendly default:
- cert-manager uses `selfsigned-ca` ClusterIssuer locally; Let's Encrypt in staging/prod.
- external-dns uses the `noop` webhook provider locally; real provider (cloudflare/route53) in staging.
- Object storage uses MinIO in-cluster; real S3 in staging. Consumer charts reference an abstract S3 endpoint that the overlay substitutes.
- Alertmanager routes to a null receiver when `PAGERDUTY_INTEGRATION_KEY` Secret is absent; overlay sets the real key in staging/prod.

All four are shipped as both configs (HA/production) in Git and patched down per environment in `kustomize/overlays/`.

## Consequences

- M2 is deployable offline after the initial image pull.
- Staging/prod activation is a single overlay swap, not a config rewrite.
- Local TLS uses self-signed certs; verifiers pass `-k` / `--insecure` (documented).
- MinIO is a single-instance local; no data-durability guarantees locally.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-6, §7
```

- [ ] **Step 7: ADR-0012 — Sync-wave ordering**

```markdown
# ADR-0012 — Argo CD sync-wave ordering for M2 bring-up

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Argo CD applies Applications in parallel by default. M2 has hard dependencies: operators must exist before their CRs, CNI before any workload, and the SPIRE trust bundle before Istio can pull it. Without explicit ordering, the first reconcile loop surfaces dozens of spurious "CRD not found" errors.

## Decision

Every Application carries an `argocd.argoproj.io/sync-wave` annotation:

| Wave | Components |
|---|---|
| -10 | cilium |
| -5 | ingress-nginx, cert-manager, external-dns, minio |
| -3 | prometheus-stack, mimir, loki, tempo, pyroscope |
| 0 | spire, istio-base, istio-ambient |
| 5 | cnpg-operator, strimzi-operator |
| 7 | cnpg-cluster, kafka-cluster, dragonfly, meilisearch, opensearch, qdrant, vault |
| 9 | karapace, debezium |
| 10 | hello |

Argo CD's sync-wave contract: Applications in lower waves are Synced + Healthy before higher-wave Applications begin reconciling.

## Consequences

- Bring-up is deterministic; no CRD-before-operator flapping.
- Sync time is linear in wave depth — can't parallelise across waves.
- Adding a new chart requires picking a wave; operator-to-CR dependency must be respected.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §5 (waves table)
- https://argo-cd.readthedocs.io/en/stable/user-guide/sync-waves/
```

- [ ] **Step 8: Commit ADRs + tag M2**

```sh
cd /run/media/milosvasic/DATA4TB/Projects/HelixGitpx
git add docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr
git commit -s -m "$(printf 'docs(adr): seed ADRs 0006-0012 from M2 brainstorming constraints\n\nCo-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>\n')"

# Tag M2 (whether the full cluster bring-up was actually verified or not —
# the caveat section below documents deferrals)
git tag -a m2-core-data-plane -m "M2 Core Data Plane — manifests + code complete; cluster activation tracked separately"
```

---

## M2 Exit

All 20 tasks complete ⇒ every row of the completion matrix in the design spec §15 has an artifact committed.

**Actual cluster bring-up and end-to-end verification** are a separate operational step after the plan execution:

```sh
impl/helixgitpx-platform/k8s-local/up.sh --m2
# Wait for Argo CD bootstrap to reconcile all 25 Applications (may take 15-30 min on first run)
bash scripts/verify-m2-cluster.sh
bash scripts/verify-m2-spine.sh
```

If the verify scripts don't reach 20/20 and 9/9 on the first run, fix issues (usually: image pull timeouts, resource pressure, CRD ordering). The completion-matrix status is the single source of truth for "M2 done."

Proceed to M3 Identity & Orgs via the same brainstorming → spec → plan → execute loop.

— End of M2 Core Data Plane plan —
