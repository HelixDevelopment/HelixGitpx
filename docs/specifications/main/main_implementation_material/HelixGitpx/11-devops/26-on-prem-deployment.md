# 26 — On-Premise & Air-Gapped Deployment

> **Document purpose**: Deploy HelixGitpx **inside the customer's infrastructure** — a VPC, a private data centre, or a fully air-gapped environment — with the same guarantees as the SaaS offering.

---

## 1. Deployment Topologies

| Topology | Description | Typical customer |
|---|---|---|
| **Managed Dedicated** | HelixGitpx runs the stack in the customer's chosen cloud; we operate it. | Regulated SaaS, enterprises |
| **Customer-Managed** | Customer's SRE runs the stack; we support and provide updates. | Large enterprises with platform teams |
| **Air-Gapped** | No outbound internet; artefacts + updates shipped via signed bundles. | Defence, intelligence, critical infrastructure |
| **Hybrid** | Some planes in cloud, others on-prem (e.g. compute on-prem, object store in cloud). | Complex regulatory regimes |

---

## 2. Sizing

### 2.1 Minimal (pilot — up to 50 users / 500 repos)

- **Kubernetes**: 3 nodes × (8 vCPU, 32 GB RAM, 500 GB NVMe).
- **Postgres**: 1 primary + 1 replica, 8 vCPU / 32 GB / 500 GB SSD each.
- **Kafka**: 3 brokers × (4 vCPU / 16 GB / 500 GB SSD).
- **Object store**: 2 TB.
- **GPU (optional AI)**: 1 × NVIDIA L4 or equivalent.

### 2.2 Standard (1000 users / 10 k repos)

- K8s: 6 app + 3 data + 2 GPU nodes.
- PG: HA cluster (3 instances, each 16 vCPU / 64 GB / 2 TB).
- Kafka: 5 brokers (8 vCPU / 32 GB / 2 TB).
- Object store: 20 TB replicated.

### 2.3 Large (10 000 users / 100 k repos)

- K8s: 30+ nodes auto-scaled.
- PG: 5-instance HA + read replicas per region.
- Kafka: 9+ brokers with tiered storage.
- Object store: 200 TB + erasure coded.
- GPU: fleet of 4× A100/L40S per region.

---

## 3. Prerequisites

### 3.1 Infrastructure

- Kubernetes 1.29+ (vanilla, OpenShift, EKS/GKE/AKS, or Talos Linux on bare metal).
- Container runtime: containerd or CRI-O.
- Ingress: Envoy Gateway (bundled) or compatible (NGINX, Traefik, F5).
- Storage classes:
  - `fast-ssd` (IOPS ≥ 10 k, for PG/Kafka).
  - `standard-ssd` (for everything else).
  - Object storage: S3-compatible (MinIO, Ceph RGW, NetApp StorageGRID, Dell ECS, AWS S3, Azure Blob via MinIO gateway).
- DNS: ability to create at least 5 A/AAAA records + one wildcard for the deployment domain.
- TLS: ACME issuer with HTTP-01 or DNS-01, or existing enterprise PKI.
- NTP: accurate within 10 ms.

### 3.2 Identity

Choose one:
- **OIDC** (Okta, Ping, Azure AD, Keycloak, Authentik, Google Workspace…).
- **SAML 2.0** (for SSO-only shops).
- **LDAP / AD** (via bundled bridge).

### 3.3 Networking

- Egress allowlists (unless air-gapped):
  - `*.githubusercontent.com`, `api.github.com` (for GitHub upstream).
  - Same per upstream provider.
  - Only if AI-cloud is enabled: the configured LLM provider endpoints.
- No inbound from public internet unless ingress/edge is in scope.

---

## 4. Installation

### 4.1 Online (internet-connected)

Using our Helm umbrella chart:

```bash
helm repo add helixgitpx https://charts.helixgitpx.example.com
helm repo update

# Render values from our wizard
helixctl platform values-wizard --output values.yaml

# Install
helm install helixgitpx helixgitpx/helixgitpx \
  --namespace helixgitpx --create-namespace \
  --values values.yaml \
  --version 1.0.0

# Verify
helixctl platform status
```

