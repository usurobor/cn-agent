// Package cell runs the kernel under a telos (the CC/PC/WC cells) and drives
// the CC→PC→WC→CC loop. A "cell" here is nothing more than the kernel FSM
// driven with a telos-specific Actor bound; the loop chains four such cell
// runs and shows the measured incoherence going 1 → 0. Spike code; see
// ../DESIGN.md.
package cell

import (
	"fmt"
	"io"
	"path/filepath"

	"cnos.dev/cnos/cell-runner/kernel"
	"cnos.dev/cnos/cell-runner/model"
)

// runCell drives one cell (one kernel run to a terminal state), then writes the
// computed receipt to <workspace>/receipts/<receiptName>.receipt.json. It
// returns the RunContext (carrying CM/wave/judgment outputs) and the terminal
// state ("decided" on accept, "held" on fail-closed).
func runCell(t *kernel.Table, a kernel.Actor, rc *kernel.RunContext, receiptName string) (string, error) {
	rc.Facts.ContractPresent = true // the contract is the α-input, present at handover
	rc.Log("[%s.%s] telos=%s contract=%s", rc.Cell, verb(rc.Cell), rc.Telos, rc.Contract.ID)

	final, err := kernel.Drive(t, a, rc, "init")
	if err != nil {
		return final, err
	}

	if rc.Receipt != nil {
		path := filepath.Join(rc.Workspace, "receipts", receiptName+".receipt.json")
		if werr := model.WriteJSON(path, rc.Receipt); werr != nil {
			return final, fmt.Errorf("write receipt: %w", werr)
		}
		rc.Log("  → receipt: receipts/%s.receipt.json (verdict=%s, transmissibility=%s)",
			receiptName, rc.Receipt.Verdict, rc.Receipt.Transmissibility)
	}
	return final, nil
}

// verb labels a cell's action for the transcript.
func verb(cell string) string {
	switch cell {
	case "CC":
		return "measure"
	case "PC":
		return "plan"
	case "WC":
		return "execute"
	default:
		return "run"
	}
}

// newLog wraps an io.Writer into a nil-safe transcript sink.
func newLog(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}
