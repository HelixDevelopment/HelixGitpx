# Script — 03 Your first repository binding

**Track:** Getting started · **Length:** 5 min · **Goal:** viewer binds one repo to two upstreams and pushes.

## Cold open (0:00 – 0:10)
Four browser tabs open (GitHub, GitLab, GitFlic, GitVerse). "One repo. Four hosts. Twenty seconds."

## Body

1. **Create the repo in HelixGitpx** — 0:10 – 0:50.
   Web app → New repository → `myorg/hello`. Private, EU residency.
2. **Add a binding** — 0:50 – 2:10.
   Repo page → Bindings → + GitHub. Paste PAT, verify scope is `repo`.
   Direction: `write`. Show the binding appears green.
3. **Add a second binding** — 2:10 – 2:50.
   + GitLab. Same flow. Direction: `write`.
4. **Clone + push** — 2:50 – 4:10.
   Terminal: `git clone https://git.helixgitpx.io/myorg/hello.git`.
   `echo hi > README.md && git add . && git commit -m first && git push`.
   Show commit arriving on GitHub AND GitLab within seconds.
5. **Verify the audit log** — 4:10 – 4:45.
   Repo → Activity. Two fan-out events, both green.

## Wrap-up (4:45 – 5:00)
Recap: "Bind, push, done."

## Companion doc
`docs/manuals/src/user-guide/00-introduction.md` · `docs.helixgitpx.io/getting-started/first-repo`
