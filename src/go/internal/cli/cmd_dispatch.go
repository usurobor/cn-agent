package cli

import (
	"context"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/dispatch"
)

// DispatchCmd implements the "dispatch" command for role identity rotation.
//
// Usage: cn dispatch --role α|β --branch cycle/N [--issue N] [--project NAME] [--backend claude|stub|print]
//
// This command enables CDD triadic coordination by dispatching isolated
// role identity sessions while preserving branch/artifact continuity.
type DispatchCmd struct{}

func (c *DispatchCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "dispatch",
		Summary:  "Role identity rotation primitive for CDD triadic coordination",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: true, // requires .cn/ for project detection and .cdd/ artifact coordination
	}
}

// Help returns detailed help text for the dispatch command
func (c *DispatchCmd) Help() string {
	return `cn dispatch - Role identity rotation primitive for CDD

USAGE:
  cn dispatch --role α|β --branch cycle/N [OPTIONS]

DESCRIPTION:
  Dispatches isolated role identity sessions for CDD triadic coordination.
  Enables γ to spawn α and β roles while preserving cognitive isolation
  and branch/artifact continuity.

REQUIRED FLAGS:
  --role α|β        Role to dispatch (alpha or beta)
  --branch cycle/N  Branch to operate on (must be cycle/N format)

OPTIONAL FLAGS:
  --issue N         Issue number (auto-detected from branch if not provided)
  --project NAME    Project name (auto-detected from hub if not provided)
  --backend TYPE    Backend to use: claude|stub|print (default: claude)

BACKENDS:
  claude   Claude CLI backend (requires claude binary and authentication)
           Uses subscription credentials or API key. Stream-json output
           for observability.
  stub     Test backend that records prompts without execution
  print    Print backend that outputs constructed prompt to stdout

ENVIRONMENT:
  CN_DISPATCH_BACKEND  Default backend (overridden by --backend flag)

AUTHENTICATION:
  Claude backend authenticates via subscription (Pro/Max) or API key.
  When using subscription credentials, usage draws from plan allocation.
  If ANTHROPIC_API_KEY is set, API billing applies.

EXIT CODES:
  0  Success
  1  Backend execution failed
  2  Invalid arguments or worktree unsafe

EXAMPLES:
  cn dispatch --role α --branch cycle/295
  cn dispatch --role β --branch cycle/295 --backend stub
  cn dispatch --role α --branch cycle/295 --backend print

SAFETY:
  - Verifies current branch matches --branch argument
  - Checks for uncommitted changes
  - Validates branch exists on origin
  - Records attempt descriptor to .cdd/unreleased/{N}/dispatch/

For more information, see gamma/SKILL.md §2.5 in the CDD documentation.`
}

func (c *DispatchCmd) Run(ctx context.Context, inv Invocation) error {
	args, err := dispatch.ParseArgs(inv.Args)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ %s\n", err)
		return err
	}

	backend := dispatch.ResolveBackend(args.Backend)

	if !dispatch.IsValidBackend(backend) {
		err := fmt.Errorf("invalid backend %q (use claude|stub|print)", backend)
		fmt.Fprintf(inv.Stderr, "✗ %s\n", err)
		return err
	}

	project := args.Project
	if project == "" {
		project = dispatch.DetectProject(inv.HubPath)
	}

	dispatcher := &dispatch.Dispatcher{
		Backend: backend,
		Project: project,
		HubPath: inv.HubPath,
		Stdout:  inv.Stdout,
		Stderr:  inv.Stderr,
	}

	result, err := dispatcher.Dispatch(ctx, args)
	if err != nil {
		return err
	}

	if err := dispatch.RecordAttempt(inv.HubPath, args.Issue, result); err != nil {
		fmt.Fprintf(inv.Stderr, "⚠ Failed to record attempt: %v\n", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("dispatch failed with exit code %d", result.ExitCode)
	}

	return nil
}
