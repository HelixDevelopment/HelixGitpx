## 2. Your first repository binding

You'll create an organization, make a repository, bind it to two Git hosts,
and push a first commit. End to end this takes about 10 minutes.

### 2.1 Create an organization

Open the web app at your tenant URL and click **New organization**. Give
it a name and pick a data residency zone. EU, UK, and US are the GA
choices; see [Administrator Guide §5](../administrator-guide/00-introduction.md)
for the tradeoffs.

> 💡 Residency is per-org, set at creation time, and can be changed later
> by the org owner via `SetOrgResidency`. Existing data migrates in a
> maintenance window; new writes land in the new zone immediately.

### 2.2 Create a repository

Inside your org, **New repository**. Pick a name. Choose *private* unless
you intend the repository to be public on all of its bound upstreams.

### 2.3 Add upstream bindings

In the repository's **Bindings** tab, click **+ Add upstream** and pick
GitHub. Paste a Personal Access Token with the `repo` scope. Set
*direction* to `write`. Save.

Repeat for a second upstream — GitLab is a sensible default so you get
visible federation on day one.

> ⚠️ Tokens stay inside the HelixGitpx secret store. We **never** log
> them. See [Security Handbook §4](../security-handbook/00-introduction.md).

### 2.4 Clone and push

Grab the HelixGitpx clone URL from the repo page header. It looks like:

```
https://git.helixgitpx.io/<org>/<repo>.git
```

```bash
git clone https://git.helixgitpx.io/<org>/<repo>.git
cd <repo>
echo "# Hello" > README.md
git add README.md
git commit -m "feat: initial commit"
git push
```

### 2.5 Verify the fanout

Refresh GitHub and GitLab in your browser. The commit should be there,
on both, within a few seconds. On the HelixGitpx repo page, the
**Activity** tab shows two `repo.pushed` events — one per upstream —
each with a link to the mirrored commit.

### 2.6 What just happened

1. Your `git push` hit **git-ingress**, the HelixGitpx Git smart-HTTP
   server.
2. git-ingress recorded the push and enqueued a `fanout_push_wf` Temporal
   workflow.
3. The workflow asked **adapter-pool** to push the ref to each bound
   upstream.
4. adapter-pool called the GitHub and GitLab APIs with the tokens you
   stored.
5. **audit-service** recorded every step with a signed, Merkle-verified
   log entry.

If any upstream push fails, it retries with exponential backoff. When an
upstream is completely offline, the push is queued and drained when the
upstream is back. Your local push never waits.

---
