# 24 — Threat Model

> **Document purpose**: Systematically enumerate **threats against HelixGitpx**, the assets they target, the mitigations already in place, and residual risk. Uses **STRIDE** per trust boundary plus **MITRE ATT&CK** mapping for TTPs. Reviewed every 6 months and on every significant architectural change.

---

## 1. Scope

In scope:
- HelixGitpx control plane (services, data stores, CI/CD).
- Data plane (events, Git objects, AI inference).
- Client applications (web, mobile, desktop, CLI).
- Third-party integrations (Git upstreams, OIDC, billing, notification channels).
- Plugin runtime (WASM).

Out of scope:
- Vulnerabilities in third-party services themselves (tracked separately).
- End-user devices (covered by customer-side policy).
- Physical security of customer premises (for on-prem deployments).

---

## 2. Assets

| Asset | Sensitivity | Owner | Where it lives |
|---|---|---|---|
| Source code (customer repos) | **Confidential** | Customer | PG + object store |
| Upstream credentials (tokens, SSH keys) | **Secret** | Customer | Vault (encrypted, SPIRE-bound) |
| User PII (email, name, IP logs) | **Sensitive** | Customer | PG + Loki |
| Session tokens / PATs | **Secret** | User | Redis + PG (hash) |
| AI training data / feedback | **Sensitive** | Customer | Object store + PG |
| Audit log | **Sensitive** (integrity-critical) | HelixGitpx + customer | PG append-only + Rekor |
| Billing / payment info | **Secret** | Customer + processor | Processor only |
| Platform secrets (DB passwords, mTLS roots) | **Secret** | HelixGitpx | Vault |
| Plugin binaries | **Integrity-critical** | Author + HelixGitpx | OCI registry + Cosign |
| Event stream | **Confidential** | Customer | Kafka (encrypted in transit + at rest) |

---

## 3. Actors

| Actor | Motivation | Capability |
|---|---|---|
| External attacker (opportunistic) | Credentials, crypto, data theft | Internet reconnaissance, CVE exploitation |
| External attacker (targeted APT) | Supply chain, IP theft | Resourced, patient, multi-stage |
| Disgruntled user | Retaliation, sabotage | Valid creds, limited scope |
| Malicious insider (engineer) | Exfiltration, sabotage | Elevated internal access |
| Compromised upstream | Data corruption | Valid webhook signatures |
| Malicious plugin author | Backdoor, data exfil | Code execution inside WASM sandbox |
| Curious / nuisance bot | Scanning, scraping | Automated, low sophistication |
| Nation-state (for regulated customers) | Intelligence | Unlimited resources |

---

## 4. Trust Boundaries

```
┌─────────────────────────────────────────────────┐
│ T0  Public Internet                             │
└─────────────────────────────────────────────────┘
              │ HTTPS / Git over SSH
              ▼
┌─────────────────────────────────────────────────┐
│ T1  Cloudflare Edge (WAF, DDoS)                 │
└─────────────────────────────────────────────────┘
              │ mTLS
              ▼
┌─────────────────────────────────────────────────┐
│ T2  Envoy Gateway (TLS termination)             │
└─────────────────────────────────────────────────┘
              │ mTLS (Istio Ambient)
              ▼
┌─────────────────────────────────────────────────┐
│ T3  Service Mesh (SPIFFE SVIDs)                 │
│     ┌──────┐ ┌──────┐ ┌──────┐                  │
│     │ API  │ │ Repo │ │ ...  │                  │
│     └──────┘ └──────┘ └──────┘                  │
└─────────────────────────────────────────────────┘
              │ mTLS (SASL/OAUTHBEARER on Kafka)
              ▼
┌─────────────────────────────────────────────────┐
│ T4  Data Plane (PG, Kafka, Redis, OS, …)        │
└─────────────────────────────────────────────────┘
              │ Explicit NAT gateway
              ▼
┌─────────────────────────────────────────────────┐
│ T5  Upstream Git Hosts (GitHub, etc.)           │
└─────────────────────────────────────────────────┘
```

Additional boundary: **T6 — Plugin Runtime** (Wasmtime sandbox within service pod).

---

## 5. STRIDE per Component

### 5.1 API Gateway (T2 → T3)

| Threat | Attack | Mitigation | Residual |
|---|---|---|---|
| **S**poofing | Forged client cert | mTLS + SPIFFE identity verification | Low |
| **T**ampering | Modified HTTP body in transit | TLS 1.3 + HSTS | Very low |
| **R**epudiation | User denies action | Audit log with trace id, source IP, device id | Low |
| **I**nformation disclosure | Verbose error messages | Sanitised error mapping (see error-catalog.md) | Low |
| **D**oS | Connection exhaustion | Cloudflare + Envoy rate limits + per-token bucket | Medium (sophisticated DDoS) |
| **E**oP | Token downgrade | Strict JWT validation; no cookies for API paths | Low |

