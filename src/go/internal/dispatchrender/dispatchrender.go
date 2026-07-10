// Package dispatchrender invokes the vendored cn-install-wake renderer
// (cnos#609) to render (or re-render) a cds-dispatch workflow file.
//
// Extracted out of internal/repoinstall (cnos#656) so the exact same
// argument-building logic serves two callers without drifting apart:
// internal/repoinstall's `cn repo install --dispatch cds` (the original
// render) and internal/repostatus's `cn repo status` (a fresh re-render
// into a temp file, compared against the live workflow + the
// .cn/repo.state.json ledger sha256 to classify drift as matches-ledger /
// user-edit / renderer-moved — design doc A2). If these two call sites
// built renderer args independently, a future flag added to one but not
// the other would silently make A2's drift classification produce false
// "user_edit" verdicts for every installed repo — this package is the
// single source of truth for "how do we invoke the renderer" so that
// class of bug is structurally impossible.
package dispatchrender

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/usurobor/cnos/src/go/internal/pkg"
)

// Options carries the renderer invocation inputs — the same fields a
// .cn/repo.state.json workflow managed_files entry records (P2's render
// contract: schemas/repo_state.cue's #ManagedFile), enough to reproduce
// the render byte-for-byte given the same vendored renderer version.
type Options struct {
	// RepoRoot is the repository root; the renderer is resolved from
	// RepoRoot's vendored RendererPackage and runs with RepoRoot as its
	// working directory.
	RepoRoot string
	// RendererPackage is the vendored package providing
	// commands/install-wake/cn-install-wake — "cnos.core" today.
	RendererPackage   string
	Tier              string // "agent" | "engine"
	Agent             string
	WorkflowPatSecret string
	BotName           string
	BotID             string
	// OutPath is where the renderer writes the workflow file — absolute,
	// or resolved relative to RepoRoot by the caller.
	OutPath string

	// Stdout/Stderr receive the renderer subprocess's own output.
	// Default io.Discard when nil — callers that want the renderer's
	// interactive install-time output (internal/repoinstall) pass the
	// real stdout/stderr; callers doing a silent background comparison
	// (internal/repostatus's fresh-render drift check) leave these nil.
	Stdout io.Writer
	Stderr io.Writer
}

// RendererPath returns the path to the vendored cn-install-wake script
// Render would invoke for opts.RepoRoot/opts.RendererPackage.
func RendererPath(repoRoot, rendererPackage string) string {
	return filepath.Join(pkg.VendorPath(repoRoot, rendererPackage), "commands", "install-wake", "cn-install-wake")
}

// Render invokes the vendored renderer to write opts.OutPath. Returns an
// error naming the renderer path if it is not present (the caller's
// package must be installed first) or if the renderer process itself
// fails (renderer stderr is included in the returned error).
func Render(ctx context.Context, opts Options) error {
	rendererPath := RendererPath(opts.RepoRoot, opts.RendererPackage)
	if _, err := os.Stat(rendererPath); err != nil {
		return fmt.Errorf("dispatch renderer not found at %s (the %s package must be installed for --dispatch cds): %w", rendererPath, opts.RendererPackage, err)
	}

	args := []string{"cds-dispatch", "--out", opts.OutPath, "--agent", opts.Agent}
	if opts.Tier == "engine" {
		// Engine tier: the only renderer arg beyond the output path and the
		// concurrency-group agent name is --tier engine. The identity flags
		// (workflow-pat-secret / bot-name / bot-id) are deliberately not
		// forwarded — the engine tier binds nothing to them.
		args = append(args, "--tier", "engine")
	} else {
		if opts.WorkflowPatSecret != "" {
			args = append(args, "--workflow-pat-secret", opts.WorkflowPatSecret)
		}
		if opts.BotName != "" {
			args = append(args, "--bot-name", opts.BotName)
		}
		if opts.BotID != "" {
			args = append(args, "--bot-id", opts.BotID)
		}
	}

	stdout, stderr := opts.Stdout, opts.Stderr
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	cmd := exec.CommandContext(ctx, rendererPath, args...)
	cmd.Dir = opts.RepoRoot
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("render dispatch workflow: %w", err)
	}
	return nil
}
