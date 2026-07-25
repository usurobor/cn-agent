// Package kernel is the package-agnostic cell FSM: a declarative transition
// table (transitions.json) plus a generic evaluator that never hardcodes a
// state name in a switch/if-chain. It mirrors the cnos.issues/issues-fsm
// idiom: guards are a fixed generic predicate vocabulary; which guard applies
// to which state, and what step/target it drives, is entirely table-driven.
//
// The kernel is telos-agnostic. A "cell" is this kernel run under a telos —
// CC, PC, and WC all drive THIS table; they differ only by which Actor is
// bound (see actor.go, drive.go). Spike code; see ../DESIGN.md.
package kernel

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

// transitionsJSON is the kernel transition table, embedded so the binary is
// self-contained regardless of the working directory (criterion 2 runs
// `go run ./cmd/cell-runner` from the module root).
//
//go:embed transitions.json
var transitionsJSON []byte

// Rule is one entry in a state's rule list. A rule matches a FactSnapshot when
// every AllTrue guard is true, every AllFalse guard is false, and (if AnyTrue
// is non-empty) at least one AnyTrue guard is true. Rules are evaluated in
// file order; the first match wins. Rule is pure data decoded from the table.
type Rule struct {
	AllTrue  []string `json:"all_true,omitempty"`
	AnyTrue  []string `json:"any_true,omitempty"`
	AllFalse []string `json:"all_false,omitempty"`

	// To is the target state this rule transitions into.
	To string `json:"to"`
	// Step is the actor-step id run while taking this transition (resolved
	// against the step registry in drive.go). The step produces the evidence
	// (facts) the NEXT state's guard reads.
	Step string `json:"step"`
	// Reason is the human-readable rationale printed for this transition.
	Reason string `json:"reason"`
}

// StateTransition holds every rule for one state, evaluated when Trigger
// fires (the kernel has exactly one trigger: "advance").
type StateTransition struct {
	State   string `json:"state"`
	Trigger string `json:"trigger"`
	Rules   []Rule `json:"rules"`
}

// Table is the full declarative transition table — the single source of truth
// for the kernel FSM. It is loaded from JSON, never constructed as Go literals.
type Table struct {
	Doc         []string          `json:"_doc,omitempty"`
	States      []string          `json:"states"`
	Terminal    []string          `json:"terminal"`
	Guards      map[string]string `json:"guards,omitempty"`
	Transitions []StateTransition `json:"transitions"`
}

// DefaultTable parses the embedded transition table.
func DefaultTable() (*Table, error) { return parseTable(transitionsJSON) }

// LoadTable reads and parses a transition table from path (used by tests that
// want to exercise an alternate table).
func LoadTable(path string) (*Table, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read transition table: %w", err)
	}
	return parseTable(data)
}

func parseTable(data []byte) (*Table, error) {
	var t Table
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parse transition table JSON: %w", err)
	}
	return &t, nil
}

// IsTerminal reports whether state is a terminal state (no outgoing edges).
func (t *Table) IsTerminal(state string) bool {
	for _, s := range t.Terminal {
		if s == state {
			return true
		}
	}
	return false
}

// guardFuncs is the fixed, generic predicate vocabulary the evaluator
// understands — the ONLY place FactSnapshot fields are read as booleans. It is
// never switched on a state name; it only names reusable evidence predicates
// over the fact model, exactly like issues-fsm's guard registry.
var guardFuncs = map[string]func(FactSnapshot) bool{
	"contract_present": func(f FactSnapshot) bool { return f.ContractPresent },
	"artifact_present": func(f FactSnapshot) bool { return f.ArtifactPresent },
	"review_approved":  func(f FactSnapshot) bool { return f.ReviewApproved },
	"review_rejected":  func(f FactSnapshot) bool { return f.ReviewRejected },
	"receipt_bound":    func(f FactSnapshot) bool { return f.ReceiptBound },
	"v_pass":           func(f FactSnapshot) bool { return f.VPass },
	"v_fail":           func(f FactSnapshot) bool { return f.VFail },
}

// evalGuard evaluates a single named guard. It errors if the table references
// a guard the engine does not recognize — a table-authoring error.
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
