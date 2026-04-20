terraform {
  required_version = ">= 1.10"
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 4.5"
    }
  }
}

provider "vault" {}

resource "vault_jwt_auth_backend" "github_actions" {
  description        = "GitHub Actions OIDC"
  path               = "github-actions"
  type               = "jwt"
  oidc_discovery_url = "https://token.actions.githubusercontent.com"
  bound_issuer       = "https://token.actions.githubusercontent.com"
}

resource "vault_jwt_auth_backend_role" "gha_deploy" {
  backend           = vault_jwt_auth_backend.github_actions.path
  role_name         = "gha-deploy"
  token_policies    = ["deploy"]
  bound_audiences   = ["https://github.com/helixgitpx"]
  bound_claims_type = "string"
  bound_claims = {
    repository = "helixgitpx/*"
    workflow   = "deploy"
  }
  user_claim    = "actor"
  token_ttl     = 900
  token_max_ttl = 1800
  role_type     = "jwt"
}

resource "vault_policy" "deploy" {
  name   = "deploy"
  policy = file("${path.module}/../policies/deploy.hcl")
}

resource "vault_policy" "ci" {
  name   = "ci"
  policy = file("${path.module}/../policies/ci.hcl")
}
