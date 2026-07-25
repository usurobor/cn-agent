// Command cell-runner is the walking-skeleton driver for the three-cell
// agentic loop (CC→PC→WC→CC). It runs one turn of the loop over a toy task,
// prints a transcript showing the measured incoherence going 1 → 0, and writes
// CUE-valid CM / wave / receipt objects under the workspace.
//
// Exit codes: 0 on a converged loop; non-zero when the loop fails closed (a
// cell reached the kernel's `held` state) or does not converge. This is spike
// code (an exploratory proof-of-life), not canonical runtime — see
// ../../DESIGN.md and ../../README.md.
package main

import (
	"flag"
	"fmt"
	"os"

	"cnos.dev/cnos/cell-runner/cell"
)

func main() {
	workspace := flag.String("workspace", "", "writable run root for CM/wave/receipt outputs (required)")
	task := flag.String("task", "", "task directory holding the target markdown (required, e.g. testdata/toy)")
	failInjection := flag.String("fail-injection", "", "inject a WC fault: 'beta-reject' or 'v-fail' (drives the kernel to 'held')")
	flag.Parse()

	switch *failInjection {
	case "", "beta-reject", "v-fail":
	default:
		fmt.Fprintf(os.Stderr, "cell-runner: unknown --fail-injection %q (want 'beta-reject' or 'v-fail')\n", *failInjection)
		os.Exit(2)
	}

	res, err := cell.RunLoop(cell.Options{
		Workspace:     *workspace,
		TaskDir:       *task,
		FailInjection: *failInjection,
		Log:           os.Stdout,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\ncell-runner: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nOK: incoherence %d → %d (%s)\n", res.InitialIncoherence, res.FinalIncoherence, res.FinalJudgment)
}
