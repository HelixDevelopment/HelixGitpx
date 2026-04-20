# sdk-examples/python/example.py
# HelixGitpx Python SDK — practical usage examples.
#
# Install:
#   pip install helixgitpx-sdk
#
# Async-first; a sync wrapper is available via `helixgitpx.sync`.
# Internally uses grpclib / httpx with retries + idempotency + OTel.
#
# Run:  HGX_PAT=hpxat_xxx python sdk-examples/python/example.py
import asyncio
import logging
import os
from datetime import timedelta

from helixgitpx import (
    HelixClient,
    Visibility,
    CaseStatus,
    MergeStrategy,
    HelixError,
    UnavailableError,
)

logging.basicConfig(level=logging.INFO)
log = logging.getLogger("helixgitpx.example")


async def main() -> None:
    pat = os.environ.get("HGX_PAT")
    if not pat:
        raise SystemExit("HGX_PAT env var required")

    # Client:
    # - baseUrl selects transport (https → Connect-JSON by default; grpc+tls → gRPC).
    # - retries={max:5, on:[UNAVAILABLE, ABORTED], jitter:True}
    # - idempotency_key auto-generated for writes.
    # - OTel spans emitted if an SDK is configured.
    async with HelixClient(
        base_url="https://api.helixgitpx.example.com",
        pat=pat,
        timeout=timedelta(seconds=30),
        user_agent="helixgitpx-example-py/1.0",
    ) as hc:
        # 1. Who am I?
        me = await hc.auth.get_me()
        log.info("Hello, %s (%s)", me.display_name, me.email)

        # 2. Create a repo
        org = await hc.org.get_org(slug="acme")
        repo = await hc.repo.create_repo(
            org_id=org.id,
            slug=f"demo-py-{asyncio.get_event_loop().time():.0f}",
            display_name="Python SDK demo",
            visibility=Visibility.INTERNAL,
            default_branch="main",
            auto_bind_all_enabled_upstreams=True,
            init_with_readme=True,
        )
        log.info("Created repo %s (id=%s)", repo.slug, repo.id.value)

        # 3. Watch events in the background
        watch_task = asyncio.create_task(_watch(hc, repo.id))

        # 4. Open + merge a PR
        try:
            pr = await hc.pr.create_pr(
                repo_id=repo.id,
                title="Hello, world",
                body="Created via Python SDK example",
                head_ref="feature/hello",
                base_ref="main",
                labels=["example"],
            )
            log.info("Opened PR #%d", pr.number)

            await hc.pr.merge_pr(
                id=pr.id,
                strategy=MergeStrategy.SQUASH,
                commit_title="Hello, world!",
                delete_source_branch=True,
            )
            log.info("Merged.")
        except HelixError as e:
            log.warning("PR flow failed: %s", e)

        # 5. Paginate escalated conflicts and resolve one with AI
        async for c in hc.conflict.list_cases_iter(
            repo_id=repo.id, status=CaseStatus.ESCALATED
        ):
            log.info("Conflict %s kind=%s", c.id.value, c.kind.name)

            proposals = await hc.conflict.propose_resolutions(
                case_id=c.id, max_proposals=3, use_ai=True
            )
            if not proposals.items:
                continue
            best = proposals.items[0]

            try:
                await hc.conflict.apply_resolution(
                    case_id=c.id,
                    strategy=best.strategy,
                    apply_plan=best.apply_plan,
                    comment="Applied via Python SDK example",
                )
            except UnavailableError:
                log.warning("Transient error, retry driven by client automatically")
            break  # demo

        await asyncio.sleep(30)
        watch_task.cancel()
        try:
            await watch_task
        except asyncio.CancelledError:
            pass

    log.info("Done.")


async def _watch(hc: HelixClient, repo_id) -> None:
    """Subscribe to repo events with automatic resume token handling.

    The SDK's watch_repo is an async generator; simply consume it.
    Reconnection + resume is built in (token persisted in-memory).
    """
    async for ev in hc.repo.watch_repo(
        repo_id=repo_id,
        event_types=["ref.*", "pr.*", "issue.*"],
        max_reconnect_backoff=timedelta(seconds=30),
    ):
        log.info("event %s @ %s", ev.event_type, ev.occurred_at.isoformat())


if __name__ == "__main__":
    asyncio.run(main())

# ---- Sync alternative (for scripts) --------------------------
# from helixgitpx.sync import HelixClient as SyncClient
# with SyncClient(base_url=..., pat=...) as hc:
#     me = hc.auth.get_me()
#     print(me.display_name)
