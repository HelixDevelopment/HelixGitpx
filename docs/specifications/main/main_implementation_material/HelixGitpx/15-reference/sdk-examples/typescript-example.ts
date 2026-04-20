// sdk-examples/typescript/example.ts
// HelixGitpx TypeScript SDK — practical usage examples.
//
// Install:
//   pnpm add @helixgitpx/sdk
//
// The SDK targets Node.js 20+, modern browsers, Deno, and Bun. It uses
// Connect-ES under the hood so the same code speaks gRPC-Web, Connect,
// or Connect-JSON depending on transport.
//
// Run:
//   HGX_PAT=hpxat_xxx tsx sdk-examples/typescript/example.ts
import {
  createClient,
  type HelixGitpxClient,
  type RepoEvent,
  Visibility,
  CaseStatus,
  MergeStrategy,
} from "@helixgitpx/sdk";
import { setTimeout as wait } from "node:timers/promises";

async function main() {
  const pat = process.env.HGX_PAT;
  if (!pat) throw new Error("HGX_PAT env var required");

  // Shared client:
  // - Transport chosen automatically (gRPC node, Connect browser).
  // - Retries on UNAVAILABLE / ABORTED with jittered backoff.
  // - Idempotency-Key generated for writes.
  // - OpenTelemetry spans emitted if an SDK is installed.
  const hc: HelixGitpxClient = createClient({
    baseUrl: "https://api.helixgitpx.example.com",
    pat,
    timeoutMs: 30_000,
    userAgent: "helixgitpx-example-ts/1.0",
  });

  // 1. Who am I?
  const me = await hc.auth.getMe({});
  console.log(`Hello, ${me.displayName} (${me.email})`);

  // 2. Resolve org id.
  const org = await hc.org.getOrg({ slug: "acme" });

  // 3. Create a repo with fan-out on.
  const repo = await hc.repo.createRepo({
    orgId: org.id,
    slug: `demo-${Date.now()}`,
    displayName: "TS SDK demo",
    visibility: Visibility.INTERNAL,
    defaultBranch: "main",
    autoBindAllEnabledUpstreams: true,
    initWithReadme: true,
  });
  console.log(`Created ${repo.slug} (id=${repo.id.value})`);

  // 4. Watch repo events — async iterator handles reconnect + resume.
  const watchAbort = new AbortController();
  void watchRepo(hc, repo.id, watchAbort.signal);

  // 5. Open + merge a PR (assumes branch already exists).
  try {
    const pr = await hc.pr.createPR({
      repoId: repo.id,
      title: "Hello, world",
      body: "Created via TypeScript SDK example",
      headRef: "feature/hello",
      baseRef: "main",
      labels: ["example"],
    });
    console.log(`Opened PR #${pr.number}`);

    await hc.pr.mergePR({
      id: pr.id,
      strategy: MergeStrategy.SQUASH,
      commitTitle: "Hello, world!",
      deleteSourceBranch: true,
    });
    console.log("Merged.");
  } catch (err: unknown) {
    console.warn("PR flow failed:", err);
  }

  // 6. Paginate escalated conflicts, apply best proposal.
  for await (const c of hc.conflict.listCasesIter({
    repoId: repo.id,
    status: CaseStatus.ESCALATED,
  })) {
    console.log(`Conflict ${c.id.value} kind=${c.kind}`);

    const proposals = await hc.conflict.proposeResolutions({
      caseId: c.id,
      maxProposals: 3,
    });
    const best = proposals.items[0];
    if (!best) continue;

    await hc.conflict.applyResolution({
      caseId: c.id,
      strategy: best.strategy,
      applyPlan: best.applyPlan,
      comment: "Applied via TS SDK example",
    });
    break; // demo
  }

  // Let the watcher run for 30s then exit.
  await wait(30_000);
  watchAbort.abort();
  console.log("Done.");
}

async function watchRepo(
  hc: HelixGitpxClient,
  repoId: { value: string },
  signal: AbortSignal
): Promise<void> {
  let resumeToken: string | undefined;

  while (!signal.aborted) {
    try {
      // The async iterator handles streaming and surfaces events as they arrive.
      const stream = hc.repo.watchRepo(
        {
          repoId,
          resumeToken,
          eventTypes: ["ref.*", "pr.*", "issue.*"],
        },
        { signal }
      );

      for await (const ev of stream as AsyncIterable<RepoEvent>) {
        console.log(
          `event ${ev.eventType} @ ${ev.occurredAt?.toDate().toISOString()}`
        );
        resumeToken = ev.resumeToken;
      }
    } catch (err) {
      if (signal.aborted) return;
      console.warn("watch error, reconnecting in 2s:", err);
      await wait(2000);
    }
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});

// Browser WebSocket fallback (for completeness):
//
// import { createWebClient } from "@helixgitpx/sdk/web";
//
// const hc = createWebClient({
//   baseUrl: "https://api.helixgitpx.example.com",
//   token: localStorage.getItem("hgx_access_token")!,
// });
//
// // When the browser blocks gRPC, the client transparently uses
// // Connect-JSON over HTTPS and WebSocket for streaming.
