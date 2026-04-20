# 25 — Migration Guide

> **Document purpose**: Help teams **migrate onto HelixGitpx** from GitHub, GitLab, Bitbucket, Gitea, or self-hosted Git — without disrupting ongoing work. Also covers migration *off* HelixGitpx for portability.

---

## 1. Migration Strategy Overview

HelixGitpx's federation model makes migration fundamentally non-destructive: you don't have to "cut over" — you **add HelixGitpx as an upstream**, let it synchronise, and then optionally make HelixGitpx your primary.

Three strategies, tailored to organisational readiness:

| Strategy | When to use | Cutover time |
|---|---|---|
| **Federation-first** | Default. You keep using existing host; HelixGitpx mirrors and adds value incrementally. | Zero — never cut over |
| **Gradual cutover** | You want HelixGitpx as primary over time. Feature flags per repo. | Weeks per repo |
| **Big-bang migration** | Consolidating multiple hosts onto HelixGitpx | Planned maintenance window |

---

## 2. Prerequisites

Before starting any migration:

- [ ] HelixGitpx org created with billing plan matching your needs.
- [ ] At least one org admin with MFA enabled.
- [ ] Upstream credentials available (OAuth app install, PAT with appropriate scopes, or SSH key).
- [ ] Network access to HelixGitpx from your CI runners and dev machines.
- [ ] Notification: announce the migration plan to your team 2 weeks ahead.

---

## 3. Provider-Specific Playbooks

### 3.1 From GitHub

**Capabilities mapped 1:1**: repos, branches, tags, PRs, issues, releases, labels, milestones, protected branches, webhooks, deploy keys, LFS.

**Near 1:1**: Actions (status checks mirror; actual workflow runs remain on GitHub unless moved).

**Manual bridge**: Projects (GitHub Projects v2) — exported as JSON; no direct equivalent.

#### Step-by-step

1. **Install the HelixGitpx GitHub App** on your org (preferred) or create a fine-grained PAT with `repo`, `admin:org`, `workflow` scopes.
2. In HelixGitpx: **Org → Upstreams → Add upstream → GitHub** → paste credentials → choose **shadow mode** (strongly recommended for first 7 days).
3. **Bulk import**: Select repos → "Import". HelixGitpx does a `git clone --mirror`, pulls PRs/issues/releases/webhooks via API, and creates per-repo bindings.
4. Observe for 7 days:
   - Check `Sync History` on each repo for anomalies.
   - Review a handful of PRs/issues for faithful mirroring.
5. **Enable bidirectional**: Toggle each upstream from shadow → enabled. A push to GitHub now fans out everywhere; a push to HelixGitpx fans back to GitHub.
6. (Optional) Add other providers as additional upstreams; HelixGitpx federates them too.

#### Preserving PR/issue numbers

HelixGitpx preserves numbers when importing. During parallel operation, new numbers are coordinated across upstreams to prevent collisions (see [09-conflict-resolution.md §7.4]).

#### Actions / Checks

HelixGitpx receives check statuses as webhooks and aggregates them. Branch-protection rules evaluated centrally — one amber check on any upstream fails the protection globally.

### 3.2 From GitLab (SaaS or Self-Hosted)

**1:1**: repos, MRs (mapped to PRs), issues, labels, milestones, releases, protected branches, webhooks, LFS.

**Near 1:1**: CI/CD pipelines → checks; GitLab Pages stays on GitLab until manual migration.

**Manual bridge**: Snippets (exported), Boards (labels map), Wikis (Git-backed — import as separate repo).

**Special scopes**: Your PAT needs `api`, `read_repository`, `write_repository`.

Otherwise same flow as GitHub.

### 3.3 From Bitbucket (Cloud or DC)

**1:1**: repos, branches, tags, PRs, webhooks, LFS.

**Near 1:1**: Pipelines (webhook-level).

**Manual bridge**: Snippets (no direct equivalent; exported as gists/files).

**Quirks**:
- Bitbucket DC needs admin-generated "application password" for PAT auth.
- Bitbucket numeric PR IDs are preserved.

### 3.4 From Gitea / Forgejo / Codeberg

Near-complete parity. Single adapter handles all three (they share API). Webhooks supported. Import is straightforward.

### 3.5 From Self-Hosted Git (Gogs, cgit, raw SSH)

- Treat as **Generic Git**: HelixGitpx is the only place with PRs/issues.
- Import via SSH URL + key; periodic pull if no webhook.
- PRs/issues/releases created in HelixGitpx are first-class; they cannot federate back to a plain Git server (no API there).

### 3.6 From Azure DevOps

**1:1**: repos, branches, tags, PRs, policies, webhooks.

**Partial**: Work Items → Issues (labels/milestones mapped; some Azure-specific fields go into extended metadata).

**Manual bridge**: Pipelines — status only.

### 3.7 From AWS CodeCommit

