# 37 — Design System & Tokens

> **Document purpose**: Source of truth for the HelixGitpx visual language: color tokens, typography, spacing, elevations, motion. Shared across web (Angular + Tailwind), mobile/desktop (Compose Multiplatform), and marketing surfaces. Changes here propagate to every client.

---

## 1. Brand Essentials

- **Name**: HelixGitpx (always capitalised as written; no hyphen).
- **Primary mark**: the green chameleon spiral. Never rotate; respect the clear-space (0.5× the logo's short side).
- **Voice**: clear, candid, technically serious but warm. Avoid jargon unless it's what the user searched for.

---

## 2. Colour Tokens

We maintain tokens at three layers, per [Design Tokens Community Group](https://design-tokens.github.io/community-group/format/):

```
primitive → semantic → component
```

### 2.1 Primitives (brand palette)

| Token | Hex | Notes |
|---|---|---|
| `color.primitive.green.50`  | `#F2FAE3` | Background tint |
| `color.primitive.green.100` | `#E2F2C8` | |
| `color.primitive.green.200` | `#CBE8A1` | |
| `color.primitive.green.300` | `#B4DE7A` | |
| `color.primitive.green.400` | `#9ACD32` | **Brand green** (logo) |
| `color.primitive.green.500` | `#85B724` | Hover / pressed |
| `color.primitive.green.600` | `#6F9D17` | Active / dark-mode accent |
| `color.primitive.green.700` | `#588011` | |
| `color.primitive.green.800` | `#3E600B` | |
| `color.primitive.green.900` | `#25400B` | Text-on-light accent |
| `color.primitive.teal.400`  | `#7EC8C2` | Logo secondary |
| `color.primitive.teal.600`  | `#54A39B` | |
| `color.primitive.gray.0`    | `#FFFFFF` | Paper |
| `color.primitive.gray.50`   | `#FAFAFA` | |
| `color.primitive.gray.100`  | `#F3F4F6` | |
| `color.primitive.gray.200`  | `#E5E7EB` | |
| `color.primitive.gray.300`  | `#D1D5DB` | |
| `color.primitive.gray.400`  | `#9CA3AF` | |
| `color.primitive.gray.500`  | `#6B7280` | |
| `color.primitive.gray.600`  | `#4B5563` | |
| `color.primitive.gray.700`  | `#374151` | |
| `color.primitive.gray.800`  | `#1F2937` | |
| `color.primitive.gray.900`  | `#0F172A` | Ink |
| `color.primitive.red.500`   | `#EF4444` | Error |
| `color.primitive.amber.500` | `#F59E0B` | Warning |
| `color.primitive.sky.500`   | `#0EA5E9` | Info |

### 2.2 Semantic (light + dark modes)

| Token | Light | Dark |
|---|---|---|
| `color.bg.page`          | `gray.0`     | `gray.900` |
| `color.bg.surface`       | `gray.50`    | `gray.800` |
| `color.bg.elevated`      | `gray.0`     | `gray.700` |
| `color.bg.hover`         | `gray.100`   | `gray.700` |
| `color.fg.default`       | `gray.900`   | `gray.50`  |
| `color.fg.muted`         | `gray.500`   | `gray.400` |
| `color.border.default`   | `gray.200`   | `gray.700` |
| `color.border.strong`    | `gray.300`   | `gray.600` |
| `color.accent.default`   | `green.600`  | `green.400` |
| `color.accent.hover`     | `green.700`  | `green.300` |
| `color.accent.on`        | `gray.0`     | `gray.900` |
| `color.danger.default`   | `red.500`    | `red.500`  |
| `color.warning.default`  | `amber.500`  | `amber.500`|
| `color.info.default`     | `sky.500`    | `sky.500`  |
| `color.success.default`  | `green.600`  | `green.400` |

### 2.3 Contrast

Every pairing (`fg` × `bg`) respects **WCAG AA**: ≥ 4.5:1 for body text; ≥ 3:1 for large text / non-text. Automated contrast checks run in CI against our defined pairings.

---

## 3. Typography

- **Sans**: `Inter` (self-hosted, variable font). Fallback stack: `-apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`.
- **Mono**: `JetBrains Mono` (self-hosted). Fallback: `ui-monospace, "SF Mono", Menlo, Consolas, monospace`.
- **Emoji**: system native (Apple Color Emoji / Segoe UI Emoji / Noto Color Emoji).

### Scale

| Token | Size | Line | Use |
|---|---|---|---|
| `text.xs`   | 12px | 16px | Labels, meta |
| `text.sm`   | 14px | 20px | Secondary body |
| `text.base` | 16px | 24px | Body (minimum) |
| `text.lg`   | 18px | 28px | Emphasis |
| `text.xl`   | 20px | 28px | Subheads |
| `text.2xl`  | 24px | 32px | H3 |
| `text.3xl`  | 30px | 36px | H2 |
| `text.4xl`  | 36px | 44px | H1 (page) |
| `text.5xl`  | 48px | 52px | Display |
| `text.6xl`  | 60px | 64px | Hero |

Weights: `400`, `500`, `600`, `700`. Letter-spacing tightens at larger sizes.

---

## 4. Spacing & Layout

4px base grid.

| Token | Value |
|---|---|
| `space.0` | 0 |
| `space.1` | 4px |
| `space.2` | 8px |
| `space.3` | 12px |
| `space.4` | 16px |
| `space.5` | 20px |
| `space.6` | 24px |
| `space.8` | 32px |
| `space.10` | 40px |
| `space.12` | 48px |
| `space.16` | 64px |
| `space.20` | 80px |

### Container widths

| Token | Max-width |
|---|---|
| `container.sm` | 640px |
| `container.md` | 768px |
| `container.lg` | 1024px |
| `container.xl` | 1280px |
| `container.2xl` | 1536px |

---

## 5. Radii & Borders

| Token | Value |
|---|---|
| `radius.none` | 0 |
| `radius.sm`   | 4px |
| `radius.md`   | 6px |
| `radius.lg`   | 8px  |
| `radius.xl`   | 12px |
| `radius.2xl`  | 16px |
| `radius.full` | 9999px |

Border width defaults: `1px`; focus rings `2px`; emphasised dividers `2px`.

---

## 6. Elevation (shadows)

| Token | Usage |
|---|---|
| `shadow.sm` | low emphasis cards |
| `shadow.md` | popovers, dropdowns |
| `shadow.lg` | modals |
| `shadow.xl` | coach marks, hero cards |

Shadows mirror-compensate in RTL. Dark mode shadows include a subtle light border to keep edges readable.

---

## 7. Motion

- **Duration**: `fast` 120 ms, `base` 200 ms, `slow` 300 ms.
- **Easing**: `ease-out` for enter; `ease-in` for exit; `ease-in-out` for transforms.
- All non-essential motion respects `prefers-reduced-motion` (replace with opacity fade only).
- Reserved for purposeful motion: state changes, positional transitions, loading feedback. No ambient motion.

---

## 8. Iconography

- Base library: Lucide (Apache-2.0) mirrored into our monorepo.
- Stroke width `1.5`, size units of 4px (16, 20, 24).
- Avoid metaphors ambiguous across cultures.
- Direction-indicating icons flipped in RTL.

---

## 9. Imagery & Illustration

- Illustrations in the HelixGitpx palette only.
- No stock imagery with people in feature screenshots (privacy, longevity).
- Product screenshots kept current (auto-refreshed via Playwright).
- Accessible alt text mandatory.

---

## 10. Cross-Platform Implementation

### Web (Angular + Tailwind)

- Generated `tailwind.config.ts` from the design-token source of truth (`design-tokens/tokens.json`).
- CSS custom properties exposed per theme (`--hgx-color-bg-page`) for runtime theming.
- Storybook instance at `storybook.helixgitpx.example.com` shows every component + state.

### Mobile / Desktop (KMP + Compose Multiplatform)

- Tokens codegen'd into `shared/ui/src/commonMain/.../Tokens.kt` as Compose `Color` / `TextStyle` / `Shape` constants.
- Material 3 `ColorScheme` derived from semantic tokens.
- Dark mode automatically follows system; manual override available.

### Native specifics

- iOS: Dynamic Type honoured.
- Android: Material You optional opt-in.
- Desktop: native menubar styling respects OS theme.

---

## 11. Component Library

Every component has:
- Design spec (Figma).
- Accessible implementation per platform.
- Storybook story for web; Compose preview for mobile/desktop.
- Documented states (default / hover / focus / pressed / disabled / loading / error).
- a11y audit notes (keyboard, screen reader, contrast).
- Unit + visual-regression + axe tests.

Core components: Button, IconButton, Input, Select, Checkbox, Radio, Switch, Tabs, Menu, Dialog, Popover, Toast, Breadcrumb, Card, Avatar, Badge, Alert, ProgressBar, Skeleton, Table, Pagination, DatePicker, CodeBlock, DiffViewer, ConflictResolver.

---

## 12. Token Pipeline

- **Source**: `design-tokens/tokens.json` (Style-Dictionary format).
- **Generators**: `pnpm build-tokens` emits:
  - `web/styles/tokens.css` (CSS custom properties).
  - `web/tailwind.generated.ts`.
  - `clients/shared/ui/src/commonMain/kotlin/.../Tokens.kt`.
  - `ios/HelixGitpxTokens.swift` (via KMP cinterop; optional).
  - `docs/assets/token-samples.html` for reference.

PRs to tokens.json require design + frontend + mobile approval.

---

## 13. Theming & Customisation

- **System theme** (default).
- **Light / Dark** explicit.
- **High contrast** mode (≥ 7:1).
- **Solarized** (opt-in, OSS tradition).
- Enterprise customers: limited brand overlay (logo + primary accent) on the web app.

---

## 14. Governance

- Design system council: product design lead, frontend lead, mobile lead, accessibility expert.
- Meets monthly.
- Proposals go through `design-tokens/` PRs with before/after screenshots.
- Breaking token changes follow a deprecation cycle (one major release).

---

## 15. Versioning & Migration

- Semver for the published token package.
- Breaking changes: 90-day deprecation window, codemods where possible.
- Release notes mention visual changes prominently.

---

*— End of Design System & Tokens —*