### 5.2 Auth Service

| Threat | Attack | Mitigation |
|---|---|---|
| Credential stuffing | Automated login attempts | Rate limits, anomaly detection (RB-200), MFA enforcement for admins |
| Refresh-token theft | Cookie/local-storage theft | Rotating refresh + reuse detection → family revocation |
| OIDC MITM | Malicious IdP or hijacked redirect | PKCE, strict redirect-uri allowlist, `nonce` check |
| MFA bypass | SIM-swap / SMS intercept | FIDO2/TOTP only; SMS deprecated; recovery codes printed one-time |
| Account takeover via password reset | Predictable tokens | Random 256-bit tokens, short TTL, rate-limited, one-use |
| PAT abuse | Leaked token | Prefix `hpxat_` triggers secret scanners; auto-revoke on leak detection |

### 5.3 Git Ingress

| Threat | Attack | Mitigation |
|---|---|---|
| Unauthorized push | Token without scope | Per-scope PAT; scope check on every ref |
| Malicious pack | Zip bomb / resource exhaustion | Size limits, streaming parse, sandbox `git-receive-pack` |
| Ref injection | Force-push over protection | Branch protection enforced centrally (not delegated to upstream) |
| LFS exfiltration | Large file abuse | Per-repo/per-org quotas; content-type allowlist |

### 5.4 Adapter Pool

| Threat | Attack | Mitigation |
|---|---|---|
| Credential leak into logs | Verbose upstream error | Redaction rules on structured logger |
| SSRF via user-supplied upstream URL | Attacker adds upstream pointing at internal IP | Egress allowlist on `adapter-pool` namespace; DNS resolution constrained |
| Replay attack on webhooks | Reused signed payload | HMAC verify + `dedup` in Redis (event_id) |
| Upstream credential compromise | Token from leak | Per-upstream rotation API, anomaly alerts |

### 5.5 Event Plane (Kafka)

| Threat | Attack | Mitigation |
|---|---|---|
| Unauthorized produce/consume | Missing auth | SASL/OAUTHBEARER tied to SPIFFE SVID; per-topic ACLs |
| Schema injection | Compat-breaking producer | Karapace strict compatibility |
| Log tampering | Direct broker access | Strict tenant isolation on brokers; disk encryption |

### 5.6 Conflict Resolver & AI

| Threat | Attack | Mitigation |
|---|---|---|
| Prompt injection via user code/comment | Model executes attacker's instructions | System-prompt layering, guardrails, structured output enforcement, no tool-use without validation |
| Data exfiltration via AI outputs | Model leaks training data | On-org boundary; no cross-tenant training; DPO training eval filters PII |
| Model poisoning | Attacker feeds bad feedback | Weighted by reputation; outlier detection in curator; shadow evaluation before promotion |
| Sandbox escape | Proposed patch runs code | Ephemeral sandbox pod with no network, no mounts, CPU+mem caps, kill timeout |

### 5.7 Plugin Runtime

| Threat | Attack | Mitigation |
|---|---|---|
| Sandbox escape | WASM module exploits Wasmtime | Wasmtime hardened config; regular upgrades; seccomp; per-plugin resource caps |
| Capability overreach | Plugin requests broader net access | Manifest declared capabilities verified at install; drift check at runtime |
| Supply-chain compromise | Signed malicious plugin | Cosign + Rekor; SBOM scan; revocation via Rekor checkpoints |
| DoS on host | Infinite loops | Fuel metering + epoch interruption |

### 5.8 Web / Mobile / Desktop Clients

| Threat | Attack | Mitigation |
|---|---|---|
| XSS on web | Crafted Markdown in issue body | CSP strict-dynamic, DOMPurify, sandboxed iframe for preview |
| CSRF on web | State-changing GET/POST from other origin | SameSite=Lax cookies; double-submit token for forms; API uses bearer tokens |
| Clickjacking | Iframe wrap | X-Frame-Options DENY + CSP frame-ancestors 'none' |
| Local token theft (mobile) | Malware / rooted device | Hardware-backed Keystore / Keychain; optional biometric unlock gate; detect jailbreak/root |
| Code tampering (desktop) | Malicious MSIX | Code-signed + auto-update via signed manifests |

### 5.9 CI/CD

