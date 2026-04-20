# Runbook — Certificate near expiry

**Alert:** `CertManagerCertificateExpiringIn7Days` or `CertManagerCertificateExpiryError`.

## Immediate checks

- `kubectl get certificate -A | grep -E 'False|True.*[0-6]d'`.
- cert-manager controller logs: `kubectl -n cert-manager logs -l app=cert-manager --tail=200`.

## Remediation

1. **Let's Encrypt rate limit** — switch to staging ACME issuer, verify, switch back.
2. **DNS-01 solver failing** — check DNS provider credentials secret (Route53/Cloudflare token).
3. **Private CA (SPIFFE/SPIRE)** — verify spire-server is healthy; rotate upstream CA if needed.

## Escalation

Expiry < 24h + no path to renewal → page SRE lead to issue static cert as stop-gap.
