# HelixGitpx Troubleshooting Guide

## 1. Introduction

Solutions to problems you will actually hit. Ordered from "first day with
HelixGitpx" to "an incident at 03:00".

### 1.1 Audience

- End users troubleshooting their own workflow.
- Operators debugging production incidents.
- Support engineers assisting customers.

### 1.2 What you'll find

- **Chapter 2:** my push is slow or hanging.
- **Chapter 3:** my upstream shows stale commits.
- **Chapter 4:** I see a conflict in the inbox — now what?
- **Chapter 5:** AI responses are nonsense or refused.
- **Chapter 6:** OPA denied something that should have passed.
- **Chapter 7:** the web app won't load or shows 502.
- **Chapter 8:** mobile push notifications aren't arriving.
- **Chapter 9:** CI runner can't authenticate.
- **Chapter 10:** operator: pod CrashLoopBackOff.
- **Chapter 11:** operator: Kafka consumer lag spike.
- **Chapter 12:** operator: CNPG failover loop.
- **Chapter 13:** operator: cert-manager issuance failure.

### 1.3 How to ask for help

When filing an issue or ticket, attach:

1. Your HelixGitpx version: `helixgitpx --version` (CLI) or the footer
   of any web page.
2. The request ID from the failing operation (X-Request-Id header or
   "Copy request ID" button in the web app).
3. The time window in UTC.
4. The exact commands, URLs, or inputs that failed.
5. Any relevant OPA decision ID.

### 1.4 Links

- Status page: [status.helixgitpx.io](https://status.helixgitpx.io).
- Runbooks (operators): `docs/operations/runbooks/`.
- DR runbook: `tools/dr/dr-drill-runbook.md`.
- Support: `support@helixgitpx.io`.

---
