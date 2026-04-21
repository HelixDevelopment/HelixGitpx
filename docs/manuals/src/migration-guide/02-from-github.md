## 2. Migrating from GitHub.com (SaaS)

Moving a team from GitHub.com to HelixGitpx is a four-phase process. You
never have to finish; each phase is safe to pause.

### 2.1 Phase 1 — bind as an additional upstream

```bash
helixgitpx org create --name acme --residency EU
helixgitpx auth login    # OIDC flow
helixgitpx repo import --from-github acme/my-repo --token $GH_TOKEN
```

The `repo import` command:

1. Clones the repo via the GitHub API.
2. Reads PRs, issues, labels via REST + GraphQL.
3. Mirrors the full history into HelixGitpx storage.
4. Adds GitHub as a *write* upstream binding.

From this moment, developer workflows work against HelixGitpx OR
GitHub — both are active.

### 2.2 Phase 2 — mirror-verify (1–2 weeks)

Enable dual-push in your CI by adding a second remote:

```bash
git remote set-url --add --push origin git@github.com:acme/my-repo.git
git remote set-url --add --push origin git@hx.acme.io:acme/my-repo.git
```

Verify parity daily with:

```bash
helixgitpx repo compare acme/my-repo --against github
```

Any divergence surfaces in the Conflicts inbox.

### 2.3 Phase 3 — flip primary

Update team docs to use the HelixGitpx clone URL. Add HelixGitpx PATs
to developer `~/.config/helixgitpx/tokens`. Rotate CI tokens.

GitHub remains bound as a mirror; pushes continue to fan out there.

### 2.4 Phase 4 — graduate or keep as insurance

Decide whether to keep the GitHub mirror:

- **Keep** — as an offsite replica. HelixGitpx handles the sync; no
  maintenance burden.
- **Remove** — `helixgitpx repo unbind --repo acme/my-repo --upstream github`.
  The repo stops mirroring to GitHub; all future state lives in
  HelixGitpx only.

### 2.5 What travels across

| GitHub concept | HelixGitpx concept | Fidelity |
|---|---|---|
| Commits, branches, tags | Same | 1:1 |
| Pull requests | Pull requests | High (title, body, comments, reviewers) |
| Issues | Issues | High |
| Labels | Labels | 1:1; org-level label catalog |
| Milestones | Milestones | 1:1 |
| Projects (v2) | Projects | Partial; board view re-imported |
| Actions workflows | Not mirrored | Run your own CI against HelixGitpx |
| Marketplace apps | Not mirrored | Use HelixGitpx webhooks + your own automation |

### 2.6 Known caveats

- **GitHub Apps** don't carry across. If you use an App that requires
  installation events, re-install against HelixGitpx's webhook URL.
- **Branch protections** from GitHub are exported as OPA Rego rules.
  Review them before enforcing — GitHub's model is stricter on
  "required status checks" than ours.

---
