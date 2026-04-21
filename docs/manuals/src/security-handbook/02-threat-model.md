## 2. Threat model

This chapter enumerates the attacker capabilities HelixGitpx defends
against, plus the assumptions under which those defences hold.

### 2.1 Attacker personas

| Persona | Capability | Motivation |
|---------|-----------|------------|
| **Opportunistic script-kiddie** | Scans for open surfaces; sprays credentials harvested elsewhere. | Crypto miners, credential-stuffing payoff. |
| **Malicious upstream** | Controls one of the bound Git providers; can push, mutate metadata, issue webhooks. | Data exfiltration, cross-tenant leakage. |
| **Compromised CI runner** | Holds a PAT; can push to any repo the PAT has access to. | Lateral movement into production. |
| **Insider (rogue member)** | Has an active HelixGitpx session + membership. | Theft, sabotage, blackmail. |
| **Nation-state** | Deep resources, long dwell time, pre-positioned implants. | Strategic Git targets (critical infrastructure, defence). |
| **Provider outage + bad merge** | Accidental adversary. | No malice — but still drops data. |

### 2.2 Assets in order of loss severity

1. Customer source code (HIGH — IP, secrets).
2. Audit log integrity (HIGH — compliance, forensic response).
3. Org tenancy boundaries (HIGH — cross-tenant leak is a DSR + press event).
4. Upstream provider tokens (MEDIUM — revoke-and-rotate feasible).
5. Service availability (MEDIUM — SLO breach, revenue hit).
6. Build-time secrets (MEDIUM — typically short-lived).

### 2.3 Defences, mapped to attackers

| Attack | Control |
|--------|---------|
| Credential spray | WebAuthn enforced for staff; TOTP for customers; Keycloak rate-limit. |
| Upstream pushes malicious ref | `conflict-resolver` surfaces divergence; human gate on apply. |
| Upstream replays old webhook | HMAC-SHA256 signature + unique delivery ID + 5-minute window. |
| Compromised CI PAT | PAT scopes; per-org residency; OPA deny on out-of-scope repo. |
| Insider `rm -rf` | Every mutating action emits `audit.events`; Merkle-anchored every hour; soft-delete window. |
| OPA bundle tamper | Bundle signed (`cosign sign-blob`); OPA agents verify signature on pull. |
| Image-based supply chain | Cosign signature + SBOM attestation + Trivy gate in `supply-chain.yml`. |
| mTLS downgrade | SPIFFE + Istio Ambient enforce; no plain HTTP east-west. |
| AI prompt-injection exfil | NeMo Guardrails on every ai-service call; decisions in audit log with prompt hash. |

### 2.4 Out of scope

- Physical security of the cluster (operator's responsibility).
- Supply chain of the distro base images (we pin and scan; upstream
  mis-signing is outside our control).
- Compromise of the OIDC provider (if Keycloak is root-owned, all bets
  off).

### 2.5 Trust boundaries

The HelixGitpx cluster has four:

1. **Public internet ↔ edge** — Cloudflare + Istio gateway + cert-manager TLS.
2. **Edge ↔ workload** — Istio Ambient mTLS; each workload has a SPIFFE SVID.
3. **Workload ↔ data plane** — CNPG + Strimzi TLS; per-DB role per service.
4. **Cluster ↔ external LLM provider (optional)** — outbound-only; per-org policy; never receives customer source code unless the org opted in.

### 2.6 Residual risk

Most attackers above can be caught, blocked, or mitigated. **What
remains**: compromise of the operator's laptop (SSH keys + kubeconfig)
would grant an attacker root-equivalent cluster access. Mitigation is
operational: hardware-backed SSH, short-lived kubeconfigs, PAM + MFA
on the workstation. See the
[Operator Guide chapter 4](../operator-guide/04-identity-mtls.md).

---
