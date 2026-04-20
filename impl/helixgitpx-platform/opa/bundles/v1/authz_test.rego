package helixgitpx.authz_test

import data.helixgitpx.authz

test_owner_can_delete_org if {
    authz.allow with input as {
        "user": {"effective_role": "owner"},
        "action": {"op": "org.delete"},
    }
}

test_viewer_cannot_delete_org if {
    not authz.allow with input as {
        "user": {"effective_role": "viewer"},
        "action": {"op": "org.delete"},
    }
}

test_admin_can_add_member if {
    authz.allow with input as {
        "user": {"effective_role": "admin"},
        "action": {"op": "team.member.add"},
    }
}

test_member_can_read if {
    authz.allow with input as {
        "user": {"effective_role": "member"},
        "action": {"op": "read.org"},
    }
}

test_empty_role_cannot_read if {
    not authz.allow with input as {
        "user": {"effective_role": ""},
        "action": {"op": "read.org"},
    }
}
