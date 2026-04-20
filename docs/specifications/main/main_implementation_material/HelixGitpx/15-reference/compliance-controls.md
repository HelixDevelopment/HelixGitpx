# Compliance Control Mapping

> Cross-reference of **HelixGitpx security controls** to major compliance frameworks. Used internally by the compliance team to respond to questionnaires and prepare for audits. Customers on Enterprise plans receive the relevant subsets on request (under NDA).

---

## Legend

- **Control ID**: internal HelixGitpx reference, stable.
- **Implementation**: how / where the control is implemented + a link to the authoritative doc.
- **Frameworks**: which external control maps to this one (non-exhaustive).
- **Evidence**: what we collect to prove it's operating.

Frameworks covered:
- **SOC 2** (AICPA TSC 2017 / revised 2022).
- **ISO/IEC 27001:2022** (Annex A controls).
- **NIST CSF 2.0** (Categories).
- **PCI DSS 4.0** (selected; HelixGitpx itself is not in scope for card data, but relevant for vendor customers).
- **CIS Controls v8**.
- **GDPR** (relevant articles).

---

## Access Control

### HGX-AC-001 — Unique User Identity
- **Implementation**: OIDC-only login with `sub` uniquely mapped to internal `users.id`; local password login disabled by default (see [08-security/11-security-compliance.md]).
- **SOC 2**: CC6.1, CC6.2
- **ISO 27001**: A.5.16, A.5.17
- **NIST CSF**: PR.AA-01
- **CIS**: 5.2
- **Evidence**: Auth audit log; MFA enrolment reports.

### HGX-AC-002 — MFA Enforcement for Privileged Access
- **Implementation**: admin roles require MFA; step-up enforced at session level; [27-data-retention-privacy.md §DSR].
- **SOC 2**: CC6.1
- **ISO 27001**: A.5.17, A.8.5
- **NIST CSF**: PR.AA-03
- **CIS**: 6.3
- **Evidence**: Enforcement logs; periodic coverage report.

### HGX-AC-003 — Role-Based + Attribute-Based Access Control
- **Implementation**: OPA (Rego) policies; per-repo roles; branch protection rules (see `15-reference/policies/repo-authz.rego`).
- **SOC 2**: CC6.1, CC6.3
- **ISO 27001**: A.5.15, A.5.18
- **NIST CSF**: PR.AA-05
- **Evidence**: OPA decision logs; bundle deployment history.

### HGX-AC-004 — Just-in-Time & Break-Glass Access
- **Implementation**: `helixctl break-glass` time-bound elevation; Vault dynamic secrets; SPIFFE SVIDs rotated hourly.
- **SOC 2**: CC6.1
- **ISO 27001**: A.5.15, A.8.2
- **NIST CSF**: PR.AA-05
- **Evidence**: Break-glass audit entries; Vault audit log.

### HGX-AC-005 — Workload Identity
- **Implementation**: SPIFFE/SPIRE per-pod X.509 SVIDs; mTLS via Istio Ambient; see ADR-0010.
- **ISO 27001**: A.5.16
- **NIST CSF**: PR.AA-03

---

## System Operations

### HGX-SO-001 — Change Management via GitOps
- **Implementation**: Argo CD reconciling signed commits; PR review + 2 approvals on production branches; see [31-release-management.md].
- **SOC 2**: CC8.1
- **ISO 27001**: A.8.32
- **CIS**: 16.1

### HGX-SO-002 — Vulnerability Management
- **Implementation**: Snyk + Grype + Semgrep in CI; SLA: Critical 7 d / High 30 d / Med 90 d; [11-security-compliance.md].
- **SOC 2**: CC7.1
- **ISO 27001**: A.8.8
- **NIST CSF**: ID.RA-01, PR.IR-01
- **CIS**: 7.4, 7.7

### HGX-SO-003 — Patch Management
- **Implementation**: OS via Talos immutable images; container images re-pulled per release; auto-updated dependencies via Renovate.
- **ISO 27001**: A.8.8
- **CIS**: 7.3

### HGX-SO-004 — Incident Response
- **Implementation**: Runbook [19-operations-runbook.md]; PagerDuty; post-mortems within 5 business days.
- **SOC 2**: CC7.3, CC7.4
- **ISO 27001**: A.5.24-A.5.26
- **NIST CSF**: RS.MA-01 through RS.CO-05

### HGX-SO-005 — Business Continuity & DR
- **Implementation**: Multi-region active-active/passive; quarterly DR drills; RPO/RTO documented in [16-infrastructure-scaling.md §9].
- **SOC 2**: A1.2, A1.3
- **ISO 27001**: A.5.29, A.5.30
- **NIST CSF**: RC.RP-01

