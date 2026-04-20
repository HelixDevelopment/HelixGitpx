# 31 — Release Management & Branching Strategy

> **Document purpose**: Define how HelixGitpx **moves code from PR to production**: branching, tagging, release trains, hotfix flow, version pinning, and customer communication. Complements [17-devops-cicd.md] with the release-specific playbook.

---

## 1. Principles

1. **Trunk-based development.** `main` is always deployable. Short-lived branches. Merge via queue.
2. **Continuous delivery.** Services ship to staging multiple times per day automatically.
3. **Release trains, not release gates.** A scheduled train leaves even if some cargo isn't ready; late changes wait for the next one.
4. **Feature flags decouple deploy from release.** Code can be live in production while invisible to users.
5. **Every artefact signed & traceable.** Commit → build → image → Helm chart → Argo CD sync → running pod, linked by digests.

---

## 2. Branch Model

### 2.1 Branches

- `main` — trunk. Always green. Always deployable.
- `release/X.Y` — active release train; minor releases cut from here.
- `release/X.Y.Z-hotfix` — emergency hotfix branch off a release tag.
- Feature branches — created from `main`, short-lived (ideally < 3 days).

No long-lived feature branches (`develop`, `next`, etc.). Merge small and often.

### 2.2 Branch Protection

- `main`: require PR, ≥ 2 approvals, all status checks green, merge queue.
- `release/*`: require PR, ≥ 2 approvals from release managers.
- Force-push blocked. Deletions blocked.
- Signed commits required on `main` and `release/*`.

### 2.3 Merge Queue

- Queue batches PRs; re-runs CI against HEAD before merging.
- Prevents semantic conflicts that pass individually but fail together.
- Merged via squash by default; merge commit for release-train cuts.

---

## 3. Versioning

### 3.1 Services

Each service is semver-independent: `<service>@<MAJOR>.<MINOR>.<PATCH>`.

Example tags: `repo-service@1.4.2`, `auth-service@2.0.1`.

### 3.2 Platform

The **platform release** aggregates a snapshot of service versions and the docs/charts at that moment:

`helixgitpx@YYYY.MM.<n>` (CalVer), e.g. `helixgitpx@2026.05.1`.

Customers using our SaaS see the platform version in the footer of the web app.

### 3.3 SDKs / CLI

Each SDK and the CLI follow semver independently, tracking the underlying API major.

### 3.4 Helm Chart

Chart version follows platform CalVer; `appVersion` reflects the headline platform release.

---

## 4. Release Cadence

| Line | Cadence | Support |
|---|---|---|
| **Continuous** (services → staging) | many times/day | rolling |
| **Stable platform release** | weekly | N and N-1 |
| **LTS release** | quarterly | 18 months |
| **Mobile** | bi-weekly stable, nightly channel | N and N-1 |
| **Desktop** | bi-weekly stable | N and N-1 |
| **Plugins** | owner-driven | per plugin |

Mondays are train-cut days. Fridays are freeze days (no new trains).

---

## 5. Release Train Workflow

1. **T-3 days**: PRs targeting next train labelled `target/2026.05.1`.
2. **T-1**: release manager on rotation triages labelled PRs; decides which will make the cut.
3. **T (Monday)**:
   - Cut `release/2026.05` from `main` (if not yet) or tag `release-2026.05.1` on existing branch.
   - Run **release-hardening pipeline** (superset of normal CI: extended test suites, DAST, chaos run, perf baseline check, supply-chain audit).
   - Auto-generate release notes from commit messages + curated highlights.
4. **T+1**: deploy to canary (5 % of prod traffic).
5. **T+2**: if canary metrics green, promote to 100 %.
6. **T+3**: publish release notes, update customer docs, send announcement.

Hotfix and security patches bypass the train; see §8.

---

## 6. Deploy Path

```
PR merged to main
    ↓
Image built + signed (Cosign keyless) + SBOM attached
    ↓
Service test suite passes → staging deploy (auto, Argo CD)
    ↓
Observed in staging for ≥ 30 min (staging SLOs green)
    ↓
Release manager approves promotion (for prod)
    ↓
GitOps PR updating image tag in helixgitpx-platform/envs/prod-eu
    ↓
Argo CD reconciles → Argo Rollouts canary:
  5% → observe → 25% → 50% → 100%
    ↓
Promoted → post-deploy checks → done
```

Rollback is a revert of the GitOps PR; Argo CD re-syncs.

---

## 7. Canary Analysis

