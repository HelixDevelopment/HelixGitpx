## 2. My push is slow or hanging

Triage in three minutes.

### 2.1 Is it you or the system?

```bash
curl -fsS https://api.helixgitpx.io/healthz
# {"status":"ok"}
```

If that fails, the service itself is down — see
[operator-guide chapter 5: incident response](../operator-guide/00-introduction.md)
or email `support@helixgitpx.io`.

If `healthz` returns 200, the issue is either network-local or
operation-specific.

### 2.2 Network-local?

```bash
curl -fsS -w "\nTTFB: %{time_starttransfer}s\n" \
    https://api.helixgitpx.io/api/v1/orgs \
    -H "Authorization: Bearer $HGX_PAT" >/dev/null
```

TTFB > 1 s from your location usually means your ISP is routing
strangely. Try a different network or a VPN as a diagnostic.

### 2.3 Large-push-specific

Smart-HTTP pack uploads can take minutes for a fresh clone of a large
repo. Expected throughput:

- First push of a multi-GB repo: 5-50 MB/s depending on your uplink.
- Subsequent pushes (delta-only): should complete in seconds.

Check current transfer size:

```bash
git push --progress
```

If progress stalls for >30 s with no output, see §2.4.

### 2.4 Partial push stuck in fan-out

If your local push returns success but commits aren't visible on one
upstream, the fan-out workflow is in progress or in DLQ:

```bash
helixgitpx sync list --repo $org/$repo --status running
helixgitpx sync list --repo $org/$repo --status dlq
```

To retry a DLQ job:

```bash
helixgitpx sync retry --job $id
```

### 2.5 Authentication-specific

```
error: unable to access 'https://git.helixgitpx.io/…': HTTP/2 stream 1 was not closed cleanly: PROTOCOL_ERROR (err 1)
```

This usually means the PAT you're using has been revoked. Re-issue:

```bash
helixgitpx auth pat create --name workstation --scopes repo:write
```

### 2.6 When to escalate

If §2.1 + §2.2 are both clean but your push keeps timing out after
60 s, file a support ticket with:

- Output of `git push --progress` (redact any URLs containing tokens).
- `helixgitpx --version`.
- Your org slug.
- The time (UTC) of the failing attempt (for log correlation).

---
