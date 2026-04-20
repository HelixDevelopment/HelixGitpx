# Script — 11 Day-0 cluster bring-up (GitOps)

**Track:** Operators · **Length:** 15 min · **Goal:** viewer brings up a full HelixGitpx cluster from scratch.

## Body

1. **Prereqs** — 0:30 – 2:00.
   6-node K8s, StorageClass, LoadBalancer, DNS, KMS.
2. **Bootstrap Argo CD** — 2:00 – 3:30.
   Helm install argocd; apply the root Application.
3. **Watch the app-of-apps sync** — 3:30 – 6:30.
   Sync-waves: Cilium → cert-manager → SPIRE → CNPG → Strimzi → … → services.
4. **Apply seeds** — 6:30 – 8:30.
   Keycloak realm import; initial OPA bundle; first org.
5. **Smoke test** — 8:30 – 11:00.
   `helixgitpx login`, create repo, bind upstream, push, verify fanout.
6. **Observability** — 11:00 – 13:00.
   Grafana dashboards per service; Tempo traces; Loki logs.
7. **Alerts + runbook smoke** — 13:00 – 14:30.
   Trigger a synthetic broker kill; alert fires; runbook link.

## Wrap-up (14:30 – 15:00)
"Cluster to production traffic in under 2 hours."

## Companion doc
`docs/manuals/src/deployment-cookbook/` · `impl/helixgitpx-platform/argocd/`
