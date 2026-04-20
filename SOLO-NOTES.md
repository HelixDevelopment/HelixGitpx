# Solo-Operation Notes

This file documents deviations from `CONTRIBUTING.md` that apply while HelixGitpx is maintained by a single engineer. Each deviation is temporary; re-enable the rule once the constraint no longer applies.

## Active deviations

| Rule | Deviation | Re-enable when |
|---|---|---|
| 2 approvers required for PRs to `main` | 1 self-review | team size ≥ 2 |
| CODEOWNERS-enforced reviews | configured but not enforced in GitHub branch protection | team size ≥ 2 |
| DCO `Signed-off-by` | enforced | always |
| Conventional Commits | enforced | always |
| Signed commits (GPG/SSH) | enforced | always |

## How this file is used

- New engineers read this first after `README.md`.
- Every deviation lists the exact condition under which it is lifted.
- When a deviation is lifted, this file is updated in the same PR that enables the enforcement.
