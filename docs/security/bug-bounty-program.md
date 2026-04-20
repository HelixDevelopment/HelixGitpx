# Bug Bounty Program

**Status:** Launching at GA on `<GA-DATE>`.
**Platform:** HackerOne (managed) — public program.
**Contact:** `security@helixgitpx.io` (for questions); submissions via H1.

## Scope

### In-scope

- `*.helixgitpx.io` production services.
- Official Android app, iOS app, Desktop apps.
- Any repository at `github.com/HelixGitpx/*`.

### Out-of-scope

- Staging / preview environments.
- DoS / spam / volumetric attacks.
- Social engineering.
- Physical attacks.
- Automated scanner output without manual validation.
- Third-party services (e.g. upstream Git hosts, Stripe).

## Rewards (USD)

| Severity | Range | Examples |
|----------|-------|----------|
| Critical | $10,000 – $25,000 | RCE, full auth bypass, tenant-to-tenant data leak |
| High | $2,500 – $10,000 | SSRF to internal metadata, stored XSS with session theft, privilege escalation |
| Medium | $500 – $2,500 | Reflected XSS, CSRF on state-changing endpoints, auth logic flaws |
| Low | $100 – $500 | Info disclosure, missing security headers with demonstrable impact |

Duplicates, informational, and out-of-scope reports are not rewarded.

## Rules

- Do not access, modify, or delete data that is not yours.
- Use test accounts you control.
- Coordinated disclosure: 90 days, extendable by mutual agreement.
- No public disclosure without written authorization.
- Provide reproducible PoC.

## Safe harbor

We will not pursue legal action against researchers acting in good faith within these rules.