---

## Cryptography

### HGX-CR-001 — Encryption in Transit
- **Implementation**: TLS 1.3 everywhere; mTLS between services; strict cipher list; OCSP stapling at edge.
- **SOC 2**: CC6.7
- **ISO 27001**: A.8.24
- **PCI DSS**: 4.2.1
- **CIS**: 3.10

### HGX-CR-002 — Encryption at Rest
- **Implementation**: PG / Kafka / object store with AES-GCM-256; keys in Vault with customer-provided KMS option for Enterprise.
- **SOC 2**: CC6.7
- **ISO 27001**: A.8.24
- **PCI DSS**: 3.5, 3.6

### HGX-CR-003 — Key Management
- **Implementation**: Vault + Shamir unseal; rotation automated for dynamic secrets; root keys in HSM where available.
- **SOC 2**: CC6.1
- **ISO 27001**: A.8.24
- **PCI DSS**: 3.7

### HGX-CR-004 — Signing & Verification
- **Implementation**: Cosign keyless + Rekor; SLSA L3 provenance; Kyverno admission verifies.
- **SOC 2**: CC8.1
- **ISO 27001**: A.8.28, A.8.30, A.8.32
- **NIST CSF**: PR.DS-06

---

## Network Security

### HGX-NS-001 — Zero-Trust Network
- **Implementation**: Cilium default-deny NetworkPolicies; Istio Ambient mTLS; explicit egress allowlists (see `18-manifests/network-policy-samples.yaml`).
- **SOC 2**: CC6.6
- **ISO 27001**: A.8.20, A.8.22
- **NIST CSF**: PR.AA-05, PR.DS-02

### HGX-NS-002 — DDoS Protection
- **Implementation**: Cloudflare edge; Envoy rate limiting; WAF rules; [11-devops/16 §7.2].
- **ISO 27001**: A.8.20

### HGX-NS-003 — Secure Configuration Baselines
- **Implementation**: CIS K8s benchmark (Level 2); Kyverno-enforced PSS restricted.
- **SOC 2**: CC6.6
- **ISO 27001**: A.8.9
- **CIS**: 4.1, 4.2

---

## Application Security

### HGX-AS-001 — Secure Development Lifecycle
- **Implementation**: Threat modelling for new components ([24-threat-model.md]); SAST (Semgrep, CodeQL, Snyk Code) in CI; dependency scanning; secrets scanning.
- **SOC 2**: CC8.1
- **ISO 27001**: A.8.25, A.8.28
- **NIST CSF**: PR.IR-01

### HGX-AS-002 — Input Validation & Output Encoding
- **Implementation**: Typed APIs (proto); CSP strict-dynamic; HTML sanitiser for rendered markdown; SQL via prepared statements (sqlc).
- **ISO 27001**: A.8.25, A.8.26

### HGX-AS-003 — Supply-Chain Security
- **Implementation**: SBOM CycloneDX 1.5; SLSA L3; Cosign signatures; Rekor transparency log; Connaisseur admission.
- **SOC 2**: CC8.1
- **NIST CSF**: ID.SC-03, PR.DS-06

### HGX-AS-004 — AI Safety & Sandbox
- **Implementation**: Guardrails, structured output enforcement, sandbox validation; [07-ai/10-llm-self-learning.md].
- **ISO 27001**: A.5.37

---

## Data Protection

### HGX-DP-001 — Data Classification
- **Implementation**: Labels on K8s workloads; `data_classification` JSON on entities; [27-data-retention-privacy.md §1].
- **ISO 27001**: A.5.12
- **NIST CSF**: ID.AM-05

### HGX-DP-002 — Data Retention & Minimisation
- **Implementation**: Retention matrix enforced by nightly purge jobs; PII masked in logs.
- **ISO 27001**: A.5.14, A.5.34
- **GDPR**: Art. 5(1)(c), 5(1)(e)

### HGX-DP-003 — Data Subject Rights
- **Implementation**: DSR workflows (access / rectify / erase / port); [27-data-retention-privacy.md §4].
- **GDPR**: Art. 15-22
- **CCPA**: §1798.100-1798.130

### HGX-DP-004 — Breach Notification
- **Implementation**: 72-hour regulator notification procedure; customer notification via status page + email.
- **GDPR**: Art. 33, 34
- **ISO 27001**: A.6.3

### HGX-DP-005 — Sub-Processor Management
- **Implementation**: Public sub-processor list; 30-day change notice; DPAs with every sub-processor.
- **ISO 27001**: A.5.19-A.5.22
- **GDPR**: Art. 28

---

## Logging & Monitoring

