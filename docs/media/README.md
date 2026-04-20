# HelixGitpx Video Curriculum

Per [Constitution Article III §3](../../CONSTITUTION.md#3-article-iii--documentation),
every major documentation section has a parallel video lesson.

This directory holds **scripts** and **brand assets**. Recorded video
binaries live outside git (too large) and are uploaded to the video host
defined in `brand/hosting.md`.

## Course curriculum (initial cut)

### Track 1 — Getting started (users)

1. **What is HelixGitpx?** — 3 min, animated explainer.
2. **Install the desktop app** — 2 min, live screencast.
3. **Your first repository binding** — 5 min, live screencast.
4. **Pushing, pulling, PRs** — 6 min, screencast.
5. **Resolving a conflict (with AI)** — 5 min, screencast.
6. **Using the mobile app** — 4 min, device capture.

### Track 2 — For developers

7. **Service architecture tour** — 10 min, diagrams + code walk.
8. **Writing a new adapter-pool provider** — 15 min, live code.
9. **Extending the OPA bundle** — 8 min, Rego walk.
10. **Webhook HMAC verification** — 6 min, Go code walk.

### Track 3 — For operators

11. **Day-0 cluster bring-up (GitOps)** — 15 min, screencast.
12. **Multi-region failover drill** — 18 min, controlled failure demo.
13. **Chaos engineering walkthrough** — 12 min, Litmus demo.
14. **Perf budget enforcement** — 8 min, k6 + CI walk.
15. **Incident response in production** — 14 min, scenario replay.

### Track 4 — For security and compliance

16. **OPA policy-as-code end-to-end** — 12 min, rego + diff review.
17. **SOC 2 evidence collection** — 8 min, tooling walk.
18. **Bug bounty triage** — 10 min, workflow.

### Track 5 — For admins

19. **Org management, RBAC, residency** — 8 min.
20. **Billing and plan management** — 6 min.
21. **Migrating from GitHub / GitLab** — 12 min.

## Production pipeline

- **Scripts** — Markdown, one file per lesson, under `video-scripts/`.
  Each file follows the [script template](./video-scripts/_TEMPLATE.md).
- **Brand** — fonts, colours, logo, end-card, lower-thirds in `brand/`.
- **Recording** — OBS Studio scenes configured per track.
- **Editing** — DaVinci Resolve free tier; project files under `editing/`.
- **Publishing** — Vimeo (primary), YouTube (mirror), self-hosted MinIO
  (archive). Publishing workflow under
  [`.github/workflows/video-publish.yml`](../../.github/workflows/video-publish.yml).

## Status

At GA: scripts are the deliverable. Recording + editing is a separate
project budgeted post-GA. See
[`docs/marketing/launch-checklist.md`](../marketing/launch-checklist.md)
for the production schedule.
