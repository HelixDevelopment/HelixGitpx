# HelixGitpx Migration Guide

## 1. Introduction

Moving teams from a single Git host to HelixGitpx federation is a
migration you can do incrementally — no big-bang cutover, no team-wide
disruption.

### 1.1 Audience

- Team leads planning a migration.
- Engineering managers estimating effort.
- DevOps engineers executing the migration.

### 1.2 Migration paths

- **Chapter 2:** from GitHub.com (SaaS).
- **Chapter 3:** from GitHub Enterprise Server.
- **Chapter 4:** from GitLab.com (SaaS).
- **Chapter 5:** from self-hosted GitLab.
- **Chapter 6:** from Bitbucket.
- **Chapter 7:** from Azure DevOps.
- **Chapter 8:** from AWS CodeCommit.
- **Chapter 9:** from Gitea / Forgejo / Gitee / GitFlic / GitVerse /
  SourceHut.
- **Chapter 10:** multi-provider consolidation (you use two or more today).

### 1.3 Migration phases

1. **Bind** — add HelixGitpx as an additional upstream alongside your
   current host. No change for developers.
2. **Mirror-verify** — push/pull goes through HelixGitpx; verify parity
   with your existing host over 1–2 weeks.
3. **Flip** — switch the primary remote in developer workflows. Old
   host remains bound as a mirror.
4. **Graduate** — when comfortable, remove the old host binding, or keep
   it as insurance.

You can stop at any step. The whole point of federation is that you
never have to "finish" a migration.

### 1.4 What comes across

| Artifact | Migrated | Notes |
|----------|----------|-------|
| Git history | Yes | Every object, signed commits preserved. |
| Branches + tags | Yes | Identical names and refs. |
| PRs / MRs | Yes | Metadata normalized per adapter. |
| Issues | Yes | Canonicalised schema with provider-specific fields preserved. |
| Labels | Yes | Mapped via per-org rules. |
| Comments | Yes | With author attribution. |
| Webhooks | Re-issued | HelixGitpx issues HMAC-signed webhooks. |
| CI integrations | Update | Point runners at HelixGitpx URLs. |
| Personal access tokens | Re-issued | Managed in HelixGitpx per-user. |

### 1.5 What doesn't come across

- Provider-specific marketplace apps (GitHub Actions marketplace, GitLab
  apps). These remain on their provider; HelixGitpx emits equivalent
  events so you can trigger external CI from there.
- Provider billing. You are billed by HelixGitpx for HelixGitpx; the
  upstream host billing is independent.

### 1.6 Support

Email `support@helixgitpx.io` with the subject "Migration from <host>"
and we'll assign a migration engineer for Team+ plans.

---
