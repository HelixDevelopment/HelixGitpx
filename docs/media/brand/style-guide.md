# HelixGitpx Video Brand Style Guide

## Identity

- **Name:** HelixGitpx (never "Helix Gitpx" or "Helixgitpx").
- **Tagline:** "One namespace, many hosts."

## Colours

| Role | Hex | Notes |
|------|-----|-------|
| Primary | `#2563eb` | Matches the docs-site `custom.css`. |
| Secondary | `#0f172a` | Dark backgrounds. |
| Accent | `#22d3ee` | Highlights, CTA buttons. |
| Success | `#10b981` | Passing tests, successful pushes. |
| Warning | `#f59e0b` | Conflicts, attention. |
| Error | `#ef4444` | Failures, rate-limits. |
| Text on dark | `#f8fafc` | — |
| Text on light | `#0f172a` | — |

## Typography

- **Headings:** Inter SemiBold.
- **Body:** Inter Regular.
- **Code:** JetBrains Mono.
- Include Latin Extended + Cyrillic subsets — the project has Russian
  upstreams (GitFlic, GitVerse) and bilingual docs are expected.

## Audio

- **Music:** royalty-free ambient from Epidemic Sound or Musicbed. No
  vocals during screencasts. Duck music under narration by −18 dB.
- **Narration:** −14 LUFS integrated, −2 dBTP peak. Noise floor ≤ −60 dB.

## Safe frames

- 16:9, 1920×1080 or 3840×2160.
- Title safe: 5 % margin.
- Caption safe: bottom third, plus alt captions in SRT/VTT.

## End card

Logo centered, URL underneath, "Next: <lesson>" below. 3-second hold.

## Intro

8-second animated logo reveal. Shared file: `brand/intro.mov`.

## Accessibility

- All lessons ship captions (`.vtt`) in English by GA; Russian follows.
- Contrast ≥ 4.5:1 for all on-screen text.
- No flashing lights > 3 Hz.
