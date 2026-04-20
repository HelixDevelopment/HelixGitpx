# HelixGitpx

**Helix Git Proxy eXtended** — a federated, privacy-preserving Git proxy that mirrors a single source of truth across multiple upstream Git hosts (GitHub, GitLab, GitFlic, GitVerse, …) and resolves the inevitable conflicts.

## Where to go next

| If you want to… | Read |
|---|---|
| Understand what HelixGitpx is | [`docs/specifications/main/main_implementation_material/HelixGitpx/00-core/01-vision-scope-constraints.md`](docs/specifications/main/main_implementation_material/HelixGitpx/00-core/01-vision-scope-constraints.md) |
| Browse the full spec suite | [`docs/specifications/main/main_implementation_material/HelixGitpx/README.md`](docs/specifications/main/main_implementation_material/HelixGitpx/README.md) |
| Read it as a website | `make docs && open http://localhost:3001` |
| Start hacking | [`impl/helixgitpx/README.md`](impl/helixgitpx/README.md) |
| Understand the roadmap | [`docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md`](docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md) |
| See current milestone plans | [`docs/superpowers/plans/`](docs/superpowers/plans/) |

## Quick start

```sh
mise install                    # pin toolchain per mise.toml
make bootstrap                  # fetch deps for every sub-project
make dev                        # bring up compose stack + hello service
curl "http://localhost:8001/v1/hello?name=world"
```

## Repository layout

- `docs/` — authoritative specifications (immutable input; edit via PR with spec review)
- `impl/` — implementation sub-projects (Go monorepo, Angular web, KMP clients, GitOps, Docusaurus)
- `Upstreams/` — scripts configuring this repo's own federation across GitHub/GitLab/GitFlic/GitVerse
- `.github/workflows/` — CI pipelines (all manually-triggered; see ADR-0001)

## License

Apache-2.0 (code) / CC-BY-SA-4.0 (documentation). See `LICENSE`.
