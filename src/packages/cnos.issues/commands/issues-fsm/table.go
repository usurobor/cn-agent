package issuesfsm

import (
	"encoding/json"
	"fmt"
	"os"
)

// Rule is one entry in a StateTransition's rules list. A rule matches a
// FactSnapshot when every AllTrue guard evaluates true, every AllFalse guard
// evaluates false, and (if AnyTrue is non-empty) at least one AnyTrue guard
// evaluates true. Rules within a StateTransition are evaluated in file
// order; the first matching rule wins.
//
// Rule is pure data (json-decoded from the package-owned transition table —
// AC1); it carries no CDS-specific behavior itself.
type Rule struct {
	AllTrue  []string `json:"all_true,omitempty"`
	AnyTrue  []string `json:"any_true,omitempty"`
	AllFalse []string `json:"all_false,omitempty"`

	// Outcome classifies the rule's result: "valid" (state is fine as-is),
	// "blocked" (a transition is disallowed pending missing evidence), or
	// "proposed" (a transition is enabled and proposed — never executed;
	// Phase 1 is read-only, see AC8).
	Outcome string `json:"outcome"`

	// Action is a short machine-readable action id, e.g. "none", "block",
	// "propose_status_todo", "propose_delta_recovery",
	// "propose_repair_dispatch".
	Action string `json:"action"`

	// TargetState is the proposed next CDS status (label suffix, no
	// "status:" prefix), or "" if the rule proposes no direct status change
	// (e.g. delta-recovery, which proposes opening a PR, not a label move).
	TargetState string `json:"target_state,omitempty"`

	// RepairPass marks a proposal as a repair re-entry (cnos#516) rather
	// than a first pass.
	RepairPass bool `json:"repair_pass,omitempty"`

	// Reason is the human-readable explanation printed for this rule's
	// outcome (used as both the "reason" and, when Outcome == "blocked",
	// the blocked_reason).
	Reason string `json:"reason"`

	// EvidenceGuards names the guards whose false subset is reported as the
	// missing-evidence list when this rule's Outcome is "blocked" (AC3).
	EvidenceGuards []string `json:"evidence_guards,omitempty"`
}

// StateTransition holds every rule for one CDS state, evaluated when Trigger
// fires (Phase 1 has exactly one trigger: "evaluate").
type StateTransition struct {
	State   string `json:"state"`
	Trigger string `json:"trigger"`
	Rules   []Rule `json:"rules"`
}

// Table is the full declarative transition table: the single source of
// truth for CDS issue-status transitions (AC1). It is loaded from JSON data
// (src/packages/cnos.cds/skills/cds/fsm/transitions.json) — never
// constructed as Go literals in production code (AC1's negative case).
type Table struct {
	Doc         []string          `json:"_doc,omitempty"`
	States      []string          `json:"states"`
	Guards      map[string]string `json:"guards,omitempty"`
	Transitions []StateTransition `json:"transitions"`
}

// LoadTable reads and parses a transition table from path.
func LoadTable(path string) (*Table, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read transition table: %w", err)
	}
	var t Table
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parse transition table JSON: %w", err)
	}
	return &t, nil
}

