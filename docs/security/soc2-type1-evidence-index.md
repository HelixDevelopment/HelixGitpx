# SOC 2 Type I Evidence Index

**Audit window:** Point-in-time at GA + 30 days.
**Auditor:** `[VERIFY-AT-INTEGRATION]` — shortlist: Prescient Assurance, Johanson Group, A-LIGN.
**Trust Services Criteria:** Security (required), Availability, Confidentiality.

## Common Criteria (CC) mapping

| TSC | Control | Evidence artifact |
|-----|---------|-------------------|
| CC1.1 | Org structure documented | `docs/company/org-chart.md`; board minutes |
| CC1.4 | Personnel policies | `docs/people/handbook.md`; onboarding logs |
| CC2.1 | Information systems boundary defined | `docs/specifications/.../00-core/01-vision-scope-constraints.md` |
| CC2.2 | Internal communication | Slack retention policy; #all-helixgitpx channel logs |
| CC3.1 | Risk assessment process | `docs/security/risk-register.md`; quarterly review notes |
| CC5.1 | Control activities | All policies in `docs/security/policies/` |
| CC6.1 | Logical access — identification & authentication | Keycloak config export; OIDC flow doc |
| CC6.2 | Access provisioning / deprovisioning | Ticket IDs from JIRA for last 6 onboards/offboards |
| CC6.3 | Role-based access | OPA bundle v2 (`enforcement.rego`); RBAC policies |
| CC6.6 | Logical access monitoring | `audit-service` logs retained 90d+; SIEM screenshots |
| CC6.7 | Data transmission encryption | TLS configs (cert-manager issuers); mTLS via SPIFFE/SPIRE |
| CC6.8 | Malicious software | Falco runtime rules; container image scans (Trivy) |
| CC7.1 | Monitoring | Prometheus + alertmanager config; Grafana dashboards |
| CC7.2 | Anomaly detection | Falco + OTel trace sampling; incident review logs |
| CC7.3 | Evaluation of security events | Incident response playbook `docs/security/incident-response.md` |
| CC7.4 | Security incident response | Post-mortems from M8 chaos drills |
| CC7.5 | Disaster recovery | `tools/dr/dr-drill-runbook.md`; DR drill results |
| CC8.1 | Change management | GitHub PR history; Argo CD sync logs |
| CC9.1 | Risk mitigation | Risk register with mitigations column |
| CC9.2 | Vendor management | Vendor inventory + review cadence |

## Availability (A1)

- A1.1 — SLO definitions: `docs/specifications/.../12-operations/slos.md`.
- A1.2 — Monitoring + alerting: Prometheus rules committed.
- A1.3 — Recovery: `tools/dr/dr-drill-runbook.md` quarterly.

## Confidentiality (C1)

- C1.1 — Data classification: `docs/security/data-classification.md`.
- C1.2 — Retention: `docs/security/retention-policy.md` (defaults: 2y audit, 90d PII).

## Evidence collection

Automated where possible. `tools/soc2/collect.sh` pulls artifacts into `/tmp/soc2-evidence-<date>/`.
