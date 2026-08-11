package cell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cnos.dev/cnos/cell-runner/actor/stub"
	"cnos.dev/cnos/cell-runner/kernel"
	"cnos.dev/cnos/cell-runner/model"
)

// Options configures one loop run.
type Options struct {
	// Workspace is the writable run root (CM/wave/receipts are written here).
	Workspace string
	// TaskDir is the source task directory (holds the target markdown). Its
	// target is copied into the workspace so the source is never mutated.
	TaskDir string
	// FailInjection injects a fault into the WC turn: "" (none), "beta-reject",
	// or "v-fail".
	FailInjection string
	// Log receives the loop transcript (nil → discarded).
	Log io.Writer
}

// Result summarizes a loop run.
type Result struct {
	InitialIncoherence int
	FinalIncoherence   int
	FinalJudgment      string
	Held               bool
	HeldCell           string
	TargetPath         string
}

// RunLoop executes one CC→PC→WC→CC turn of the agentic loop. It returns a
// non-nil error when the loop fails closed (a cell reached `held`) or does not
// converge — the caller maps that to a non-zero exit.
func RunLoop(opts Options) (Result, error) {
	log := newLog(opts.Log)
	var res Result

	target, err := prepareWorkspace(opts.Workspace, opts.TaskDir)
	if err != nil {
		return res, err
	}
	res.TargetPath = target

	t, err := kernel.DefaultTable()
	if err != nil {
		return res, fmt.Errorf("load kernel table: %w", err)
	}

	fmt.Fprintf(log, "=== cnos cell-runner — CC→PC→WC→CC walking skeleton ===\n")
	fmt.Fprintf(log, "workspace: %s\ntarget:    %s\n", opts.Workspace, target)
	if opts.FailInjection != "" {
		fmt.Fprintf(log, "fault:     %s (injected into the WC turn)\n", opts.FailInjection)
	}
	fmt.Fprintln(log)

	// --- Turn 1: CC.measure -------------------------------------------------
	fmt.Fprintln(log, "── turn 1: CC.measure ──────────────────────────────────")
	cc1 := kernel.NewRunContext("CC", "judgment", "ref:turn-1/CC", opts.Workspace, target,
		model.Contract{ID: "cc-measure", Telos: "judgment", Target: target,
			Instruction: "measure the target's coherence", Acceptance: []string{"cm_has_target", "cm_has_provenance"}},
		"", log)
	if _, err := runCell(t, stub.NewCC(), cc1, "cc-1"); err != nil {
		return res, err
	}
	if err := writeCCOutputs(cc1, "cc-1"); err != nil {
		return res, err
	}
	res.InitialIncoherence = cc1.CM.Incoherence()
	fmt.Fprintf(log, "  incoherence: %d   judgment: %s\n\n", res.InitialIncoherence, cc1.Judgment.Verdict)

	if cc1.Judgment.Verdict == "coherent" {
		res.FinalIncoherence = res.InitialIncoherence
		res.FinalJudgment = "coherent"
		fmt.Fprintln(log, "target already coherent; nothing to do.")
		return res, nil
	}

	// --- Turn 2: PC.plan ----------------------------------------------------
	fmt.Fprintln(log, "── turn 2: PC.plan ─────────────────────────────────────")
	pc := kernel.NewRunContext("PC", "relation_graph", "ref:turn-2/PC", opts.Workspace, target,
		model.Contract{ID: "pc-plan", Telos: "relation_graph", Target: target,
			Instruction: "plan a wave that clears the measured defect", Acceptance: []string{"wave_has_exactly_one_contract"}},
		"", log)
	pc.CM = cc1.CM // the PC plans from the CM the CC produced
	if _, err := runCell(t, stub.NewPC(), pc, "pc"); err != nil {
		return res, err
	}
	if err := model.WriteJSON(filepath.Join(opts.Workspace, "wave", "pc.wave.json"), pc.Wave); err != nil {
		return res, fmt.Errorf("write wave: %w", err)
	}
	pc.Log("  → wave: wave/pc.wave.json (%d contract)\n", len(pc.Wave.Contracts))

	// --- Turn 3: WC.execute -------------------------------------------------
	fmt.Fprintln(log, "── turn 3: WC.execute ──────────────────────────────────")
	wcContract := pc.Wave.Contracts[0]
	wc := kernel.NewRunContext("WC", "artifact", "ref:turn-3/WC", opts.Workspace, target,
		wcContract, opts.FailInjection, log)
	wcFinal, err := runCell(t, stub.NewWC(), wc, "wc")
	if err != nil {
		return res, err
	}
	if wcFinal == "held" {
		res.Held = true
		res.HeldCell = "WC"
		fmt.Fprintf(log, "\nLOOP HELD at WC (fail closed): stop receipt written to receipts/wc.receipt.json\n")
		return res, fmt.Errorf("loop failed closed: WC reached 'held' (%s)", describeHold(opts.FailInjection))
	}
	fmt.Fprintln(log)

	// --- Turn 4: CC.re-measure ---------------------------------------------
	fmt.Fprintln(log, "── turn 4: CC.re-measure ───────────────────────────────")
	cc2 := kernel.NewRunContext("CC", "judgment", "ref:turn-4/CC", opts.Workspace, target,
		model.Contract{ID: "cc-remeasure", Telos: "judgment", Target: target,
			Instruction: "re-measure the target after the WC edit", Acceptance: []string{"cm_has_target", "cm_has_provenance"}},
		"", log)
	if _, err := runCell(t, stub.NewCC(), cc2, "cc-2"); err != nil {
		return res, err
	}
	if err := writeCCOutputs(cc2, "cc-2"); err != nil {
		return res, err
	}
	res.FinalIncoherence = cc2.CM.Incoherence()
	res.FinalJudgment = cc2.Judgment.Verdict
	fmt.Fprintf(log, "  incoherence: %d   judgment: %s\n\n", res.FinalIncoherence, cc2.Judgment.Verdict)

	fmt.Fprintf(log, "=== LOOP COMPLETE: incoherence %d → %d, judgment=%s, receipt chain closed ===\n",
		res.InitialIncoherence, res.FinalIncoherence, res.FinalJudgment)

	if res.FinalJudgment != "coherent" || res.FinalIncoherence != 0 {
		return res, fmt.Errorf("loop did not converge: final judgment=%s incoherence=%d", res.FinalJudgment, res.FinalIncoherence)
	}
	return res, nil
}

