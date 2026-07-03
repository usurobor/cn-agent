package cli

import (
	"context"

	issuesmap "github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map"
)

// IssuesMapCmd implements `cn issues map` — generate the interactive
// issue-board (Voronoi + Pivot views) from the repository's open issues.
//
// Per the dispatch boundary (INVARIANTS.md T-002, eng/go §2.18), this file
// owns only the dispatch wiring; all domain logic (fetch, taxonomy parse,
// effort computation, render, write) lives in the issuesmap package, which
// is Go-source co-located under src/packages/cnos.issues/commands/issues-map/
// per cnos#556 (mirroring the cnos#392 cdd-verify precedent). `cn issues
// map` remains a compiled-in kernel command (SourceKernel/TierKernel,
// registered in src/go/cmd/cn/main.go) — this is Go-source co-location,
// not package-command exec-dispatch (PACKAGE-SYSTEM.md §7).
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
