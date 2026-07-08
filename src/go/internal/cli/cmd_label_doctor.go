package cli

import (
	"context"

	labeldoctor "github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor"
)

// LabelDoctorCmd implements `cn label-doctor` (cnos#493) — audits (and,
// unless --dry-run, repairs) the target repo's GitHub labels against
// src/packages/cnos.core/labels.json. Name is registered as the single
// string "label-doctor", matching cell.go's existing runtime error text
// ("Run label-doctor before retrying") and cnos.core's label-doctrine
// naming exactly.
//
// Invocation shape note: because Name contains a hyphen,
// cli/dispatch.go's ResolveCommand routes it through the noun-verb path
// ("cn label doctor", two argv tokens) rather than a flat single
// hyphenated token — this is the same, already-documented
// ("Hyphenated flat forms ... are NOT supported", dispatch.go)
// mechanism every other hyphenated kernel command in this registry goes
// through (issues-fsm -> "issues fsm", cdd-verify -> "cdd verify",
// repo-install -> "repo install"), not a choice made here. See
// .cdd/unreleased/493/self-coherence.md §ACs (AC4) for the full
// write-up.
//
// Per the dispatch boundary (INVARIANTS.md T-002, eng/go §2.18), this
// file owns only the dispatch wiring; all domain logic (manifest
// loading, GitHub REST primitives, audit/repair, flag parsing) lives in
// the labeldoctor package, imported in-process from this go.work-linked
// module — mirroring cmd_issues_fsm.go's issuesfsm import exactly.
type LabelDoctorCmd struct{}

func (c *LabelDoctorCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "label-doctor",
		Summary: "Audit and repair the repo's GitHub labels against cnos.core's labels.json (cnos#493)",
		Source:  SourceKernel,
		Tier:    TierKernel,
		// NeedsHub is false: this must work in a freshly-cloned repo
		// with no .cn/ yet, same reasoning as RepoInstallCmd.
	}
}

func (c *LabelDoctorCmd) Run(ctx context.Context, inv Invocation) error {
	return labeldoctor.Run(ctx, inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)
}
