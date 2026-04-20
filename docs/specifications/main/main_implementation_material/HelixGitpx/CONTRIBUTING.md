# Contributing to HelixGitpx

Thanks for your interest in HelixGitpx. We welcome contributions — code, docs, translations, bug reports, design feedback, community help. This document explains how to get involved effectively.

---

## Before You Start

1. **Read the [Developer Guide](12-operations/20-developer-guide.md)** — local setup, coding standards, testing expectations.
2. **Read the [Code of Conduct](CODE_OF_CONDUCT.md)** — it applies to every interaction in our spaces.
3. **Search existing issues and discussions** before opening new ones.

---

## Ways to Contribute

### Report a Bug

1. Check if the issue is already filed.
2. File at `github.com/vasic-digital/helixgitpx/issues/new/choose` with the **Bug** template.
3. Include:
   - HelixGitpx version (footer of the web app or `helixctl version`).
   - Steps to reproduce (short + concrete).
   - What you expected vs. what happened.
   - Logs / screenshots (redact secrets).
   - Environment (OS, browser, client, region).

### Propose a Feature

- **Small, well-scoped features**: open an issue with the **Feature** template; include use cases, alternatives, and a rough design.
- **Large or architectural changes**: start an **RFC** in `docs/rfcs/` — see [RFC template](docs/rfcs/0000-template.md). RFCs get discussed publicly, sometimes leading to an ADR.

Feature requests are welcome but not all are accepted. We'll explain why if we pass.

### Fix a Bug / Build a Feature

1. Comment on the issue to say you're taking it on; get rough agreement with a maintainer first — saves wasted effort.
2. Fork, branch from `main`, make your change.
3. Write tests — units + integration as appropriate. Coverage ≥ 100 % gate is real.
4. Follow our [coding standards](12-operations/20-developer-guide.md#5-coding-standards).
5. Open a PR to `main`. Fill in the template.

### Improve Docs

Docs live in this repo. Small edits can be PR'd directly; larger restructures benefit from a quick issue first.

### Translate

Crowdin project: `https://crowdin.com/project/helixgitpx`. Join, request a role for your locale, translate. See [32-i18n-l10n.md](05-frontend/32-i18n-l10n.md).

### Report a Security Issue

**Don't open a public issue.** Email **security@helixgitpx.example.com** with details; PGP key at `docs/security/pgp.asc`. See our [responsible disclosure policy](08-security/36-trust-center.md#13-responsible-disclosure).

---

## Pull Request Requirements

- **Branch naming**: `feat/<short-desc>`, `fix/<short-desc>`, `docs/<short-desc>`, `chore/<short-desc>`.
- **Commit messages**: [Conventional Commits](https://www.conventionalcommits.org/). Example:
  - `feat(repo): add transfer API endpoint`
  - `fix(conflict): correct CRDT merge when labels empty`
  - `docs(api): clarify idempotency key TTL`
- **Signed commits** required on merges to `main` (Gitsign or GPG).
- **Scope**: keep PRs small (< 400 lines of diff ideal; > 1000 requires justification).
- **Tests**: include; adjust coverage gates if needed.
- **Docs**: update relevant docs + ADRs if architectural.
- **CI**: must be green. Flaky? Flag it, don't ignore it.

### Review Process

- Two approvals required on `main`.
- Reviewers respond within one business day.
- Address feedback in new commits (don't force-push until requested); squash at merge.
- Reviewer trust is earned — comments, not orders. Disagreements discussed; escalate to maintainers if needed.

### Merge Queue

We use a merge queue. After approval, add the `ready-to-merge` label. The queue batches PRs, runs CI against HEAD, and merges in order.

---

## CLA / DCO

Contributions are licensed under **Apache-2.0**. We require sign-off:

```
Signed-off-by: Your Name <your.email@example.com>
```

(Equivalent of DCO; add via `git commit -s`.)

If contributing on behalf of a company, ensure you're authorised and mention it in the PR description. A Corporate CLA may be required for substantial contributions; we'll guide.

---

## First-Time Contributors

Look for issues labelled:
- `good-first-issue` — small, well-scoped, good for onboarding.
- `help-wanted` — we'd welcome help.
- `docs` — documentation improvements.

Not sure where to start? Ask in `#helixgitpx-community` on our Slack or the Discussion forum; maintainers will suggest.

---

## Community Norms

- **Be kind.** Be specific. Be brief.
- **Credit others.** Name co-authors in commit messages.
- **Assume good faith.** Read charitably before responding.
- **Disagreements are fine.** Ad hominem is not.
- **English is our working language.** If you're more comfortable in another language, do your best in English and we'll help smooth it.

Harassment or abuse → see [Code of Conduct](CODE_OF_CONDUCT.md). Reporting: `conduct@helixgitpx.example.com`.

---

## Governance

HelixGitpx is stewarded by **vasic-digital**. Technical decisions go through:

1. **Issue / RFC** — public discussion.
2. **Maintainer review** — consensus-seeking.
3. **ADR** — recorded decisions, immutable.
4. **Implementation** — PR, review, merge.

Maintainers:
- Project maintainers — listed in `MAINTAINERS.md`.
- Area maintainers — per module; listed in each module's README.

Becoming a maintainer: sustained, high-quality contributions + culture fit. Existing maintainers nominate; consensus to add.

---

## Release Notes

Changelogs are generated from commits. Your PR should have a clear title — it becomes the changelog entry. Substantive changes can also add a "release notes" section to the PR body.

---

## Recognition

All contributors are listed on our `CONTRIBUTORS.md` and, for significant contributions, credited in release notes. We also run a quarterly "community highlights" post on the blog.

---

## Questions?

- **Discussion forum**: for general questions.
- **Slack (`helixgitpx-community`)**: real-time chat, office hours.
- **Office hours**: Thursdays 15:00 CET.
- **Email**: `community@helixgitpx.example.com` for anything that doesn't fit above.

Welcome aboard. 🦎
