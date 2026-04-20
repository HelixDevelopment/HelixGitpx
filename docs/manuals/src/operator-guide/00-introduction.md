# HelixGitpx Operator Guide

## 1. Introduction

This manual covers operating a HelixGitpx installation in production: day-0
bring-up, day-1 hardening, day-2 daily operations, and disaster recovery.

It assumes you have read the
[Architecture overview](../../../specifications/main/main_implementation_material/HelixGitpx/01-architecture/02-system-architecture.md)
and are comfortable with Kubernetes, Helm, and GitOps.

### 1.1 Audience

- Platform SREs who own the HelixGitpx cluster(s).
- On-call engineers responding to alerts.
- Release engineers coordinating deploys.

### 1.2 What you'll learn

- **Chapter 2:** cluster prerequisites (K8s, networking, storage).
- **Chapter 3:** installing via Argo CD app-of-apps.
- **Chapter 4:** configuring OIDC (Keycloak) and mTLS (SPIFFE/SPIRE).
- **Chapter 5:** SLOs, alerts, and the runbook catalogue.
- **Chapter 6:** backup, restore, and DR drills.
- **Chapter 7:** multi-region active-passive failover.
- **Chapter 8:** performance tuning and capacity planning.
- **Chapter 9:** security hardening and compliance evidence.
- **Chapter 10:** upgrades and rollbacks.

### 1.3 Minimum cluster profile

| Resource | Minimum | Recommended |
|---------|---------|-------------|
| Nodes | 6 (2 control, 4 worker) | 9+ (3 control, 6+ worker) |
| vCPU per worker | 8 | 16 |
| RAM per worker | 32 GiB | 64 GiB |
| Storage class | RWX, snapshots | Ceph RBD or equivalent |
| Ingress | Any L7 with cert-manager | Istio Ambient + gateway API |

### 1.4 Reference implementation

The GitOps source of truth is `impl/helixgitpx-platform/argocd/` — the
app-of-apps file bootstraps every component with sync-wave ordering. Do
not deploy manually; every change goes through GitOps.

---
