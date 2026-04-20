# HelixGitpx Deployment Cookbook

## 1. Introduction

Recipes for standing up HelixGitpx in different environments, from a
single-node k3d laptop cluster to a multi-region production fleet.

### 1.1 Audience

- Operators bringing up a new environment.
- Consultants automating customer installs.
- Developers reproducing a bug in a fresh cluster.

### 1.2 Recipes

- **Chapter 2:** local laptop via k3d and Argo CD (10 minutes).
- **Chapter 3:** single-region EKS install (1 hour).
- **Chapter 4:** single-region Hetzner install (1 hour).
- **Chapter 5:** single-region bare-metal (kubeadm + Ceph).
- **Chapter 6:** air-gapped install (internal mirrors).
- **Chapter 7:** multi-region active-passive.
- **Chapter 8:** GPU pool for on-prem LLM inference.
- **Chapter 9:** upgrade between minor versions.
- **Chapter 10:** teardown and data export.

### 1.3 Source of truth

Every recipe generates resources from `impl/helixgitpx-platform/` — the
Helm charts, Argo apps, and Kustomize overlays in this repo are the
canonical deployment artefacts. The cookbook just composes them for
specific environments.

### 1.4 Architecture primer

Read
[`01-architecture/02-system-architecture.md`](../../../../docs/specifications/main/main_implementation_material/HelixGitpx/01-architecture/02-system-architecture.md)
before any non-trivial install. Every recipe assumes you already know
the component list and sync-wave ordering.

### 1.5 What's NOT in this cookbook

- Kubernetes basics. We assume you know `kubectl`, Helm, and Argo CD.
- Kernel tuning. If you need `sysctl` on hosts, talk to
  `support@helixgitpx.io` for the current recommendation.
- Customer-specific policy. Bring your OPA overlay; we provide the base
  bundle.

---
