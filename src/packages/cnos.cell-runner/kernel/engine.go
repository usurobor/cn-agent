package kernel

import (
	"fmt"
	"strings"
)

// FactSnapshot is the evaluator's entire input: an explicit, inference-free
// observation of one cell's evidence. Each boolean is produced by an actor
// step (or, for ContractPresent, by the cell handing over the contract). The
// evaluator reads these; it never re-derives a fact by guessing.
type FactSnapshot struct {
	// State is the current kernel state.
	State string

	ContractPresent bool
	ArtifactPresent bool
	ReviewApproved  bool
	ReviewRejected  bool
	ReceiptBound    bool
	VPass           bool
	VFail           bool
}

// Decision is the evaluator's structured output for one (state, facts) pair:
// the matched transition, or Matched=false when no rule fires (either a
// terminal state, or a genuinely stuck state). It is a pure projection of a
// Table + FactSnapshot — read-only and idempotent.
type Decision struct {
	From    string
	To      string
	Step    string
	Guard   string
	Reason  string
	Matched bool
}

// Evaluate is the generic evaluator: it looks up state's entry in t and
// returns the Decision from the first matching rule. It performs no
// cell-specific reasoning — every state name, guard, step, and target comes
// from t. This is the "never hardcode a state name in a switch" invariant.
func Evaluate(t *Table, state string, f FactSnapshot) (Decision, error) {
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
			return Decision{
				From:    state,
				To:      r.To,
				Step:    r.Step,
				Guard:   guardLabel(r),
				Reason:  r.Reason,
				Matched: true,
			}, nil
		}
		// State entry exists but no rule matched: not an error here — the
		// caller (Drive) decides whether that is terminal or stuck.
		return Decision{From: state, Matched: false}, nil
	}
	// No entry for this state at all (e.g. a terminal state).
	return Decision{From: state, Matched: false}, nil
}

// guardLabel renders a rule's guard set for the transcript/receipt.
func guardLabel(r Rule) string {
	parts := append([]string{}, r.AllTrue...)
	for _, g := range r.AnyTrue {
		parts = append(parts, "any:"+g)
	}
	for _, g := range r.AllFalse {
		parts = append(parts, "!"+g)
	}
	if len(parts) == 0 {
		return "(none)"
	}
	return strings.Join(parts, "&")
}

// GuardsKnown reports whether every guard referenced anywhere in the table is
// implemented by the engine's guard registry — a table-authoring lint used by
// tests to keep transitions.json and the guard vocabulary in sync.
func (t *Table) GuardsKnown() error {
	var unknown []string
	for _, tr := range t.Transitions {
		for _, r := range tr.Rules {
			for _, g := range append(append(append([]string{}, r.AllTrue...), r.AnyTrue...), r.AllFalse...) {
				if _, ok := guardFuncs[g]; !ok {
					unknown = append(unknown, g)
				}
			}
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("table references unknown guards: %s", strings.Join(unknown, ", "))
	}
	return nil
}
