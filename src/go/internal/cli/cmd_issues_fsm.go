package cli

import (
	"context"

	issuesfsm "github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm"
)

// IssuesFsmCmd implements `cn issues fsm` — the cnos#568 Phase 1 read-only
// issue-state reconciler. Phase 1 supports exactly one sub-verb,
// "evaluate" (parsed by the issuesfsm package itself from the first
// remaining arg); there is no "apply" sub-verb — label-mutation authority
// is Phase 2 (cnos#569), out of scope here.
//
// Per the dispatch boundary (INVARIANTS.md T-002, eng/go §2.18), this file
// owns only the dispatch wiring; all domain logic (fact-snapshot assembly,
// transition-table evaluation, decision rendering) lives in the issuesfsm
// package, which is Go-source co-located under
// src/packages/cnos.issues/commands/issues-fsm/ per cnos#568 (mirroring the
// cnos#556 issues-map / cnos#392 cdd-verify precedent). `cn issues fsm`
// remains a compiled-in kernel command (SourceKernel/TierKernel,
// registered in src/go/cmd/cn/main.go) — this is Go-source co-location,
// not package-command exec-dispatch (PACKAGE-SYSTEM.md §7).
type IssuesFsmCmd struct{}

func (c *IssuesFsmCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "issues-fsm",
		Summary: "Evaluate CDS issue state from observed facts (read-only; Phase 1, cnos#568)",
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
