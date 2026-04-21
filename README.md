# HelixGitpx

**Helix Git Proxy eXtended** — a federated Git proxy that mirrors a single source of truth across multiple upstream Git hosts (GitHub, GitLab, GitFlic, GitVerse, Gitea, Gitee, Bitbucket, Azure DevOps, AWS CodeCommit, Forgejo, SourceHut, and generic Git-over-HTTPS) and resolves the inevitable conflicts with policy- and AI-assisted flows.

**Status:** v1.0.0 GA. Milestones `m1-foundation` through `m8-ga` tagged.

## Governance

| If you want to… | Read |
|---|---|
| The supreme policy | [`CONSTITUTION.md`](CONSTITUTION.md) |
| Rules for AI contributors | [`AGENTS.md`](AGENTS.md) |
| Orientation for Claude | [`CLAUDE.md`](CLAUDE.md) |
| How to contribute | [`CONTRIBUTING.md`](CONTRIBUTING.md) |

The Constitution is the load-bearing document. The mandatory testing matrix (seven test types, 100 % coverage per type per module touched, mocks only in unit tests, no skipped tests ever) is Article II.

## Where to go next

| If you want to… | Read |
|---|---|
| Understand what HelixGitpx is | [`docs/specifications/main/main_implementation_material/HelixGitpx/00-core/01-vision-scope-constraints.md`](docs/specifications/main/main_implementation_material/HelixGitpx/00-core/01-vision-scope-constraints.md) |
| Browse the full spec suite | [`docs/specifications/main/main_implementation_material/HelixGitpx/README.md`](docs/specifications/main/main_implementation_material/HelixGitpx/README.md) |
| Read it as a website | `make docs && open http://localhost:3001` |
| Start hacking | [`impl/helixgitpx/README.md`](impl/helixgitpx/README.md) |
| Understand the roadmap | [`docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md`](docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md) |
| See milestone plans | [`docs/superpowers/plans/`](docs/superpowers/plans/) |
| User manuals (all formats) | [`docs/manuals/`](docs/manuals/) |
| Video curriculum | [`docs/media/`](docs/media/) |
| Integration plans | [`docs/integrations/`](docs/integrations/) |
| **What is NOT finished (and why)** | [`docs/UNFINISHED.md`](docs/UNFINISHED.md) |

## Quick start

```sh
mise install                    # pin toolchain per mise.toml
make bootstrap                  # fetch deps for every sub-project
make dev                        # bring up compose stack + hello service
curl "http://localhost:8001/v1/hello?name=world"
```

## Test matrix (mandatory)

```bash
make test-all         # seven types; CI enforces
make coverage-audit   # per-package coverage report
make upstream-push    # push main + tags to every configured upstream
```

Individual types: `make test-unit`, `test-integration`, `test-e2e`, `test-security`, `test-stress`, `test-ddos`, `test-benchmark`.

## One-shot green suite

```bash
bash scripts/verify-everything.sh
```

Runs every artifact verifier (M1-M8 artifacts, Argo paths, Helm chart lint,
Rego syntax) plus `go vet` + `go test` across the workspace. Cluster-probe
verifiers short-circuit cleanly when no cluster is reachable.

## Repository layout

- `docs/` — specifications, integrations, manuals, media, operations, security, superpowers plans.
- `impl/helixgitpx/` — Go monorepo (platform + 18 services + gen + tools/scaffold).
- `impl/helixgitpx-web/` — Angular 19 + Nx web app.
- `impl/helixgitpx-clients/` — KMP + Compose shells (Android/iOS/Desktop).
- `impl/helixgitpx-platform/` — Helm charts, Argo CD apps, Kustomize overlays, SQL, OPA.
- `impl/helixgitpx-docs-site/` — Docusaurus public docs (`docs.helixgitpx.io`).
- `impl/helixgitpx-website/` — Astro marketing site (`helixgitpx.io`).
- `test/` — cross-service integration / e2e / security / stress / ddos / benchmark / chaos suites.
- `tools/` — coverage-audit, perf, fuzz, chaos, dr, docs-export.
- `scripts/` — milestone verifiers + upstream sync.
- `Upstreams/` — per-upstream scripts configuring this repo's federation targets.
- `.github/workflows/` — CI pipelines (all `workflow_dispatch`-only; ADR-0001).

## License

Apache-2.0 (code) / CC-BY-SA-4.0 (documentation). See `LICENSE`.
