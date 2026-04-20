# Vault + OIDC configuration

Activation is M2 (requires a running Vault cluster). Artifacts ship here for
GitOps from day one.

## Activation checklist (M2)

1. `export VAULT_ADDR=https://vault.internal`
2. `terraform init && terraform apply` in `terraform/`.
3. Verify: `vault read auth/github-actions/role/gha-deploy`.
4. Flip `.github/workflows/deploy.yml` from the VAULT_ADDR-gated skip into the real fetch step.

## Safety

Tokens issued by `gha-deploy` are TTL-bound (15 min default). Never widen scope
without an ADR and a security review.
