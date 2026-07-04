package cli

import (
	"context"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/cell"
)

// CellReturnCmd implements "cn cell return".
//
// Usage: cn cell return --issue N --verdict V --review path
//
// Reads a cn.operator-review.v1 artifact and, on iterate/reject verdicts,
// transitions the GitHub issue from status:review to status:changes.
// On converge, reports the converge path; no label action taken.
type CellReturnCmd struct{}

func (c *CellReturnCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "cell-return",
		Summary:  "Deliver operator verdict; transition status:review → status:changes on iterate/reject",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: false, // operates via gh CLI; no hub required
	}
}

func (c *CellReturnCmd) Help() string {
	return `cn cell return - Deliver operator verdict for a cell at status:review

USAGE:
  cn cell return --issue N --verdict V --review path

DESCRIPTION:
  Reads a cn.operator-review.v1 artifact and routes the operator's verdict
  back into the existing cell lifecycle.

  On iterate or reject: removes the status:review label and adds status:changes
  on the GitHub issue. Use 'cn cell resume' afterward to re-arm the cycle.

  On converge: reports the converge path; no label action is taken.
  The operator merges the PR to complete the cell lifecycle.

  On clarify: reports that clarification is needed; no label action is taken.

REQUIRED FLAGS:
  --issue N     GitHub issue number (positive integer)
  --verdict V   One of: converge | iterate | reject | clarify
  --review path Path to a cn.operator-review.v1 YAML-frontmatter artifact

EXAMPLES:
  cn cell return --issue 500 --verdict iterate \
    --review .cdd/unreleased/500/operator-review.md

  cn cell return --issue 497 --verdict converge \
    --review .cdd/unreleased/497/operator-review.md

SCHEMA:
  The artifact at --review must begin with YAML frontmatter containing:
    schema: cn.operator-review.v1
  See src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md.

EXIT CODES:
  0  Success (transition applied or converge/clarify path reported)
  1  Error (artifact missing, schema mismatch, gh CLI failure)

DEPENDENCIES:
  gh CLI must be installed and authenticated (existing runtime dependency).`
}

func (c *CellReturnCmd) Run(ctx context.Context, inv Invocation) error {
	args, err := cell.ParseReturnArgs(inv.Args)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn cell return: %s\n\n", err)
		fmt.Fprintf(inv.Stderr, "Run 'cn help cell' to see usage.\n")
		return err
	}

	returner := &cell.Returner{
		Stdout: inv.Stdout,
		Stderr: inv.Stderr,
	}

	if err := returner.Return(ctx, args); err != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn cell return: %s\n", err)
		return err
	}
	return nil
}

// CellResumeCmd implements "cn cell resume".
//
// Usage: cn cell resume --issue N
//
// Re-arms an existing cycle on cycle/{N}: verifies the branch and artifact
// directory exist, and appends an R[N+1] section header to self-coherence.md
// so the next dispatched round knows it is a continuation.
type CellResumeCmd struct{}

func (c *CellResumeCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "cell-resume",
		Summary:  "Re-arm an existing cycle after status:changes (preserves branch and artifacts)",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: true, // needs repo root to locate .cdd/unreleased/{N}/
	}
}

func (c *CellResumeCmd) Help() string {
	return `cn cell resume - Re-arm an existing cycle after it returns from status:changes

USAGE:
  cn cell resume --issue N

PRECONDITION:
  The current local git branch MUST be cycle/{N}. cn cell resume is a
  local-only helper (Option B v0; see cycle/500 R1 self-coherence §R1
  for the design choice). The command refuses with review_resume_wrong_branch
  if the current branch is anything other than cycle/{N}.

  Run 'git checkout cycle/{N}' (or 'git fetch origin && git checkout -B
  cycle/{N} origin/cycle/{N}') before invoking.

DESCRIPTION:
  Verifies the cycle/{N} branch exists on origin and the .cdd/unreleased/{N}/
  artifact directory is present, then appends an R[N+1] section header to
  self-coherence.md so the next dispatched round proceeds as a continuation.

  Existing artifacts are preserved unchanged. The cycle branch is not
  recreated. No new cycle directory is created. The caller is responsible
  for the commit and push of the §R[N+1] marker (Option A — fetch + auto-
  commit + push — is deferred to a future cycle).

  After resume, δ routes α R[N+1] → β R[N+1] → γ closeout amendment per the
  resumed-from-changes shape (delta/SKILL.md §9.10).

REQUIRED FLAGS:
  --issue N   GitHub issue number (positive integer)

EXAMPLE:
  git checkout cycle/500
  cn cell resume --issue 500
  git add .cdd/unreleased/500/self-coherence.md
  git commit -m "alpha-500: open §R2"
  git push origin cycle/500

EXIT CODES:
  0  Success (cycle re-armed; R[N+1] header appended)
  1  Error (wrong branch, branch not found, artifact directory missing, write failure)

INVARIANTS:
  - cycle/{N} branch on origin: preserved (not recreated)
  - .cdd/unreleased/{N}/ directory: preserved (no new directory created)
  - self-coherence.md: §R[N+1] header appended (NOT replaced)
  - All existing artifacts in .cdd/unreleased/{N}/: unchanged`
}

