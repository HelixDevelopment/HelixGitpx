#!/usr/bin/env node
// Copy (or symlink under POSIX) the spec tree into Docusaurus-friendly docs/.
// Rewrites intra-spec links to Docusaurus routes.

import fs from "node:fs/promises";
import path from "node:path";
import { existsSync } from "node:fs";

const SRC = path.resolve("../../docs/specifications/main/main_implementation_material/HelixGitpx");
const DST = path.resolve("./docs");

async function copyTree(src, dst) {
  const entries = await fs.readdir(src, { withFileTypes: true });
  await fs.mkdir(dst, { recursive: true });
  for (const entry of entries) {
    const s = path.join(src, entry.name);
    const d = path.join(dst, entry.name);
    if (entry.isDirectory()) {
      await copyTree(s, d);
    } else if (entry.name.endsWith(".md")) {
      const raw = await fs.readFile(s, "utf8");
      const rewritten = raw
        .replace(/]\(\.\.\/(\d\d-[a-z-]+)\//g, "](../../$1/")
        .replace(/]\(docs\/specifications\/[^)]+\/HelixGitpx\//g, "](/");
      await fs.writeFile(d, rewritten);
    } else {
      await fs.copyFile(s, d);
    }
  }
}

if (!existsSync(SRC)) {
  console.error(`sync-docs: source not found at ${SRC}`);
  process.exit(1);
}
await fs.rm(DST, { recursive: true, force: true });
await copyTree(SRC, DST);
console.log(`sync-docs: copied spec tree → ${DST}`);
