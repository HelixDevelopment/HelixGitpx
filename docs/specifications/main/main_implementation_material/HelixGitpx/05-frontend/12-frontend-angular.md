# 12 — Frontend (Angular Web Application)

> **Document purpose**: Define the **architecture, tooling, state management, styling, accessibility, and testing** of the HelixGitpx web application.

---

## 1. Overview

| Attribute | Value |
|---|---|
| Framework | **Angular 19+** (standalone components, signals, deferred loading) |
| Language | TypeScript 5.5+ |
| State | **NgRx Signal Store** + per-feature component state; RxJS for event streams |
| Styling | **Tailwind CSS 3.4** + shadcn-for-Angular patterns + CSS variables for theming |
| Build / workspace | **Nx** monorepo |
| RPC | **Connect-Go** clients (gRPC-Web transparent) via `@connectrpc/connect-web` |
| Testing | Jest (unit), Playwright (E2E), Storybook (visual) |
| PWA | `@angular/pwa` — offline shell, background sync |
| a11y | Angular CDK a11y; axe-core checks in CI; WCAG 2.2 AA |
| i18n | Angular `@angular/localize` + ICU MessageFormat via Weblate |

---

## 2. Repository Layout (Nx Monorepo)

```
apps/
├── web/                        # main Angular app
│   └── src/app/
│       ├── app.config.ts
│       ├── app.routes.ts
│       ├── layout/
│       ├── features/
│       │   ├── dashboard/
│       │   ├── repos/
│       │   ├── pull-requests/
│       │   ├── issues/
│       │   ├── conflicts/
│       │   ├── upstreams/
│       │   ├── settings/
│       │   └── admin/
│       └── shell/
├── storybook/
└── e2e/                        # Playwright

libs/
├── ui/                         # shadcn-style design system components
├── data-access/
│   ├── auth/
│   ├── orgs/
│   ├── repos/
│   ├── events/
│   ├── conflicts/
│   └── ai/
├── feature/                    # smart feature components (optional layer)
├── util/                       # pure utils
├── sdk/                        # generated Connect clients + proto types
└── assets/                     # theme, icons, fonts
```

Nx enforces **module boundaries** via `eslint-plugin-boundaries`; e.g., `features/*` may depend on `data-access/*`, not the other way around.

---

## 3. Rendering Strategy

- **Client-side rendering** by default (SPA).
- **SSR** for marketing/landing pages and public repo views (`@angular/ssr` on Cloudflare Workers).
- **View Transitions API** for route changes where supported.
- **Deferred loading** (`@defer`) for below-the-fold features.
- **Service Worker** for offline shell (read-only view of last-visited repos).

---

## 4. State Management

### 4.1 Per-Feature Signal Store

Each feature owns a **Signal Store** (NgRx Signal Store) that holds:

- Local UI state (filters, selections).
- Derived computed signals.
- Effects for async ops.
- Subscriptions to live events.

### 4.2 Shared Stores

- `authStore` — current principal, tokens, refresh logic.
- `orgsStore` — active org, available orgs.
- `prefsStore` — UI preferences.
- `themeStore` — dark/light/system.
- `eventsStore` — multiplexed subscription broker; features subscribe to filtered sub-streams.

### 4.3 Live Events Integration

The global `eventsStore` opens a single gRPC-Web stream and fans events out to interested signals. Feature stores register selectors like:

```ts
private readonly events = inject(EventsStore);
private readonly repoId = signal<string>('');

constructor() {
  this.events
    .scoped({ repo_id: this.repoId() }, ['ref.*', 'pr.*'])
    .subscribe(ev => this.ingest(ev));
}
```

Stream is persistent across route navigations; scope changes update the subscription server-side.

### 4.4 Optimistic Updates

Writes apply immediately to local state, with an in-flight marker. On server confirmation → clear marker; on error → rollback + toast.

### 4.5 Caching & ETag

- HTTP cache via `HttpInterceptor` that honours `ETag` / `If-None-Match`.
- In-memory cache (LRU) for list views; invalidated by relevant events.

---

## 5. Routing & Navigation

- **Standalone routes** with `loadComponent` for code-splitting.
- Routes guarded by `canActivate`:
  - `authGuard` — redirect to login if no session.
  - `orgGuard` — ensure active org in path.
  - `permissionGuard` — OPA decision cached per route.
- Deep-linkable state (filters, sort, etc.) via query params.

### 5.1 Example Route Tree

```
/                                    → dashboard
/login
/callback
/orgs/:orgSlug
/orgs/:orgSlug/repos
/orgs/:orgSlug/repos/:repoSlug
/orgs/:orgSlug/repos/:repoSlug/code
/orgs/:orgSlug/repos/:repoSlug/prs
/orgs/:orgSlug/repos/:repoSlug/prs/:num
/orgs/:orgSlug/repos/:repoSlug/issues
/orgs/:orgSlug/repos/:repoSlug/conflicts
/orgs/:orgSlug/repos/:repoSlug/settings
/orgs/:orgSlug/settings
/orgs/:orgSlug/members
/orgs/:orgSlug/upstreams
/settings/profile
/settings/security
/admin                               (role=admin)
```

---

## 6. Design System

