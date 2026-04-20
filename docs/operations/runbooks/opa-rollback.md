# Runbook — OPA bundle rollback

**Alert:** `OPAEvaluationErrorRateHigh` or mass `403` responses correlated with bundle revision change.

## Immediate checks

- Active bundle revision: `curl opa-bundle-server/version`.
- Per-service OPA sidecar logs: `kubectl logs -l app=<svc> -c opa --tail=50`.

## Remediation

1. Identify last-known-good revision from CI log (tag-based).
2. Flip pointer: `kubectl -n policy patch configmap opa-bundle-pointer -p '{"data":{"active":"v2.3.4"}}'`.
3. Force OPA agents to re-poll: `kubectl rollout restart deploy -l opa.bundle-reload=true` or wait `bundle.polling_min` seconds.

## Post-incident

- File policy regression test capturing the bad scenario.
- Add diff-review gate to CI if the bad bundle made it past review.

## Escalation

Rollback fails / pointer not propagating → page security + platform lead.
