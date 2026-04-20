# 35 — Support Handbook

> **Document purpose**: Guide the support and customer-success functions in responding to customer requests quickly, consistently, and with the right level of technical depth. Pairs with [30-sla.md] which sets the external commitments.

---

## 1. Support Tiers

| Tier | Staffed by | Handles |
|---|---|---|
| **T1** | Generalist support engineers | Onboarding, how-to, account issues, basic troubleshooting |
| **T2** | Senior support + subject-matter experts | Complex issues; escalates to engineering |
| **T3** | Engineering on-call (product team that owns the component) | Incidents, data issues, deep bugs |
| **Customer Success Manager (CSM)** | Dedicated per Enterprise account | Proactive account health, training, upsell |
| **Technical Account Manager (TAM)** | Per Enterprise account | Technical advisory, roadmap alignment |

---

## 2. Channels

| Channel | Audience |
|---|---|
| `support@helixgitpx.example.com` | Team / Business |
| In-app support chat | Team+ (business hours + 24/7 on Enterprise) |
| Enterprise portal | Enterprise only |
| Community forum | Free tier + anyone |
| Status page / RSS | All |
| Phone (Enterprise) | P1 only, 24/7 |
| Slack Connect (Enterprise) | Dedicated channels |

---

## 3. Intake & Triage

When a ticket arrives, T1 classifies:

1. **Severity** (P1–P4) per [30-sla.md §2].
2. **Plan** (Free / Team / Business / Enterprise).
3. **Category** (Account, Billing, Git ingress, Adapter, Conflict, AI, API, Security, Other).
4. **Is this a known incident?** If yes, reference status page.

Ticket is auto-routed based on (severity × plan × category). Response SLA timers start.

---

## 4. Canned Templates (edit before sending)

### 4.1 Welcome / Acknowledge

> Hi {first_name},
>
> Thanks for reaching out. I'm {agent_name} from the HelixGitpx support team.
>
> I'm looking into this now and will get back to you by {sla_time}. In the meantime, could you share:
> - The affected org / repo slug
> - Your approximate time of occurrence (with timezone)
> - Any error messages you've seen
>
> Thanks!

### 4.2 Known Incident

> Hi {first_name},
>
> Thanks for reporting this. We're currently investigating an incident that matches what you're seeing: {incident_link}
>
> We'll update the status page as we know more. You can subscribe to updates there. I'll reach back out personally once resolved, too.

### 4.3 How-To / Documentation Pointer

> Hi {first_name},
>
> What you're describing is covered in our docs here: {docs_link}
>
> The short version: {one-paragraph summary}
>
> Let me know if any of it is unclear or if you hit a snag while trying this!

### 4.4 Need More Info

> Hi {first_name},
>
> To help debug, could you run this and paste the output?
>
> ```
> helixctl doctor org {org}
> ```
>
> That captures environment info and recent activity — nothing sensitive is included. If you'd rather send it privately, reply to this email as an attachment.

### 4.5 Escalating to Engineering

> Hi {first_name},
>
> I'm escalating this to our engineering team. Based on what you've described, it's beyond routine troubleshooting.
>
> I'll stay as your point of contact; engineering will reply through me. Expected response: within {sla}.
>
> Do you have additional context or a time window the issue is most reproducible? Any info helps.

### 4.6 Feature Request

> Hi {first_name},
>
> Thanks for this idea — I can see how it would help.
>
> I've logged it in our product backlog under {ticket_ref}. Our product team reviews new ideas every Thursday. Not all ideas are pursued, but all are considered. You'll get an email if it's picked up or if we decide to pass.

### 4.7 Billing Dispute

> Hi {first_name},
>
> I understand — let me dig into your invoice and usage.
>
> I've pulled the meters for {period} and I see the following breakdown: {details}. If this doesn't match what you were expecting, happy to walk through it or adjust where appropriate. Our billing team has the authority to issue credits up to {amount} without escalation.

### 4.8 Resolved

> Hi {first_name},
>
> Glad we got this sorted. Summary for your records:
> - Reported: {summary_of_issue}
> - Cause: {root_cause}
> - Fix: {resolution}
>
> If anything like it happens again, referencing ticket {id} will speed things up. Have a great week!

---

## 5. Troubleshooting Playbooks

For each common category, T1 has a step-by-step playbook.

### 5.1 "My push isn't appearing on upstream X"

