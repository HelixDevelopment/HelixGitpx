# Feature Flag Catalog

> Every flag in HelixGitpx lives here while it's active — purpose, owner, rollout status, targeting, expected removal date. Flags without a **removal plan** accumulate tech debt and confuse operators.

Tooling: **OpenFeature** client + **Unleash** self-hosted as the evaluator. Flags are YAML-managed in `platform/flags/` and deployed via Argo CD.

---

## Conventions

- Name: `kebab-case`, prefixed by domain.
- Types: `release`, `experiment`, `ops` (kill-switch), `permission`, `plan-gate`.
- Every flag has an **owner**, **created-on**, and **sunset target**.
- **Sunset target = expected removal date** — hard prompt in CI if overdue.
- Strategies: percentage rollout, org allowlist, user allowlist, environment-only, always-on (stale).

---

## Active Flags

### Release rollouts

| Flag | Owner | Target | Status | Notes |
|---|---|---|---|---|
| `repo.search-hybrid-rrf` | search-team | 100 % | stabilising | RRF fusion of Meilisearch + Qdrant; remove flag 2026-07-01 |
| `ai.conflict-auto-apply-v2` | ai-team | 25 % prod, 100 % staging | experiment | v2 model with new confidence heuristic |
| `live.ws-backpressure-v2` | platform | 10 % | experiment | Adaptive window sizing |
| `ui.conflict-resolver-three-pane` | frontend | opt-in beta | release | New UI; replaces old two-pane |
| `kmp.mobile-push-resume` | mobile | 50 % | experiment | Resume-on-FCM-push for Android |
| `release.canary-analysis-multi-metric` | sre | 100 % | stabilising | New Argo Rollouts analysis template |

### Operational kill-switches (ops)

| Flag | Owner | Default | Effect |
|---|---|---|---|
| `ai.enabled` | ai-team | on | Global AI circuit breaker; turn off in outages |
| `ai.cloud-routing-allowed` | security | off | Whether any cloud-LLM fallback may run |
| `adapter.github-disabled` | platform | off | Stops calls to GitHub temporarily (upstream outage) |
| `adapter.gitlab-disabled` | platform | off | — |
| `adapter.bitbucket-disabled` | platform | off | — |
| `live-events.emit-disabled` | platform | off | Stops live-event delivery; events still persisted |
| `sync.fanout-paused` | platform | off | Pause all fan-out to upstreams; local state updates continue |
| `ci.block-deploys` | sre | off | Global deploy freeze |
| `writes.read-only` | sre | off | Reject all write operations; reads + exports ok |
| `billing.meter-disabled` | billing | off | Stops metering (safety against billing bug) |
| `plugin.loading-disabled` | security | off | Prevents any WASM plugin load |

### Permission / policy toggles

| Flag | Owner | Scope | Description |
|---|---|---|---|
| `policy.require-mfa-on-admin` | security | org | Enforce MFA for admin-class actions (default on for Business+) |
| `policy.require-signed-commits-main` | security | org | Branch-protection default for `main` |
| `policy.require-signed-tags` | security | org | Signed tags required; default on Enterprise |
| `plugin.allow-unsigned` | security | global | Default off; enabling requires audit approval (dev-only) |
| `ai.allow-user-data-training` | ai/privacy | org | Default off in EU; on in other regions with opt-out |

### Plan gates

| Flag | Owner | Gate |
|---|---|---|
| `feature.audit-export-long-retention` | billing | Business+ only |
| `feature.residency-pinning` | billing | Enterprise only |
| `feature.dedicated-inference-pool` | billing | Enterprise only |
| `feature.sso-saml` | billing | Business+ only |
| `feature.cloud-ai-routing-opt-in` | billing | paid plans |
| `feature.on-device-ai-preview` | billing | desktop stable ≥ 1.2 |

### Experiments

| Flag | Owner | Hypothesis | Exit criteria |
|---|---|---|---|
| `conflict.ai-rationale-display` | research | "Showing AI rationale increases accept rate ≥ 5 %" | Compare 90 % vs 10 % for 30 days |
| `ui.pr-summary-sticky-card` | frontend-research | "Sticky summary reduces review time 15 %" | Eye-tracking sample + qual feedback |
| `onboarding.guided-first-repo` | growth | "Guided onboarding reduces time-to-first-push" | T1→T2 conversion +10 % |
| `notify.mobile-collapsed-batch` | mobile | "Digest reduces churn" | Uninstall rate ↓ 5 % |

---

## Retired Flags (last 90 days, for reference)

| Flag | Removed | Outcome |
|---|---|---|
| `ui.dark-mode` | 2026-03-04 | Default on |
| `repo.uuidv7-ids` | 2026-02-11 | Default on |
| `kmp.network-change-reconnect` | 2026-02-18 | Default on |
| `ai.lora-per-task` | 2026-03-22 | Default on |

Retired flags stay listed for 90 days then move to `docs/flags/archive/`.

---

## Lifecycle

1. **Propose**: open PR adding flag definition + initial targeting. Reviewer checks purpose, owner, sunset target.
2. **Develop**: feature gated by flag; default off in prod; on in dev/staging for implementers.
3. **Launch internal**: enable for HelixGitpx's own dogfood org.
4. **Early access**: opt-in beta for design partners.
5. **Rollout**: percentage ramp; analysis via Grafana annotations.
6. **100 %**: monitor for 14 days.
7. **Remove flag** before sunset; both `true` and `false` branches removed from code; release note captures the change.

Tooling:

- `helixctl flag get <n>` — current evaluation for a principal.
- `helixctl flag set <n> --strategy=percent --value=10` — interactive change (audit-logged).
- `helixctl flag stale` — list flags past sunset; CI warning.

---

## Evaluation Consistency

- Flag decisions are **stable per principal** within a session (hash-based).
- Same principal gets the same decision across services (avoid UI / backend mismatch).
- Decisions are propagated in request context so every service sees the same variant.

---

## Audit

- Every flag change audited.
- Ops flags flipping from default fires `ops_flag_flipped` event at P2 severity → Slack.
- Kill-switch flipping (e.g. `writes.read-only`) is a P1-worthy event; incident ticket opened automatically.

---

## Discipline

- ≤ 100 active flags globally; above that triggers a cleanup sprint.
- Sunset discipline enforced: overdue flags block PRs in affected service after 30 days past sunset.
- Quarterly flag-hygiene review by platform + product leads.

---

## Template (new flag)

```yaml
name: ui.conflict-resolver-three-pane
type: release
owner: frontend
created: 2026-04-01
sunset: 2026-08-01
description: >
  New three-pane conflict resolver UI replacing the two-pane. Measures
  reduction in resolution time and accept rate.
strategies:
  default:
    variation: "off"
  environments:
    staging: { variation: "on", percent: 100 }
    prod:    { variation: "on", percent: 10, strategy: "gradual" }
variations:
  - key: "on"
    weight: 10
  - key: "off"
    weight: 90
metrics:
  primary: "ui_conflict_time_to_resolve_seconds"
  guardrails: ["ui_error_rate", "slo_burn_rate_1h"]
```
