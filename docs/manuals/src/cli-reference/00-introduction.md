# HelixGitpx CLI Reference

## 1. Introduction

The `helixgitpx` CLI is a single static binary that talks to the
HelixGitpx API. It's the fastest way to bind repos, inspect conflicts,
and drive CI workflows.

### 1.1 Install

- **Homebrew:** `brew install helixgitpx/tap/helixgitpx`
- **Debian/Ubuntu:** `sudo apt install ./helixgitpx_<version>_amd64.deb`
- **Fedora/RHEL:** `sudo rpm -i helixgitpx-<version>.x86_64.rpm`
- **Windows (winget):** `winget install HelixGitpx.CLI`
- **Direct:** download from the releases page on each upstream.

The binary is reproducibly built; the SBOM and Cosign signature ship
alongside every release.

### 1.2 First run

```bash
helixgitpx login https://app.helixgitpx.io   # OIDC browser flow
helixgitpx org list
helixgitpx repo bind --repo myorg/myrepo --upstream github,gitlab
```

### 1.3 Global flags

| Flag | Purpose |
|------|---------|
| `--endpoint <url>` | Override the API base URL |
| `--profile <name>` | Select a named profile (default `default`) |
| `--output <fmt>` | `text` (default), `json`, `yaml` |
| `--verbose` | Noisy logs |

### 1.4 Configuration

The CLI looks at `~/.config/helixgitpx/config.yaml`. One profile per
environment is typical.

```yaml
profiles:
  default:
    endpoint: https://api.helixgitpx.io
  staging:
    endpoint: https://api.staging.helixgitpx.io
    token_file: ~/.config/helixgitpx/staging.token
```

### 1.5 Command topics

- `helixgitpx login` / `logout` — session lifecycle.
- `helixgitpx org` — organization management.
- `helixgitpx team` — teams and membership.
- `helixgitpx repo` — repository CRUD and bindings.
- `helixgitpx upstream` — manage upstream connections.
- `helixgitpx sync` — inspect and retry sync operations.
- `helixgitpx conflict` — inbox + resolve flow.
- `helixgitpx ai` — summarize, propose, chat.
- `helixgitpx search` — semantic and code search.
- `helixgitpx policy` — download and simulate OPA bundles.
- `helixgitpx billing` — plan and invoices.

Every command supports `--help` with full flag reference.

### 1.6 Scripting

CLI output on `--output json` is stable across minor versions; scripts
can rely on it. For long-running operations, `helixgitpx … --watch` tails
event streams and exits on terminal state.

---