### HGX-LM-001 — Centralised, Tamper-Evident Audit Log
- **Implementation**: `audit.events` append-only table + Kafka topic; Merkle anchored to Rekor; [18-observability.md].
- **SOC 2**: CC7.1, CC7.2, CC4.1
- **ISO 27001**: A.8.15
- **NIST CSF**: DE.AE-03, DE.AE-08

### HGX-LM-002 — Alerting on Security Events
- **Implementation**: Alerts catalog (§Security); on-call rotation.
- **ISO 27001**: A.8.16, A.5.25

### HGX-LM-003 — Log Retention
- **Implementation**: Audit 365 d hot + long-term cold; security logs 1 y minimum; [27-data-retention-privacy.md §2].
- **SOC 2**: CC7.2
- **ISO 27001**: A.8.15

### HGX-LM-004 — Time Synchronisation
- **Implementation**: NTP on all hosts; timestamps in UTC; variance < 10 ms.
- **CIS**: 8.4

---

## Physical / Environmental (for on-prem customers)

Mostly inherited from customer data centre; HelixGitpx SaaS relies on sub-processor controls.

- **HGX-PE-001** — Data centre certifications (AWS / GCP / Azure + bare-metal DCs) — ISO 27001, SOC 2 Type II, PCI DSS, FedRAMP where applicable.
- **HGX-PE-002** — Media sanitisation on hardware decommissioning.

---

## Third-Party Management

### HGX-TP-001 — Sub-Processor Due Diligence
- **Implementation**: Annual security review + SOC 2 / ISO 27001 verification; on-change re-assessment.
- **SOC 2**: CC9.2
- **ISO 27001**: A.5.19-A.5.22

### HGX-TP-002 — Contractual Safeguards
- **Implementation**: DPAs with SCCs; security clauses in MSAs.
- **GDPR**: Art. 28

---

## Governance

### HGX-GV-001 — Information Security Policy
- **Implementation**: Maintained in this repo under `docs/policies/`; reviewed annually.
- **ISO 27001**: A.5.1

### HGX-GV-002 — Security Awareness Training
- **Implementation**: Onboarding + annual refresh; phishing simulations quarterly.
- **ISO 27001**: A.6.3
- **SOC 2**: CC1.4

### HGX-GV-003 — Risk Management
- **Implementation**: Risk register in this repo; quarterly review; threat model review.
- **ISO 27001**: Clause 6.1.2
- **NIST CSF**: GV.RM

### HGX-GV-004 — Internal Audit
- **Implementation**: Compliance team; annual external audit.
- **ISO 27001**: Clause 9.2

### HGX-GV-005 — Management Review
- **Implementation**: Quarterly exec review of security posture.
- **ISO 27001**: Clause 9.3

---

## Evidence Catalogue (highlights)

| Evidence | Location | Cadence |
|---|---|---|
| Access review (privileged users) | Spreadsheet-of-truth + audit log dashboard | Quarterly |
| Patch SLA reports | Vuln dashboard export | Monthly |
| DR drill run book + results | `postmortems/dr-drills/` | Quarterly |
| Pen-test report | GRC portal | Annual |
| Backup restore test | `runbooks/backup-restore-test.md` | Monthly |
| Security awareness completion | LMS export | Semi-annual |
| Incident responses | `postmortems/` | Per incident |
| Change records | GitHub + Argo CD | Continuous |
| Risk register | `docs/risks/` | Quarterly |
| Sub-processor list | Trust page | Continuous |

---

## Framework Cheat Sheet

Map external auditors expect:

**SOC 2 (Trust Services Criteria) → HelixGitpx**
- CC1 Control Environment → HGX-GV-001, HGX-GV-002
- CC2 Communication → HGX-GV-001
- CC3 Risk Assessment → HGX-GV-003
- CC4 Monitoring → HGX-LM-001
- CC5 Control Activities → entire control set
- CC6 Logical Access → HGX-AC-*
- CC7 System Operations → HGX-SO-*, HGX-LM-001
- CC8 Change Management → HGX-SO-001, HGX-AS-003
- CC9 Risk Mitigation → HGX-TP-*
- A1 Availability → HGX-SO-005
- C1 Confidentiality → HGX-CR-*, HGX-DP-*
- P1-P8 Privacy → HGX-DP-003, HGX-DP-004

**ISO 27001:2022 Annex A ↔ HelixGitpx** — see inline references throughout this doc.

---

## Update Discipline

Every material change to security posture updates this file. PRs that add/remove controls must:

1. Update the control table.
2. Update the affected control's implementation notes.
3. Link to the PR implementing the change.
4. Add / adjust evidence requirement.

Reviewed by the compliance team + at least one engineer on the affected system.
