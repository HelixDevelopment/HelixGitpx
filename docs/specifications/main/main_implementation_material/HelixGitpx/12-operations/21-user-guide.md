# 21 — User Guide

> **Document purpose**: Help end-users, org admins, and repo maintainers **use HelixGitpx effectively**. Covers account setup, upstream connections, repo management, conflict resolution, AI features, and everyday workflows.

---

## 1. Getting Started

### 1.1 Create an Account

Visit `app.helixgitpx.example.com` and sign in with your identity provider (Google, GitHub, Microsoft, or your corporate SSO via OIDC/SAML). First login automatically provisions your user profile.

### 1.2 Create or Join an Organisation

- **Create**: enter a unique slug, display name, default visibility, and region.
- **Join**: accept an invitation email or use an invite link.

Orgs are the container for people, teams, repos, and upstream connections.

### 1.3 Install Clients

- **Web**: any modern browser. PWA-installable for desktop-like experience.
- **Mobile**: Google Play, App Store, F-Droid.
- **Desktop**: download installers from `helixgitpx.example.com/download` (Windows, macOS, Linux).
- **CLI**: `helixctl` — `curl -fsSL helixgitpx.example.com/install.sh | sh`.

All clients share your account and data.

---

## 2. Connect Upstream Git Services

HelixGitpx synchronises your repositories with the Git hosts you already use. You keep working on GitHub, GitLab, Gitee, etc. — HelixGitpx federates.

### 2.1 Add an Upstream

1. **Org → Upstreams → Add upstream**.
2. Pick a provider (GitHub, GitLab, Gitee, GitFlic, GitVerse, Bitbucket, Codeberg, Gitea, Forgejo, SourceHut, Azure DevOps, AWS CodeCommit, Generic Git).
3. Authenticate: OAuth (preferred), Personal Access Token, App token, or SSH key — as appropriate for the provider.
4. Choose **shadow mode** (recommended first time): HelixGitpx will read but not push for a soak period, so you can verify mirroring is correct before going live.
5. Save. The provider is listed as **Connected**.

### 2.2 Enable / Disable / Disconnect

- **Disable**: pauses synchronisation both ways. Your data is preserved. Re-enable any time.
- **Disconnect**: removes credentials and stops tracking entirely. Webhooks are de-registered.
- **Rotate credentials**: rotates the token without interrupting flow.

### 2.3 Provider-Specific Notes

- **GitHub**: we recommend using a GitHub App (more granular scopes, higher rate limits).
- **GitLab self-hosted**: enter your base URL.
- **Azure DevOps**: both cloud and Server/TFS supported.
- **Generic Git**: covers any SSH/HTTPS Git host without an API. PRs/issues will be simulated inside HelixGitpx since the upstream doesn't support them.

---

## 3. Create a Repository

1. **Repos → New repo**.
2. Slug, display name, description.
3. Visibility (public / private / internal).
4. Toggle **"create on all enabled upstreams"** (default on).
5. Choose primary upstream (optional) — conflict resolution prefers this source.

HelixGitpx provisions the repo across every enabled upstream in parallel. You get a list of clone URLs you can use from any tool.

### 3.1 Migrate an Existing Repo

- **Repos → Import** → paste a Git URL → we clone mirror, optionally pull issues/PRs/releases via the provider API, and then fan out to other upstreams.

### 3.2 Archive / Delete

- **Archive**: makes the repo read-only everywhere; no more sync. Reversible.
- **Delete**: tombstones locally; requires explicit flag to also delete on upstreams.

---

## 4. Everyday Workflows

### 4.1 Pushing Code

Clone from any upstream you like — e.g. `git clone git@github.com:acme/infra.git` — and push as usual. HelixGitpx detects the push (via webhook or periodic poll) and replicates it to all other enabled upstreams within seconds.

You can also push directly to HelixGitpx's ingress:

```
git remote add helixgitpx git@push.helixgitpx.example.com:acme/infra.git
git push helixgitpx
```

### 4.2 Branches, Tags, Releases

- All created the same way as on any Git host.
- Tags / releases are mirrored with assets.
- Release assets larger than 100 MB are stored on the HelixGitpx asset store and referenced on each upstream.

### 4.3 Pull Requests / Merge Requests

Open a PR on any upstream. HelixGitpx creates **mirror PRs** on every other upstream with a note linking back. Reviews, comments, and merge actions federate.

When a PR merges anywhere, it's merged everywhere. "Once merged, always merged" — HelixGitpx enforces this invariant.

### 4.4 Issues

Open on any upstream; mirrored across all others. Comments replay in order. Labels, milestones, and assignees use CRDT merging — concurrent edits never lose data.

### 4.5 Branch Protection

- **Repo → Settings → Branch Protection**.
- Rules: required reviews, required status checks (aggregated across upstreams), signed commits, up-to-date branch, force-push prevention.
- HelixGitpx enforces these **centrally** and, where supported, also pushes the rules to upstreams.

---

## 5. Conflicts — How HelixGitpx Handles Divergence

Occasionally, the same ref ends up pointing at different commits on different upstreams, or the same label is added/removed concurrently. HelixGitpx detects this and offers resolution.

### 5.1 Types of Conflicts

- **Ref divergence** (most common): two upstreams have different head SHAs.
- **Rename collision**: same file renamed to different paths.
- **PR state mismatch**: merged on one, closed on another.
- **Metadata race**: labels/milestones concurrent edits.
- **Tag collision**: same tag, different commits.
- **LFS divergence**: extremely rare; quarantined for manual review.

### 5.2 How They Resolve

