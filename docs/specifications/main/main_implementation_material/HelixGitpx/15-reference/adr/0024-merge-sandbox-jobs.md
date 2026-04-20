# ADR-0024 — Three-way merges run in ephemeral Kubernetes Jobs

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

conflict-resolver performs three-way merges on user content. A compromised merge tool or malicious repo contents could exfiltrate credentials or exploit libc bugs. Running merges in the main service pod shares its privilege surface.

## Decision

Each merge is a short-lived Kubernetes Job:
- seccomp: `RuntimeDefault`
- capabilities: drop `ALL`
- network policy: egress=none (no DNS, no HTTP — the Job reads input via mounted configmap, writes output to a mounted emptyDir)
- service account: dedicated, no API-server access
- image: alpine + git only

The controller (conflict-resolver) waits on Job completion, reads the output, emits `conflict.resolved`.

## Consequences

- One Job per merge — ~200 ms startup overhead, dwarfed by merge time on non-trivial repos.
- Exploit blast radius is one emptyDir + CPU time; no persistent state reachable.
- Scale tied to Kubernetes scheduler throughput; M8 may introduce a Job-pool pre-warmer.

## Links

- `docs/superpowers/specs/2026-04-20-m5-federation-conflict-engine-design.md` §2 C-5
