# HelixGitpx Manuals

In-depth user & operator guides, exported in every major format per
[Constitution Article III §2](../../CONSTITUTION.md#3-article-iii--documentation).

## Content index

Source lives in [`src/`](./src/). One Markdown file per chapter.

| Manual | Audience | Source |
|--------|----------|--------|
| User Guide | End users (web/mobile/desktop) | [src/user-guide/](./src/user-guide/) |
| Operator Guide | SREs, platform operators | [src/operator-guide/](./src/operator-guide/) |
| Developer Guide | Service-level contributors | [src/developer-guide/](./src/developer-guide/) |
| Administrator Guide | Org admins / tenant owners | [src/administrator-guide/](./src/administrator-guide/) |
| API Reference | API consumers | [src/api-reference/](./src/api-reference/) |
| CLI Reference | CLI users | [src/cli-reference/](./src/cli-reference/) |
| Security Handbook | Customer CISOs / auditors | [src/security-handbook/](./src/security-handbook/) |
| Deployment Cookbook | Self-hosters | [src/deployment-cookbook/](./src/deployment-cookbook/) |
| Troubleshooting | All | [src/troubleshooting/](./src/troubleshooting/) |
| Migration Guide | Adopters from GitHub/GitLab/etc. | [src/migration-guide/](./src/migration-guide/) |

## Output formats

All ten manuals are published in **all** of the following formats on every
release. The export pipeline lives at [`tools/docs-export/`](../../tools/docs-export/).

| Format | Tool | Output path |
|--------|------|-------------|
| HTML (site) | Docusaurus | `impl/helixgitpx-docs-site/build/` |
| PDF | `pandoc` + weasyprint | `docs/manuals/dist/<name>.pdf` |
| ePub | `pandoc` | `docs/manuals/dist/<name>.epub` |
| MOBI | `ebook-convert` (Calibre) | `docs/manuals/dist/<name>.mobi` |
| DOCX | `pandoc` | `docs/manuals/dist/<name>.docx` |
| Markdown | native source | `docs/manuals/src/<name>/*.md` |
| Plain text | `pandoc --to plain` | `docs/manuals/dist/<name>.txt` |
| Zipped bundle | `zip` | `docs/manuals/dist/<name>.zip` |

All exports are regenerated from the Markdown source — single source of
truth, no format drift.

## Building locally

```bash
# HTML (inside Docusaurus)
cd impl/helixgitpx-docs-site && npx docusaurus build

# All offline formats (PDF, ePub, DOCX, txt)
bash tools/docs-export/build-all.sh
```

Output lands in `docs/manuals/dist/`.

## Status

This directory ships as scaffolding at GA. Chapter-by-chapter content is a
dedicated post-GA effort tracked in `docs/marketing/launch-checklist.md`
(§7 days before). Each chapter gets its own PR; expected cadence is one
manual per fortnight through Year 2.
