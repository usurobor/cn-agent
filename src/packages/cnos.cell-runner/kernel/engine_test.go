package kernel

import (
	"io"
	"testing"

	"cnos.dev/cnos/cell-runner/model"
)

// TestEvaluateFirstMatchWins checks the generic evaluator picks the correct
// transition per state from the default (embedded) table, purely from facts.
func TestEvaluateFirstMatchWins(t *testing.T) {
	tab, err := DefaultTable()
	if err != nil {
		t.Fatalf("DefaultTable: %v", err)
	}
	cases := []struct {
		name     string
		state    string
		facts    FactSnapshot
		wantTo   string
		wantStep string
		matched  bool
	}{
		{"init->produced", "init", FactSnapshot{ContractPresent: true}, "produced", "alpha_produce", true},
		{"produced->reviewed", "produced", FactSnapshot{ArtifactPresent: true}, "reviewed", "beta_review", true},
		{"reviewed->closed", "reviewed", FactSnapshot{ReviewApproved: true}, "closed", "gamma_close", true},
		{"reviewed->held", "reviewed", FactSnapshot{ReviewRejected: true}, "held", "delta_hold", true},
		{"closed->validated", "closed", FactSnapshot{ReceiptBound: true}, "validated", "v_validate", true},
		{"validated->decided", "validated", FactSnapshot{VPass: true}, "decided", "delta_accept", true},
		{"validated->held", "validated", FactSnapshot{VFail: true}, "held", "delta_hold", true},
		{"no facts, no match", "init", FactSnapshot{}, "", "", false},
		{"terminal decided", "decided", FactSnapshot{}, "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, err := Evaluate(tab, c.state, c.facts)
			if err != nil {
				t.Fatalf("Evaluate: %v", err)
			}
			if d.Matched != c.matched {
				t.Fatalf("matched=%v want %v", d.Matched, c.matched)
			}
			if d.To != c.wantTo || d.Step != c.wantStep {
				t.Fatalf("got (to=%q step=%q) want (to=%q step=%q)", d.To, d.Step, c.wantTo, c.wantStep)
			}
		})
	}
}

// TestUnknownGuardErrors ensures a table referencing an unknown guard is a
// clean error, not a silent pass.
func TestUnknownGuardErrors(t *testing.T) {
	tab := &Table{
		States:   []string{"init"},
		Terminal: []string{},
		Transitions: []StateTransition{{
			State:   "init",
			Trigger: "advance",
			Rules:   []Rule{{AllTrue: []string{"no_such_guard"}, To: "x", Step: "y"}},
		}},
	}
	if _, err := Evaluate(tab, "init", FactSnapshot{}); err == nil {
		t.Fatal("expected error for unknown guard, got nil")
	}
	if err := tab.GuardsKnown(); err == nil {
		t.Fatal("GuardsKnown should reject an unknown guard")
	}
}

// TestDefaultTableGuardsKnown checks the shipped table only uses implemented guards.
func TestDefaultTableGuardsKnown(t *testing.T) {
	tab, err := DefaultTable()
	if err != nil {
		t.Fatal(err)
	}
	if err := tab.GuardsKnown(); err != nil {
		t.Fatalf("default table has unknown guards: %v", err)
	}
}

// scriptedActor is a minimal fact-only actor used to exercise Drive without
// any file IO: each step sets the fact its transition produces. reject/vFail
// select the held branches.
type scriptedActor struct {
	reject bool
	vFail  bool
}

func (a scriptedActor) Alpha(rc *RunContext) error { rc.Facts.ArtifactPresent = true; return nil }
func (a scriptedActor) Beta(rc *RunContext) error {
	if a.reject {
		rc.Facts.ReviewRejected = true
	} else {
		rc.Facts.ReviewApproved = true
	}
	return nil
}
func (a scriptedActor) Gamma(rc *RunContext) error { rc.Facts.ReceiptBound = true; return nil }
func (a scriptedActor) Validate(rc *RunContext) error {
	if a.vFail {
		rc.Facts.VFail = true
	} else {
		rc.Facts.VPass = true
	}
	return nil
}
func (a scriptedActor) DeltaAccept(rc *RunContext) error {
	rc.Receipt = &model.Receipt{Verdict: "accept"}
	return nil
}
func (a scriptedActor) DeltaHold(rc *RunContext) error {
	rc.Receipt = &model.Receipt{Verdict: "hold"}
	return nil
}

func drive(t *testing.T, a scriptedActor) string {
	t.Helper()
	tab, err := DefaultTable()
	if err != nil {
		t.Fatal(err)
	}
	rc := NewRunContext("WC", "artifact", "ref:test", t.TempDir(), "", model.Contract{}, "", io.Discard)
	rc.Facts.ContractPresent = true
	final, err := Drive(tab, a, rc, "init")
	if err != nil {
		t.Fatalf("Drive: %v", err)
	}
	return final
}

// TestDriveReachesDecided drives the happy path to the decided terminal.
func TestDriveReachesDecided(t *testing.T) {
	if got := drive(t, scriptedActor{}); got != "decided" {
		t.Fatalf("final=%q want decided", got)
	}
}

// TestDriveBetaRejectHolds drives the β-reject branch to held.
func TestDriveBetaRejectHolds(t *testing.T) {
	if got := drive(t, scriptedActor{reject: true}); got != "held" {
		t.Fatalf("final=%q want held", got)
	}
}

// TestDriveVFailHolds drives the V-fail branch to held.
func TestDriveVFailHolds(t *testing.T) {
	if got := drive(t, scriptedActor{vFail: true}); got != "held" {
		t.Fatalf("final=%q want held", got)
	}
}
