package kernel

import (
	"fmt"
	"io"

	"cnos.dev/cnos/cell-runner/model"
)

// RunContext is the mutable state of one kernel run (one cell turn). Actor
// steps read and mutate it: they set facts, append evidence refs, and build
// the CM / wave / judgment / receipt outputs. The kernel Drive loop owns the
// state machine; the RunContext owns the accumulated evidence.
type RunContext struct {
	// Cell is the telos-role label: "CC", "PC", or "WC".
	Cell string
	// Telos is the output kind: "judgment", "relation_graph", or "artifact".
	Telos string
	// Ref is a synthetic, git-native-style reference string for this turn
	// (e.g. "ref:turn-1/CC"); used as measured_at_ref / produced_at_ref. The
	// spike does not touch git, so refs are stable synthetic identifiers.
	Ref string

	// Workspace is the run's writable root.
	Workspace string
	// Target is the path (inside the workspace) of the file the cell acts on.
	Target string
	// Contract is the α-input for this run.
	Contract model.Contract
	// FailInjection selects a fault to inject ("", "beta-reject", "v-fail").
	FailInjection string

	// Facts is the evaluator's observation, mutated by actor steps.
	Facts *FactSnapshot
	// Evidence accumulates typed evidence refs for the receipt.
	Evidence map[string]string
	// Transitions records every kernel edge traversed (receipt evidence).
	Transitions []model.TransitionRecord

	// Outputs (set by the relevant actor step for this cell's telos).
	CM       *model.CM
	Wave     *model.Wave
	Judgment *model.Judgment
	Receipt  *model.Receipt

	logw io.Writer
}

// NewRunContext builds a RunContext with an initialized fact snapshot and
// evidence map, logging the transcript to logw.
func NewRunContext(cell, telos, ref, workspace, target string, contract model.Contract, failInjection string, logw io.Writer) *RunContext {
	return &RunContext{
		Cell:          cell,
		Telos:         telos,
		Ref:           ref,
		Workspace:     workspace,
		Target:        target,
		Contract:      contract,
		FailInjection: failInjection,
		Facts:         &FactSnapshot{},
		Evidence:      map[string]string{},
		logw:          logw,
	}
}

// Log writes a transcript line.
func (rc *RunContext) Log(format string, args ...any) {
	fmt.Fprintf(rc.logw, format+"\n", args...)
}

// AddEvidence records a typed evidence ref for the receipt.
func (rc *RunContext) AddEvidence(key, ref string) { rc.Evidence[key] = ref }
