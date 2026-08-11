package cell

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"cnos.dev/cnos/cell-runner/model"
)

const toyTask = "../testdata/toy"

// TestHappyLoop runs the full CC→PC→WC→CC loop and asserts incoherence 1 → 0,
// a coherent final judgment, and that every expected object was written.
func TestHappyLoop(t *testing.T) {
	ws := t.TempDir()
	var log bytes.Buffer
	res, err := RunLoop(Options{Workspace: ws, TaskDir: toyTask, Log: &log})
	if err != nil {
		t.Fatalf("RunLoop: %v\n%s", err, log.String())
	}
	if res.InitialIncoherence != 1 || res.FinalIncoherence != 0 {
		t.Fatalf("incoherence %d -> %d, want 1 -> 0", res.InitialIncoherence, res.FinalIncoherence)
	}
	if res.FinalJudgment != "coherent" {
		t.Fatalf("final judgment=%q want coherent", res.FinalJudgment)
	}
	transcript := log.String()
	for _, want := range []string{"incoherence: 1", "incoherence: 0", "LOOP COMPLETE"} {
		if !strings.Contains(transcript, want) {
			t.Fatalf("transcript missing %q", want)
		}
	}
	for _, f := range []string{
		"cm/cc-1.cm.json", "cm/cc-2.cm.json",
		"wave/pc.wave.json",
		"receipts/cc-1.receipt.json", "receipts/pc.receipt.json",
		"receipts/wc.receipt.json", "receipts/cc-2.receipt.json",
	} {
		if _, err := os.Stat(filepath.Join(ws, f)); err != nil {
			t.Fatalf("expected emitted file %s: %v", f, err)
		}
	}

	// The first CM records the defect; the second clears it.
	cm1 := readCM(t, filepath.Join(ws, "cm/cc-1.cm.json"))
	cm2 := readCM(t, filepath.Join(ws, "cm/cc-2.cm.json"))
	if len(cm1.Defects) != 1 || cm1.Defects[0].Kind != "missing-required-section" {
		t.Fatalf("cm1 defects=%v want one missing-required-section", cm1.Defects)
	}
	if len(cm2.Defects) != 0 {
		t.Fatalf("cm2 defects=%v want none", cm2.Defects)
	}

	// The WC receipt must be an accepted, transmissible artifact.
	rc := readReceipt(t, filepath.Join(ws, "receipts/wc.receipt.json"))
	if rc.Verdict != "accept" || rc.Transmissibility != "accepted" {
		t.Fatalf("wc receipt verdict=%q transmissibility=%q want accept/accepted", rc.Verdict, rc.Transmissibility)
	}
}

// TestFailInjection covers both fail-closed branches: β-reject and V-fail both
// drive WC to `held`, error out, and write a not_transmissible stop receipt.
func TestFailInjection(t *testing.T) {
	cases := []struct {
		fault     string
		wantFinal string // From-state of the held edge
	}{
		{"beta-reject", "reviewed"},
		{"v-fail", "validated"},
	}
	for _, c := range cases {
		t.Run(c.fault, func(t *testing.T) {
			ws := t.TempDir()
			var log bytes.Buffer
			res, err := RunLoop(Options{Workspace: ws, TaskDir: toyTask, FailInjection: c.fault, Log: &log})
			if err == nil {
				t.Fatalf("expected fail-closed error, got nil\n%s", log.String())
			}
			if !res.Held || res.HeldCell != "WC" {
				t.Fatalf("res.Held=%v cell=%q want held WC", res.Held, res.HeldCell)
			}
			rc := readReceipt(t, filepath.Join(ws, "receipts/wc.receipt.json"))
			if rc.Verdict != "hold" {
				t.Fatalf("verdict=%q want hold", rc.Verdict)
			}
			if rc.Transmissibility != "not_transmissible" {
				t.Fatalf("transmissibility=%q want not_transmissible", rc.Transmissibility)
			}
			if rc.BoundaryDecision.Action != "reject" {
				t.Fatalf("action=%q want reject", rc.BoundaryDecision.Action)
			}
			last := rc.Transitions[len(rc.Transitions)-1]
			if last.To != "held" || last.From != c.wantFinal {
				t.Fatalf("held edge %s->%s want %s->held", last.From, last.To, c.wantFinal)
			}
			// A held run must NOT reach turn 4.
			if _, err := os.Stat(filepath.Join(ws, "cm/cc-2.cm.json")); err == nil {
				t.Fatal("re-measure should not run after a held WC")
			}
		})
	}
}

// TestEmittedObjectsCueVet validates every emitted object against schema/*.cue
// with the real `cue` binary. Skipped if cue is not on PATH.
func TestEmittedObjectsCueVet(t *testing.T) {
	cueBin, err := exec.LookPath("cue")
	if err != nil {
		t.Skip("cue not on PATH; skipping schema-vet integration test")
	}
	ws := t.TempDir()
	if _, err := RunLoop(Options{Workspace: ws, TaskDir: toyTask}); err != nil {
		t.Fatalf("RunLoop: %v", err)
	}
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	objs := []struct{ file, def string }{
		{"cm/cc-1.cm.json", "#CM"},
		{"cm/cc-2.cm.json", "#CM"},
		{"wave/pc.wave.json", "#Wave"},
		{"receipts/cc-1.receipt.json", "#Receipt"},
		{"receipts/pc.receipt.json", "#Receipt"},
		{"receipts/wc.receipt.json", "#Receipt"},
		{"receipts/cc-2.receipt.json", "#Receipt"},
		{"judgment/cc-1.judgment.json", "#Judgment"},
		{"judgment/cc-2.judgment.json", "#Judgment"},
	}
	for _, o := range objs {
		cmd := exec.Command(cueBin, "vet", "./schema/", filepath.Join(ws, o.file), "-d", o.def)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("cue vet %s -d %s failed: %v\n%s", o.file, o.def, err, out)
		}
	}
}

func readCM(t *testing.T, path string) model.CM {
	t.Helper()
	var cm model.CM
	readJSON(t, path, &cm)
	return cm
}

func readReceipt(t *testing.T, path string) model.Receipt {
	t.Helper()
	var r model.Receipt
	readJSON(t, path, &r)
	return r
}

func readJSON(t *testing.T, path string, v any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
}
