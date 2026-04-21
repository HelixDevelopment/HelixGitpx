## 2. Members, roles, and teams

An HelixGitpx org has three role tiers. Every user has exactly one org
role. Repository-level overrides are additive.

### 2.1 Org roles

| Role | Can do |
|------|--------|
| Owner | Everything, including delete-org and residency change. |
| Admin | Everything except delete-org and billing. |
| Member | Read + write per repo ACL; cannot invite or change settings. |

There is always ≥ 1 Owner. The system refuses to demote the last Owner.

### 2.2 Inviting members

Settings → Members → **Invite**. Supply an email address and a role.
An invite email with an OIDC-acceptance link is sent. Invites expire
after 7 days.

**Bulk invite via CLI:**

```bash
helixgitpx org members invite \
  --org acme \
  --csv team.csv \
  --default-role member
```

CSV format: `email,role,teams` (teams comma-separated).

### 2.3 Teams

Teams group members and apply shared repo permissions. A user can be
on any number of teams. Team access is *additive* — the most permissive
access wins.

```bash
helixgitpx team create --org acme --name backend
helixgitpx team add-member --team backend --user alice@acme.com
helixgitpx team grant --team backend --repo acme/api --permission admin
```

### 2.4 SSO provisioning (Scale+)

When SSO is configured (Keycloak realm federation), inviting becomes
optional — users can sign in directly and get provisioned as a Member
of the org whose email domain matches. Owners/Admins still require
explicit promotion.

### 2.5 De-provisioning

Removing a user removes their org membership and revokes every PAT
they issued. Their **commits remain** (attribution is preserved). Their
**comments remain** (under a "deactivated account" badge).

Full data export is available on request via
[`privacy@helixgitpx.io`](mailto:privacy@helixgitpx.io).

---
