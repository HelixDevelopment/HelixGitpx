// sdk-examples/go/main.go
// HelixGitpx Go SDK — practical usage examples.
//
// Install:
//   go get github.com/vasic-digital/helixgitpx-go@v1
//
// The SDK wraps the generated gRPC clients with Connect interceptors
// for auth, retries, idempotency, and OTel instrumentation.
//
// Examples in this file:
//   1. Connect, authenticate with a PAT, fetch current user.
//   2. Create a repo with fan-out to all enabled upstreams.
//   3. Watch repo events (bidi stream) and print ref updates.
//   4. Open a PR and merge it.
//   5. Paginate through conflicts and resolve one with human strategy.
//   6. Resume a watch after a reconnect using a resume token.
//
// Run:  HGX_PAT=hpxat_xxxxx go run sdk-examples/go/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/vasic-digital/helixgitpx-go/v1/client"
	helixgitpxv1 "github.com/vasic-digital/helixgitpx-go/v1/gen/helixgitpx/v1"
)

func main() {
	pat := os.Getenv("HGX_PAT")
	if pat == "" {
		log.Fatal("HGX_PAT env var required")
	}

	// ---- 1. Client construction ----------------------------
	// The shared client has sensible defaults:
	// - HTTP/2 with keep-alives
	// - gRPC + Connect protocol negotiation
	// - Exponential backoff for retryable codes (UNAVAILABLE, ABORTED)
	// - OTel instrumentation (propagates trace context)
	// - Automatic PAT header injection
	// - Per-request idempotency key generation for writes
	hc, err := client.New(
		client.WithBaseURL("https://api.helixgitpx.example.com"),
		client.WithPAT(pat),
		client.WithTimeout(30*time.Second),
		client.WithUserAgent("helixgitpx-example/1.0"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	me, err := hc.Auth.GetMe(ctx, connect.NewRequest(&helixgitpxv1.GetMeRequest{}))
	if err != nil {
		log.Fatalf("GetMe: %v", err)
	}
	fmt.Printf("Hello, %s (%s)\n", me.Msg.DisplayName, me.Msg.Email)

	// ---- 2. Create a repo ---------------------------------
	org := "acme"                       // your org slug
	repoReq := &helixgitpxv1.CreateRepoRequest{
		OrgId:                      mustOrgID(ctx, hc, org),
		Slug:                       fmt.Sprintf("demo-%d", time.Now().Unix()),
		DisplayName:                "Demo from the Go SDK",
		Visibility:                 helixgitpxv1.Visibility_VISIBILITY_INTERNAL,
		DefaultBranch:              "main",
		AutoBindAllEnabledUpstreams: true,
		InitWithReadme:             true,
	}
	repo, err := hc.Repo.CreateRepo(ctx, connect.NewRequest(repoReq))
	if err != nil {
		log.Fatalf("CreateRepo: %v", err)
	}
	fmt.Printf("Created repo %s (id=%s)\n", repo.Msg.Slug, repo.Msg.Id.Value)

	// ---- 3. Watch repo events (bidi stream) ---------------
	// The client handles reconnects with resume tokens automatically.
	watchCtx, watchCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer watchCancel()

	go watchRepo(watchCtx, hc, repo.Msg.Id)

	// ---- 4. Open a PR + merge ------------------------------
	// Assumes there's a branch `feature/hello` pointing at a real commit.
	pr, err := hc.PR.CreatePR(ctx, connect.NewRequest(&helixgitpxv1.CreatePRRequest{
		RepoId:   repo.Msg.Id,
		Title:    "Hello, world",
		Body:     "Demo PR created via the Go SDK.",
		HeadRef:  "feature/hello",
		BaseRef:  "main",
		Draft:    false,
		Labels:   []string{"example"},
	}))
	if err != nil {
		log.Printf("CreatePR: %v (skipping merge)", err)
	} else {
		fmt.Printf("Opened PR #%d\n", pr.Msg.Number)
		_, merr := hc.PR.MergePR(ctx, connect.NewRequest(&helixgitpxv1.MergePRRequest{
			Id:                 pr.Msg.Id,
			Strategy:           helixgitpxv1.MergeStrategy_SQUASH,
			CommitTitle:        "Hello, world!",
			DeleteSourceBranch: true,
		}))
		if merr != nil {
			// FAILED_PRECONDITION commonly means protection rules unmet.
			log.Printf("MergePR: %v", merr)
		} else {
			fmt.Println("Merged.")
		}
	}

	// ---- 5. Paginate conflicts + resolve one --------------
	iter := hc.Conflict.ListCasesIter(ctx, &helixgitpxv1.ListCasesRequest{
		RepoId: repo.Msg.Id,
		Status: helixgitpxv1.CaseStatus_ESCALATED,
	})
	for iter.Next() {
		c := iter.Value()
		fmt.Printf("Conflict %s kind=%s subject=%s\n", c.Id.Value, c.Kind, c.Subject)

		// Ask the service for proposals, pick the highest-confidence one.
		pr, err := hc.Conflict.ProposeResolutions(ctx,
			connect.NewRequest(&helixgitpxv1.ProposeRequest{CaseId: c.Id, MaxProposals: 3}))
		if err != nil {
			log.Printf("propose: %v", err)
			continue
		}
		if len(pr.Msg.Items) == 0 {
			continue
		}
		best := pr.Msg.Items[0]
		_, err = hc.Conflict.ApplyResolution(ctx, connect.NewRequest(&helixgitpxv1.ApplyRequest{
			CaseId:    c.Id,
			Strategy:  best.Strategy,
			ApplyPlan: best.ApplyPlan,
			Comment:   "Applied via Go SDK example",
		}))
		if err != nil {
			log.Printf("apply: %v", err)
		}
		break // demo: only one
	}
	if err := iter.Err(); err != nil {
		log.Printf("iter err: %v", err)
	}

	fmt.Println("Done.")
}

// watchRepo streams repo events and keeps a resume token for recovery.
func watchRepo(ctx context.Context, hc *client.Client, repoID *helixgitpxv1.UUID) {
	var resume string
	backoff := 1 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		stream, err := hc.Repo.WatchRepo(ctx, connect.NewRequest(&helixgitpxv1.WatchRepoRequest{
			RepoId:      repoID,
			ResumeToken: resume,
			EventTypes:  []string{"ref.*", "pr.*", "issue.*"},
		}))
		if err != nil {
			log.Printf("watch start: %v (retry in %s)", err, backoff)
			time.Sleep(backoff)
			backoff = min(backoff*2, 30*time.Second)
			continue
		}
		backoff = 1 * time.Second

		for stream.Receive() {
			ev := stream.Msg()
			fmt.Printf("event %s @ %s\n", ev.EventType, ev.OccurredAt.AsTime().Format(time.RFC3339))
			resume = ev.ResumeToken
		}
		if err := stream.Err(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
			log.Printf("watch err: %v (reconnecting)", err)
		}
	}
}

// mustOrgID resolves an org slug to an id, panicking on error.
func mustOrgID(ctx context.Context, hc *client.Client, slug string) *helixgitpxv1.UUID {
	resp, err := hc.Org.GetOrg(ctx, connect.NewRequest(&helixgitpxv1.GetOrgRequest{Slug: slug}))
	if err != nil {
		log.Fatalf("resolve org: %v", err)
	}
	return resp.Msg.Id
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
