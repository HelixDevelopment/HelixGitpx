# 28 — Accessibility Standards (WCAG 2.2 AA)

> **Document purpose**: Define the **accessibility commitments** for HelixGitpx across web, mobile, desktop, and documentation; how we design for inclusion; how we test; and how users can report issues.

---

## 1. Commitment

HelixGitpx targets **WCAG 2.2 Level AA** conformance as a minimum. Level AAA where feasible. This applies to every client surface (web, iOS, Android, Windows, macOS, Linux, CLI where it applies) and documentation.

Equivalent regional standards:
- **EN 301 549** (EU).
- **Section 508** (US Federal).
- **ADA Title III** practical application (US).
- **AODA** (Ontario, Canada).

---

## 2. Principles (WCAG 2.2 POUR)

| Principle | What it means for us |
|---|---|
| **Perceivable** | Text alternatives, captions, contrast ≥ 4.5:1, scalable text |
| **Operable** | Full keyboard access, no time-outs without extension, no flashing |
| **Understandable** | Consistent navigation, predictable behaviour, clear error messages |
| **Robust** | Valid HTML/native components, compatible with assistive tech |

---

## 3. Design Foundations

### 3.1 Colour & Contrast

- Text contrast ≥ **4.5:1** (AA) for body, ≥ **3:1** for large text.
- Non-text contrast ≥ **3:1** for UI elements and graphical indicators.
- Colour never used as the only means of conveying information — always paired with icon or text.
- **High-contrast theme** (≥ 7:1) available.
- Dark / light / system modes.
- Tested with common colour-vision deficiencies (deuter-, protan-, tritan-opia) via tooling.

### 3.2 Typography

- Default body ≥ 14 px on web; scales with browser / OS settings.
- Minimum font-size respected (no clamping that fights user OS setting).
- Line-height ≥ 1.5× body on prose; paragraphs 2× body apart.
- No text as image (except logos).

### 3.3 Spacing & Targets

- Touch targets ≥ 44 × 44 pt (iOS HIG) / 48 × 48 dp (Material).
- Hover / focus visual targets ≥ 24 × 24 CSS px (WCAG 2.2 SC 2.5.8).

### 3.4 Motion

- All non-essential animations respect `prefers-reduced-motion`.
- No auto-playing video / audio.
- No content that flashes > 3×/s (seizure safety).

---

## 4. Keyboard & Focus

- **Every interactive element** reachable and operable by keyboard.
- Logical tab order matching visual order.
- **Visible focus indicator** with ≥ 3:1 contrast; never removed.
- **Skip links** (Skip to content, Skip to navigation) on every page.
- **Shortcuts** discoverable via `?` — can be disabled to prevent collisions (WCAG 2.1.4).
- Traps avoided; modals return focus on close.
- **Roving tabindex** for composite widgets (tabs, menus).

---

## 5. Screen Readers

- **Semantic HTML first**; ARIA only to patch gaps.
- Every image has an `alt` (empty for decorative).
- Every form control has an associated `<label>`.
- Live regions for async updates (toast, build status, new comment).
- Landmarks: `banner`, `nav`, `main`, `contentinfo`, `complementary`, `search`.
- Headings in order — no skipped levels.
- Links descriptive — avoid "click here".

Tested with:

- **VoiceOver** (macOS, iOS).
- **NVDA** (Windows).
- **JAWS** (Windows).
- **TalkBack** (Android).
- **Orca** (Linux).

---

## 6. Mobile-Specific

- Respect OS text-size settings up to 200 %.
- Support grayscale and high-contrast modes.
- VoiceOver rotor / TalkBack reading controls available.
- Gestures have non-gesture alternatives (swipe-to-delete always has a button).
- Haptic feedback for critical actions; disable-able.

---

## 7. Desktop-Specific