Or, preferred: **Argo CD**. Fork our `helixgitpx-platform-template` repo, adjust values, let Argo CD reconcile.

### 4.2 Air-Gapped

1. On an internet-connected machine, download the offline bundle:

   ```bash
   helixctl platform download-bundle \
     --version 1.0.0 \
     --arch linux/amd64 \
     --include ai-models \
     --output helixgitpx-1.0.0-bundle.tar
   ```

2. Transfer via approved mechanism (signed USB, diode, etc.).
3. On-prem: import into the internal OCI registry and chart museum.

   ```bash
   helixctl platform load-bundle \
     --file helixgitpx-1.0.0-bundle.tar \
     --registry registry.internal:5000 \
     --chartmuseum https://charts.internal
   ```

4. Install normally with `values.airgapped.yaml` (rewrites image refs to the internal registry; disables all outbound calls).

Bundle contents:

- All container images (Cosign-signed; verified on load).
- Helm charts.
- CRDs, Kyverno policies.
- LLM model weights (optional subset).
- OPA bundles.
- Documentation snapshot.
- SBOM and SLSA provenance for every artefact.

Updates: monthly signed patch bundles; auto-apply in staging, manual approval for prod.

---

## 5. AI in On-Prem / Air-Gapped

- **Default: fully local**. vLLM or Ollama on GPU nodes; models shipped with the bundle.
- Shipped models (options): Llama 3.x (Meta community license), Mistral/Codestral, Qwen2.5-Coder, DeepSeek-Coder, StarCoder, Phi-3.
- Embedding model: BGE, E5, nomic-embed.
- Self-learning pipeline runs entirely in-cluster using customer's own feedback (never leaves).
- For air-gapped: base model updates via signed bundles only.

Customers with extreme sensitivity (or model-ban policies) can disable AI entirely:

```yaml
services:
  aiService:
    enabled: false
featureFlags:
  aiConflictAutoApply: false
```

The rest of HelixGitpx functions normally — auto-resolution falls back to policy + CRDT only.

---

## 6. Observability Export

On-prem deployments almost always already have a monitoring stack. HelixGitpx integrates:

- **Metrics**: Prometheus remote-write to customer's Mimir / Thanos / VictoriaMetrics.
- **Logs**: Loki or OpenSearch; alternatively, stream to customer's SIEM (Splunk, Elastic, Chronicle).
- **Traces**: OTLP to customer's collector (Datadog / New Relic / Jaeger / Tempo).
- **Dashboards**: provided as Grafana JSON; import to customer's Grafana.

All exports tagged with tenant labels for segregation.

---

## 7. Backups

- **Postgres** PITR to customer's object store (default `backups/postgres/`).
- **Kafka** tiered storage + daily metadata snapshot.
- **Vault** Raft snapshot daily.
- **OpenSearch** daily snapshots.
- Retention per customer policy (default 30 d hot, 1 y cold).
- **Restore drills**: monthly-automated; quarterly-rehearsed.

---

## 8. Disaster Recovery

| Scenario | Strategy |
|---|---|
| Single node failure | K8s reschedules; no action |
| AZ failure | Auto-failover for stateful services |
| Region failure (multi-region deployment) | Promote secondary (see 11-devops/16) |
| Cluster corruption | Restore from backup + Kafka replay |
| Software-level corruption | `helixctl fsck --deep` + event replay into fresh projectors |

RPO/RTO defaults:
- Postgres: RPO ≤ 30 s, RTO ≤ 15 min.
- Event store: RPO ≤ 30 s, RTO ≤ 5 min.
- Object store: RPO depends on customer's storage backend.

Customers can tune aggressiveness (sync-replication adds latency, async reduces RPO risk).

---

## 9. Identity Integration

