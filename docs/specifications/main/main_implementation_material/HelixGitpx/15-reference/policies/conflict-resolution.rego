# policies/conflict/resolution.rego
# Conflict-resolution strategy selection policy.
# Given a conflict case, emit the preferred strategy, whether AI may be
# invoked, and the minimum confidence needed for auto-apply.
#
# Called by the conflict-resolver service.
#
# Input shape:
#   { "case":     { id, kind, repo_id, subject, severity, ... },
#     "repo":     { id, org_id, primary_upstream, branch_protection_strict, ... },
#     "org":      { tier, risk_profile, require_human_review_on_high },
#     "system":   { model_name, model_version, model_confidence_baseline, ... } }
#
# Output shape:
#   { "strategy":          "prefer_primary" | "prefer_newer" | "three_way_merge"
#                          | "crdt_merge"   | "policy_deny"  | "escalate_to_human",
#     "use_ai":            bool,
#     "auto_apply":        bool,
#     "min_confidence":    number,
#     "reason":            string,
#     "notify_roles":      [ ... ],
#     "sandbox_required":  bool }

package helixgitpx.conflict.resolution

import future.keywords.if
import future.keywords.in

default decision := {
  "strategy": "escalate_to_human",
  "use_ai": false,
  "auto_apply": false,
  "min_confidence": 0.99,
  "reason": "default_escalate",
  "notify_roles": ["maintainer"],
  "sandbox_required": true,
}

# ---- Kind-driven defaults ---------------------------------
decision := d if {
  input.case.kind == "metadata_concurrent"
  d := {
    "strategy": "crdt_merge",
    "use_ai": false,
    "auto_apply": true,
    "min_confidence": 1.0,                 # CRDT is deterministic
    "reason": "metadata_uses_crdt",
    "notify_roles": [],
    "sandbox_required": false,
  }
}

decision := d if {
  input.case.kind == "ref_divergence"
  input.repo.primary_upstream != ""
  not input.repo.branch_protection_strict
  d := {
    "strategy": "prefer_primary",
    "use_ai": false,
    "auto_apply": true,
    "min_confidence": 1.0,
    "reason": "primary_upstream_wins",
    "notify_roles": [],
    "sandbox_required": false,
  }
}

decision := d if {
  input.case.kind == "ref_divergence"
  input.repo.branch_protection_strict
  d := {
    "strategy": "three_way_merge",
    "use_ai": true,
    "auto_apply": auto_ok,
    "min_confidence": ai_threshold,
    "reason": "three_way_with_ai_assist",
    "notify_roles": ["maintainer"],
    "sandbox_required": true,
  }
}

decision := d if {
  input.case.kind == "rename_collision"
  d := {
    "strategy": "three_way_merge",
    "use_ai": true,
    "auto_apply": false,             # always human review for renames
    "min_confidence": 1.0,
    "reason": "renames_always_reviewed",
    "notify_roles": ["maintainer", "author"],
    "sandbox_required": true,
  }
}

decision := d if {
  input.case.kind == "pr_state"
  input.case.severity in {"low", "normal"}
  d := {
    "strategy": "prefer_newer",
    "use_ai": false,
    "auto_apply": true,
    "min_confidence": 1.0,
    "reason": "pr_state_mismatch_prefer_newer",
    "notify_roles": ["author"],
    "sandbox_required": false,
  }
}

decision := d if {
  input.case.kind == "lfs_divergence"
  d := {
    "strategy": "escalate_to_human",
    "use_ai": false,
    "auto_apply": false,
    "min_confidence": 1.0,
    "reason": "lfs_divergence_quarantined",
    "notify_roles": ["maintainer", "org_admin"],
    "sandbox_required": false,
  }
}

decision := d if {
  input.case.kind == "tag_collision"
  input.case.metadata.signed == true
  d := {
    "strategy": "policy_deny",
    "use_ai": false,
    "auto_apply": false,
    "min_confidence": 1.0,
    "reason": "signed_tag_immutable",
    "notify_roles": ["maintainer"],
    "sandbox_required": false,
  }
}

# ---- Org / risk overrides ---------------------------------
decision := d if {
  base := base_decision
  input.case.severity == "critical"
  input.org.require_human_review_on_critical
  d := merge(base, {"auto_apply": false, "strategy": "escalate_to_human"})
}

decision := d if {
  base := base_decision
  input.org.risk_profile == "regulated"
  base.use_ai == true
  d := merge(base, {"auto_apply": false})
}

# ---- Helpers ----------------------------------------------
auto_ok if {
  input.system.model_confidence_baseline >= 0.92
  not input.org.require_human_review_on_high
}

ai_threshold := 0.92 if {
  input.org.risk_profile == "standard"
}
ai_threshold := 0.97 if {
  input.org.risk_profile == "elevated"
}
ai_threshold := 1.0 if {
  input.org.risk_profile == "regulated"
}

base_decision := decision

merge(a, b) := c if {
  c := object.union(a, b)
}
