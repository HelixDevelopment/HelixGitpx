// sdk-examples/kotlin/Example.kt
// HelixGitpx Kotlin SDK — practical usage examples.
//
// Works on JVM, Android, and as part of a Kotlin Multiplatform
// commonMain module. Uses kotlinx.coroutines + Flow for streaming.
//
// Install (JVM / Android):
//   implementation("io.helixgitpx.sdk:helixgitpx-sdk:1.0.0")
//
// Install (KMP):
//   implementation("io.helixgitpx.sdk:helixgitpx-sdk-core:1.0.0")
//
// Run (JVM):
//   HGX_PAT=hpxat_xxx ./gradlew :examples:jvmRun

package io.helixgitpx.examples

import io.helixgitpx.sdk.HelixClient
import io.helixgitpx.sdk.auth.Pat
import io.helixgitpx.sdk.model.CaseStatus
import io.helixgitpx.sdk.model.ConflictKind
import io.helixgitpx.sdk.model.MergeStrategy
import io.helixgitpx.sdk.model.Visibility
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.onEach
import kotlinx.coroutines.flow.retryWhen
import kotlin.time.Duration.Companion.minutes
import kotlin.time.Duration.Companion.seconds

fun main(): Unit = runBlocking {
    val pat = System.getenv("HGX_PAT") ?: error("HGX_PAT env var required")

    // Shared HelixClient:
    // - Uses OkHttp on JVM/Android, Ktor elsewhere
    // - Automatic retries on UNAVAILABLE/ABORTED
    // - Coroutine-friendly everywhere; Flow for streams
    // - OpenTelemetry Kotlin SDK integrated when present
    val client = HelixClient.builder()
        .baseUrl("https://api.helixgitpx.example.com")
        .auth(Pat(pat))
        .userAgent("helixgitpx-example-kt/1.0")
        .build()

    client.use { hc ->
        // 1. Who am I?
        val me = hc.auth.getMe()
        println("Hello, ${me.displayName} (${me.email})")

        // 2. Create a repo
        val org = hc.org.getOrg(slug = "acme")
        val repo = hc.repo.createRepo {
            orgId = org.id
            slug = "demo-${System.currentTimeMillis()}"
            displayName = "KMP SDK demo"
            visibility = Visibility.INTERNAL
            defaultBranch = "main"
            autoBindAllEnabledUpstreams = true
            initWithReadme = true
        }
        println("Created ${repo.slug} (id=${repo.id.value})")

        // 3. Watch repo events in the background with auto-resume
        val watchJob = launch {
            watchRepoWithResume(hc, repo.id)
        }

        // 4. Open + merge a PR
        try {
            val pr = hc.pr.createPR {
                repoId = repo.id
                title = "Hello, world"
                body = "Created via Kotlin SDK example"
                headRef = "feature/hello"
                baseRef = "main"
                labels = listOf("example")
            }
            println("Opened PR #${pr.number}")

            hc.pr.mergePR {
                id = pr.id
                strategy = MergeStrategy.SQUASH
                commitTitle = "Hello, world!"
                deleteSourceBranch = true
            }
            println("Merged.")
        } catch (e: Exception) {
            println("PR flow failed: $e")
        }

        // 5. Paginate & apply best conflict proposal
        hc.conflict.listCasesFlow(
            repoId = repo.id,
            status = CaseStatus.ESCALATED,
        ).onEach { c ->
            println("Conflict ${c.id.value} kind=${c.kind}")
            if (c.kind == ConflictKind.LFS_DIVERGENCE) return@onEach   // always human

            val proposals = hc.conflict.proposeResolutions(caseId = c.id, maxProposals = 3)
            val best = proposals.items.firstOrNull() ?: return@onEach
            hc.conflict.applyResolution {
                caseId = c.id
                strategy = best.strategy
                applyPlan = best.applyPlan
                comment = "Applied via Kotlin SDK example"
            }
        }.collect()

        delay(30.seconds)
        watchJob.cancel()
        println("Done.")
    }
}

private suspend fun watchRepoWithResume(hc: HelixClient, repoId: io.helixgitpx.sdk.model.UUID) {
    var resumeToken: String? = null
    hc.repo.watchRepoFlow(
        repoId = repoId,
        resumeToken = { resumeToken },
        eventTypes = listOf("ref.*", "pr.*", "issue.*"),
    )
        .retryWhen { cause, attempt ->
            val delay = (1L shl minOf(attempt.toInt(), 5)).coerceAtMost(30) * 1000L
            println("watch err $cause (attempt=$attempt, retrying in ${delay}ms)")
            delay(delay)
            true
        }
        .catch { println("fatal watch error: $it") }
        .collect { ev ->
            println("event ${ev.eventType} @ ${ev.occurredAt}")
            resumeToken = ev.resumeToken
        }
}

/* ===== Android-specific notes ================================
 * On Android, inject HelixClient via Hilt / Koin:
 *
 *   @Provides @Singleton
 *   fun helixClient(authTokenStore: AuthTokenStore): HelixClient =
 *       HelixClient.builder()
 *           .baseUrl(BuildConfig.HGX_BASE_URL)
 *           .authProvider { authTokenStore.currentBearer() }
 *           .applicationContext(context)    // enables WorkManager-backed retries
 *           .build()
 *
 * Use WorkManager for the offline outbox replay:
 *
 *   HelixOutboxWorker.enqueuePeriodic(context)
 *
 * On iOS (via KMP), the same HelixClient works; hook WatchRepoFlow to
 * a SwiftUI view model with Flow.asAsyncSequence() (via KMP-NativeCoroutines).
 * ============================================================ */
