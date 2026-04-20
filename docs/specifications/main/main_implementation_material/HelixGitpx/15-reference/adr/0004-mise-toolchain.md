# ADR-0004 — mise for toolchain management

**Status**: Accepted
**Date**: 2026-04-20
**Deciders**: Милош Васић

---

## Context and Problem Statement

The monorepo spans Go, Node, JVM, and a dozen CLI tools (buf, sqlc, kubectl, helm, kind, k3d, skaffold, tilt, cosign, syft, grype, kyverno-cli, checkov, gofumpt, golangci-lint, gremlins, delve, goose). Each developer needs identical versions to reproduce builds. A hermetic Nix-based environment is the strictest option; a simpler PATH-managed pinning is sufficient for M1 and lower cost to adopt.

## Decision Drivers

- Reduce "works on my machine" failures due to version drift
- Zero per-developer setup overhead (one `mise install` after checkout)
- Low friction compared to container-based or Nix-based approaches
- Quick adoption by existing team members

## Considered Options

- **A** `devbox` (Nix-based) — stricter, slower first-run, larger downloads
- **B** `mise` (former `rtx`) — simpler, non-hermetic, PATH-based, rapid adoption
- **C** No pinning — risk silent version drift

## Decision

**We will adopt `mise` (Option B) with `mise.toml` at the repo root to pin every language runtime and CLI.** Developers run `mise install` after checkout; the shell auto-activates pinned versions in any subdirectory of the repo.

## Consequences

### Positive

- Single source of truth for all tool versions; no `asdf`/`nvm`/`sdkman` sprawl.
- Automatic shell integration: no need to manually set PATH.
- Non-hermetic approach acceptable for M1; tight reproducibility can be re-evaluated at M8 when on-prem deployment requires stricter guarantees.

### Negative / Trade-offs

- Non-hermetic: tools install into `~/.local/share/mise/` and rely on system libc/openssl.
- macOS/Linux parity requires testing on both platforms (not Windows).
- Version selection is snapshot-in-time; no automatic security patches (mitigated by scheduled updates and pinning).

### Risks & Mitigations

- Risk: tool version installs fail on unsupported platforms → Mitigation: document supported platforms (Linux x86_64, macOS arm64/x86_64); use GitHub Actions runners as source of truth
- Risk: system libc incompatibility → Mitigation: Nix hermetic container as fallback for problematic hosts (M8+)

## Implementation Notes

- Root `mise.toml` file pins:
  - Go 1.23.x, Node 20.x, JVM (17+), Rust (if needed)
  - All CLI tools (buf, sqlc, kubectl, helm, kind, k3d, skaffold, tilt, cosign, syft, grype, kyverno-cli, checkov, golangci-lint, gofumpt, gremlins, delve, goose)
- Go toolchain auto-management (GOTOOLCHAIN) coexists cleanly: we pin `go 1.23` in go.mod and optionally `toolchain go1.23.4` in go.work.
- Pre-commit hook (future M2): warn if `mise.toml` is edited but `mise install` not run.

## Validation & Exit Criteria

- `mise install` succeeds and auto-activates on shell startup.
- `go version`, `node --version`, `buf --version` report pinned versions.
- CI environment (GitHub Actions) runs same version matrix (tracked in `.github/workflows/ci-*.yml`).

## References

- `mise.toml` (root of repo)
- https://mise.jdx.dev/
- `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §4 C-5
