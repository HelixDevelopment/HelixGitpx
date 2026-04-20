# ADR-0002 — Portable container runtime (docker OR podman) via wrapper

**Status**: Accepted
**Date**: 2026-04-20
**Deciders**: Милош Васић

---

## Context and Problem Statement

Primary dev machines have `podman`; many ecosystem templates and tools assume `docker`. Hardcoding either breaks one class of machines, fragmenting the developer experience.

## Decision Drivers

- Support both podman and docker without per-machine configuration
- Preserve existing tooling patterns (docker-compose syntax)
- Minimize platform-specific workarounds in build scripts

## Considered Options

- **A** Hardcode `docker` — breaks podman-first machines
- **B** Hardcode `podman` — breaks docker-first machines and breaks many SaaS CI templates
- **C** Wrapper script that detects available runtime — works everywhere, adds one indirection layer

## Decision

**We will adopt a wrapper script (Option C).** All local-dev tooling (Makefiles, scripts) invokes `impl/helixgitpx-platform/compose/bin/compose`, which detects the available runtime at invocation time (preferring `docker compose` → `podman compose` → `podman-compose`). Scripts never call `docker` or `podman` directly.

## Consequences

### Positive

- Works on both docker and podman hosts with no per-machine setup.
- Developers still see the actual runtime in logs via `compose` output.
- Integrates cleanly with CI (GitHub Actions runners have docker preinstalled).

### Negative / Trade-offs

- Features that differ between runtimes (selinux labels, some network modes) are constrained to the common subset.
- One additional shell indirection for every compose invocation (negligible performance cost).

### Risks & Mitigations

- Risk: wrapper logic drifts from actual runtime behavior → Mitigation: regular testing on both runtimes; wrapper is minimal and reviewed carefully
- Risk: feature parity breaks (e.g., podman adds a flag that docker doesn't have) → Mitigation: document in runbooks; use lowest-common-denominator features

## Implementation Notes

- Wrapper location: `impl/helixgitpx-platform/compose/bin/compose`
- Detection logic (bash script, ~30 lines): try `docker compose`, then `podman compose`, then `podman-compose`, fail with a helpful error if none found
- Used by all Makefiles and scripts in `impl/helixgitpx-platform/`

## Validation & Exit Criteria

- `compose up -d postgres` succeeds on both docker and podman machines
- Wrapper script has zero external dependencies (pure bash)

## References

- `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §4 C-3, §8
- `impl/helixgitpx-platform/compose/bin/compose` (implementation)
