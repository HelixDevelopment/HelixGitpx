## 3. Argo CD bootstrap

HelixGitpx is deployed exclusively via GitOps. All 53 applications live
under `impl/helixgitpx-platform/argocd/applications/`. This chapter walks
the first-time bring-up.

### 3.1 Install Argo CD

```bash
kubectl create namespace argocd
helm repo add argo https://argoproj.github.io/argo-helm
helm install argocd argo/argo-cd -n argocd --version 7.6.12 \
    -f impl/helixgitpx-platform/helm/argocd/values.yaml
kubectl -n argocd wait deploy/argocd-server --for=condition=Available --timeout=5m
```

### 3.2 Apply the app-of-apps

```bash
kubectl apply -f impl/helixgitpx-platform/argocd/app-of-apps.yaml
```

This registers the root Application. Argo CD discovers every child
Application under `impl/helixgitpx-platform/argocd/applications/` and
syncs them in wave order (see §3.3).

### 3.3 Sync-wave order

| Wave | Components | Why first |
|-----:|-----------|-----------|
| −10  | Cilium | Cluster can't talk to itself without a CNI. |
| −5   | cert-manager | Every subsequent chart needs TLS. |
|  0   | SPIRE, Keycloak, MinIO, observability stack (Loki/Tempo/Mimir) | Platform prerequisites. |
|  1-3 | CNPG + Strimzi + Dragonfly + Vault + OPA bundle server | Data plane. |
|  5-7 | Debezium, Karapace, search stack (Meili/Qdrant/OpenSearch/Zoekt) | Eventing + search. |
|  7   | Ollama, vLLM, NeMo Guardrails, mirrormaker2 | AI + cross-region. |
|  10  | HelixGitpx services (hello, auth, orgteam, repo, …) + website + docs-site | Workloads; depend on everything above. |

### 3.4 Verify sync

```bash
kubectl -n argocd get applications
argocd app list | grep -vi healthy    # anything non-healthy is a red flag
```

Expected: 53 applications, every one `Synced + Healthy`. If a chart
stays `Missing`, run `scripts/verify-argo-paths.sh` locally — a broken
path is the most common cause.

### 3.5 Known bring-up pitfalls

- **CNPG initializes slowly.** First `Cluster` can take 3-5 min. That's
  normal; don't re-apply while it's initialising.
- **Strimzi Kafka** requires a working StorageClass with RWO + snapshot.
  If snapshots unsupported, backups will be disabled.
- **SPIRE + Istio Ambient** can race. If workloads can't get SVIDs,
  cordon the affected node and re-start the spire-agent DaemonSet.

---