1. Confirm upstream X is **enabled** and not in shadow mode: UI → Upstreams.
2. Check Sync History on the repo (`Repo → Sync History`).
3. Look for rate-limit or auth errors.
4. If rate-limited → RB-130 guidance → suggest credential rotation or upstream concurrency reduction.
5. If auth-failed → ask customer to rotate credentials.
6. If no error: check Kafka lag dashboard for the sync-orchestrator group (T2 if lag > 5 min).

### 5.2 "A conflict keeps reappearing"

1. Identify the conflict case (`helixctl conflict show <id>`).
2. Look at the set of upstreams involved; find which one keeps producing conflicting updates.
3. Usually: parallel pushes to the same branch on two upstreams.
4. Recommend: pin a primary upstream; educate on workflow.
5. If genuinely stuck → escalate.

### 5.3 "AI proposal was wrong"

1. Collect: case id, proposal chosen, actual correct result.
2. File feedback via the system (`helixctl conflict feedback`).
3. Reassure: 5-min undo window; outside window reverse-PRs flow.
4. Escalate to AI team so the example is added to evals.

### 5.4 "Login is stuck"

1. Have customer open incognito + clear cookies.
2. Verify IdP (OIDC issuer in `auth.identity_providers`) is healthy.
3. Check status page for recent identity provider outages.
4. If PAT-based: verify PAT scope.

### 5.5 "How do I…?"

1. Search docs first; point to the relevant page.
2. If docs are insufficient, write a quick how-to and file a ticket for the docs team to absorb.

---

## 6. Escalation Triggers (T1 → T2 / T3)

Escalate immediately when:

- Any P1 ticket (skip T1 triage queue; straight to ops).
- Data integrity concerns (lost commits, missing events, corrupted refs).
- Security concerns (credential exposure, suspicious activity).
- Customer reports "our prod is down because of HelixGitpx" (regardless of our view).
- > 3 hops of back-and-forth without resolution.
- Legal / privacy / regulatory question.

---

## 7. Tooling

- **Ticket system**: Zendesk / Freshdesk / custom.
- **helixctl**: primary debugging tool; support agents have restricted-scope elevated access (audit-logged).
- **Grafana read-only** dashboards tagged for support.
- **Runbook index** pinned in the support chat workspace.
- **helixctl customer impersonate --read-only** — assumes a customer's read-only view to reproduce UI issues; audited.

Never use customer impersonation for write actions without explicit consent.

---

## 8. Customer Data Handling

- Minimum necessary principle: don't pull what you don't need.
- Log everything you do: the audit trail is there for the customer's protection and yours.
- Never share customer data in community forums; only in the private ticket.
- PII scrubbed from shared screenshots before attaching.
- Sensitive content: use the secure portal, not email.

---

## 9. Working with Engineering

- T2 maintains a weekly "escalations" meeting with each product team.
- Escalation format: short Markdown doc with customer context, impact, timeline, what we've tried, ask.
- Engineering on-call response time matches severity:
  - P1: 15 min
  - P2: 2 h
- Engineering owns the root-cause; support owns the customer relationship and comms.

---

## 10. Quality & Metrics

- **FRT (First Response Time)**: plan-SLA adherence.
- **TTR (Time to Resolution)**: measured per severity.
- **CSAT** (post-ticket survey): target ≥ 4.7/5.
- **Deflection**: % of tickets resolved without engineering.
- **Repeat-customer ratio**: low is better (fewer repeating issues).
- **Knowledge-base age**: any article > 1 y old reviewed.

Weekly dashboard review.

---

## 11. Knowledge Base

- Every ticket that could have been a docs read is a candidate for a new KB article.
- Quarterly "doc sprint" — support writes, docs team edits.
- KB links included in AI-drafted responses (with human review before sending).

---

## 12. Proactive Support (Enterprise / CSM)

- Monthly health check call.
- Usage + margin review.
- Roadmap preview (NDA'd).
- Renewal planning 120 days out.
- Executive-sponsor program for strategic accounts.

---

## 13. Crisis Communication

During big incidents:

- Single source of truth = status page.
- Support writes customer-facing updates that track the incident channel internally.
- Legal reviews any statement about data loss / breach.
- CS manages exec-to-exec outreach where appropriate.

---

## 14. Continuous Improvement

- Monthly retro: top support categories, proposed product / docs fixes.
- Quarterly: incoming-volume forecast vs. actual; staffing planning.
- Annual: full handbook review.

---

*— End of Support Handbook —*
