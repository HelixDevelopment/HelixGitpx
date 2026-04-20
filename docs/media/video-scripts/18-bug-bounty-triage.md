# Script — 18 Bug bounty triage

**Track:** Security & compliance · **Length:** 10 min · **Goal:** viewer knows how a report becomes a fix.

## Body

1. **Report intake** — 0:30 – 1:30.
   HackerOne → internal queue. Auto-acknowledge.
2. **Severity assessment** — 1:30 – 3:00.
   CVSSv3.1 + business impact. Severity matrix.
3. **Reproduce** — 3:00 – 4:30.
   Safe test env; if can't reproduce, ask for more info.
4. **Fix + regression test** — 4:30 – 6:30.
   Patch, add test covering the attack, OPA rule if applicable.
5. **Reward decision** — 6:30 – 7:30.
   Follow the published schedule.
6. **Coordinated disclosure** — 7:30 – 9:00.
   90 days; earlier if reporter and we agree.
7. **Public CVE** — 9:00 – 9:45.
   File CVE when applicable; publish advisory.

## Wrap-up (9:45 – 10:00)
"Researchers catch the bugs we didn't think of. Pay them."

## Companion doc
`docs/security/bug-bounty-program.md`
