# Script — 07 Service architecture tour

**Track:** Developers · **Length:** 10 min · **Goal:** viewer can point at any HelixGitpx component and explain its role.

## Cold open (0:00 – 0:15)
C4 L1 diagram appears. "Every line here is a real RPC. Let's walk them."

## Body

1. **Ingress edge** — 0:15 – 1:30.
   Istio Ambient gateway, mTLS termination, rate-limiter, OPA decision.
2. **Identity layer** — 1:30 – 2:30.
   Keycloak OIDC → auth-service → JWT + SPIFFE handoff.
3. **Tenancy layer** — 2:30 – 3:30.
   orgteam-service, repo-service, outbox pattern into Kafka.
4. **Federation layer** — 3:30 – 5:00.
   adapter-pool (12 providers), webhook-gateway, sync-orchestrator (Temporal), conflict-resolver.
5. **Realtime layer** — 5:00 – 6:00.
   live-events-service (gRPC streaming + SSE + WS fallback).
6. **AI + search** — 6:00 – 7:30.
   ai-service, LiteLLM, NeMo Guardrails, search-service with RRF.
7. **Data plane** — 7:30 – 8:30.
   Postgres (CNPG), Kafka (Strimzi), Dragonfly, MinIO, Meilisearch, Qdrant, OpenSearch, Zoekt.
8. **Observability** — 8:30 – 9:30.
   OTel → Tempo/Loki/Mimir/Pyroscope.

## Wrap-up (9:30 – 10:00)
"Every service is on the diagram for a reason. Your job is to know why."

## Companion doc
`docs/specifications/main/main_implementation_material/HelixGitpx/01-architecture/02-system-architecture.md`