- Full keyboard navigation on all windows.
- Respect OS-level accessibility: Windows high contrast / Narrator, macOS VoiceOver, Linux ATs.
- Menu bar (native) reachable by keyboard.
- Multiple-monitor-safe focus behaviour.

---

## 8. Forms & Errors

- Inline validation with clear error messages (text + icon + colour).
- Errors never rely on colour alone.
- `aria-live="polite"` or `assertive` depending on urgency.
- Required fields marked textually (not only `*`).
- Auto-complete hints (`autocomplete="email"`, `autocomplete="one-time-code"`, etc.).

---

## 9. Rich UI Patterns

### 9.1 Code Viewer

- Monospaced font; no reliance on whitespace characters for meaning without announced alternatives.
- Copy button with visible keyboard focus.
- Line numbers navigable via JAWS/NVDA table mode.

### 9.2 Diff Viewer

- `role="region"` with descriptive label.
- ARIA describedby states "added" / "removed" in summary.
- Side-by-side and inline modes; user preference remembered.

### 9.3 Conflict Resolver

- Panels announced as regions.
- Actions (Accept / Edit / Reject) have keyboard shortcuts.
- Three-way merge indicators semantic, not only colour-coded.

### 9.4 Interactive Charts

- Accessible data table companion view for every chart.
- Keyboard navigation of data points.

---

## 10. Docs & Marketing Site

- Same WCAG 2.2 AA bar as the app.
- Alt text curated on every image.
- Code blocks have language attributes for screen-reader pronunciation hints.
- Tables have captions + `scope` attributes.
- PDFs (if any) are tagged PDFs; runnable with screen readers.

---

## 11. Testing

### 11.1 Automated

- **axe-core** (Playwright) on every PR — CI blocks on critical/serious issues.
- **jest-axe** in component unit tests.
- **Lighthouse** accessibility score ≥ 95 target.
- **pa11y-ci** on the docs site.
- Mobile: **accessibility-test-framework** (Android) / **XCTest Accessibility** (iOS).

### 11.2 Manual

- Weekly rotating component reviews with screen-reader passes.
- Per-release manual regression suite (20-minute run-through).
- Before each major release: full WCAG 2.2 AA audit by a qualified tester.

### 11.3 User Testing

- Quarterly sessions with users of assistive technologies (paid, with consent).
- Feedback logged in a dedicated queue with SLA.

---

## 12. Third-Party Components

- Prefer components with documented accessibility (Radix UI, Angular CDK, Material, shadcn).
- Wrap / audit any component whose accessibility is unproven.
- Never adopt a dependency whose AT-story is unclear.

---

## 13. Content Guidelines

- Plain language; short sentences; common words.
- Readability target: Grade 8 for product copy, Grade 10 for docs.
- Acronyms expanded on first use.
- Instructions not relying on sensory cues ("the red button" → "the Delete button").

---

## 14. Issue Intake & Response

- **Help → Accessibility Feedback** (web, mobile, desktop) and `accessibility@helixgitpx.example.com`.
- Triage within 2 business days.
- Severity SLA:
  - **Critical** (blocks core task for a user with disability): fix within 7 d.
  - **High**: 30 d.
  - **Medium**: 90 d.
  - **Low**: next release.
- Public VPAT (Voluntary Product Accessibility Template) published and updated per release.

---

## 15. Internationalisation Alignment

- RTL (Arabic, Hebrew) mirrors layout appropriately.
- Language attribute set on root; specific content tagged when different from page language.
- Screen-reader pronounceability tested per locale.

---

## 16. Training

- Onboarding includes accessibility module for every engineer and designer.
- Checklists used during design review + code review.
- Yearly refresher.

---

## 17. Roadmap

- **Short term**: achieve WCAG 2.2 AA on every surface; publish VPAT.
- **Medium term**: AAA for code viewer, diff viewer, conflict resolver.
- **Long term**: research cognitive-accessibility improvements beyond WCAG (clearer flows, AI-powered assist).

---

*— End of Accessibility Standards —*
