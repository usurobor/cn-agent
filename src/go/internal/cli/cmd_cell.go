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

DESCRIPTION:
  Verifies the cycle/{N} branch exists on origin and the .cdd/unreleased/{N}/
  artifact directory is present, then appends an R[N+1] section header to
  self-coherence.md so the next dispatched round proceeds as a continuation.

  Existing artifacts are preserved unchanged. The cycle branch is not
  recreated. No new cycle directory is created.

  After resume, δ routes α R[N+1] → β R[N+1] → γ closeout amendment per the
  resumed-from-changes shape (delta/SKILL.md §9.10).

REQUIRED FLAGS:
  --issue N   GitHub issue number (positive integer)

EXAMPLE:
  cn cell resume --issue 500

EXIT CODES:
  0  Success (cycle re-armed; R[N+1] header appended)
  1  Error (branch not found, artifact directory missing, write failure)

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
