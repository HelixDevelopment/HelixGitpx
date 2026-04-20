# 23 — WASM Plugin SDK

> **Document purpose**: Enable third parties and enterprises to ship **custom Git-provider adapters, webhook normalisers, and policy validators** without forking HelixGitpx. Plugins are WebAssembly components compiled from any language, signed with Cosign, and sandboxed at runtime.

---

## 1. Why WASM

| Option | Verdict |
|---|---|
| Fork the codebase | ✗ unmaintainable, security boundary vague |
| Native Go plugins (`plugin` pkg) | ✗ ABI fragile, Linux-only, same process |
| External sidecars (gRPC) | ✗ operational overhead, per-tenant isolation hard |
| **WASM component model** | ✓ cross-language, sandboxed, signed, fast-warm, hot-reloadable |

ADR-0015 records the decision.

---

## 2. Plugin Shapes

| Shape | Interface | Example |
|---|---|---|
| **Git Adapter** | `helixgitpx:adapter/universal-git@1.0.0` | Custom corporate Git host |
| **Webhook Normaliser** | `helixgitpx:webhook/normaliser@1.0.0` | Proprietary upstream's payload format |
| **Policy Validator** | `helixgitpx:policy/validator@1.0.0` | Regulatory content screening before push |
| **Conflict Heuristic** | `helixgitpx:conflict/heuristic@1.0.0` | Domain-specific rename detection |
| **AI Post-Processor** | `helixgitpx:ai/post-processor@1.0.0` | Redact or transform model outputs |

Interfaces are **WIT** (WebAssembly Interface Types) under `proto/wit/helixgitpx/`.

---

## 3. Runtime — Wasmtime

- Host: `wasmtime` (component model enabled).
- **Epoch interruption** enforces wall-clock time budgets.
- **Fuel metering** caps CPU per call.
- **Memory limits** per plugin instance (default 128 MiB).
- Parallelism: pool of pre-compiled module instances per plugin + per-tenant partition.
- Cold-start budget: ≤ 50 ms.
- Hot call budget: ≤ 2 ms typical.

---

## 4. Sandbox Model

- **Deny by default**: no filesystem, no network, no clocks beyond host-approved imports.
- **Explicit host imports**: HTTP fetch (allow-listed FQDNs), logger, metrics, OTel span emitter.
- **No ambient authority**: every capability must be granted in the plugin manifest.
- Per-tenant resource limits enforced by host.

### 4.1 Host Imports (selected)

```wit
// helixgitpx:host/net@1.0.0
interface net {
    record http-request {
        method: string,
        url: string,
        headers: list<tuple<string, string>>,
        body: option<list<u8>>,
    }
    record http-response {
        status: u16,
        headers: list<tuple<string, string>>,
        body: list<u8>,
    }
    http-fetch: func(req: http-request) -> result<http-response, http-error>
}

// helixgitpx:host/log@1.0.0
interface log {
    log: func(level: level, msg: string, attrs: list<tuple<string, string>>)
}

// helixgitpx:host/metrics@1.0.0
interface metrics {
    counter-add: func(name: string, value: u64, labels: list<tuple<string,string>>)
    histogram-observe: func(name: string, value: f64, labels: list<tuple<string,string>>)
}
```

All host imports are tenant-scoped (labels injected by host; plugin cannot forge).

---

## 5. Plugin Manifest (`helixgitpx-plugin.toml`)

```toml
[plugin]
id          = "com.example.custom-host"
name        = "Custom Corporate Git Adapter"
version     = "1.2.3"
shape       = "git-adapter"
interface   = "helixgitpx:adapter/universal-git@1.0.0"
min_host    = "1.0.0"
license     = "Apache-2.0"
homepage    = "https://example.com/helixgitpx-plugin"

[capabilities]
net.http    = ["https://git.example.com", "https://api.example.com"]
log         = true
metrics     = true
secrets     = ["vault:plugins/com.example.custom-host/*"]

[limits]
memory_mib       = 128
call_timeout_ms  = 5000
cold_start_ms    = 100
max_concurrency  = 8

[signing]
cosign_identity  = "https://github.com/vasic-digital/helixgitpx-plugin-example/.github/workflows/release.yaml@refs/heads/main"
cosign_issuer    = "https://token.actions.githubusercontent.com"
sbom             = "sbom.cdx.json"
slsa_provenance  = "provenance.intoto.jsonl"
```

---

## 6. Example — Minimal Git Adapter (Rust)