Straightforward Git mirror. No PR/issue concept in CodeCommit → HelixGitpx becomes the collaboration layer.

---

## 4. `helixctl migrate` CLI

For bulk operations:

```bash
# Import all repos from a GitHub org that match a filter
helixctl migrate github \
  --target-org acme \
  --from-org acme-public \
  --filter "topic:public AND NOT archived" \
  --shadow \
  --import-history \
  --include issues,prs,releases,labels,webhooks \
  --concurrency 8

# Dry run
helixctl migrate github --target-org acme --from-org acme --dry-run

# Migrate a curated list from a CSV
helixctl migrate gitlab --file repos.csv --shadow

# Status
helixctl migrate status
helixctl migrate report --format csv > migration.csv
```

CSV format:

```
upstream_owner,upstream_repo,target_slug,visibility,import_issues,import_prs
acme,infra,infra,private,true,true
acme,www,website,public,true,false
```

---

## 5. Mapping Nuances

### 5.1 Identities

When a GitHub user commented on an issue, their comment is imported with their GitHub username and an immutable `imported_from` attribute. When the same person later joins HelixGitpx with a matching verified email, their identity is **merged** (audit event emitted) — historical entries update to show the real HelixGitpx user.

### 5.2 Timestamps

All imported timestamps preserve original UTC values. HelixGitpx never rewrites history.

### 5.3 Large Histories

For repos > 10 GB, use the **progressive import** mode: HelixGitpx streams packs, persists checkpoints, and resumes on failure. Default concurrency scales to the slowest of (upstream rate limit, local network, object-store throughput).

### 5.4 LFS Objects

LFS objects are deduplicated across repos using content addresses — migrating a large monorepo doesn't balloon storage.

### 5.5 Webhooks

Existing webhooks on the source remain in place. HelixGitpx additionally registers its own webhook so it receives events. If an existing webhook duplicates functionality you'd rather route through HelixGitpx, you can decommission it post-cutover.

### 5.6 Permissions

Teams/permissions are **not** imported automatically — too many organisational nuances. Instead, the migration wizard presents suggested mappings and lets admins tweak before applying.

---

## 6. Rollback

Federation-first means rollback is trivial: set the HelixGitpx upstream to **disabled**; traffic continues on the source host unchanged. Any in-flight mirror writes complete; no new ones start.

A full **disconnect** removes credentials and stops webhook registrations — no data is deleted on source.

To roll back data in HelixGitpx (remove the imported repo), use **Archive** then **Delete** with a confirmation token.

---

## 7. Post-Migration Hardening

- Rotate all upstream credentials.
- Enable MFA enforcement in HelixGitpx org settings.
- Define branch protection rules at the HelixGitpx level — they are authoritative.
- Configure AI model budgets appropriately.
- Enable audit log export to your SIEM if required.
- Run the **Compatibility Report**: `helixctl doctor org acme` — surfaces mismatches between source and HelixGitpx state.

---

## 8. Migrating **off** HelixGitpx

We believe users should be free to leave. `helixctl export` produces a complete portable archive:

```bash
helixctl export org acme --out ./acme-export
```

Contents:
- Per-repo: Git bundle + all refs.
- Per-repo: `issues.json`, `prs.json`, `releases.json`, `comments.json`, `labels.json`, `milestones.json`, `reviews.json`.
- Per-repo: `protection.json`, `webhooks.json`, `lfs/` directory with raw objects.
- Per-org: `org.json`, `teams.json`, `memberships.json`, `upstreams.json`.
- Audit export (CSV + JSONL) per retention policy.
- Human-readable `README.md` and machine-readable `manifest.json`.

Formats are interoperable with the importers on other hosts (a pre-built GitHub Actions workflow bundles a round-trip). We also publish a spec so tools can write conversion scripts.

Billing records are retained per legal requirement independent of export.

---

## 9. Common Gotchas

- **Rate limits**: upstream APIs are not unlimited. `helixctl migrate` throttles appropriately, but a 50 k-repo import spans days.
- **Merge commits on import**: we preserve history bit-for-bit. Strange-looking old merges remain strange.
- **Huge binaries outside LFS**: consider running `git-filter-repo` to migrate into LFS *before* importing; cheaper than post-hoc cleanup.
- **Webhooks at upstream**: some providers cap webhook count per repo — verify before enabling many integrations.

---

## 10. Success Checklist

- [ ] Shadow mode ran for ≥ 7 days without divergence alerts.
- [ ] Sample PR/issue round-trip tested end-to-end.
- [ ] Team training completed (user guide §1–4).
- [ ] CI/CD pointing at HelixGitpx as status target (if using native CI).
- [ ] IDE / editor integrations updated to HelixGitpx endpoints.
- [ ] Internal wiki / bookmarks updated.
- [ ] Credentials rotated.
- [ ] On-call rotation aware of new system.

---

*— End of Migration Guide —*
