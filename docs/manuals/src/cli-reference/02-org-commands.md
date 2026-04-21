## 2. Org commands

The `helixgitpx org` subtree manages organizations — tenancy's top-level
container.

### 2.1 Create

```bash
helixgitpx org create \
  --name acme \
  --residency EU \
  --owner-email you@acme.com
```

Options:

- `--residency` — `EU` (default), `UK`, or `US`.
- `--plan` — `free` (default), `team`, `scale`, `enterprise`. Paid
  plans require an active Stripe subscription; unpaid plans trial for
  14 days.

### 2.2 List

```bash
helixgitpx org list
helixgitpx org list --output json
```

### 2.3 Show

```bash
helixgitpx org show --org acme
```

Output includes: residency, plan, seat count, owner(s), created-at.

### 2.4 Set residency

```bash
helixgitpx org set-residency --org acme --to UK
```

Requires Owner. Triggers an async migration workflow; the CLI returns
the workflow ID. Monitor with `helixgitpx sync show <id>`.

### 2.5 Delete

```bash
helixgitpx org delete --org acme --confirm "acme"
```

Requires Owner + the `--confirm <name>` safety string matching the org
name. Deletion is soft for 30 days, hard after.

### 2.6 Members

See [§2 of the Administrator Guide](../administrator-guide/02-members-rbac.md).
CLI equivalents:

```bash
helixgitpx org members list --org acme
helixgitpx org members invite --org acme --email bob@acme.com --role member
helixgitpx org members remove --org acme --user bob@acme.com
```

---
