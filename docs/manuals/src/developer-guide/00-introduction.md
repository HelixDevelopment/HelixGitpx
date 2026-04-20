# HelixGitpx Developer Guide

## 1. Introduction

This manual is for engineers who ship code into the HelixGitpx platform —
service authors, library authors, adapter authors, and anyone whose PR
lands in `impl/`.

It assumes you have read the
[Constitution](../../../CONSTITUTION.md) and the
[Architecture overview](../../../../docs/specifications/main/main_implementation_material/HelixGitpx/01-architecture/02-system-architecture.md).

### 1.1 Audience

- Backend engineers writing Go services.
- Frontend engineers writing Angular or Compose UI.
- Adapter engineers integrating new Git providers.
- AI engineers extending the `ai-service` use-case catalogue.

### 1.2 What you'll learn

- **Chapter 2:** repo layout, workspace, and the tools/scaffold generator.
- **Chapter 3:** coding conventions per language (Go, TS, Kotlin).
- **Chapter 4:** proto-driven contracts (`buf generate`).
- **Chapter 5:** local development loop (compose up, migrate, hot reload).
- **Chapter 6:** writing tests in all seven required types.
- **Chapter 7:** extending OPA bundles safely.
- **Chapter 8:** submitting a new provider adapter.
- **Chapter 9:** merging a change (PR etiquette, reviews, CI gates).
- **Chapter 10:** release and the upstream federation cadence.

### 1.3 Minimum toolchain

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.23.4 | `mise` or direct |
| Node | ≥20 | `mise` |
| Gradle | ≥8.14 | `mise` |
| buf | latest | `go install github.com/bufbuild/buf/cmd/buf@latest` |
| Docker or Podman | latest | OS package |
| kubectl | 1.31 | OS package |
| helm | 3.16 | OS package |

The `mise.toml` at the repo root pins every version; `mise install`
gets you a working environment.

### 1.4 First 30 minutes

```bash
git clone https://github.com/HelixGitpx/HelixGitpx.git
cd HelixGitpx
mise install                 # toolchain
make bootstrap               # deps across all sub-projects
make dev                     # compose up
cd impl/helixgitpx && go test ./... && echo "ready"
```

---