// guardFuncs is the generic, fixed predicate vocabulary the evaluator
// understands. This is the only place FactSnapshot fields are interpreted as
// booleans; which guards apply to which CDS state — and what to do when they
// pass or fail — is entirely table-driven (AC1). This map is never switched
// on a CDS-specific *state* name (e.g. "review", "changes", "in-progress");
// it only names generic, reusable evidence predicates over the fact model.
var guardFuncs = map[string]func(FactSnapshot) bool{
	"run_active":              func(f FactSnapshot) bool { return f.RunState == "queued" || f.RunState == "in_progress" },
	"branch_exists":           func(f FactSnapshot) bool { return f.BranchExists },
	"branch_has_commits":      func(f FactSnapshot) bool { return f.CommitsBeyondBase > 0 },
	"pr_exists":               func(f FactSnapshot) bool { return f.PRExists },
	"pr_has_commits":          func(f FactSnapshot) bool { return f.PRCommitCount > 0 },
	"review_request_present":  func(f FactSnapshot) bool { return f.ReviewRequestPresent },
	"repair_contract_present": func(f FactSnapshot) bool { return f.RepairContractPresent },
	"cdd_artifacts_present":   func(f FactSnapshot) bool { return len(f.CDDArtifacts) > 0 },
	"checks_passing":          func(f FactSnapshot) bool { return f.ChecksState == "success" },
	// cnos#575: claim / hard-block / release-back-to-queue guards, added
	// following the review_request_present/REVIEW-REQUEST.yml precedent
	// exactly (marker-file presence, not comment-text parsing).
	"claim_request_present":   func(f FactSnapshot) bool { return f.ClaimRequestPresent },
	"block_request_present":   func(f FactSnapshot) bool { return f.BlockRequestPresent },
	"release_request_present": func(f FactSnapshot) bool { return f.ReleaseRequestPresent },
}

// evalGuard evaluates a single named guard against f. Returns an error if
// the table references a guard name the engine doesn't recognize — a table
// authoring error, not a fact-model gap.
func evalGuard(name string, f FactSnapshot) (bool, error) {
	fn, ok := guardFuncs[name]
	if !ok {
		return false, fmt.Errorf("transition table references unknown guard %q", name)
	}
	return fn(f), nil
}

// ruleMatches reports whether r's conditions all hold against f.
func ruleMatches(r Rule, f FactSnapshot) (bool, error) {
	for _, g := range r.AllTrue {
		v, err := evalGuard(g, f)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	for _, g := range r.AllFalse {
		v, err := evalGuard(g, f)
		if err != nil {
			return false, err
		}
		if v {
			return false, nil
		}
	}
	if len(r.AnyTrue) > 0 {
		any := false
		for _, g := range r.AnyTrue {
			v, err := evalGuard(g, f)
			if err != nil {
				return false, err
			}
			if v {
				any = true
				break
			}
		}
		if !any {
			return false, nil
		}
	}
	return true, nil
}

// Evaluate is the generic evaluator engine's entry point: it reads the
// current state from f (via CurrentState), looks up that state's
// StateTransition in t, and returns the Decision produced by the first
// matching rule. The engine performs no CDS-specific reasoning of its own —
// every state name, guard combination, and proposed action comes from t.
func Evaluate(t *Table, f FactSnapshot) (Decision, error) {
	state := CurrentState(f)
	d := Decision{
		Issue:        f.Issue,
		CurrentState: state,
		Facts:        f,
	}

	for _, tr := range t.Transitions {
		if tr.State != state {
			continue
		}
		for _, r := range tr.Rules {
			ok, err := ruleMatches(r, f)
			if err != nil {
				return Decision{}, err
			}
			if !ok {
				continue
			}
			d.Outcome = r.Outcome
			d.Action = r.Action
			d.TargetState = r.TargetState
			d.RepairPass = r.RepairPass
			d.Reason = r.Reason
			if r.TargetState != "" {
				d.EnabledTransition = fmt.Sprintf("%s -> %s", state, r.TargetState)
			}
			if r.Outcome == "blocked" {
				d.BlockedReason = r.Reason
				for _, g := range r.EvidenceGuards {
					v, err := evalGuard(g, f)
					if err != nil {
						return Decision{}, err
					}
					if !v {
						d.MissingEvidence = append(d.MissingEvidence, g)
					}
				}
			}
			return d, nil
		}
		// State matched but no rule matched — a table-authoring gap, not a
		// crash: report it as blocked so the operator sees it explicitly.
		d.Outcome = "blocked"
		d.BlockedReason = fmt.Sprintf("no rule in the transition table matched state %q for issue #%d", state, f.Issue)
		return d, nil
	}

	// No transition entry at all for this state (including "" — unlabeled).
	d.Outcome = "blocked"
	d.BlockedReason = fmt.Sprintf("no transition entry for state %q in the CDS transition table", state)
	return d, nil
}
