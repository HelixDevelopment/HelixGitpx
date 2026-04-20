// examples/temporal/fanout_push.go
// Reference implementation of the FanOutPush Temporal workflow.
// This is the authoritative example of how HelixGitpx structures
// durable workflows:
//   - Deterministic workflow functions (no wall-clock, no I/O).
//   - Activities wrap all side effects.
//   - Retries + heartbeats configured per activity.
//   - Child workflows for per-upstream operations (parallel + independent).
//   - Signals for cancellation / resume.
//   - Queries for status.
//
// Package path: github.com/vasic-digital/helixgitpx/services/sync-orchestrator/workflow
package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	helixv1 "github.com/vasic-digital/helixgitpx/gen/helixgitpx/v1"
	"github.com/vasic-digital/helixgitpx/services/sync-orchestrator/activities"
)

// FanOutPushInput is the workflow input. Serializable; keep stable.
type FanOutPushInput struct {
	JobID        string             `json:"job_id"`
	RepoID       string             `json:"repo_id"`
	OrgID        string             `json:"org_id"`
	RefName      string             `json:"ref_name"`
	OldSHA       string             `json:"old_sha"`
	NewSHA       string             `json:"new_sha"`
	Origin       string             `json:"origin"`
	Bindings     []Binding          `json:"bindings"`
	CorrelationID string            `json:"correlation_id"`
	IdempotencyKey string           `json:"idempotency_key"`
}

// Binding describes a target upstream.
type Binding struct {
	UpstreamID   string `json:"upstream_id"`
	Provider     string `json:"provider"`
	RemoteOwner  string `json:"remote_owner"`
	RemoteName   string `json:"remote_name"`
	PushEnabled  bool   `json:"push_enabled"`
}

// StepResult captures per-upstream outcome. Persisted in workflow history.
type StepResult struct {
	UpstreamID   string    `json:"upstream_id"`
	Status       string    `json:"status"`       // succeeded|failed|skipped
	ErrorCode    string    `json:"error_code,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	RefsUpdated  int       `json:"refs_updated"`
	BytesOut     int64     `json:"bytes_out"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
}

// FanOutPushResult is the workflow output.
type FanOutPushResult struct {
	JobID      string       `json:"job_id"`
	Status     string       `json:"status"` // succeeded|partial|failed|cancelled
	Steps      []StepResult `json:"steps"`
	StartedAt  time.Time    `json:"started_at"`
	FinishedAt time.Time    `json:"finished_at"`
}

// ------------------------------------------------------------
// FanOutPush — replicate a local ref update to every enabled
// upstream binding in parallel. Writes events + updates sync.jobs.
// ------------------------------------------------------------

const (
	cancelSignalName = "cancel"
	statusQueryName  = "status"
)

func FanOutPush(ctx workflow.Context, in FanOutPushInput) (*FanOutPushResult, error) {
	log := workflow.GetLogger(ctx)
	log.Info("FanOutPush starting", "repo_id", in.RepoID, "ref", in.RefName)

	// Record start to the DB (Activity so it survives replay).
	setupOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    500 * time.Millisecond,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    5,
		},
	}
	setupCtx := workflow.WithActivityOptions(ctx, setupOpts)
	if err := workflow.ExecuteActivity(setupCtx,
		activities.RecordJobStart, in.JobID).Get(setupCtx, nil); err != nil {
		return nil, fmt.Errorf("record start: %w", err)
	}

	// Set up query + signal handlers.
	var (
		result = &FanOutPushResult{
			JobID:     in.JobID,
			Status:    "running",
			Steps:     make([]StepResult, 0, len(in.Bindings)),
			StartedAt: workflow.Now(ctx),
		}
		cancelled bool
	)
	if err := workflow.SetQueryHandler(ctx, statusQueryName, func() (*FanOutPushResult, error) {
		return result, nil
	}); err != nil {
		return nil, fmt.Errorf("set query: %w", err)
	}

	cancelCh := workflow.GetSignalChannel(ctx, cancelSignalName)

	// Watch for cancel signal in parallel with work.
	workflow.Go(ctx, func(innerCtx workflow.Context) {
		var reason string
		cancelCh.Receive(innerCtx, &reason)
		cancelled = true
		log.Warn("cancel signal received", "reason", reason)
	})

	// Fan out: one child workflow per binding for independence + replay safety.
	futures := make([]workflow.Future, 0, len(in.Bindings))
	for _, b := range in.Bindings {
		if !b.PushEnabled {
			result.Steps = append(result.Steps, StepResult{
				UpstreamID: b.UpstreamID,
				Status:     "skipped",
			})
			continue
		}

		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:        fmt.Sprintf("push-%s-%s", in.JobID, b.UpstreamID),
			TaskQueue:         "sync-orchestrator",
			ParentClosePolicy: temporal.PARENT_CLOSE_POLICY_TERMINATE,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    1 * time.Second,
				BackoffCoefficient: 2.0,
				MaximumInterval:    5 * time.Minute,
				MaximumAttempts:    6,
				NonRetryableErrorTypes: []string{
					"AdapterAuthFailed",
					"UpstreamCapabilityUnsupported",
					"BindingDisabled",
				},
			},
			WorkflowTaskTimeout: 30 * time.Second,
		}
		childCtx := workflow.WithChildOptions(ctx, cwo)
		f := workflow.ExecuteChildWorkflow(childCtx, PushToUpstream, PushToUpstreamInput{
			JobID:       in.JobID,
			RepoID:      in.RepoID,
			OrgID:       in.OrgID,
			Binding:     b,
			RefName:     in.RefName,
			OldSHA:      in.OldSHA,
			NewSHA:      in.NewSHA,
			Origin:      in.Origin,
		})
		futures = append(futures, f)
	}

	for _, f := range futures {
		if cancelled {
			break
		}
		var step StepResult
		err := f.Get(ctx, &step)
		if err != nil {
			// Convert to a canonical StepResult; never let a child kill the parent.
			step.Status = "failed"
			step.ErrorMessage = err.Error()
		}
		result.Steps = append(result.Steps, step)
	}

	// Compute overall status.
	result.FinishedAt = workflow.Now(ctx)
	result.Status = decideStatus(result.Steps, cancelled)

	// Persist + emit event.
	finalOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: 500 * time.Millisecond,
			MaximumAttempts: 10,            // important: we don't want to lose the final state
		},
	}
	finalCtx := workflow.WithActivityOptions(ctx, finalOpts)
	if err := workflow.ExecuteActivity(finalCtx,
		activities.RecordJobEnd, in.JobID, result).Get(finalCtx, nil); err != nil {
		return result, fmt.Errorf("record end: %w", err)
	}
	// Fire sync.completed event (outbox-backed in the activity).
	_ = workflow.ExecuteActivity(finalCtx, activities.EmitSyncCompleted, result).Get(finalCtx, nil)

	return result, nil
}

