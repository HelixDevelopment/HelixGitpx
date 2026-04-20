# ADR-0022 — wazero for WASM plugin runtime

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

adapter-pool needs a WASM plugin host to extend provider support without recompiling the service. Options: wasmtime-go (CGO, bigger image), wazero (pure Go, zero CGO), wasmer-go (dev discontinued).

## Decision

`tetratelabs/wazero` for the plugin host. Pure Go means adapter-pool's distroless image stays small; no libc/libssl dependencies leak into the sandbox. Example plugin written in TinyGo (WASI target).

## Consequences

- adapter-pool Docker image stays < 50 MiB.
- Plugin ABI is text/binary JSON across the WASM boundary; documented in `adapter-pool/examples/plugin-hello/README.md`.
- Performance ~2× slower than native Go for heavy workloads; acceptable for glue code calling remote APIs.

## Links

- `docs/superpowers/specs/2026-04-20-m5-federation-conflict-engine-design.md` §2 C-2
- https://github.com/tetratelabs/wazero
