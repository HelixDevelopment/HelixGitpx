## 3. Pushing, pulling, and pull requests

This chapter walks the day-to-day developer workflow. It assumes you've
completed [Chapter 2](./02-first-repo.md) and have a bound repository.

### 3.1 Clone

```bash
git clone https://git.helixgitpx.io/<org>/<repo>.git
cd <repo>
```

The HelixGitpx Git URL is the only URL you ever need. You never type
`github.com` or `gitlab.com` again — bindings handle the fan-out.

### 3.2 Branching conventions

From `CONTRIBUTING.md`:

- `feat/<topic>` — new functionality.
- `fix/<topic>` — bug fix.
- `docs/<topic>` — docs-only changes.
- `chore/<topic>` — toolchain, CI, or housekeeping.

The web app enforces these prefixes by default; the org admin can
override in **Settings → Branching**.

### 3.3 Push a change

```bash
git switch -c feat/readme-update
echo "## Features" >> README.md
git add README.md
git commit -m "docs: expand readme"
git push -u helixgitpx feat/readme-update
```

The ref appears on every bound upstream within seconds. Open a PR from
the HelixGitpx web app, the CLI, or the mobile app.

### 3.4 Opening a pull request

From the web app: click **New PR**, pick target branch, describe the
change. The PR is mirrored to every upstream in both directions.

From the CLI:

```bash
helixgitpx pr create \
  --base main \
  --head feat/readme-update \
  --title "docs: expand readme" \
  --body "Adds a features section." \
  --reviewers @alice,@bob
```

### 3.5 Reviews

Inline comments land on whichever host the reviewer is on — GitHub,
GitLab, or HelixGitpx web — and propagate. If two reviewers race on the
same line from different hosts, the conflict-resolver tags both threads
into the conflict inbox.

### 3.6 Merging

Merging from HelixGitpx triggers one fan-out merge per bound upstream.
Options:

- **Squash** — default for feature branches.
- **Rebase** — default for linear histories (per org setting).
- **Merge commit** — only when the history has review-value structure.

Branch protection rules are enforced once centrally by OPA, not
per-provider.

### 3.7 When an upstream is down

- Push to HelixGitpx always succeeds (we write to our own storage).
- Fan-out queues the upstream; retries with exponential back-off.
- Once the upstream is back, the queue drains. Your commit shows up
  there when it can.

You never wait on a slow upstream. See
[Troubleshooting Chapter 3](../troubleshooting/00-introduction.md) if the
queue takes longer than 15 minutes to drain.

---
