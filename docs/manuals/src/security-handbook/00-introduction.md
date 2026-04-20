# HelixGitpx Security Handbook

## 1. Introduction

This handbook is the single document to hand to a customer CISO or an
auditor who asks "how do you secure HelixGitpx?". It cross-references the
public security docs (`docs/security/`) and the operator runbooks
(`docs/operations/runbooks/`).

### 1.1 Audience

- Customer CISOs and security engineers performing due diligence.
- SOC 2 / ISO 27001 auditors.
- Procurement teams requiring a security attestation.
- Internal security engineers.

### 1.2 What you'll find

- **Chapter 2:** threat model.
- **Chapter 3:** data protection (at rest, in transit, in use).
- **Chapter 4:** identity and access management.
- **Chapter 5:** authorization model (OPA + NeMo Guardrails).
- **Chapter 6:** supply chain integrity (SLSA, Cosign, SBOM).
- **Chapter 7:** monitoring, detection, and response.
- **Chapter 8:** disaster recovery.
- **Chapter 9:** compliance posture (SOC 2, ISO 27001, GDPR).
- **Chapter 10:** responsible disclosure, bug bounty.

### 1.3 Short version

HelixGitpx uses:

- TLS 1.3 everywhere ingress; SPIFFE/SPIRE mTLS everywhere east-west.
- AES-256 at rest; KMS-managed keys; 90-day rotation.
- OPA Rego bundles for authorization; signed, diff-reviewed, versioned.
- NeMo Guardrails around every LLM surface.
- Cosign-signed images, SBOMs, SLSA Level 3 at GA.
- Falco + OpenTelemetry + Prometheus for runtime detection.
- Annual external pen-test; public bug bounty on HackerOne.
- SOC 2 Type I at GA; Type II tracked post-GA.
- ISO 27001 gap-analysed at GA; certification tracked post-GA.

### 1.4 Evidence packet

For audits, we provide on request under NDA:

- SOC 2 Type I report.
- ISO 27001 gap analysis.
- Latest pen-test executive summary.
- SBOM and Cosign attestation chain.
- DR drill results (last 4 quarters).

Contact: `security@helixgitpx.io`.

---
