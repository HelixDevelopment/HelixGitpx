# 30 — Service Level Agreements (SLA)

> **Document purpose**: Formalise **external service-level commitments** HelixGitpx makes to paying customers across plans — availability, performance, support response, incident communication, remedies. Internal SLOs that drive engineering practice live in [18-observability.md].

---

## 1. Plan → SLA Matrix

| Capability | Free | Team | Business | Enterprise |
|---|---|---|---|---|
| **Monthly uptime** | best-effort | **99.5 %** | **99.9 %** | **99.95 %** (99.99 % option) |
| **API latency p99** (read) | — | 500 ms | 300 ms | **300 ms**, audited monthly |
| **API latency p99** (write) | — | 1 s | 500 ms | **500 ms** |
| **Git push replication** p99 to first upstream | — | 15 s | 5 s | **5 s** |
| **Live-event delivery** p99 | — | 1 s | 500 ms | **500 ms** |
| **Conflict escalation** notification | — | 1 h | 15 min | **1 min** |
| **Scheduled maintenance notice** | — | 48 h | 5 d | 14 d |
| **RPO** (data loss window) | — | 5 min | 30 s | 30 s; contractual |
| **RTO** (recovery time) | — | 1 h | 15 min | 15 min |
| **Support response — P1** | community | 4 h | 1 h | **15 min** 24/7 |
| **Support response — P2** | community | next business day | 4 h | 2 h |
| **Support channels** | forum | email | email + chat | email + chat + dedicated TAM |
| **Incident review reports** | — | summary | detailed | full RCA + action items |
| **Credits on breach** | — | optional | yes (see §5) | contractual |
| **Dedicated success manager** | — | — | optional | yes |

Free tier has no SLA — community best-effort.

---

## 2. Definitions

- **Downtime**: any 1-minute interval where the HelixGitpx service is unavailable as measured by our external synthetic probes across ≥ 3 independent vantage points.
- **Monthly Uptime Percentage**: `(TotalMinutes − Downtime) / TotalMinutes × 100 %`.
- **Excluded Downtime**: scheduled maintenance announced per the notice window; emergency security maintenance; force majeure; customer-caused outages; issues outside HelixGitpx's control (customer IdP down, upstream Git provider outages unless they're within our direct control).
- **Incident**: any event customer-visible enough to appear on the public status page.
- **Severity** (external-facing):
  - **P1**: Complete or majority unavailability; data at risk.
  - **P2**: Significant functionality impaired; workaround exists.
  - **P3**: Minor degradation.
  - **P4**: Cosmetic / low impact.

---

## 3. Measurement

- Uptime measured by HelixGitpx-operated synthetic probes from ≥ 3 regions every 60 s covering the public API and web shell.
- Latency measured server-side at the API gateway via OTel histograms; rolled up per calendar month.
- Independent third-party monitoring (e.g. DataDog Synthetic / Catchpoint) validates our numbers; customers can subscribe to these probes for their own verification.
- Reports published monthly on the status page and available via API.

---

## 4. Exclusions

