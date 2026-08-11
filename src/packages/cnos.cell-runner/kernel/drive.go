package kernel

import (
	"fmt"

	"cnos.dev/cnos/cell-runner/model"
)

// Drive runs the kernel FSM to a terminal state. Starting at `initial`, it
// repeatedly: (1) evaluates the table for the current state, (2) records the
// matched transition, (3) runs the transition's step actor (which mutates the
// facts that gate the NEXT transition), and (4) advances to the target state —
// until a terminal state is reached.
//
// Drive is fully generic: it contains no cell-specific and no state-specific
// logic. Every state name, guard, and step comes from the table + the step
// registry. It returns the terminal state reached, plus an error if a step
// fails or the machine gets stuck (no guard satisfied at a non-terminal
// state — the fail-closed signal).
func Drive(t *Table, a Actor, rc *RunContext, initial string) (string, error) {
	state := initial
	rc.Facts.State = state
	for {
		if t.IsTerminal(state) {
			return state, nil
		}
		d, err := Evaluate(t, state, *rc.Facts)
		if err != nil {
			return state, err
		}
		if !d.Matched {
			return state, fmt.Errorf("kernel stuck at state %q: no transition guard satisfied", state)
		}

		rc.Transitions = append(rc.Transitions, model.TransitionRecord{
			From:  d.From,
			To:    d.To,
			Step:  d.Step,
			Guard: d.Guard,
		})
		rc.Log("  kernel: %-9s --%-16s--> %-9s  [%s]", d.From, d.Guard, d.To, d.Step)

		step, ok := stepFuncs[d.Step]
		if !ok {
			return state, fmt.Errorf("transition table references unknown step %q", d.Step)
		}
		if err := step(a, rc); err != nil {
			return state, fmt.Errorf("step %q failed: %w", d.Step, err)
		}

		state = d.To
		rc.Facts.State = state
	}
}
