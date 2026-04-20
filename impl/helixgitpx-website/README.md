# HelixGitpx marketing website

Astro + Tailwind static site that ships at `helixgitpx.io`.

Distinct from the docs site (`docs.helixgitpx.io`, Docusaurus). This site
carries the landing page, pricing, trust center, customers, careers,
legal, and the marketing blog.

## Why Astro

- Content-first, zero-JS-by-default, fast Core Web Vitals.
- Component islands let us pull in small interactive pieces (forms, search)
  without shipping a framework for the whole site.
- Build output is plain HTML + CSS + a sprinkle of JS — perfect behind
  NGINX / edge caching.

## Develop

```bash
cd impl/helixgitpx-website
npm install
npm run dev          # http://localhost:4321
```

## Build

```bash
npm run build        # output in dist/
```

## Deploy

Same pattern as `impl/helixgitpx-docs-site/`: build to static assets,
containerize, deploy via the Argo CD `website` Application.

## Structure

- `src/layouts/BaseLayout.astro` — shared shell, nav, footer.
- `src/pages/` — one file per URL.
- `src/components/` — reusable islands.
- `src/styles/` — Tailwind globals.
- `public/` — static assets (logos, OG images, favicons).

## Brand

Colour palette and typography are synced with
[`docs/media/brand/style-guide.md`](../../docs/media/brand/style-guide.md).
The Tailwind config mirrors the palette.

## Content ownership

- Marketing-owned: `/`, `/why`, `/features`, `/pricing`, `/customers`,
  `/careers`, `/press`, `/contact`, blog, legal.
- Product-owned: `/trust`, `/security`, `/subprocessors`.
- Engineering-owned: `/api` summary (deep ref lives on docs).