```rust
// src/lib.rs
wit_bindgen::generate!({ world: "plugin" });

use exports::helixgitpx::adapter::universal_git::*;
use helixgitpx::host::net;
use helixgitpx::host::log;

struct Adapter;

impl Guest for Adapter {
    fn provider_info() -> ProviderInfo {
        ProviderInfo {
            name: "Custom Corporate Host".into(),
            kind: "wasm".into(),
            capabilities: Capabilities {
                supports_pull_requests: true,
                supports_issues: true,
                supports_releases: false,
                supports_webhooks: true,
                max_file_size_bytes: 50 * 1024 * 1024,
                ..Default::default()
            },
        }
    }

    fn push_ref(ctx: CallContext, req: PushRefRequest) -> Result<PushRefResponse, Error> {
        log::log(log::Level::Info, &format!("pushing {} to {}", req.ref_name, req.url), vec![]);
        let resp = net::http_fetch(net::HttpRequest {
            method: "POST".into(),
            url: format!("{}/push", req.url),
            headers: vec![("authorization".into(), req.credential_ref.clone())],
            body: Some(req.pack_bytes),
        }).map_err(to_err)?;

        if resp.status != 200 {
            return Err(Error { code: "adapter.upstream_error".into(), message: format!("status {}", resp.status) });
        }
        Ok(PushRefResponse { new_sha: req.new_sha, replicated_at_ms: ctx.now_ms })
    }

    // ... list_refs, create_pr, etc.
}

export!(Adapter);
```

Build & publish:

```bash
cargo component build --release --target wasm32-wasip2
cosign sign-blob --keyless target/wasm32-wasip2/release/plugin.wasm > plugin.wasm.sig
helixctl plugin publish --file plugin.wasm --manifest helixgitpx-plugin.toml --signature plugin.wasm.sig
```

Supported toolchains:

- Rust (`cargo-component`) — recommended.
- Go (`TinyGo` with component support).
- C / C++ (wit-bindgen-cpp).
- JavaScript / TypeScript (`jco` / componentize-js).
- Python (wasmtime-py + componentize-py).

---

## 7. Admission & Verification

At install time, the plugin host performs:

1. **Integrity**: SHA-256 match against manifest.
2. **Signature**: Cosign keyless verify → identity + issuer in manifest → Rekor transparency log entry exists and is non-revoked.
3. **SBOM scan**: Syft → Grype; deny on High/Critical CVEs.
4. **SLSA level**: require ≥ L2 (L3 for Enterprise customers).
5. **Static analysis**: wasmtime compile + WIT conformance check.
6. **Capability check**: requested capabilities do not exceed org's plugin policy.
7. **Version / host compat**.

Installation emits `helixgitpx.audit.events{action=plugin.install}`.

---

## 8. Runtime Safety

- Plugin calls run on a dedicated worker pool; failures do not crash the host.
- Circuit breakers per plugin: if error rate > 10 % for 5 min → disable; alert.
- Per-tenant fuel budget: plugins cannot monopolise a shared host; exhaustion → interruption → error to caller.
- Logs and metrics emitted by plugin are tagged `plugin_id`.
- Panic / trap → plugin disabled for 60 s backoff.

---

## 9. Data Handling

- Plugins do not persist data; state must be handed back to host or stored through a granted backing (KV table) with quota.
- PII / secrets: plugin can request a secret reference (`vault:...`); host fetches on-demand and injects into a memory region zeroed after use.
- Egress data to external endpoints is audited.

---

## 10. Testing & CI

- Unit tests in plugin's language.
- **Contract tests** via the host's `helixctl plugin conformance` suite — runs every required RPC with canned inputs and asserts output shape.
- **Fuzz tests** on every input surface using `cargo-fuzz` (or equivalent).
- **Performance benchmarks**: cold/warm call, p50/p99.
- **Security**: `cargo audit`, Semgrep, Grype on the WASM artefact (binary scanner).
- Golden master against staging HelixGitpx with shadow-mode traffic.

---

## 11. Marketplace (Future)

- Central catalog at `plugins.helixgitpx.example.com`.
- Verified publishers (Anthropic-esque identity verification).
- Star rating, download count, health signals.
- Install is one click; enterprises can lock to their internal catalog.

---

## 12. Versioning & Compatibility

- WIT interfaces follow semver.
- Host supports the last two majors of every interface.
- Plugins declare `min_host` and `max_host` (both inclusive).
- Breaking interface changes get a new major with a 12-month deprecation window.

---

## 13. Operations

Runbooks:
- Plugin disabled due to error rate → `RB-202`.
- Plugin sandbox violation → immediate disable, audit, notify org owner.
- Fleet-wide rollback if an installed plugin is flagged retroactively (Rekor revocation).

---

## 14. Compliance Notes

- Code-signing identity required (no anonymous plugins).
- For public plugins, source code availability preferred; for private, SBOM mandatory.
- Export-controlled customers: plugin installation requires additional review gate.

---

*— End of WASM Plugin SDK —*
