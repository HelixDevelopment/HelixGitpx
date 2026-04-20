// sdk-examples/swift/Example.swift
// HelixGitpx Swift SDK — practical usage examples for iOS / macOS.
//
// Install (SwiftPM):
//   .package(url: "https://github.com/vasic-digital/helixgitpx-swift", from: "1.0.0")
//
// Uses Connect-Swift under the hood. Swift Concurrency (async/await)
// is first class; AsyncSequence for streaming.

import Foundation
import HelixGitpxSDK

@main
struct Example {
    static func main() async throws {
        guard let pat = ProcessInfo.processInfo.environment["HGX_PAT"] else {
            fatalError("HGX_PAT env var required")
        }

        // Client
        // - URLSession-backed transport with automatic retries
        // - Keychain-backed token storage (on Apple platforms)
        // - Idempotency-Key auto-generated for writes
        // - OSLog integration + OTel bridge when configured
        let client = HelixClient(
            baseURL: URL(string: "https://api.helixgitpx.example.com")!,
            pat: pat,
            userAgent: "helixgitpx-example-swift/1.0",
            timeout: 30
        )

        // 1. Who am I?
        let me = try await client.auth.getMe()
        print("Hello, \(me.displayName) (\(me.email))")

        // 2. Create a repo
        let org = try await client.org.getOrg(slug: "acme")
        let repo = try await client.repo.createRepo(
            orgId: org.id,
            slug: "demo-swift-\(Int(Date().timeIntervalSince1970))",
            displayName: "Swift SDK demo",
            visibility: .internal,
            defaultBranch: "main",
            autoBindAllEnabledUpstreams: true,
            initWithReadme: true
        )
        print("Created repo \(repo.slug) (id=\(repo.id.value))")

        // 3. Watch repo events with auto-resume
        let watchTask = Task {
            try await watchRepo(client: client, repoId: repo.id)
        }

        // 4. Open + merge a PR
        do {
            let pr = try await client.pr.createPR(
                repoId: repo.id,
                title: "Hello, world",
                body: "Created via Swift SDK example",
                headRef: "feature/hello",
                baseRef: "main",
                labels: ["example"]
            )
            print("Opened PR #\(pr.number)")

            _ = try await client.pr.mergePR(
                id: pr.id,
                strategy: .squash,
                commitTitle: "Hello, world!",
                deleteSourceBranch: true
            )
            print("Merged.")
        } catch {
            print("PR flow failed: \(error)")
        }

        // 5. Apply best AI proposal for one escalated conflict
        for try await conflict in client.conflict.listCases(
            repoId: repo.id,
            status: .escalated
        ) {
            print("Conflict \(conflict.id.value) kind=\(conflict.kind)")

            let proposals = try await client.conflict.proposeResolutions(
                caseId: conflict.id,
                maxProposals: 3
            )
            guard let best = proposals.items.first else { continue }

            _ = try await client.conflict.applyResolution(
                caseId: conflict.id,
                strategy: best.strategy,
                applyPlan: best.applyPlan,
                comment: "Applied via Swift SDK example"
            )
            break  // demo
        }

        try? await Task.sleep(for: .seconds(30))
        watchTask.cancel()

        print("Done.")
    }

    static func watchRepo(client: HelixClient, repoId: UUID) async throws {
        var resume: String? = nil
        while !Task.isCancelled {
            do {
                for try await ev in client.repo.watchRepo(
                    repoId: repoId,
                    resumeToken: resume,
                    eventTypes: ["ref.*", "pr.*", "issue.*"]
                ) {
                    print("event \(ev.eventType) @ \(ev.occurredAt)")
                    resume = ev.resumeToken
                }
            } catch {
                if Task.isCancelled { return }
                print("watch error: \(error) — reconnecting")
                try await Task.sleep(for: .seconds(2))
            }
        }
    }
}

/* ===== SwiftUI integration =================================
 * @MainActor
 * final class RepoViewModel: ObservableObject {
 *     @Published var events: [RepoEvent] = []
 *     private let client: HelixClient
 *     private var watchTask: Task<Void, Never>?
 *
 *     init(client: HelixClient) { self.client = client }
 *
 *     func start(repoId: UUID) {
 *         watchTask = Task {
 *             for await ev in client.repo.watchRepo(repoId: repoId, resumeToken: nil, eventTypes: ["*"]) {
 *                 events.append(ev)
 *             }
 *         }
 *     }
 *     func stop() { watchTask?.cancel() }
 * }
 * ============================================================ */
