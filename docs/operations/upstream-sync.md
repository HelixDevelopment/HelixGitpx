# Runbook — Upstream federation sync

This project practices the federation pattern it specifies. Every change on
`main` is mirrored to **all** configured upstreams.
([Constitution Article IV §2](../../CONSTITUTION.md#4-article-iv--versioning-and-distribution).)

## Configured upstreams

Inventory of target hosts — each is a `Upstreams/<name>.sh` that exports
`UPSTREAMABLE_REPOSITORY` to a Git URL.

| Target   | URL                                                   |
|----------|-------------------------------------------------------|
| GitHub   | `git@github.com:HelixDevelopment/HelixGitpx.git`       |
| GitLab   | `git@gitlab.com:helixdevelopment1/helixgitpx.git`      |
| GitFlic  | `git@gitflic.ru:helixdevelopment/helixgitpx.git`       |
| GitVerse | `git@gitverse.ru:helixdevelopment/HelixGitpx.git`      |

Adding a new target: drop a new `Upstreams/<Target>.sh` script that exports
`UPSTREAMABLE_REPOSITORY`. No other config changes needed — the push and
status scripts auto-discover.

## Cadence

- **Every merge to `main`** — the scheduled workflow fires within the hour
  (or manually via `make upstream-push`).
- **Every tag** — the same workflow pushes `--tags` alongside the branch.
- **Daily anchor** — a manual-only `workflow_dispatch` CI job at
  `.github/workflows/upstream-sync.yml` is expected to be triggered at
  least once per day per Constitution §IV §3.

## Manual sync (local)

```bash
make upstream-push       # push main + tags to all upstreams
make upstream-status     # show ahead/behind vs each upstream
```

Under the hood:

- `scripts/push-to-all-upstreams.sh` — adds each remote transiently, pushes
  `main` and tags, removes the remote alias.
- `scripts/upstream-status.sh` — fetches each upstream and reports the
  `ahead N` / `behind N` vs local HEAD.

## CI (manual-trigger)

Per the `workflow_dispatch`-only mandate (Constitution §CI), there is no
schedule trigger. Operators trigger the sync job manually. A typical day:

1. `gh workflow run upstream-sync.yml` (or click in the Actions UI).
2. CI runs `scripts/push-to-all-upstreams.sh` using a deploy key with push
   rights to each upstream (stored as a GitHub Actions secret per target).
3. Job summary reports each target's push result.

## SSH + credentials

Each upstream needs a working SSH identity (or HTTPS token) on the machine
that runs the push. Suggested:

- Local: dedicated SSH key per target in `~/.ssh/helixgitpx-<target>`,
  wired via `~/.ssh/config` `IdentityFile` directives.
- CI: per-target deploy key as an encrypted secret
  (`UPSTREAM_<TARGET>_SSH_KEY`), installed into `~/.ssh/` at job start.

## Failures

- **Rejected non-fast-forward:** upstream has diverged. Resolve on `main`,
  do not force-push without explicit approval. A diverged upstream is a
  governance incident — escalate in `#helixgitpx-sre`.
- **SSH auth failed:** verify the deploy key is installed and authorized on
  the target platform. Rotate if compromised.
- **Network unreachable:** retry after the incident; the sync is idempotent.

## Observability

- `make upstream-status` prints per-target divergence on demand.
- The CI job uploads a JSON summary to `tools/upstream/reports/`.
- Prometheus scrape exists at `helixgitpx_upstream_last_success_timestamp`
  (gauge per target) for alerting on stale syncs.