Most conflicts resolve automatically:

- **CRDT merging** for metadata.
- **Policy-based** for refs (prefer primary, prefer signed, prefer newer).
- **Three-way merge** where a common ancestor exists and the merge is clean.
- **AI proposal** where merging isn't obvious; each proposal comes with a rationale and confidence score.

### 5.3 Your Role

- **Inbox → Conflicts**: shows cases awaiting your attention.
- For each case: review side-by-side diff, accept / edit / reject the AI proposal, or propose your own.
- Auto-applied resolutions have a **5-minute undo window** — one click reverses on every upstream.
- You'll get notifications on desktop / mobile when a conflict needs you.

---

## 6. AI Features

### 6.1 Conflict Resolution Assistant

Already described above. The assistant explains what it proposes and why.

### 6.2 PR Summariser

Click **Summarise this PR** to get a concise, human-readable description of what changed, surfaced reviewers likely to care, and potential risks.

### 6.3 Review Assistant

Inline comments suggest common issues (unused imports, missing error handling, unclear naming). Accept / dismiss per suggestion. Dismissed suggestions feed back into the model.

### 6.4 Label & Milestone Suggestions

When you open a new issue or PR, we suggest labels based on title + body + repo history.

### 6.5 Semantic Search

Ask questions in plain language: "Where do we initialise the Redis client?" — we search the repo's code semantically.

### 6.6 Chatops

In the repo chat drawer, ask natural-language things: "open a PR from `foo` to `main` and assign @alice". The assistant confirms before executing.

### 6.7 Privacy

- AI runs on HelixGitpx's self-hosted models by default; your code never leaves our environment.
- Cloud models only used if your org explicitly opts in.
- You can opt out of AI feedback learning via **Settings → Privacy → Do not use my interactions to improve AI**.

---

## 7. Notifications

### 7.1 Channels

- In-app (mobile / desktop / web push).
- Email.
- Slack / Microsoft Teams / Discord webhooks.
- Custom webhook for your own systems.
- SMS (optional, enterprise tier).

### 7.2 Subscriptions

Configure per scope (org, repo, PR, issue). Sensible defaults; tweak via **Settings → Notifications**.

---

## 8. Personal Access Tokens

- **Settings → Developer → Tokens**.
- Scopes (fine-grained): `repo:read`, `repo:write`, `org:read`, etc.
- Expiration required (max 1 year).
- The token is shown **once** on creation.

Tokens are recommended for automation; for interactive use, OIDC via the UI/CLI is preferred.

---

## 9. `helixctl` CLI

The CLI mirrors the UI for automation and power users.

```bash
helixctl login                                       # OIDC device-flow
helixctl repo create acme/infra --visibility=private --create-on-all
helixctl repo list acme
helixctl conflict list acme/infra --status=escalated
helixctl conflict resolve <case_id> --strategy=prefer_primary
helixctl upstream add acme --provider=github --oauth
helixctl sync trigger acme/infra
helixctl pr list acme/infra --state=open
helixctl event stream --scope=repo:acme/infra      # tail live events
```

CLI configuration in `~/.helixgitpx/config.yaml`.

---

## 10. Integrations

- **GitHub Actions / GitLab CI / Jenkins**: our action/orb/plugin publishes status back to HelixGitpx, which aggregates across upstreams.
- **VSCode extension**: conflict view, PR review, event subscriber.
- **IntelliJ plugin**: same features.
- **Terraform provider**: manage orgs, repos, upstreams, branch protection as code.
- **Kubernetes Operator** (enterprise): sync repo-definitions-as-CRDs.

---

## 11. Accessibility

HelixGitpx is WCAG 2.2 AA compliant. If you find accessibility issues, please report via **Help → Accessibility Feedback**.

Keyboard shortcuts: press `?` anywhere to view.

---

## 12. Data Portability

- **Export**: full data export per repo (Git bundle + issues/PRs/releases JSON) or per org.
- **Import**: reverse; also supports direct imports from GitHub / GitLab / Bitbucket.

---

## 13. Privacy & Data Requests

- **Settings → Privacy → Data Requests**: export, rectify, erase, object, restrict.
- GDPR + CCPA + PIPEDA workflows supported.

---

## 14. Billing

- Plans: Free, Team, Enterprise. Features scale accordingly (active repos, upstream count, AI tokens, audit retention, regional pinning).
- **Usage dashboard** shows current meters.
- Hard limits warn first, throttle next, block last.

---

## 15. Support

- **Docs**: <https://docs.helixgitpx.example.com>.
- **Status**: <https://status.helixgitpx.example.com>.
- **Community forum**: <https://community.helixgitpx.example.com>.
- **Support**:
  - Free: community.
  - Team: email, 1 business day.
  - Enterprise: 24/7 with SLA.

---

## 16. Troubleshooting FAQ

**My push isn't showing up on upstream X.**
Check the upstream status in **Org → Upstreams**; the upstream may be disabled or rate-limited. Look at **Repo → Sync History** for the failed job and retry.

**I got an "AI proposal rejected by sandbox" message.**
The model suggested a patch that failed compile/lint/tests. Either accept a different proposal or resolve manually.

**I pushed, but nothing happened.**
Verify your remote URL and credentials. Your push may have been blocked by branch protection.

**A conflict keeps reappearing.**
Someone may be pushing to the same branch on multiple upstreams simultaneously. Consider pinning a primary upstream, or coordinate with collaborators.

**The app shows old data after reconnecting.**
Pull to refresh (mobile) or hit the reload icon (web); local cache is rebuilt on reconnect.

---

*— End of User Guide —*
