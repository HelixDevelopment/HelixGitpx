package helixgitpx.authz

default allow := false

# Owners do anything in their scope.
allow if {
    input.user.effective_role == "owner"
}

# Admins manage members in their team or descendants.
allow if {
    startswith(input.action.op, "team.member.")
    input.user.effective_role == "admin"
}

# Viewers can read within their team/ancestors.
allow if {
    startswith(input.action.op, "read.")
    input.user.effective_role != ""
}
