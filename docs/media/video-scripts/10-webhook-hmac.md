# Script — 10 Webhook HMAC verification

**Track:** Developers · **Length:** 6 min · **Goal:** viewer understands why constant-time comparison matters and how we enforce it.

## Body

1. **Threat model** — 0:30 – 1:30.
   An attacker with the URL but not the secret. Why naive bytes.Equal is a timing oracle.
2. **The code** — 1:30 – 3:00.
   Walk `platform/webhook/hmac.go` — `hmac.Equal` + header parsing.
3. **Provider quirks** — 3:00 – 4:30.
   GitHub `sha256=...`, GitLab `X-Gitlab-Token`, Bitbucket UUIDs. All end with a constant-time compare.
4. **Fuzz corpus** — 4:30 – 5:30.
   `tools/fuzz/corpora/webhook/` + `go test -fuzz`.

## Wrap-up
"Two lines of code, one real attacker, your job."

## Companion doc
`impl/helixgitpx/platform/webhook/`
