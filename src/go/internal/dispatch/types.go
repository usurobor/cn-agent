// Package dispatch implements role identity rotation for CDD triadic coordination.
//
// The dispatch command enables γ to spawn isolated α and β role sessions
// while preserving cognitive isolation and branch/artifact continuity.
package dispatch

import (
	"time"
)

// Args represents parsed command-line arguments for dispatch
type Args struct {
	Role     string // "alpha" or "beta" (α|β normalized to full names)
	Branch   string // e.g. "cycle/295"
	Issue    int    // issue number
	Project  string // project name
	Backend  string // "claude"|"stub"|"print"
	PromptFile string // optional prompt file path
}

// Result represents the outcome of a dispatch operation
type Result struct {
	AttemptID      string    `json:"attempt_id"`      // e.g. "dispatch-295-alpha-001"
	Role           string    `json:"role"`            // "alpha"|"beta"
	Issue          int       `json:"issue"`           // issue number
	Branch         string    `json:"branch"`          // e.g. "cycle/295"
	Backend        string    `json:"backend"`         // "claude"|"stub"|"print"
	Status         string    `json:"status"`          // "completed"|"failed"|"timeout"|"cancelled"|"backend_unavailable"|"partial"
	ExitCode       int       `json:"exit_code"`       // backend exit code
	StartedAt      time.Time `json:"started_at"`      // when dispatch began
	EndedAt        time.Time `json:"ended_at"`        // when dispatch completed
	DurationMs     int64     `json:"duration_ms"`     // duration in milliseconds
	PromptHash     string    `json:"prompt_hash"`     // sha256 of the generated prompt
	BranchHeadBefore string  `json:"branch_head_before"` // commit SHA before dispatch
	BranchHeadAfter  string  `json:"branch_head_after"`  // commit SHA after dispatch
	ArtifactRefs   []string  `json:"artifact_refs"`   // paths to artifacts created/modified
	FailureKind    *string   `json:"failure_kind"`    // failure category when status != "completed"
	Diagnostic     *string   `json:"diagnostic"`      // human-readable error message
	Retryable      bool      `json:"retryable"`       // whether this attempt can be retried
	ResumeHint     string    `json:"resume_hint"`     // "none"|"reinvoke_same_prompt"|"read_partial_artifact"|"operator_required"
	ResultRef      string    `json:"result_ref"`      // path or reference to detailed result
}

// Backend represents a dispatch backend (Claude CLI, stub, print)
type Backend interface {
	// Name returns the backend name
	Name() string

	// Dispatch executes the role session with the given prompt
	Dispatch(args *Args, prompt string) (*BackendResult, error)

	// IsAvailable checks if the backend can be used
	IsAvailable() error
}

// BackendResult represents the result from a backend execution
type BackendResult struct {
	ExitCode   int    // exit code from the backend
	Stdout     string // captured stdout
	Stderr     string // captured stderr
	Duration   time.Duration // execution duration
	FailureKind *string // failure category
	Diagnostic *string // error details
}