func decideStatus(steps []StepResult, cancelled bool) string {
	if cancelled {
		return "cancelled"
	}
	var succ, fail, skip int
	for _, s := range steps {
		switch s.Status {
		case "succeeded":
			succ++
		case "failed":
			fail++
		case "skipped":
			skip++
		}
	}
	switch {
	case fail == 0 && succ > 0:
		return "succeeded"
	case succ > 0 && fail > 0:
		return "partial"
	case succ == 0 && fail > 0:
		return "failed"
	default:
		return "succeeded" // all skipped counts as success (nothing to do)
	}
}

// ------------------------------------------------------------
// PushToUpstream — child workflow. Uses activities that hit
// adapter-pool with retries, circuit breakers, and heartbeats.
// ------------------------------------------------------------

type PushToUpstreamInput struct {
	JobID   string  `json:"job_id"`
	RepoID  string  `json:"repo_id"`
	OrgID   string  `json:"org_id"`
	Binding Binding `json:"binding"`
	RefName string  `json:"ref_name"`
	OldSHA  string  `json:"old_sha"`
	NewSHA  string  `json:"new_sha"`
	Origin  string  `json:"origin"`
}

func PushToUpstream(ctx workflow.Context, in PushToUpstreamInput) (StepResult, error) {
	start := workflow.Now(ctx)
	step := StepResult{
		UpstreamID: in.Binding.UpstreamID,
		Status:     "pending",
		StartedAt:  start,
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout:  2 * time.Minute,
		HeartbeatTimeout:     15 * time.Second,
		ScheduleToCloseTimeout: 10 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    60 * time.Second,
			MaximumAttempts:    6,
			NonRetryableErrorTypes: []string{"AdapterAuthFailed", "UpstreamCapabilityUnsupported"},
		},
	}
	actx := workflow.WithActivityOptions(ctx, ao)

	// Pre-flight: check circuit + rate limits for this provider.
	var ok bool
	if err := workflow.ExecuteActivity(actx, activities.PreflightAdapter, in.Binding).
		Get(actx, &ok); err != nil || !ok {
		step.Status = "failed"
		step.ErrorCode = "adapter.circuit_open"
		step.FinishedAt = workflow.Now(ctx)
		return step, nil                   // non-retry from parent; child will run again per parent retry
	}

	// Push.
	var pushResult activities.PushResult
	if err := workflow.ExecuteActivity(actx, activities.AdapterPushRef, in).Get(actx, &pushResult); err != nil {
		step.Status = "failed"
		step.ErrorCode = "adapter.push_failed"
		step.ErrorMessage = err.Error()
		step.FinishedAt = workflow.Now(ctx)
		return step, err
	}

	step.Status = "succeeded"
	step.RefsUpdated = pushResult.RefsUpdated
	step.BytesOut = pushResult.BytesOut
	step.FinishedAt = workflow.Now(ctx)
	return step, nil
}

// ------------------------------------------------------------
// Activities — one-line summaries; full impls elsewhere.
// ------------------------------------------------------------
// RecordJobStart(ctx, jobID)                  → updates sync.jobs row → running
// RecordJobEnd(ctx, jobID, result)            → persists steps + status
// EmitSyncCompleted(ctx, result)              → writes event to outbox
// PreflightAdapter(ctx, binding)              → checks circuit + rate limits
// AdapterPushRef(ctx, in)                     → calls adapter-pool; heartbeats during long pushes
// ------------------------------------------------------------