func (c *CellResumeCmd) Run(ctx context.Context, inv Invocation) error {
	args, err := cell.ParseResumeArgs(inv.Args)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn cell resume: %s\n\n", err)
		fmt.Fprintf(inv.Stderr, "Run 'cn help cell' to see usage.\n")
		return err
	}

	resumer := &cell.Resumer{
		RepoRoot: inv.HubPath,
		Stdout:   inv.Stdout,
		Stderr:   inv.Stderr,
	}

	if err := resumer.Resume(ctx, args); err != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn cell resume: %s\n", err)
		return err
	}
	return nil
}

// CellFinalizeCmd implements "cn cell finalize".
//
// Usage: cn cell finalize [--issue N] [--base-sha SHA]
//
// Mechanical checkpoint + PR-open/update finalizer (cnos#591). Detects
// "matter" (uncommitted changes, commits beyond base, CDD artifacts, or an
// existing cycle/{N} branch) and, when present, commits/pushes and ensures
// a draft PR exists — idempotently. Never writes a status label and never
// calls `cn issues fsm evaluate`; the FSM remains the sole label-writer.
type CellFinalizeCmd struct{}

func (c *CellFinalizeCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "cell-finalize",
		Summary:  "Mechanical checkpoint + idempotent draft-PR open/update when a cell has matter",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: true, // needs repo root to locate .cdd/unreleased/{N}/ and run git
	}
}

func (c *CellFinalizeCmd) Help() string {
	return `cn cell finalize - Mechanical checkpoint + PR-open/update finalizer

USAGE:
  cn cell finalize [--issue N] [--base-sha SHA]

DESCRIPTION:
  Runtime/finalizer responsibility (cnos#591): if a run ends with matter,
  there must be a branch/commit/PR, not a cell stranded at status:in-progress
  with no PR. cn cell finalize checkpoints matter (commit + push) and
  opens (or idempotently reuses) a draft pull request. It does not decide
  cognition and does not own status labels — the FSM ('cn issues fsm
  evaluate --apply') remains the only mechanism that applies a lifecycle
  status label; this command has no label-mutation code path at all.

  Matter is detected as the OR of four independent signals:
    1. uncommitted file changes in the working tree
    2. commits on cycle/{N} beyond its base (main)
    3. CDD artifacts present under .cdd/unreleased/{N}/
    4. an existing cycle/{N} branch

  If none of the four signals hold, this is a clean no-op: exit 0, no
  branch created, no commit, no gh pr create call.

  When matter exists: ensures cycle/{N} exists, commits any uncommitted
  changes, pushes, then ensures a draft PR exists for cycle/{N} — creating
  one (base: main; body references Refs #{N}) if none is open, or leaving
  an already-open PR as a clean no-op (idempotent; never a duplicate).

  On checkpoint or PR-open failure, this command never claims success: it
  writes .cdd/unreleased/{N}/FINALIZE-STOP.md naming the failure and
  returns a precise error (cell_finalize_checkpoint_failed /
  cell_finalize_pr_open_failed / cell_finalize_pr_list_failed).

FLAGS (both optional):
  --issue N      GitHub issue number. If omitted, self-detected from the
                 current branch (cycle/{N}) or, failing that, by scanning
                 .cdd/unreleased/{N}/ for the issue with matter.
  --base-sha SHA Baseline SHA (e.g. $CN_WAKE_BASE_SHA from a dispatch-shaped
                 wake's write-fence baseline step) used only to widen the
                 self-detection scan to newly-committed CDD paths since
                 that baseline when the current branch is not cycle/{N}.

EXAMPLES:
  cn cell finalize
  cn cell finalize --base-sha "$CN_WAKE_BASE_SHA"
  cn cell finalize --issue 591

EXIT CODES:
  0  Success (matter checkpointed + PR ensured, or clean no-matter no-op)
  1  Error (checkpoint failed, PR list/create failed — see FINALIZE-STOP.md)

DEPENDENCIES:
  git and gh CLI must be installed and authenticated (existing runtime
  dependencies of this package).`
}

func (c *CellFinalizeCmd) Run(ctx context.Context, inv Invocation) error {
	args, err := cell.ParseFinalizeArgs(inv.Args)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn cell finalize: %s\n\n", err)
		fmt.Fprintf(inv.Stderr, "Run 'cn help cell' to see usage.\n")
		return err
	}

	finalizer := &cell.Finalizer{
		RepoRoot: inv.HubPath,
		Stdout:   inv.Stdout,
		Stderr:   inv.Stderr,
	}

	if err := finalizer.Finalize(ctx, args); err != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn cell finalize: %s\n", err)
		return err
	}
	return nil
}
