# HelixGitpx User Guide

## 1. Introduction

HelixGitpx is a federated Git proxy. From a single namespace, you can push
code to — and pull code from — a dozen different Git hosting providers.
The same repository can live simultaneously on GitHub, GitLab, Gitea,
Gitee, Bitbucket, GitFlic, GitVerse, Azure DevOps, AWS CodeCommit, Forgejo,
SourceHut, and any other Git-over-HTTPS host.

This manual teaches you the mental model, the tools, and the common flows.

### 1.1 Who is this guide for?

- Developers pushing code from a workstation.
- Reviewers who merge pull requests.
- Team leads who bind repositories to upstreams.
- End users of the HelixGitpx web, desktop, mobile apps.

Operators and infrastructure engineers should read the
[Operator Guide](../operator-guide/00-introduction.md) instead.

### 1.2 What you'll learn

- **Chapter 2:** creating an account and your first organization.
- **Chapter 3:** binding a repository to one or more upstreams.
- **Chapter 4:** pushing, pulling, and reviewing code.
- **Chapter 5:** conflict resolution when upstreams diverge.
- **Chapter 6:** using AI assistance for PR summaries and code search.
- **Chapter 7:** mobile and desktop apps.
- **Chapter 8:** account settings and data residency.
- **Chapter 9:** troubleshooting common issues.

### 1.3 Conventions used in this guide

- `inline code` — commands, file names, exact string values.
- Code blocks show the **full** command you paste into a terminal.
- Screenshots are illustrative; the UI may differ slightly per release.
- Tips are called out:

  > 💡 **Tip** — short advice worth following.

  > ⚠️ **Warning** — read this before acting.

### 1.4 Getting help

- Public docs: [`docs.helixgitpx.io`](https://docs.helixgitpx.io)
- Status: [`status.helixgitpx.io`](https://status.helixgitpx.io)
- Bug bounty: [`hackerone.com/helixgitpx`](https://hackerone.com/helixgitpx)
- Enterprise support: [`support@helixgitpx.io`](mailto:support@helixgitpx.io)

---
