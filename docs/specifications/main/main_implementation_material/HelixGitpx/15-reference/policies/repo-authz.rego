# policies/repo/authz.rego
# Repository-scoped authorization policy for HelixGitpx.
# Loaded by every service via the in-process OPA library; also enforced
# at the gateway for defense in depth. Bundle served by policy-service.
#
# Decision shape:
#   { "allow": bool, "reason": string, "obligations": [ ... ] }
#
# Input shape:
#   { "principal": { user_id, org_ids, scopes, role, mfa, device_id, ... },
#     "action":    "repo.read" | "repo.write" | "repo.admin" | "ref.push" | ... ,
#     "resource":  { kind, id, org_id, visibility, archived, primary_upstream, ... },
#     "context":   { ip, country, method, path, trace_id, ... } }

package helixgitpx.repo.authz

import future.keywords.if
import future.keywords.in
import future.keywords.every

default decision := {
  "allow": false,
  "reason": "default_deny",
  "obligations": [],
}

# ============================================================
# Decision composition — explicit deny wins.
# ============================================================
decision := d {
  some deny in deny_set
  d := {"allow": false, "reason": deny.reason, "obligations": []}
}

decision := d {
  count(deny_set) == 0
  allow_with_reason := find_allow
  d := allow_with_reason
}

# ============================================================
# DENY rules — absolute.
# ============================================================
deny_set contains {"reason": "account_suspended"} if {
  input.principal.status == "suspended"
}

deny_set contains {"reason": "org_not_member"} if {
  input.resource.kind == "repo"
  not input.resource.org_id in input.principal.org_ids
  input.resource.visibility != "public"
  input.principal.role != "platform_admin"
}

deny_set contains {"reason": "repo_archived_write"} if {
  input.action in write_actions
  input.resource.archived == true
}

deny_set contains {"reason": "mfa_required"} if {
  input.action in privileged_actions
  not input.principal.mfa
}

deny_set contains {"reason": "residency_violation"} if {
  some region in input.resource.allowed_regions
  not input.context.region in input.resource.allowed_regions
}

deny_set contains {"reason": "legal_hold"} if {
  input.action in write_actions
  input.resource.legal_hold == true
}

deny_set contains {"reason": "signed_push_required"} if {
  input.action == "ref.push"
  input.resource.require_signed_push == true
  not input.context.signature_valid
}

deny_set contains {"reason": "branch_protection_blocks_force"} if {
  input.action == "ref.push"
  input.context.force_push == true
  input.resource.block_force_push == true
}

deny_set contains {"reason": "ai_cloud_routing_not_allowed"} if {
  input.action == "ai.inference.cloud"
  input.resource.allow_cloud_ai != true
}

# ============================================================
# ALLOW rules (checked only if no denies fired).
# ============================================================
find_allow := {"allow": true, "reason": "public_read", "obligations": []} if {
  input.resource.visibility == "public"
  input.action == "repo.read"
}

find_allow := {"allow": true, "reason": "member_read", "obligations": []} if {
  input.action in read_actions
  input.resource.org_id in input.principal.org_ids
}

find_allow := {"allow": true, "reason": "role_write", "obligations": obligations} if {
  input.action in write_actions
  input.resource.org_id in input.principal.org_ids
  input.principal.role in write_roles
  obligations := write_obligations
}

find_allow := {"allow": true, "reason": "role_admin", "obligations": obligations} if {
  input.action in admin_actions
  input.resource.org_id in input.principal.org_ids
  input.principal.role in admin_roles
  obligations := admin_obligations
}

# Fallback — explicit deny with a reason so clients can surface it.
find_allow := {"allow": false, "reason": "not_authorized", "obligations": []}

# ============================================================
# Data tables
# ============================================================
read_actions := {"repo.read", "ref.read", "pr.read", "issue.read", "release.read"}
write_actions := {"repo.write", "ref.push", "pr.write", "issue.write", "release.write"}
admin_actions := {"repo.admin", "repo.delete", "repo.transfer", "upstream.admin", "policy.write", "ai.inference.cloud"}
privileged_actions := {"repo.admin", "repo.delete", "repo.transfer", "policy.write", "org.admin", "billing.admin"}

write_roles := {"maintainer", "admin", "owner"}
admin_roles := {"admin", "owner"}

# Obligations are hints the service MUST honour (audit, masking, etc).
write_obligations := [
  {"kind": "audit", "severity": "normal"},
]

admin_obligations := [
  {"kind": "audit", "severity": "high"},
  {"kind": "notify_owner_on_delete", "applies_to": ["repo.delete", "repo.transfer"]},
]
