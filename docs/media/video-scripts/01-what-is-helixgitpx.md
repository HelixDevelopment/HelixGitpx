# Script — 01 What is HelixGitpx?

**Track:** Getting started
**Target length:** 3 min
**Audience:** developers who currently use a single Git host and are
curious about federation.
**Prerequisites:** familiarity with `git push`, `git pull`, and what
GitHub/GitLab are.
**Goal:** understand the federation idea in 3 minutes and want to try it.

## Cold open (0:00 – 0:10)

Split-screen: GitHub outage page on the left, GitLab outage page on the
right, a dev throwing up their hands. Cut to black. Title card:
**"HelixGitpx — one namespace, many hosts."**

## Introduction (0:10 – 0:35)

Webcam + screen. "Today's Git workflows depend on one provider. That's
fragile. HelixGitpx fixes it. I'm <Name>. Three minutes."

## Body

1. **The problem** — 0:35 – 0:55.
   Animated diagram: 3 devs, each locked into a different host.
   Narration: "Every host is a silo. When GitHub goes down, you stop."

2. **The HelixGitpx model** — 0:55 – 1:40.
   Animated diagram: a central HelixGitpx namespace, arrows out to 12 hosts.
   Narration: "You push to HelixGitpx. HelixGitpx mirrors to every host
   you've bound. If one is down, the others still have your code."

3. **Demo: same repo, 4 hosts** — 1:40 – 2:30.
   Screencast: `git push helixgitpx main` — then browse each of GitHub,
   GitLab, GitFlic, GitVerse to show the commit appeared on all four.

4. **Conflict resolution** — 2:30 – 2:50.
   Brief cut to the conflict inbox in the web app. "When two hosts
   diverge, you see it here. An AI proposal you approve or reject."

## Wrap-up (2:50 – 3:00)

Three bullets on screen:
- Federation, not lock-in.
- Real push, real mirror, real verification.
- Start free at helixgitpx.io.

End card: logo + URL + "Next lesson: Install the desktop app."

## Shot list

- 0:00 – 0:10: animated intro.
- 0:10 – 0:35: webcam.
- 0:35 – 1:40: animated explainer.
- 1:40 – 2:30: screencast (terminal + 4 browser tabs).
- 2:30 – 2:50: web app recording.
- 2:50 – 3:00: outro card.

## Assets

- Slides: `assets/01-slides.key`
- Diagrams: `assets/01-silos.svg`, `assets/01-hub.svg`
- B-roll: none.

## Companion documentation

- `docs/manuals/src/user-guide/00-introduction.md`
- `docs.helixgitpx.io/intro`