- **Tailwind utility-first** with custom theme tokens (CSS custom properties) for light/dark + brand palette (HelixGitpx green + teal).
- **shadcn-style primitives** adapted for Angular: Button, Input, Dialog, Dropdown, Tooltip, Toast, Tabs, Table, Card, Drawer, etc.
- **Headless accessibility** via **CDK** (a11y, overlay, drag-drop).
- **Icons**: Lucide Angular.
- **Motion**: `@angular/animations` for micro-interactions; Framer-Motion-lite style patterns.
- **Density modes**: comfortable / compact.
- **Typography**: Inter (UI) + JetBrains Mono (code).

### 6.1 Tokens (excerpt)

```css
:root {
  --color-primary-50:  #f0fbe8;
  --color-primary-500: #9acd32;
  --color-primary-600: #7fb41f;
  --color-accent-500:  #7cd1c2;
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
  --shadow-md: 0 4px 12px rgba(0,0,0,0.08);
  --font-sans: Inter, system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', ui-monospace, monospace;
}

[data-theme='dark'] {
  --color-bg: #0b0f13;
  --color-fg: #e6edf3;
  /* … */
}
```

### 6.2 Component Library

All primitives live in `libs/ui` with:
- Type-safe props (strict templates).
- ARIA attributes.
- Dark-mode aware.
- Storybook stories + Chromatic visual tests.
- Keyboard navigation tests.

---

## 7. Accessibility

- WCAG 2.2 AA at GA; AAA for critical flows by Y2.
- Automated checks:
  - `axe-core` on every page in E2E.
  - Lighthouse CI threshold: a11y ≥ 95.
- Manual checks: screen reader walk-through per release (NVDA + VoiceOver + TalkBack).
- Focus management: trap in modals; restore on close.
- Keyboard shortcuts with discoverable help (`?` opens cheatsheet).
- Reduced-motion support.

---

## 8. i18n

- Source locale: `en`.
- GA locales: `de`, `fr`, `es`, `pt`, `sr`, `ru`, `zh-CN`, `ja`.
- RTL by Y2 (`ar`, `he`).
- Managed in **Weblate**.
- ICU MessageFormat for plurals and gender.
- Fallback chain: requested → base → English.

---

## 9. Performance Budgets

| Metric | Target |
|---|---|
| First Contentful Paint (4G, Moto G5) | ≤ 1.5 s |
| Largest Contentful Paint | ≤ 2.5 s |
| Total Blocking Time | ≤ 200 ms |
| Cumulative Layout Shift | ≤ 0.1 |
| INP (Interaction to Next Paint) | ≤ 200 ms |
| Main bundle (gzip) | ≤ 250 KB |
| Per-route lazy chunks | ≤ 80 KB |

Checked in CI with Lighthouse-CI + size-limit.

---

## 10. Security on the Client

- Strict Content-Security-Policy (no inline script; nonce-based).
- Trusted Types enforced.
- Sub-Resource Integrity on third-party scripts (ideally none).
- Cookies: `HttpOnly`, `Secure`, `SameSite=Lax`.
- Refresh token stored in HttpOnly cookie; access token in-memory only.
- CSRF: double-submit token for REST mutations.
- Never render unsanitised Markdown; use DOMPurify + a curated allowlist.
- Every URL validated (no `javascript:` etc.).

---

## 11. Observability

- **OpenTelemetry Web** SDK → ingest into backend OTel collector via CORS endpoint.
- Trace context (W3C traceparent) propagated into gRPC-Web calls.
- Real User Monitoring (RUM): LCP, INP, CLS, errors; sent anonymised.
- Error reporting: Sentry (self-hosted) with source maps.

---

## 12. PWA / Offline

- Offline shell: app loads; last-seen repo and PR lists cached in IndexedDB (via Dexie).
- Background Sync: queued writes (e.g. adding a comment offline) replayed on reconnect with idempotency keys.
- Install banner with screenshots manifest.
- Push notifications (for channels the user opts in to) via Web Push with VAPID keys.

---

## 13. Testing

- **Unit** (Jest): all components, stores, services; coverage ≥ 100 %.
- **Component** (Storybook + testing-library): interaction tests on every primitive.
- **Visual** (Chromatic or Percy): catch regressions.
- **E2E** (Playwright):
  - Login → create repo → bind upstream → push via fake adapter → conflict → resolve.
  - Accessibility scan at each page.
  - Cross-browser: Chromium, Firefox, WebKit.
  - Performance assertions (LCP).
- **Contract**: TypeScript types regenerated from protobufs; `tsc --noEmit` enforces compatibility.
- **Mutation** (Stryker): coverage-quality sanity.

---

## 14. Build & Deploy

- Production build: `nx build web --configuration=production` → emits static assets.
- Served via Cloudflare / Fastly CDN with cache-busting hashes and `immutable` cache-control for hashed assets.
- Per-route chunking.
- Feature flags via **OpenFeature** + Unleash SDK; evaluated at runtime.
- Rollback: swap CDN origin version pointer (blue/green).

---

## 15. Dev Experience

- **Devbox** + **mise** for local toolchain.
- Mock server (`@connectrpc/connect-node` + gRPC-Web) for offline dev.
- Storybook on port 6006; web app on 4200.
- `nx affected` in CI to only build/test changed projects.
- Git hooks (lefthook): format, lint, commit message (conventional commits).

---

*— End of Angular Frontend —*