- Beta / preview features.
- Issues caused by customer misconfiguration (e.g. invalid IdP, broken webhooks customer-side).
- Upstream Git provider downtime (we report on it, we don't cause it). Enterprise customers can get targeted credits if we materially delayed routing around the outage.
- Force majeure events.
- Actions by the customer that exceed published limits in a plan.
- DNS issues not caused by HelixGitpx-managed DNS.

---

## 5. Service Credits

### 5.1 Formula (Business / Enterprise default)

| Monthly Uptime % | Credit |
|---|---|
| < 99.95 % and ≥ 99.9 % | 10 % of month's fee |
| < 99.9 % and ≥ 99.0 % | 25 % |
| < 99.0 % and ≥ 95.0 % | 50 % |
| < 95.0 % | 100 % |

### 5.2 Request

- Customer must request within 30 days of the end of the affected month.
- Request via `billing@helixgitpx.example.com` or support ticket with incident references.
- Credit applied to next invoice. Non-cash, non-transferable.

### 5.3 Maximum

- Credits capped at 100 % of one monthly fee per SLA event.
- Enterprise may contract alternative remedies (penalties, off-contract options).

---

## 6. Support Response & Resolution

### 6.1 Response Times

Time from ticket receipt to **substantive** (not auto) reply.

| Severity | Team | Business | Enterprise |
|---|---|---|---|
| P1 | 4 h | 1 h | **15 min 24/7** |
| P2 | next business day | 4 h | 2 h |
| P3 | 3 business days | 1 business day | 8 business hours |
| P4 | 10 business days | 5 business days | 5 business days |

### 6.2 Resolution Targets

Not guaranteed, but targeted; communicated per incident.

- P1: workaround within 4 h; permanent fix within 30 d.
- P2: workaround within 2 business days; permanent within 60 d.
- P3 / P4: next reasonable release.

### 6.3 Languages

- English supported 24/7.
- German, French, Spanish, Japanese, Serbian — business hours.
- Enterprise can contract additional languages.

---

## 7. Incident Communication

- **Status page** (`status.helixgitpx.example.com`) updated:
  - P1: initial post within 15 min of detection; updates at least every 30 min.
  - P2: initial within 30 min; updates every hour.
  - P3: single post, as needed.
- Incident channel auto-created in Slack for engaged customers.
- Webhook / RSS / email subscription to status page available to all.
- **Post-incident**:
  - Business / Enterprise: detailed incident report within 5 business days.
  - Team: summary within 10 business days.
- Blameless tone; explicit timeline, root cause (or preliminary), corrective actions, customer-impact summary.

---

## 8. Maintenance Windows

- **Routine**: announced ≥ 5 business days in advance for Team+, 14 d for Enterprise.
- **Expedited security maintenance**: announced as soon as safely possible; may be hours.
- **Emergency**: minimum feasible notice; post-incident communication expected.
- Windows chosen to minimise customer impact per region; off-peak local time.

---

## 9. Capacity & Fair Use

- Published rate limits apply (see [07-rest-api.md §10]).
- Abuse or disproportionate use outside contract triggers throttling, not an SLA breach.
- Customers can contract higher limits.

---

## 10. Security Commitments

- Vulnerability remediation SLAs:
  - Critical: 7 days.
  - High: 30 days.
  - Medium: 90 days.
  - Low: 180 days.
- Security advisory publication for critical issues affecting customer data.
- Responsible disclosure honour policy (see [08-security/24-threat-model.md §11]).

---

## 11. Data Protection Commitments

- Encryption in transit (TLS 1.3) and at rest (AES-GCM-256).
- GDPR / CCPA / PIPEDA compliance.
- Data portability via export API.
- Retention per [27-data-retention-privacy.md].
- Annual SOC 2 Type II for Business / Enterprise.
- DPA (Data Processing Addendum) signed by default for all paid plans.

---

## 12. Compliance Artefacts

Available on request (NDA may apply):

- SOC 2 Type II audit report.
- ISO 27001 certificate + SoA.
- SLSA L3 provenance.
- CSA CAIQ / STAR.
- Cyber Essentials (UK), ANSSI SecNumCloud (where applicable).
- Pen-test summary (executive-level).

Enterprise-tier customers get direct access to the trust portal.

---

## 13. Business Continuity

- Multi-region capable.
- Quarterly DR drill; results summarised to Enterprise customers.
- Business Continuity & Disaster Recovery Plan documented; excerpts shareable under NDA.

---

## 14. Escalation Path (Enterprise)

1. Primary support contact.
2. Technical Account Manager (TAM).
3. Head of Support.
4. CTO.

Agreed in onboarding; contact details in shared "Customer Escalation Doc".

---

## 15. Change Notifications

- API changes per [29-api-versioning-deprecation.md] with 12-month deprecation for majors.
- Subprocessor changes: 30-day notice.
- Material policy changes: 30-day notice.

---

## 16. Termination

- **Customer-initiated**: any time (reflected on next billing cycle for monthly); annual contracts per terms.
- **HelixGitpx-initiated** for cause: 30-day cure period where feasible.
- Data export window: 90 days post-termination.
- Data deletion: within 30 days of export window closing; confirmation on request.

---

## 17. Amendments

- Standard SLA applies unless superseded by signed MSA / enterprise agreement.
- Amendments proposed by HelixGitpx: 30-day opt-out window where changes are material.

---

## 18. Governing Terms

This SLA complements the Terms of Service and Master Services Agreement. In case of conflict, MSA > DPA > SLA > Terms of Service > Documentation.

Governing law and jurisdiction per the relevant contract. For EU customers, SCCs govern cross-border transfers.

---

*— End of SLA —*
