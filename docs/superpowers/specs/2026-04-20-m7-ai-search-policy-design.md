# M7 AI, Search, Policy — Design Spec

| Field | Value |
|---|---|
| Status | APPROVED |
| Milestone | M7 — AI, Search, Policy (Weeks 37-42) |
| Scope | Roadmap §8 items 116-137 (22 items) |

## Locked constraints

- **C-1 ai-service** — Go service wrapping LiteLLM router (as HTTP client). Routes prompts to Ollama (local), vLLM (GPU staging), or external providers (OpenAI/Anthropic) based on repo config.
- **C-2 Guardrails** — NVIDIA NeMo Guardrails via its Python service running as sidecar or dedicated pod; ai-service proxies through it.
- **C-3 Self-learning** — DPO pipeline on Ray Cluster (GPU pool); feedback collected via `/api/ai/feedback`, curated (PII scrubbed), LoRA adapters promoted via shadow eval.
- **C-4 OPA bundle server** — dedicated `opa-bundles-server` service (not just ConfigMap from M3); Kyverno admission for cluster.
- **C-5 Hybrid search** — fan-out to Meilisearch + Qdrant + OpenSearch; reciprocal-rank fusion; code search via Zoekt.

## New services

- **ai-service** — gRPC `AIService.{Summarize,ProposeConflict,SuggestLabel,Search}` + REST for ChatOps.
- **search-service** — thin aggregator that fans a query to Meilisearch (fulltext), Qdrant (vector), OpenSearch (filters), Zoekt (code); merges via RRF.
- **opa-bundle-server** — hosts OPA bundles with ETag; services pull via `httpbundle` plugin.
- **zoekt-indexer** — CronJob that indexes repo content.

## 22-item matrix

Items 116-120 (AI platform), 121-125 (use cases), 126-130 (self-learning), 131-134 (policy-as-code), 135-137 (hybrid search).

## ADRs 0030-0034

- 0030 — LiteLLM as router (pluggable providers)
- 0031 — NeMo Guardrails as out-of-process policy layer
- 0032 — DPO via Ray for self-learning
- 0033 — OPA bundle server over ConfigMap (extensibility)
- 0034 — Zoekt for code search (not Sourcegraph)

— End of M7 design —
