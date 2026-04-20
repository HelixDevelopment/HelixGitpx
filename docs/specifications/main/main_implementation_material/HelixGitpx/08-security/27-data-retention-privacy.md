# 27 — Data Retention, Privacy & Data Subject Rights

> **Document purpose**: Define how long HelixGitpx **keeps data**, how it **protects personal data**, and how we **honour data-subject rights** (GDPR, CCPA, PIPEDA, and comparable regimes). This is the authoritative reference for privacy engineering and support teams.

---

## 1. Data Categories

| Category | Examples | Classification |
|---|---|---|
| **Account data** | Email, username, display name, avatar | Personal |
| **Authentication** | Password hashes (if used), MFA secrets, session tokens (hashed) | Secret |
| **Identity linkage** | OIDC `sub`, issuer, raw claims | Personal |
| **Usage logs** | IP address, user agent, timestamps, actions | Personal + Telemetry |
| **Source code** | Git repos + LFS objects | Customer data (treated as confidential; not used for training) |
| **Metadata** | Issues, PRs, comments, labels | Personal (authorship) + Customer data |
| **Audit log** | Admin actions, logins, policy decisions | Personal + Compliance |
| **Billing** | Plan, usage meters, invoices | Personal (if individual billing) |
| **AI interactions** | Prompts, outputs, feedback | Personal (user-linked) |
| **Device identifiers** | Per-device keypair IDs | Personal |
| **Observability** | Metrics, traces (mostly aggregate) | Telemetry (low-personal) |

---

## 2. Retention Matrix

| Data | Default retention | Min | Max | Controlling customer knob |
|---|---|---|---|---|
| Account data | While active | indefinite | on erasure | user deletion |
| Sessions | 24 h idle / 14 d absolute | 24 h | 90 d | org policy |
| Refresh tokens | rotating; family TTL 14 d | 7 d | 90 d | org policy |
| PATs | per-token expiry (max 1 y) | 1 d | 1 y | per-user |
| Login attempts | 90 d | 30 d | 365 d | org policy |
| Audit log | 365 d (Team), 7 y (Enterprise) | 30 d | 10 y | plan |
| Source code (active) | indefinite | indefinite | indefinite | repo deletion |
| Soft-deleted repos | 30 d tombstone then purge | 1 d | 90 d | org policy |
| Comments / issues / PRs | indefinite (lifecycle of repo) | — | — | repo deletion |
| AI training-source feedback | 18 months | 30 d | 5 y | org opt-out |
| AI prompt-run records | 90 d (personal-linked) | 30 d | 1 y | org policy |
| Application logs | 30 d hot + 90 d cold | 7 d | 1 y | org policy |
| Metrics | 400 d downsampled | 30 d | 2 y | retention config |
| Traces | 7 d (errors 90 d) | 1 d | 90 d | retention config |
| Webhook deliveries | 30 d | 7 d | 180 d | — |
| DLQ events | 30 d | 7 d | 180 d | — |
| Backups | 35 d (daily) + 12 months (monthly) | — | — | legal |
| Billing records | 7 y (SOX / tax) | 7 y | 10 y | legal |

**Tombstone vs. purge**: deletion marks the record inaccessible ("tombstone") and triggers a purge job after the retention expires (default 30 d). Irreversible after purge.

---

## 3. Principles

1. **Data minimisation** — we collect only what is necessary. IPs / user-agents anonymised after 90 days in metric stores.
2. **Purpose limitation** — data collected for one purpose is not repurposed silently.
3. **Storage limitation** — everything has a retention; retention is enforced by automated purge jobs.
4. **Integrity & confidentiality** — encrypted at rest (AES-GCM-256) and in transit (TLS 1.3 + mTLS). Keys managed by Vault + customer-supplied KMS (enterprise).
5. **Lawful basis** — Contract for service-essential processing; Legitimate Interest for security logs; Consent for marketing and AI training (opt-out by default for training).
6. **Transparency** — this doc, the privacy policy, and the in-app privacy dashboard tell users what we store and why.

---

## 4. Data-Subject Rights (DSR)

Users and admins can exercise DSRs via **Settings → Privacy → Data Requests** or by emailing **privacy@helixgitpx.example.com**. All requests handled within 30 days (GDPR standard).

### 4.1 Right of Access

- **What happens**: a ZIP is generated containing all personal data associated with the user's account — profile, sessions, identity links, audit trail, AI interactions, billing references, and cross-references to repos/issues/PRs where they've contributed.
- **Technical**: Temporal workflow `AccessDataRequest` fans out to each service, collects, encrypts with user's PGP if provided, and delivers via signed download link.

### 4.2 Right to Rectification

- Profile self-edit for most fields. Complex fields (e.g. historical audit entries) require support intervention and are logged as a correction event — original remains for integrity with a correction annotation.

### 4.3 Right to Erasure

Two flavours:

