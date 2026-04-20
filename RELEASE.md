# HelixGitpx v1.0.0 — General Availability

**Release date:** 2026-<GA-DATE>
**Status:** GA — production-ready.

## Highlights

HelixGitpx is the first open federated Git proxy. One namespace, many Git hosts —
users work in a single UI while pushes, PRs, and issues fan out across GitHub,
GitLab, Gitea, Gitee, GitFlic, GitVerse, Bitbucket, Forgejo, SourceHut, Azure
DevOps, AWS CodeCommit, and any generic Git host over HTTPS.

### What you get at GA

- **Federation** — 12 provider adapters with webhook-gated bidirectional sync.
- **Conflict resolution** — CRDT-backed collaboration + AI-assisted conflict
  proposals; human review gate mandatory.
- **AI everywhere** — PR summaries, label suggestions, ChatOps, and semantic +
  code search. Self-learning via DPO (opt-in).
- **Policy-as-code** — OPA bundle v2; every enforcement point rendered in Rego
  with a diff-review CI gate.
- **Multi-region** — active-passive eu-central-1 → eu-west-2 with
  MirrorMaker 2 + CNPG logical replication.
- **Clients** — Web (Angular), Android, iOS, Windows, macOS, Linux desktop.

### Platform

- **K8s:** 1.31, Istio Ambient, Argo CD, Temporal, SPIFFE/SPIRE, OpenTelemetry.
- **Data:** Postgres 16 + Timescale, Dragonfly (Redis), Kafka (Strimzi),
  Meilisearch + Qdrant + OpenSearch + Zoekt.
- **LLM:** LiteLLM router, Ollama (local), vLLM (GPU pool), NeMo Guardrails.

## Upgrade path

This is the initial GA. No upgrade from pre-release tags is supported; start
from v1.0.0.

## Deprecations

None at GA.

## Known limitations

- Active-active multi-region is post-GA (Year 2).
- SOC 2 Type II and ISO 27001 certification are post-GA; SOC 2 Type I attestation
  and ISO 27001 gap analysis ship with GA.
- APAC region deployment is post-GA.
- iOS push via APNs behind TestFlight-only flag at GA day.

## Credits

- Maintainer: Милош Васић.
- External contributors listed in `CONTRIBUTORS.md`.

## License

Apache-2.0 (code), CC-BY-SA-4.0 (docs).