- **OIDC**: register HelixGitpx as an application; map groups to HelixGitpx teams via claim rules.
- **SAML 2.0**: SP metadata export from HelixGitpx; IdP config template provided per major vendor.
- **LDAP/AD bridge**: the bundled `helixgitpx-dirsync` syncs users/groups hourly; handles disabled-user propagation.
- **MFA**: preserved end-to-end; we never bypass customer's MFA.

Break-glass account: a locally-provisioned admin for cases when the IdP is unavailable. Secured by Shamir-split recovery codes.

---

## 10. Compliance Artefacts Provided

- **SOC 2 Type II** letter (annual).
- **ISO 27001** certification + Statement of Applicability.
- **SLSA L3** provenance for every artefact.
- **SBOM** (CycloneDX 1.5) per image.
- **FIPS 140-2/3** crypto mode (toggleable): uses FIPS-validated OpenSSL; disables non-approved ciphers.
- **CMMC L3** readiness guide (defence customers).
- **PCI DSS** scope reduction doc (no card data in scope).

---

## 11. Operational Handoff

For customer-managed deployments we provide:

- **Runbook**: full on-prem operations playbook (superset of the SaaS runbook, with extra sections for IaC, update bundles, DR).
- **Training**: 4 × 2-hour sessions.
- **Office hours**: weekly for the first 90 days.
- **Tiered support**: 24/7 P1, business-hours P2.

---

## 12. Update Cadence

- Stable line: monthly.
- LTS line: quarterly, supported for 18 months.
- Security patches: out-of-band, 7-day remediation SLA for Critical CVEs.

All updates:
- Tested against the customer's staging instance first.
- Signed + reproducible.
- Rollback path documented.

---

## 13. Licensing

- **Per-seat** for Enterprise customers.
- Offline license verification (signed license blob with expiry).
- Soft grace period on expiry (read-only after 30 days).
- Feature tiers expressed as license capabilities; platform respects them.

---

## 14. Security Hardening for On-Prem

- **CIS Kubernetes benchmark** — automation checks via kube-bench; our defaults pass Level 2.
- **FIPS mode** — optional.
- **Strict PSS** (Pod Security Standards `restricted` profile) enforced by Kyverno.
- **No root, no privileged** — verified at admission.
- **Read-only root fs** — default for all our services.
- **Network zero-trust** — default-deny NetworkPolicies + Cilium L7 policies (see `18-manifests/network-policy-samples.yaml`).

---

## 15. Sample `values.on-prem.yaml` (highlights)

```yaml
global:
  env: onprem
  domain: helixgitpx.corp.example.com
  imageRegistry: registry.internal:5000/helixgitpx
  features:
    aiCloudAllowed: false
    egressInternet: false

auth:
  oidc:
    enabled: true
    issuerUrl: "https://sso.corp.example.com"
    clientId: "helixgitpx"
    clientSecretRef:
      name: "oidc-client-secret"
      key: "secret"
    groupsClaim: "groups"
    requireMfa: true

postgres:
  backup:
    destinationPath: "s3://backups.internal/helixgitpx/pg/"
    encryptionKeyRef: kms://keyring-internal/key-pg

kafka:
  tieredStorage:
    enabled: true
    destination: "s3://backups.internal/helixgitpx/kafka/"

observability:
  prometheus:
    remoteWrite:
      - url: "https://mimir.internal/api/v1/push"
  loki:
    gatewayUrl: "https://logs.internal"
  tempo:
    gatewayUrl: "https://traces.internal"

compliance:
  fipsMode: true
  auditRetentionYears: 7
```

---

## 16. Going Live

- [ ] Installation verified with `helixctl platform status`.
- [ ] Smoke test suite passes (`helixctl smoke run --env production`).
- [ ] DNS records cut over.
- [ ] IdP application published.
- [ ] Backups verified (restore test completed).
- [ ] Monitoring dashboards imported.
- [ ] On-call rotation configured.
- [ ] First users invited.
- [ ] Migration playbook executed (see [25-migration-guide.md]).

---

*— End of On-Prem & Air-Gapped Deployment —*
