package cli

import (
	"context"

	issuesdispatch "github.com/usurobor/cnos/packages/cnos.issues/commands/issues-dispatch"
)

// IssuesDispatchCmd implements `cn issues dispatch` — the cnos#640
// one-verb dispatch-authorization primitive. It flips an issue's
// status:ready label to status:todo (the operator's dispatch-
// authorization event, unchanged in substance — see
// dispatch-protocol/SKILL.md §1.2) and, in the same caller-visible
// operation, strips the legacy body-hold prose the #614/#633 recurrence
// showed can drift out of sync with the label. This command is
// operator/human-invoked only: no dispatch wake or scan/reconciler process
// calls it on its own initiative.
//
// Per the dispatch boundary (INVARIANTS.md T-002, eng/go §2.18), this file
// owns only the dispatch wiring; all domain logic (fetch, body
// reconciliation, label mutation) lives in the issuesdispatch package,
// which is Go-source co-located under
// src/packages/cnos.issues/commands/issues-dispatch/ per cnos#640 —
// mirroring the cnos#568 issues-fsm / cnos#556 issues-map precedent
// exactly (a sibling of both, same registry, same noun-verb routing, same
// kernel binary — not a standalone binary).
type IssuesDispatchCmd struct{}

func (c *IssuesDispatchCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "issues-dispatch",
		Summary: "Authorize one design-first issue for dispatch: flips status:ready -> status:todo and cleans the legacy body-hold phrase in one operation (cnos#640)",
		Source:  SourceKernel,
		Tier:    TierKernel,
		// NeedsHub is false: the command targets a GitHub repo (via
		// --repo / $GITHUB_REPOSITORY), not a cnos hub.
	}
}

func (c *IssuesDispatchCmd) Run(ctx context.Context, inv Invocation) error {
	return issuesdispatch.Run(ctx, inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)
}
