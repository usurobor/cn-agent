package issuesfsm

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Decision is the evaluator's structured output for one FactSnapshot: the
// AC2-required "current state, observed facts, enabled transition, blocked
// reason, proposed action" block. It is a pure projection of a Table +
// FactSnapshot through Evaluate — read-only, and safe to compute repeatedly
// (AC7 idempotence: the same inputs always produce the same Decision).
type Decision struct {
	Issue        int
	CurrentState string
	Facts        FactSnapshot

	// Outcome is "valid" (no transition needed), "blocked" (a transition is
	// disallowed pending missing evidence), or "proposed" (a transition is
	// enabled — proposed only, never executed; see AC8).
	Outcome string
	// Action is the machine-readable action id from the matched rule.
	Action string
	// TargetState is the proposed next status label suffix, or "" if none.
	TargetState string
	// EnabledTransition is "<current> -> <target>" when TargetState != "",
	// else "".
	EnabledTransition string
	// RepairPass marks a proposed transition as a repair re-entry (cnos#516).
	RepairPass bool
	// Reason is the human-readable rationale for the matched rule.
	Reason string
	// BlockedReason is set (equal to Reason) when Outcome == "blocked".
	BlockedReason string
	// MissingEvidence lists the guard names that were false and caused a
	// "blocked" outcome (AC3's "missing-evidence list").
	MissingEvidence []string
}

// Render writes the AC2-required decision block to w: current state,
// observed facts, enabled transition, blocked reason, and proposed action.
// The format is stable, deterministic (fixed field order, sorted labels),
// and does not mutate any state — it is purely a projection of d.
func (d Decision) Render(w io.Writer) {
	fmt.Fprintf(w, "cn issues fsm evaluate — issue #%d\n\n", d.Issue)

	state := d.CurrentState
	if state == "" {
		state = "(none — no status:* label present)"
	}
	fmt.Fprintf(w, "Current state: %s\n\n", state)

	fmt.Fprintln(w, "Observed facts:")
	labels := append([]string(nil), d.Facts.Labels...)
	sort.Strings(labels)
	fmt.Fprintf(w, "  labels: %s\n", strings.Join(labels, ", "))
	fmt.Fprintf(w, "  run_state: %s\n", orNone(d.Facts.RunState))
	if d.Facts.RunConclusion != "" {
		fmt.Fprintf(w, "  run_conclusion: %s\n", d.Facts.RunConclusion)
	}
	fmt.Fprintf(w, "  branch_exists: %v\n", d.Facts.BranchExists)
	fmt.Fprintf(w, "  commits_beyond_base: %d\n", d.Facts.CommitsBeyondBase)
	fmt.Fprintf(w, "  pr_exists: %v\n", d.Facts.PRExists)
	fmt.Fprintf(w, "  pr_commit_count: %d\n", d.Facts.PRCommitCount)
	fmt.Fprintf(w, "  cdd_artifacts: %s\n", orNoneList(d.Facts.CDDArtifacts))
	fmt.Fprintf(w, "  review_request_present: %v\n", d.Facts.ReviewRequestPresent)
	fmt.Fprintf(w, "  repair_contract_present: %v\n", d.Facts.RepairContractPresent)
	fmt.Fprintf(w, "  checks_state: %s\n", orNone(d.Facts.ChecksState))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Decision:")
	fmt.Fprintf(w, "  outcome: %s\n", orNone(d.Outcome))
	fmt.Fprintf(w, "  enabled_transition: %s\n", orNoneStr(d.EnabledTransition, "(none)"))
	fmt.Fprintf(w, "  blocked_reason: %s\n", orNoneStr(d.BlockedReason, "(none)"))
	if len(d.MissingEvidence) > 0 {
		fmt.Fprintf(w, "  missing_evidence: %s\n", strings.Join(d.MissingEvidence, ", "))
	}
	fmt.Fprintf(w, "  proposed_action: %s\n", orNoneStr(d.Action, "none"))
	if d.RepairPass {
		fmt.Fprintln(w, "  repair_pass: true")
	}
	fmt.Fprintf(w, "  reason: %s\n", orNoneStr(d.Reason, "(none)"))

	// AC8: this command never mutates a label. Restate it in the output so
	// the boundary is operator-visible, not just an internal invariant.
	fmt.Fprintln(w)
	fmt.Fprintln(w, "(read-only: no label was written; this is a proposal, not an action)")
}

func orNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

func orNoneStr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func orNoneList(ss []string) string {
	if len(ss) == 0 {
		return "(none)"
	}
	sorted := append([]string(nil), ss...)
	sort.Strings(sorted)
	return strings.Join(sorted, ", ")
}