1. **Account erasure**: account deactivated, PII scrubbed across services.
   - Comments and commits authored by the user are re-attributed to a "deleted-user" handle; the substantive content remains (to preserve project history and other users' rights).
   - Exceptions: **billing records** (legal retention); **security logs** (legitimate interest, usually 1 y).
2. **Specific data erasure**: targeted removal of a piece of data (e.g. a comment).

- **Technical**: `ErasureWorkflow` walks the services; each service has an `Erase(user_id)` handler that locates, tombstones, and (post-retention) hard-deletes personal data.
- **Across backups**: erasure applied in rolling fashion as backups age out (explained clearly to the user).

### 4.4 Right to Restriction / Objection

- "Pause my processing" — temporary read-only state.
- Object to AI training: toggles a flag that exclude this user's feedback from training corpora. Historical already-trained models are not retrained for one individual (standard industry practice; communicated clearly).

### 4.5 Right to Data Portability

- See §9 (Export) in [25-migration-guide.md]. Machine-readable, interoperable formats.

### 4.6 Automated Decision-Making

- HelixGitpx does not make decisions about users using solely automated means that have legal or similarly significant effects on them.
- AI-assisted conflict resolution affects *repositories*, not *individuals*; and a 5-minute undo window plus human escalation preserves human review.

---

## 5. Special Categories & Minors

- HelixGitpx does not knowingly process special categories of personal data (health, political opinions, etc.).
- We do not knowingly target users under 16 (or other local age floors). Where we learn of underage use, the account is deactivated.

---

## 6. Regional / Residency Controls

- **EU residency**: paid orgs can pin `primary_region` to the EU. Data-at-rest stays in EU; cross-region replication disabled or EU-only.
- **UK residency**: same mechanism; post-Brexit UK-only storage optional.
- **US Public Sector**: dedicated region with FedRAMP-aligned controls (roadmap, not yet certified).
- **Custom**: enterprise customers can specify custom residency via on-prem deployment.

---

## 7. Cross-Border Transfers

- Outside chosen residency, transfers happen only with an appropriate mechanism: **SCCs**, **adequacy decisions**, or **BCRs**.
- Transfer impact assessments performed per regime.
- Sub-processors listed publicly on our trust site; changes notified 30 d in advance.

---

## 8. Sub-Processors

| Role | Sub-processor | Location |
|---|---|---|
| Email delivery | (e.g.) Postmark / SES | US / EU |
| Payments | Stripe | US |
| Identity (optional) | OIDC IdP (customer's) | — |
| CDN / WAF | Cloudflare | Global |
| Observability (hosted tier) | Grafana Cloud (optional) | EU / US |

Customers using on-prem deployments can reduce or eliminate sub-processors.

---

## 9. Logging Privacy Controls

- **Redaction**: collectors strip emails, tokens, keys, SSH private-key material, and common secret patterns.
- **Sampling**: high-volume debug logs sampled to reduce footprint.
- **Masking in traces**: attributes in OTel traces that reference personal data are tagged `pii=true` and anonymised in the default sample.
- **Row-level access** to logs: Grafana Enterprise / OSS + OPA proxy provides per-tenant scoping.

---

## 10. Breach Notification

- Internal: security on-call opens incident; legal engaged within 1 h for candidate breach.
- **72-hour** clock starts when a personal-data breach is identified (GDPR).
- Customer notification: via primary contact + status page.
- Regulator notification: handled by DPO through the required channel.
- Evidence collected and retained per legal guidance.

---

## 11. Vendor Security Assessments

- Privacy impact assessments for any new data flow.
- Sub-processor due diligence annual + on change.
- Contract clauses: SCCs / DPA baked into every contract.

---

## 12. Children's Data

We do not target users under 16. If learned, account is disabled.

---

## 13. Data Retention Enforcement

Automated purge jobs run nightly per table. Metrics visible to compliance:

- `helixgitpx_retention_jobs_total{status}`.
- `helixgitpx_retention_records_purged_total{table}`.
- Alerts if any job fails for > 7 d.

---

## 14. Customer Configuration UI

Org admins can configure (subject to plan):

- Retention shorteners (below product defaults).
- Residency pinning (enterprise).
- AI training opt-out (default on for EU; opt-in elsewhere).
- Audit retention extension (enterprise, up to 10 y).

UI surfaces a "What happens if…" preview before applying destructive policies.

---

## 15. Privacy Tests

- **Unit**: every service has an `Erase(user_id)` test.
- **Integration**: cross-service erasure integration test in CI.
- **Contract**: DSR workflow e2e monthly.
- **Fuzz**: malformed export payloads.

---

## 16. Documentation & Training

- New engineers complete a privacy module within 30 days.
- Annual refresher for all engineers with access to production.
- Legal signs off on each major release for privacy impact.

---

*— End of Data Retention, Privacy & DSR —*
