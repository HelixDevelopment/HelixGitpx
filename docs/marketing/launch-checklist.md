# Launch Checklist — GA Day

## 7 days before

- [ ] `RELEASE.md` frozen; engineering sign-off.
- [ ] Docs site deployed at `docs.helixgitpx.io` with TLS green.
- [ ] Trust center route live in web app.
- [ ] Status page `status.helixgitpx.io` configured with all services.
- [ ] Final pen-test findings triaged; High/Critical closed.
- [ ] DR drill completed within last 14 days; RTO/RPO green.
- [ ] Perf budgets green on 7-day soak.

## 3 days before

- [ ] Blog post drafted, reviewed by founder + marketing.
- [ ] HackerNews + ProductHunt submissions drafted (queued, not posted).
- [ ] Partner announcements coordinated (shortlist of 5 OSS partners).
- [ ] Demo video recorded and captioned.
- [ ] Customer success team trained on billing flows + known limitations.

## 1 day before

- [ ] Feature flags audited; GA-ready defaults.
- [ ] Argo CD sync-wave order validated.
- [ ] Emergency rollback procedure rehearsed with on-call.
- [ ] All team PTO cleared for 48h post-launch.

## Launch day

- [ ] 08:00 — Final smoke test (5-min end-to-end on prod).
- [ ] 09:00 — Flip public DNS (`helixgitpx.io` A/AAAA to production).
- [ ] 09:30 — Publish blog post.
- [ ] 10:00 — Submit to HackerNews + ProductHunt.
- [ ] 10:00 — Email announcement to beta list.
- [ ] 10:15 — Social (Mastodon, Bluesky, LinkedIn, X) threads.
- [ ] 14:00 — Mid-day metrics review: signups, errors, SLO.
- [ ] 17:00 — End-of-day metrics review; post internal recap.

## Day +1 through Day +7

- [ ] Monitor HN thread; engineering on rotation to answer.
- [ ] Daily metrics review; incident response on standby.
- [ ] Collect user feedback; open issues tagged `ga-feedback`.

## Go / No-go criteria

Any of these blocks launch:
- Open P1/P2.
- Failing SLO against baseline.
- DR drill failure in last 14 days.
- Any Critical pen-test finding unresolved.
