# ADR-0039 — Docusaurus for the public docs site

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Options: Docusaurus (React-based, Facebook-maintained), MkDocs-Material
(Python-based, excellent search), Hugo (fast builds), GitBook (SaaS).

## Decision

Docusaurus. Consistent React tooling with the web app. Good
versioning, plugins, and community.

## Consequences

- Static site built in CI, shipped as NGINX image to `docs.helixgitpx.io`.
- Search via Algolia DocSearch (free tier for OSS docs).
- Locked in to Docusaurus v3; upgrade path is annual.

## Links

- Spec §LOCKED C-10
- impl/helixgitpx-docs-site/
