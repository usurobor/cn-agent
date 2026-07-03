package cli

import (
	"context"

	issuesfsm "github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm"
)

// IssuesFsmCmd implements `cn issues fsm` — the cnos#568 Phase 1 issue-
// state reconciler, extended by cnos#569 Phase 2 with a guard-gated
// mutation path. There is still exactly one sub-verb, "evaluate" (parsed
// by the issuesfsm package itself from the first remaining arg); label-
// mutation authority is a --apply FLAG on "evaluate" (cnos#569), not a
// separate "apply" sub-verb.
//
// Per the dispatch boundary (INVARIANTS.md T-002, eng/go §2.18), this file
// owns only the dispatch wiring; all domain logic (fact-snapshot assembly,
// transition-table evaluation, decision rendering, and the cnos#569
// label-write path) lives in the issuesfsm package, which is Go-source
// co-located under src/packages/cnos.issues/commands/issues-fsm/ per
// cnos#568 (mirroring the cnos#556 issues-map / cnos#392 cdd-verify
// precedent). `cn issues fsm` remains a compiled-in kernel command
// (SourceKernel/TierKernel, registered in src/go/cmd/cn/main.go) — this is
// Go-source co-location, not package-command exec-dispatch
// (PACKAGE-SYSTEM.md §7).
type IssuesFsmCmd struct{}

func (c *IssuesFsmCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "issues-fsm",
		Summary: "Evaluate (and, with --apply, mutate) CDS issue state from observed facts (cnos#568 Phase 1 / cnos#569 Phase 2)",
		Source:  SourceKernel,
		Tier:    TierKernel,
		// NeedsHub is false: --fixture mode is fully offline, and the live
		// path targets a GitHub repo (via --repo / $GITHUB_REPOSITORY) plus
		// the current git checkout, not a cnos hub.
	}
}

func (c *IssuesFsmCmd) Run(ctx context.Context, inv Invocation) error {
	return issuesfsm.Run(ctx, inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)
}
