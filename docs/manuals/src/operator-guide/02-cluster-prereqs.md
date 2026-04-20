## 2. Cluster prerequisites

Before you install HelixGitpx, the target cluster must satisfy the
minimum profile below. Deviating invites grief we cannot debug remotely.

### 2.1 Kubernetes version

- **Minimum:** 1.29.
- **Recommended:** 1.31 (the version Argo CD and Istio Ambient are tested
  against in our CI).
- **Flavours tested:** EKS 1.31, GKE 1.31, kubeadm 1.31 on Ubuntu 24.04,
  k3d / kind 1.31 for local dev.

### 2.2 Node profile

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| Worker nodes | 4 | 6+ |
| vCPU per worker | 8 | 16 |
| RAM per worker | 32 GiB | 64 GiB |
| Disk per worker | 200 GiB NVMe | 500 GiB NVMe |
| Node pools | 1 mixed | 1 general + 1 AI/GPU + 1 ingress |

AI features (Ollama, vLLM) target a dedicated GPU pool in Recommended.
GA supports CPU-only LLM backends as a fallback for clusters without GPUs.

### 2.3 Networking

- **CNI:** Cilium ≥ 1.16 in Hubble-enabled mode. The Argo app-of-apps
  installs Cilium in sync-wave `-10`.
- **Pod CIDR:** /16 minimum; /14 recommended for multi-region expansion.
- **Service mesh:** Istio Ambient. Optional but the default in our Helm
  values. Sidecar mode also works.
- **Ingress:** L7 that supports Gateway API v1, or nginx-ingress. Must
  support cert-manager annotations.

### 2.4 Storage

- **StorageClass:** a default RWX class with volume expansion enabled.
- **Snapshots:** a VolumeSnapshotClass for CNPG backups.
- **Capacity:** allocate 500 GiB per Postgres cluster + 2 TiB per Kafka
  cluster + 1 TiB MinIO per region (Git blob storage).

### 2.5 DNS and TLS

- A public zone you control for production (`*.helixgitpx.io`).
- cert-manager issuers: Let's Encrypt production (DNS-01 or HTTP-01),
  and a private CA for east-west mTLS (SPIFFE trust domain).

### 2.6 Identity

- OIDC-capable identity source (GA ships its own Keycloak; can federate
  to Okta/Google/Azure AD).
- Workload identity: SPIRE installed; every HelixGitpx service
  authenticates east-west via SVIDs.

### 2.7 Observability stack

- Prometheus + Alertmanager (scrapes HelixGitpx + Strimzi + CNPG).
- Loki for logs, Tempo for traces, Pyroscope for profiles, Grafana for
  dashboards.
- All five install via the `prometheus-stack`, `loki`, `tempo`,
  `pyroscope`, and Grafana-as-part-of-prometheus-stack charts under
  `impl/helixgitpx-platform/helm/`.

### 2.8 KMS

- A KMS or HSM the cluster can reach for at-rest envelope encryption.
  Supported: AWS KMS, GCP KMS, HashiCorp Vault Transit, self-hosted KMS
  with PKCS#11.

### 2.9 Pre-install checklist

- [ ] `kubectl get nodes` shows ≥ 4 Ready workers matching the profile.
- [ ] A default StorageClass exists and supports expansion.
- [ ] The ingress controller is installed and a Layer-7 LoadBalancer IP
      is allocated.
- [ ] `kubectl get crd certificates.cert-manager.io` exists.
- [ ] Public DNS for your base domain resolves to the ingress LB.
- [ ] KMS credentials are installed as a cluster secret.

If any checkbox fails, don't proceed; see
[Troubleshooting §10](../troubleshooting/00-introduction.md) for
diagnostic steps.

---
