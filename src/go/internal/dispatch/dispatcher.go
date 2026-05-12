package dispatch

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Dispatcher orchestrates role identity dispatch operations
type Dispatcher struct {
	Backend  string    // backend name
	Project  string    // project name
	HubPath  string    // path to hub (.cn/ directory)
	Stdout   io.Writer // output stream
	Stderr   io.Writer // error stream
}

// Dispatch executes a role identity session
func (d *Dispatcher) Dispatch(ctx context.Context, args *Args) (*Result, error) {
	attemptID := generateAttemptID(args.Issue, args.Role)
	startTime := time.Now()

	// Create result with initial values
	result := &Result{
		AttemptID:  attemptID,
		Role:       args.Role,
		Issue:      args.Issue,
		Branch:     args.Branch,
		Backend:    args.Backend,
		StartedAt:  startTime,
		Status:     "failed", // default to failed, update on success
		Retryable:  true,
		ResumeHint: "reinvoke_same_prompt",
	}

	// Check worktree safety
	if err := d.checkWorktreeSafety(args); err != nil {
		result.EndedAt = time.Now()
		result.DurationMs = result.EndedAt.Sub(startTime).Milliseconds()
		result.ExitCode = 1
		failureKind := "worktree_unsafe"
		result.FailureKind = &failureKind
		diagnostic := err.Error()
		result.Diagnostic = &diagnostic

		fmt.Fprintf(d.Stderr, "✗ %s\n", err)
		return result, err
	}

	// Get branch HEAD before dispatch
	if branchHead, err := d.getBranchHead(args.Branch); err == nil {
		result.BranchHeadBefore = branchHead
	}

	// Construct prompt
	prompt, err := d.constructPrompt(args)
	if err != nil {
		result.EndedAt = time.Now()
		result.DurationMs = result.EndedAt.Sub(startTime).Milliseconds()
		result.ExitCode = 1
		failureKind := "prompt_construction_failed"
		result.FailureKind = &failureKind
		diagnostic := err.Error()
		result.Diagnostic = &diagnostic

		fmt.Fprintf(d.Stderr, "✗ Failed to construct prompt: %s\n", err)
		return result, err
	}

	// Hash the prompt
	hasher := sha256.New()
	hasher.Write([]byte(prompt))
	result.PromptHash = fmt.Sprintf("sha256:%x", hasher.Sum(nil))

	// Get the backend
	backend, err := d.getBackend(args.Backend)
	if err != nil {
		result.EndedAt = time.Now()
		result.DurationMs = result.EndedAt.Sub(startTime).Milliseconds()
		result.ExitCode = 1
		failureKind := "backend_unavailable"
		result.FailureKind = &failureKind
		result.Status = "backend_unavailable"
		result.Retryable = false
		diagnostic := err.Error()
		result.Diagnostic = &diagnostic

		fmt.Fprintf(d.Stderr, "✗ Backend unavailable: %s\n", err)
		return result, err
	}

	// Execute dispatch
	backendResult, err := backend.Dispatch(args, prompt)
	result.EndedAt = time.Now()
	result.DurationMs = result.EndedAt.Sub(startTime).Milliseconds()

	if backendResult != nil {
		result.ExitCode = backendResult.ExitCode
		result.FailureKind = backendResult.FailureKind
		result.Diagnostic = backendResult.Diagnostic
	}

	if err != nil {
		fmt.Fprintf(d.Stderr, "✗ Dispatch failed: %s\n", err)
		return result, err
	}

	// Get branch HEAD after dispatch
	if branchHead, err := d.getBranchHead(args.Branch); err == nil {
		result.BranchHeadAfter = branchHead
	}

	// Check for artifacts in .cdd/unreleased/{N}/
	if args.Issue > 0 {
		artifacts := d.findArtifacts(args.Issue)
		result.ArtifactRefs = artifacts
	}

	// Success case
	if result.ExitCode == 0 {
		result.Status = "completed"
		result.ResumeHint = "none"
		fmt.Fprintf(d.Stdout, "✓ %s dispatch completed\n", args.Role)
	}

	return result, nil
}

// checkWorktreeSafety validates the current worktree state
func (d *Dispatcher) checkWorktreeSafety(args *Args) error {
	// Check if we're on the correct branch
	currentBranch, err := d.getCurrentBranch()
	if err != nil {
		return fmt.Errorf("cannot determine current branch: %w", err)
	}

	if currentBranch != args.Branch {
		return fmt.Errorf("current branch %q != dispatch branch %q. Check out %s first", currentBranch, args.Branch, args.Branch)
	}

	// Check for uncommitted changes
	if hasUncommittedChanges, err := d.hasUncommittedChanges(); err != nil {
		return fmt.Errorf("cannot check worktree status: %w", err)
	} else if hasUncommittedChanges {
		return fmt.Errorf("uncommitted changes present. Commit or stash changes before dispatch")
	}

	// Verify branch is reachable from origin
	if err := d.verifyBranchReachable(args.Branch); err != nil {
		return fmt.Errorf("branch %q not reachable: %w", args.Branch, err)
	}

	return nil
}

// getCurrentBranch returns the current git branch name
func (d *Dispatcher) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = d.HubPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// hasUncommittedChanges checks if there are uncommitted changes in the worktree
func (d *Dispatcher) hasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = d.HubPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// verifyBranchReachable checks if the branch exists on origin
func (d *Dispatcher) verifyBranchReachable(branch string) error {
	// Fetch the branch
	fetchCmd := exec.Command("git", "fetch", "--quiet", "origin", branch)
	fetchCmd.Dir = d.HubPath
	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}

	// Verify the branch exists
	verifyCmd := exec.Command("git", "rev-parse", "--verify", "origin/"+branch)
	verifyCmd.Dir = d.HubPath
	if err := verifyCmd.Run(); err != nil {
		return fmt.Errorf("branch does not exist on origin: %w", err)
	}

	return nil
}

// getBranchHead returns the current HEAD SHA of the given branch
func (d *Dispatcher) getBranchHead(branch string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = d.HubPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// findArtifacts looks for artifacts in .cdd/unreleased/{N}/
func (d *Dispatcher) findArtifacts(issue int) []string {
	cycleDir := filepath.Join(d.HubPath, ".cdd", "unreleased", fmt.Sprintf("%d", issue))
	var artifacts []string

	filepath.Walk(cycleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// Convert to relative path from cycle directory
		relPath, err := filepath.Rel(cycleDir, path)
		if err == nil && relPath != "." {
			artifacts = append(artifacts, ".cdd/unreleased/"+fmt.Sprintf("%d", issue)+"/"+relPath)
		}
		return nil
	})

	return artifacts
}

// generateAttemptID creates a unique attempt ID
func generateAttemptID(issue int, role string) string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("dispatch-%d-%s-%s", issue, role, timestamp)
}