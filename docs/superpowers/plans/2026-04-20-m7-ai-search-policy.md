# M7 Implementation Plan (compact)

Tasks:
- T1 — ai-service scaffold + proto (AIService) + LiteLLM client stub
- T2 — Ollama + vLLM helm charts (infra) + NeMo Guardrails helm
- T3 — AI use-case endpoints: Summarize, ProposeConflict, SuggestLabel, Search, ChatOps
- T4 — Self-learning pipeline scaffolding (feedback endpoint, curator, DPO Ray stub, LoRA promotion)
- T5 — opa-bundle-server service scaffold + Kyverno admission bundle manifests
- T6 — Rego rules bundle (bundles/v2/ with enforcement rules)
- T7 — search-service scaffold + RRF aggregator
- T8 — Zoekt indexer CronJob helm chart
- T9 — 4 new Argo CD Applications (ai, search, opa-bundle-server, zoekt)
- T10 — ADRs 0030-0034 + m7-ai tag