Canary runs analysis templates (see `18-manifests/argo-application.yaml`):

- **success-rate** ≥ 99.5 %
- **p99 latency** within SLO
- **burn-rate** ≤ 6× baseline
- **error budgets** not depleted
- **adapter circuit** not opening

Any failure during the 10-minute analysis window → **automatic rollback** + page.

---

## 8. Hotfix Flow

When a critical bug is found in a release:

1. Open issue labelled `hotfix` with incident details.
2. Branch `release/X.Y.Z-hotfix` from the affected tag.
3. Fix + add test.
4. Fast-track review (can be one reviewer if P1).
5. Cut `release-X.Y.Z+1` tag.
6. Deploy directly to affected environments (skip the train).
7. Backport to `main` (mandatory) and other active release branches.
8. Post-mortem within 5 business days.

---

## 9. Security Releases

- Same as hotfix but with **security on-call** involved.
- Coordinated disclosure if third parties affected.
- CVE may be requested via our CNA partner.
- Advisory published on status page and through release notes.
- No public mention of the fix before patches are generally available.

---

## 10. Mobile & Desktop Release

### Mobile

- **Stable** bi-weekly: tagged from `release-mobile-X.Y.Z`; AAB / IPA built; submitted to Play Store / App Store.
- **Beta** channel: TestFlight / Play internal-testing.
- **Nightly** channel: Firebase App Distribution (Android) / TestFlight (iOS).
- Desktop updaters poll signed manifest; honour channel subscription.

### Desktop

- Conveyor-produced MSIX / DMG / AppImage / .deb / .rpm.
- Signed with Developer ID / EV cert / GPG.
- Auto-update honours enterprise MDM preferences.

---

## 11. Rollback

- **Service**: revert GitOps PR; Argo CD re-syncs within minutes.
- **Database migration**: reverse migration in next release (expand-contract); avoid destructive rollbacks.
- **Kafka schema**: backward compatibility means old consumers remain functional; broken producers rolled back.
- **Mobile / desktop**: can't un-ship to stores but can "unpublish" and push forced-upgrade minimum versions for clients.

---

## 12. Feature Flags & Dark Launches

- **Per-org / per-user / per-percentage** targeting.
- Ship dark behind a flag; enable for HelixGitpx's own dogfood org; then Enterprise design partners; then global.
- Flags outlive their rollout only when there's a reason; stale-flag reports in CI surface cleanup candidates.

Kill-switches for every risky change documented in runbooks.

---

## 13. Deprecations & Migrations

Coordinated per [29-api-versioning-deprecation.md].

Release notes call out deprecations prominently; migration guides committed in the same PR as the deprecation.

---

## 14. Release Notes Template

```markdown
## HelixGitpx 2026.05.1 (Stable)

**Released**: 2026-05-04

### Highlights
- 🚀 New: ...
- ⚙️ Improved: ...
- 🐛 Fixed: ...

### Detailed changes
- Conflict resolver: ...
- API gateway: ...

### Deprecations
- `RepoService.ListRepos` deprecated in v1; use `SearchRepos` (sunset 2027-05-01).

### Upgrade notes
- Any manual action required.

### Security advisories
- None / Links.

### Stats
- PRs merged: 87
- Contributors: 22
- Test coverage: 100.3 %
- Mutation score: 92.1 %
```

---

## 15. Internal Tooling

- `helixctl release cut --train 2026.05.1 --target staging`
- `helixctl release promote --train 2026.05.1 --target prod-eu --strategy canary`
- `helixctl release list --status active`
- `helixctl release rollback --env prod-eu --to 2026.04.3`
- `helixctl release compare --from 2026.04.3 --to 2026.05.1` (shows PRs, metrics delta)

All gated behind MFA + audit.

---

## 16. Release Captaincy

- Weekly rotation. Release captain duties:
  - Triage labelled PRs.
  - Kick off train.
  - Run checks.
  - Handle rollbacks.
  - Own communication to internal stakeholders.
  - Coordinate with support on expected customer questions.
- One-week shadow for new captains.

---

## 17. Metrics of Release Health

- **Lead time** PR open → prod.
- **Deployment frequency** per service.
- **Change failure rate** (rollbacks / total deploys).
- **Mean time to recovery** (MTTR) from P1 incidents.
- Target: **elite** DORA metrics (lead < 1 day, deploys daily+, CFR < 15 %, MTTR < 1 h).

Dashboards on the platform-wide Grafana.

---

*— End of Release Management —*