// writeCCOutputs writes a CC turn's CM and judgment objects.
func writeCCOutputs(rc *kernel.RunContext, name string) error {
	if err := model.WriteJSON(filepath.Join(rc.Workspace, "cm", name+".cm.json"), rc.CM); err != nil {
		return fmt.Errorf("write CM: %w", err)
	}
	if err := model.WriteJSON(filepath.Join(rc.Workspace, "judgment", name+".judgment.json"), rc.Judgment); err != nil {
		return fmt.Errorf("write judgment: %w", err)
	}
	rc.Log("  → cm: cm/%s.cm.json   judgment: judgment/%s.judgment.json", name, name)
	return nil
}

func describeHold(faultInjection string) string {
	switch faultInjection {
	case "beta-reject":
		return "β rejected the artifact; reviewed → held"
	case "v-fail":
		return "V rejected the artifact; validated → held"
	default:
		return "a gate rejected the artifact"
	}
}

// prepareWorkspace creates the workspace layout and copies the task's target
// markdown into <workspace>/target/, returning the copy's path. The source
// task directory is never mutated (so tests and re-runs are repeatable).
func prepareWorkspace(workspace, taskDir string) (string, error) {
	if workspace == "" {
		return "", fmt.Errorf("--workspace is required")
	}
	srcTarget, err := findTarget(taskDir)
	if err != nil {
		return "", err
	}
	dstDir := filepath.Join(workspace, "target")
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dstDir, filepath.Base(srcTarget))
	data, err := os.ReadFile(srcTarget)
	if err != nil {
		return "", fmt.Errorf("read task target: %w", err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return "", fmt.Errorf("copy task target: %w", err)
	}
	return dst, nil
}

// findTarget returns the single markdown target inside taskDir.
func findTarget(taskDir string) (string, error) {
	if taskDir == "" {
		return "", fmt.Errorf("--task is required")
	}
	matches, err := filepath.Glob(filepath.Join(taskDir, "*.md"))
	if err != nil {
		return "", err
	}
	if len(matches) != 1 {
		return "", fmt.Errorf("task %q must contain exactly one .md target, found %d", taskDir, len(matches))
	}
	return matches[0], nil
}
