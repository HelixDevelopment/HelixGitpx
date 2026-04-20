# ISO/IEC 27001:2022 Gap Analysis

**Assessor:** internal (M8 pre-certification gap).
**Target:** full Annex A coverage + ISMS maturity by end of year 2 (post-GA).

## Annex A control gap summary

| Control | Theme | Status | Gap |
|---|---|---|---|
| A.5.1 Policies for information security | Organizational | ✅ | Covered by `docs/security/policies/`. Need annual sign-off cadence. |
| A.5.2 Information security roles | Organizational | ✅ | CISO role defined. Deputy needed. |
| A.5.3 Segregation of duties | Organizational | ⚠️ | Production deploy = any 2 engineers currently; tighten to production-break-glass role. |
| A.5.7 Threat intelligence | Organizational | ⚠️ | No formal TI feed yet. Plan: integrate CISA KEV + GitHub advisory. |
| A.5.23 Cloud services security | Organizational | ✅ | K8s configs reviewed + CIS benchmarked. |
| A.5.26 Response to incidents | Organizational | ✅ | Runbooks + incident template. |
| A.5.30 ICT readiness for business continuity | Organizational | ✅ | DR drill quarterly. |
| A.6.3 Awareness, education, training | People | ⚠️ | Monthly training required; currently ad-hoc. |
| A.6.8 Confidentiality or NDAs | People | ✅ | Employee + contractor NDAs in place. |
| A.7.4 Physical security monitoring | Physical | N/A | Fully cloud-native; no HQ server room. |
| A.8.1 User endpoint devices | Technological | ⚠️ | MDM enrollment for all work devices — 85% coverage currently. |
| A.8.2 Privileged access rights | Technological | ✅ | Just-in-time via Teleport; audited. |
| A.8.5 Secure authentication | Technological | ✅ | OIDC + WebAuthn (staff); OIDC + TOTP (customers). |
| A.8.9 Configuration management | Technological | ✅ | GitOps via Argo CD; all configs in git. |
| A.8.12 Data leakage prevention | Technological | ⚠️ | Endpoint DLP not yet; Zscaler / Nightfall eval planned. |
| A.8.16 Monitoring activities | Technological | ✅ | OTel + Prometheus + Falco. |
| A.8.23 Web filtering | Technological | ⚠️ | No corp web filter; rely on endpoint + awareness. |
| A.8.24 Use of cryptography | Technological | ✅ | TLS 1.3; mTLS (SPIFFE); KMS-managed keys. |
| A.8.25 Secure development lifecycle | Technological | ✅ | CI gates: SAST (Semgrep), DAST (ZAP), SCA (Trivy), mutation testing. |
| A.8.26 Application security requirements | Technological | ✅ | ASVS L2 baseline; OWASP Top-10 tests. |
| A.8.28 Secure coding | Technological | ✅ | Conventional Commits + DCO; two-approver PRs. |
| A.8.33 Test information | Technological | ⚠️ | Production PII ban in staging — tooling exists; enforcement tightening. |

## Prioritized gap closure (year 1 post-GA)

1. **A.5.3** — split break-glass role from everyday deploy.
2. **A.6.3** — monthly awareness training (Kontra / HackEDU).
3. **A.8.1** — MDM to 100%.
4. **A.8.12** — DLP tool selection + pilot.

## Out-of-scope for gap analysis

A.7.* (physical) controls N/A; A.8.23 (web filtering) accepted risk with compensating controls (endpoint).
