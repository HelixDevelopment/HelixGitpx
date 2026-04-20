package helixgitpx.enforcement

# M7 bundle v2 — expanded from v1 (authz) with deny-list for sensitive actions.

default allow := true

# Deny force-push to protected branches
deny contains reason if {
    input.action.op == "git.push"
    input.ref.protected == true
    input.action.force == true
    reason := "force push to protected branch forbidden"
}

# Deny PATs with admin:* scope unless issuer is a realm-admin
deny contains reason if {
    startswith(input.action.op, "pat.issue")
    some scope in input.action.scopes
    startswith(scope, "admin:")
    not input.user.realm_roles[_] == "helixgitpx-admin"
    reason := "admin-scoped PATs require helixgitpx-admin"
}

allow := count(deny) == 0
