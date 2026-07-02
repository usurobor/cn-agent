package cli

import (
	"context"

	"github.com/usurobor/cnos/src/go/internal/issuesmap"
)

// IssuesMapCmd implements `cn issues map` — generate the interactive
// issue-board (Voronoi + Pivot views) from the repository's open issues.
//
// Per the dispatch boundary (INVARIANTS.md T-002, eng/go §2.18), this file
// owns only the dispatch wiring; all domain logic (fetch, taxonomy parse,
// effort computation, render, write) lives in the issuesmap package.
type IssuesMapCmd struct{}

func (c *IssuesMapCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "issues-map",
		Summary: "Generate the interactive issue-board (Voronoi + Pivot) from open issues",
		Source:  SourceKernel,
		Tier:    TierKernel,
		// NeedsHub is false: the command targets a GitHub repo (via --repo /
		// $GITHUB_REPOSITORY) and an output dir, not a cnos hub.
	}
}

func (c *IssuesMapCmd) Run(ctx context.Context, inv Invocation) error {
	return issuesmap.Run(ctx, inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)
}
