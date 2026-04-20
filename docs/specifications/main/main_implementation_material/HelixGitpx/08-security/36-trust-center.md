# 36 — Trust Center (Public Security & Privacy Posture)

> **Document purpose**: The **public-facing** security, privacy, and reliability overview — what HelixGitpx tells prospective and existing customers about how we protect them. The working source of truth lives here; the customer-facing page mirrors this content with lighter design.

---

## 1. Our Commitments

- **Security is a first-class product concern.** Every feature has a threat model; every release passes supply-chain, SAST, DAST, and dependency gates.
- **Your data is yours.** We don't sell it, don't mine it, and don't use it for purposes you haven't consented to.
- **Transparency over marketing.** Public postmortems, subprocessor list, deprecation calendar.
- **Customer control.** Data residency, retention, and AI participation are opt-able.

---

## 2. Certifications & Audits

| Framework | Status |
|---|---|
| **SOC 2 Type I** | Completed (YYYY-MM) |
| **SOC 2 Type II** | In progress; report target YYYY-MM |
| **ISO/IEC 27001:2022** | Gap analysis complete; certification YYYY |
| **ISO/IEC 27017 / 27018** | On roadmap |
| **SLSA L3** | Achieved for every production artefact |
| **CSA STAR Self-Assessment** | Published |
| **Cyber Essentials (UK)** | Certified |
| **PCI DSS** | Out of scope (we don't process card data directly; our processor is PCI certified) |
| **FedRAMP Moderate** | On roadmap for dedicated US deployment |

Reports available under NDA via the Trust Portal.

---

## 3. Data Residency & Isolation

- **EU residency**: Paid orgs can pin their primary region to the EU. Data at rest stays in the EU; cross-region replication is EU-only or disabled.
- **UK residency**: Supported.
- **US residency**: Supported.
- **Custom residency**: On-prem / dedicated deployments for customers with strict requirements.
- **Logical isolation**: Every query is tenant-scoped via Postgres Row-Level Security; Kafka ACLs per consumer group; Object storage tenant-prefixed.

---

## 4. Encryption

- **In transit**: TLS 1.3; mTLS between internal services via SPIFFE SVIDs; strict cipher list; OCSP stapling at edge.
- **At rest**: AES-GCM-256 across PG, Kafka, Redis, OpenSearch, Qdrant, and object storage.
- **Keys**: HashiCorp Vault; customer-managed KMS on Enterprise; Shamir-split root; HSMs for critical keys.
- **FIPS 140-2/3 mode** available on-prem / dedicated.

---

## 5. Identity & Access

- OIDC / SAML SSO supported; local accounts optional.
- MFA (TOTP + FIDO2) enforceable at org level; SMS deprecated.
- SCIM provisioning on Business+ (roadmap).
- Fine-grained PATs with per-scope permissions.
- Session TTLs tight (15 min access, rotating refresh, family-reuse detection).

---

## 6. Application Security

- **Threat modelling** for every new component; [threat model published](../08-security/24-threat-model.md) (excerpted).
- **Secure SDLC**: SAST (Semgrep, CodeQL, Snyk Code), DAST (OWASP ZAP), Dependency scanning (Snyk OSS, govulncheck), Secrets scanning (Gitleaks) — all blocking CI gates.
- **Peer review** required on every merge.
- **Penetration testing**: Annual external; continuous bug bounty.
- **Rapid patching**: Critical CVEs within 7 days; High 30 days.

---

## 7. Infrastructure & Operations

- **Defence in depth**: WAF + DDoS protection at edge, mTLS mesh, default-deny NetworkPolicies, PSS-restricted pods.
- **Immutable infrastructure**: GitOps via Argo CD; every production change is a signed commit.
- **Signed artefacts**: Cosign keyless + Rekor transparency log; Kyverno admission verifies.
- **SLSA L3** builds: hermetic, reproducible, provenance attached.
- **SBOM** (CycloneDX 1.5) for every image.
- **Runtime monitoring**: Continuous profile + behaviour baselining; anomalies alerted.

---

## 8. Monitoring & Incident Response

- 24/7 on-call rotation.
- PagerDuty + status page auto-incidents.
- Every P1 triggers a formal post-mortem within 5 business days.
- Historical incidents public at `status.helixgitpx.example.com/incidents`.
- Customer notification: status page + email + webhook options.

---

## 9. Availability

- Target SLOs:
  - 99.9 % (Business)
  - 99.95 % (Enterprise default)
  - 99.99 % (Enterprise premium)
- Multi-region active/active-passive.
- RPO ≤ 30 s / RTO ≤ 15 min.
- Quarterly DR drills, including full region failover.
- Backups: daily + weekly + monthly retained; monthly restore test.

See [30-sla.md] for the full SLA.

---

## 10. Privacy

- **GDPR / CCPA / PIPEDA / LGPD compliant**.
- **DPA** signed by default for every paid plan; SCCs bundled.
- **Subprocessors** listed publicly, 30-day notice on changes.
- **Data-subject rights** fully supported (access / rectification / erasure / portability / objection).
- **Retention** per published policy ([27-data-retention-privacy.md]); customer-configurable within compliance bounds.
- **AI training**: default off in EU; opt-out everywhere; customer data never used to train models for other customers.

---

## 11. AI Safety & Governance

- Default self-hosted inference (on-prem / our region). Customer code never leaves the customer-chosen region unless explicitly opted-in.
- **Guardrails**: input + output filters on every model call.
- **Sandboxed validation**: AI-proposed patches run in ephemeral pods before apply.
- **Human oversight**: undo window + escalation paths for high-stakes suggestions.
- **No automated decision-making** producing legal or similarly significant effects on individuals (GDPR Art. 22).
- **Model registry** auditable; shadow mode before promotion.

---

## 12. Software Supply Chain

- Source code review + signed commits (Gitsign) on `main`.
- Dependencies: pinned, SBOM'd, CVE-scanned.
- Build: hermetic, attested (SLSA L3, in-toto provenance on Rekor).
- Release: signed, immutable digests.
- **Admission**: Kyverno verifies at deploy time.
- Third-party plugins: require Cosign signatures; sandboxed runtime; revocation via Rekor.

---

## 13. Responsible Disclosure

- **security@helixgitpx.example.com** (PGP key below).
- Bug bounty via HackerOne; scope + safe harbour published.
- Triage SLA: 24 h.
- Remediation SLA per severity (see [24-threat-model.md]).
- Credit published on our hall-of-fame (optional for reporter).

PGP key:

```
-----BEGIN PGP PUBLIC KEY BLOCK-----
[public key material here]
-----END PGP PUBLIC KEY BLOCK-----
```

---

## 14. Policies (publicly linked)

- Terms of Service
- Privacy Policy
- Acceptable Use Policy
- Data Processing Addendum
- Subprocessor list (live)
- Cookie Policy
- Accessibility Statement
- Responsible AI Use Policy

---

## 15. Compliance Artifacts Available to Customers

Under NDA via the Trust Portal:

- SOC 2 Type I/II reports
- ISO 27001 certificate + SoA (when issued)
- CAIQ / STAR questionnaire
- Latest pen-test executive summary
- Architecture diagrams (sanitised)
- Business continuity plan excerpt
- Disaster recovery test results

---

## 16. Customer Security Controls

Available to customers to harden their usage:

- MFA enforcement.
- Branch protection (required reviews, signed commits, required checks).
- Fine-grained PATs.
- IP allowlists (Enterprise).
- SAML / SSO.
- SCIM provisioning (roadmap).
- Audit log export.
- Data residency pinning.
- AI opt-out.
- Egress allowlisting for adapter pool (on-prem).
- Plugin signature enforcement.

---

## 17. Contact

- **security@helixgitpx.example.com** — vulnerabilities, incidents.
- **privacy@helixgitpx.example.com** — DSRs, DPO contact.
- **trust@helixgitpx.example.com** — compliance, questionnaires, trust portal access.
- **status.helixgitpx.example.com** — live status & historical incidents.

---

## 18. Change Log

Every material change to security, privacy, or availability posture reflected here is logged below with date and summary.

| Date | Change |
|---|---|
| 2026-01-15 | Initial trust page published. |
| 2026-03-02 | SOC 2 Type I report available. |
| 2026-03-20 | Added subprocessor: (example). |

---

*— End of Trust Center —*