| Threat | Attack | Mitigation |
|---|---|---|
| Malicious PR runs with privileges | Fork-PR token abuse | Pull_request_target disallowed; read-only creds on fork PRs |
| Dependency confusion | Malicious public package masking internal | Private registry + scoped packages; lockfile enforcement |
| Secret exfiltration from logs | `echo $SECRET` | Masked outputs; no secrets in env; OIDC-federated short-lived creds |
| Supply-chain tampering | Image rewrite post-sign | Immutable digests; Kyverno admission verifies Cosign + Rekor + SLSA |

---

## 6. Attack Surface Inventory

- Public HTTPS endpoints: `app.*`, `api.*`, `push.*`, `webhooks.*`, `status.*`, `docs.*`.
- Public SSH: `push.helixgitpx.example.com:22` (Git over SSH).
- OIDC callback URLs.
- Webhook inbound URLs per upstream provider.
- Admin endpoints (`/api/v1/admin/*`) — policy-gated, short-lived access.
- Plugin install endpoints.
- Status / metrics endpoints — internal only.

---

## 7. MITRE ATT&CK Mapping (selected)

| Tactic | Technique | Our defence |
|---|---|---|
| Initial Access | T1190 Exploit Public-Facing App | WAF, dependency scanning, Semgrep/CodeQL in CI |
| Initial Access | T1078 Valid Accounts | MFA enforcement, anomaly detection |
| Execution | T1059 Command & Scripting | No shell execution in services; sandbox for AI patches |
| Persistence | T1098 Account Manipulation | Audit log + alert on admin changes |
| Privilege Escalation | T1078.004 Cloud Accounts | Short-lived OIDC creds; Vault dynamic secrets |
| Defense Evasion | T1070 Indicator Removal | Append-only audit + Rekor anchor |
| Credential Access | T1555 Credentials from Password Stores | No plaintext creds; Vault with SPIRE-bound auth |
| Discovery | T1087 Account Discovery | Rate limits on user/email enumeration; constant-time responses |
| Lateral Movement | T1021 Remote Services | Default-deny NetworkPolicies; Cilium L7 policy |
| Collection | T1005 Data from Local System | Disk encryption; minimal local state on pods (readOnly rootfs) |
| Exfiltration | T1041 Exfiltration Over C2 | Explicit egress allowlists; DLP on PII-sensitive egress |
| Impact | T1486 Data Encrypted for Impact | Immutable backups; versioned object storage with Object Lock |
| Impact | T1498 Network DoS | Cloudflare anycast + rate limits |

---

## 8. Data-Flow Risk Highlights

- **Write-then-event outbox** eliminates dual-write inconsistencies that could be exploited to produce false audit trails.
- **Cross-region replication** is encrypted end-to-end (MirrorMaker 2 over TLS; PG logical replication over TLS).
- **LFS** bytes go directly to object storage via presigned URLs — never through our services, limiting blast radius for malicious bytes.
- **Plugins** run in a second sandbox layer inside an already-isolated service pod.

---

## 9. Known Residual Risks

| Risk | Impact | Likelihood | Plan |
|---|---|---|---|
| Sophisticated DDoS exceeding Cloudflare throughput | Service degradation | Low | Scale plan + Anycast; on-call drill |
| Novel Kubernetes container escape (0-day) | Cluster compromise | Very low | Continuous patching; Kata runtime for highest-sensitivity workloads |
| Prompt-injection leading to undesirable AI output | User confusion / reputational | Medium | Guardrails + output validation + user review step for irreversible actions |
| Supply-chain compromise of a foundational dependency | Variable | Low (prob) × High (impact) | SBOM scans, Rekor verification, vendoring for pinned deps, rapid patch SLA |
| Regulatory change requiring feature removal | Business | Medium | Feature flags + region toggles enable rapid compliance pivots |

---

## 10. Threat Review Cadence

- **Quarterly**: lightweight review with engineering leads for new components.
- **Semi-annually**: full STRIDE review.
- **Ad hoc**: when onboarding new provider, deploying a new AI model, or in response to industry incidents.
- **External**: annual pen-test; continuous bug bounty.

Outputs from every review: changes to mitigations, new alerts, test cases, or documented accepted risks in this file.

---

## 11. Responsible Disclosure

- **security@helixgitpx.example.com** — PGP key published.
- Bug bounty via HackerOne / Intigriti with clearly scoped rules and safe-harbour language.
- SLA for triage: 24 h; remediation SLAs keyed to severity (Critical 7 d, High 30 d, Medium 90 d).

---

*— End of Threat Model —*
