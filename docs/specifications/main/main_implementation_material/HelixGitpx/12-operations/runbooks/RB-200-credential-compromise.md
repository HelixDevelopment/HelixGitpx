# RB-200 — Suspected Credential Compromise

> **Severity default**: P1 (security)
> **Owner**: Security on-call
> **Last tested**: 2026-02-15 (tabletop)
>
> **This is a security runbook. Follow carefully. Everything is audited.**

## 1. Detection

Alert: `SuspectedCredentialCompromise` — anomalous login pattern detected (geo jump, concurrent sessions from distant locations, impossible travel).

Other triggers:
- Customer report via `security@helixgitpx.example.com`.
- External credential-leak scanner match against our token prefix (`hpxat_`).
- `RB-201` sibling: secret found in a public commit.
- Anomaly score from auth telemetry above threshold.

---

## 2. First 5 Minutes — Contain

Do not delay to diagnose. Contain first.

1. **Acknowledge** alert; open **security incident** channel with restricted access.
2. Identify the user/org involved (from alert labels).
3. **Force-revoke** active sessions:
   ```bash
   helixctl auth revoke-sessions --user <user-uuid> --all --reason="RB-200 suspected compromise"
   ```
4. **Invalidate** PATs and refresh-token families:
   ```bash
   helixctl auth revoke-pats --user <user-uuid> --all
   helixctl auth revoke-refresh-family --user <user-uuid> --all
   ```
5. Preserve the session/PAT rows for forensic comparison (they're revoked, not deleted).

---

## 3. Assess Scope

```bash
# Timeline of actions by the user in last 7 days
helixctl audit trace --user <user-uuid> --last=7d --format=jsonl > incident-trace.jsonl

# Orgs + resources touched
jq -r '.resource_kind + ":" + .resource_id' incident-trace.jsonl | sort -u

# Any write actions (concerning)
jq 'select(.outcome=="success" and (.action | test("create|update|delete|merge|push")))' incident-trace.jsonl
```

Look for:
- Repo or upstream credential changes.
- PR merges that look anomalous.
- New webhooks pointed at external domains.
- PAT creation (attacker establishing persistence).

---

## 4. Contain Further

Depending on what you see:

### 4.1 Webhook hijack attempt

- Remove the suspicious webhook.
- Rotate the repo webhook secret.

### 4.2 Upstream credential touched

- Rotate the upstream credential immediately via `helixctl upstream rotate-credentials`.
- Check the upstream's recent commits/activity for attacker activity.

### 4.3 Malicious PAT created during session

- Revoke that PAT.
- Audit what operations were performed with it.

### 4.4 Malicious merge

- If the PR is still mergeable, revert the merge (coordinate with repo maintainer).
- If code reached customer upstreams, fan-out the revert via HelixGitpx.

---

## 5. Engage Partners

- **Legal**: may need to involve privacy officer (72-h GDPR clock starts now if EU user data).
- **Customer success**: notify affected org's primary contact.
- **External**: if attacker used a third-party IdP to gain access, inform that IdP's security contact.

---

## 6. Preserve Evidence

- Export relevant audit log slice (`helixctl audit export --filter=...`).
- Snapshot relevant network logs from Cloudflare + Envoy for the user's IPs.
- Ensure evidence is write-locked (our audit log is append-only; this is guaranteed).
- Hand off to security team lead.

---

## 7. Recover

- User chooses a new authentication method (MFA re-enrolment mandatory).
- Review and re-generate any org-level secrets the user had access to.
- Re-authorise sessions only after user has confirmed identity via out-of-band channel.

---

## 8. Communicate

### 8.1 To the user (templated)

> Subject: [HelixGitpx] Security notice — your account was protected
>
> We detected unusual activity on your HelixGitpx account and have pre-emptively revoked your sessions and API tokens. Please sign in again and re-enrol MFA. If you did not initiate the activity below, contact security@helixgitpx.example.com immediately.
>
> [bullet summary of observed activity]

### 8.2 To the org admin (if org-scoped resources touched)

> Subject: [HelixGitpx] Security incident involving your organisation
>
> A user in your organisation had credentials revoked after suspected compromise. We have identified the following resources that may have been affected: [list]. We have taken the following actions: [list]. No further action is required of you, but we recommend [recommendations].

Avoid details in emails that could help an attacker confirm success; prefer richer details via an authenticated support portal link.

### 8.3 Public

Typically no public notification unless the incident was widespread.

---

## 9. Regulator Notification

If personal data breach is confirmed (GDPR Art. 33):

- DPO to notify supervisory authority within 72 h of awareness.
- Document: facts, likely consequences, measures taken.
- If high risk to individuals (Art. 34), also notify affected individuals.

---

## 10. Post-Incident

- Root-cause: how did the attacker get in? (phishing, credential stuffing, token leak, IdP compromise, insider?)
- Corrective actions: MFA enforcement, leak-scanner improvements, rate-limit tightening, user education.
- Postmortem within 10 business days; sanitised version for customers.
- Update threat model if a new TTP observed.

---

## 11. Related

- RB-201 (Secret leaked in public commit)
- RB-202 (Malicious plugin behaviour)
- [24-threat-model.md] — attacker tactics
- [27-data-retention-privacy.md §10] — breach notification
- Legal retention: all artefacts retained ≥ 7 y for security incidents.

---

## 12. Drill

Tabletop quarterly. Red-team exercise semi-annually.

Scenario cards under `chaos/security/red-team-cards/`.
