package issuesfsm

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// realTablePath is the CDS transition table's location relative to this
// package's directory (the package boundary AC1 pins: cnos.cds owns the
// table as data, cnos.issues owns the generic engine that reads it).
const realTablePath = "../../../cnos.cds/skills/cds/fsm/transitions.json"

func loadRealTable(t *testing.T) *Table {
	t.Helper()
	tab, err := LoadTable(realTablePath)
	if err != nil {
		t.Fatalf("LoadTable(%s): %v", realTablePath, err)
	}
	return tab
}

// --- AC1: transition table exists, is declarative JSON data, and the
// engine reads from it (not from inline Go literals / a switch on state
// names). ---

func TestAC1_TableLoadsAsData(t *testing.T) {
	tab := loadRealTable(t)
	if len(tab.States) == 0 {
		t.Fatal("transition table has no states")
	}
	if len(tab.Transitions) == 0 {
		t.Fatal("transition table has no transitions")
	}
	for _, want := range []string{"ready", "todo", "in-progress", "review", "changes"} {
		found := false
		for _, tr := range tab.Transitions {
			if tr.State == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("transition table missing state %q", want)
		}
	}
}

// --- AC2: `cn issues fsm evaluate --issue N --fixture path` prints current
// state, observed facts, enabled transition, blocked reason, and proposed
// action — and mutates nothing. ---

func TestAC2_RunPrintsFullDecisionBlock(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(),
		[]string{"evaluate", "--issue", "601", "--fixture", "testdata/review-empty.json", "--table", realTablePath},
		nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"Current state: review",
		"Observed facts:",
		"labels:",
		"cell_kind: (none) (source: absent, defaulted_to: implementation)",
		"enabled_transition:",
		"blocked_reason:",
		"proposed_action:",
		"read-only: no label was written",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got:\n%s", want, out)
		}
	}
}

func TestAC2_UnknownSubcommandRejected(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{"apply", "--issue", "1"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error for the unsupported 'apply' subcommand (Phase 2, out of scope)")
	}
	if !strings.Contains(stderr.String(), "Phase 2") {
		t.Errorf("expected stderr to explain apply is Phase 2, got: %s", stderr.String())
	}
}

// --- AC3: status:review with no PR, no commits, no REVIEW-REQUEST.yml is
// flagged blocked/invalid with a missing-evidence list. ---

func TestAC3_EmptyReviewBlocked(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/review-empty.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "blocked" {
		t.Fatalf("outcome = %q, want blocked", dec.Outcome)
	}
	if dec.BlockedReason == "" {
		t.Error("expected a non-empty blocked_reason")
	}
	wantMissing := map[string]bool{"pr_exists": true, "branch_has_commits": true, "review_request_present": true}
	if len(dec.MissingEvidence) != len(wantMissing) {
		t.Errorf("missing_evidence = %v, want exactly %v", dec.MissingEvidence, wantMissing)
	}
	for _, m := range dec.MissingEvidence {
		if !wantMissing[m] {
			t.Errorf("unexpected missing_evidence entry %q", m)
		}
	}
}

func TestAC3_ReviewWithEvidenceIsValid(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/review-with-pr.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "valid" {
		t.Errorf("outcome = %q, want valid (review has a PR)", dec.Outcome)
	}
	if dec.BlockedReason != "" {
		t.Errorf("expected no blocked_reason, got %q", dec.BlockedReason)
	}
}

// --- AC4: status:in-progress, no active run, no/empty branch, no PR ->
// proposes status:todo + requeue reason. ---

func TestAC4_DeadInProgressNoMatterProposesRequeue(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-dead-no-matter.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" {
		t.Fatalf("outcome = %q, want proposed", dec.Outcome)
	}
	if dec.TargetState != "todo" {
		t.Errorf("target_state = %q, want todo", dec.TargetState)
	}
	if dec.EnabledTransition != "in-progress -> todo" {
		t.Errorf("enabled_transition = %q, want %q", dec.EnabledTransition, "in-progress -> todo")
	}
	if dec.Action != "propose_status_todo" {
		t.Errorf("action = %q, want propose_status_todo", dec.Action)
	}
}

// --- AC5: status:in-progress, no active run, branch has commits beyond
// base, no PR -> proposes delta-recovery, NEVER status:todo (cnos#368). ---

func TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-dead-with-commits.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" {
		t.Fatalf("outcome = %q, want proposed", dec.Outcome)
	}
	if dec.Action != "propose_delta_recovery" {
		t.Fatalf("action = %q, want propose_delta_recovery", dec.Action)
	}
	if dec.TargetState == "todo" {
		t.Fatal("cnos#368 regression: dead in-progress WITH commits must never propose status:todo (blind re-dispatch over existing work)")
	}
}

func TestAC5_HealthyActiveInProgressIsValid(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-active.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "valid" || dec.Action != "none" {
		t.Errorf("active run should be valid/none, got outcome=%q action=%q", dec.Outcome, dec.Action)
	}
}

// --- AC6: status:changes without repair context blocks the todo proposal;
// with a repair contract present, enables a repair_pass proposal. ---

func TestAC6_ChangesWithoutRepairContextBlocked(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/changes-no-repair.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "blocked" {
		t.Fatalf("outcome = %q, want blocked (no repair contract present)", dec.Outcome)
	}
	if dec.TargetState == "todo" {
		t.Fatal("must not propose todo without repair context")
	}
}

func TestAC6_ChangesWithRepairContextEnablesRepairPass(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/changes-with-repair.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" {
		t.Fatalf("outcome = %q, want proposed", dec.Outcome)
	}
	if dec.TargetState != "todo" {
		t.Errorf("target_state = %q, want todo", dec.TargetState)
	}
	if !dec.RepairPass {
		t.Error("expected repair_pass=true when a repair contract is present")
	}
}

// --- AC7: idempotence — evaluating twice on the same fixture yields an
// identical decision (same proposed action, no duplicate/different action
// on the second run). ---

func TestAC7_Idempotent(t *testing.T) {
	for _, fixture := range []string{
		"testdata/review-empty.json",
		"testdata/in-progress-dead-no-matter.json",
		"testdata/in-progress-dead-with-commits.json",
		"testdata/changes-no-repair.json",
		"testdata/changes-with-repair.json",
	} {
		t.Run(fixture, func(t *testing.T) {
			tab := loadRealTable(t)
			snap, err := LoadFixture(fixture)
			if err != nil {
				t.Fatal(err)
			}
			d1, err := Evaluate(tab, snap)
			if err != nil {
				t.Fatal(err)
			}
			d2, err := Evaluate(tab, snap)
			if err != nil {
				t.Fatal(err)
			}
			if d1.Outcome != d2.Outcome || d1.Action != d2.Action || d1.TargetState != d2.TargetState || d1.Reason != d2.Reason {
				t.Fatalf("second Evaluate diverged: first=%+v second=%+v", d1, d2)
			}

			// Also drive it through the CLI Run() path twice and diff the
			// rendered bytes directly — the operator-visible surface must
			// be byte-identical across runs, not just field-equal.
			var out1, out2, stderr bytes.Buffer
			args := []string{"evaluate", "--issue", "1", "--fixture", fixture, "--table", realTablePath}
			if err := Run(context.Background(), args, nil, &out1, &stderr); err != nil {
				t.Fatalf("first Run: %v", err)
			}
			if err := Run(context.Background(), args, nil, &out2, &stderr); err != nil {
				t.Fatalf("second Run: %v", err)
			}
			if out1.String() != out2.String() {
				t.Fatalf("Run() output not idempotent for %s", fixture)
			}
		})
	}
}

// --- AC8: zero label mutation. No --apply flag exists; no gh/GitHub
// label-write call anywhere in this package's non-test source. ---

func TestAC8_NoApplyFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	// --apply is not a registered flag: flag.ContinueOnError should reject
	// it as unknown.
	err := Run(context.Background(), []string{"evaluate", "--apply", "--issue", "1"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error: --apply is not a recognized flag in Phase 1")
	}
}

func TestAC8_NoLabelMutationCodeInSource(t *testing.T) {
	forbidden := []string{
		// A registered --apply flag (as opposed to comment prose that
		// documents its absence) is the actual AC8 regression to catch.
		`fs.Bool("apply"`,
		`fs.String("apply"`,
		"gh issue edit",
		"gh label",
		"-X PATCH",
		"-X POST",
		"/labels\"",
		"AddLabel",
		"RemoveLabel",
	}
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(e.Name(), "_test.go") {
			// This test file itself legitimately contains the forbidden
			// strings as string literals to check for; skip it.
			continue
		}
		data, err := os.ReadFile(e.Name())
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		for _, f := range forbidden {
			if strings.Contains(content, f) {
				t.Errorf("%s contains forbidden label-mutation pattern %q (AC8 violation)", e.Name(), f)
			}
		}
	}
}

// --- Proof plan positive cases: ready -> todo enablement; todo is a no-op. ---

func TestReadyEnablesTodo(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/ready.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" || dec.TargetState != "todo" {
		t.Errorf("ready should propose todo, got outcome=%q target=%q", dec.Outcome, dec.TargetState)
	}
}

func TestCurrentState(t *testing.T) {
	cases := []struct {
		labels []string
		want   string
	}{
		{[]string{"status:review", "dispatch:cell"}, "review"},
		{[]string{"dispatch:cell"}, ""},
		{nil, ""},
	}
	for _, c := range cases {
		got := CurrentState(FactSnapshot{Labels: c.labels})
		if got != c.want {
			t.Errorf("CurrentState(%v) = %q, want %q", c.labels, got, c.want)
		}
	}
}

func TestEvaluate_UnknownState_Blocked(t *testing.T) {
	tab := loadRealTable(t)
	dec, err := Evaluate(tab, FactSnapshot{Issue: 1, Labels: []string{"status:done"}})
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "blocked" {
		t.Errorf("outcome = %q, want blocked for a state absent from the table", dec.Outcome)
	}
}

// --- cell_kind seam (cnos#568 operator note / cnos#570 taxonomy).
// Contract: the field is OPTIONAL and DEFAULTED when absent, and it is NOT
// ENFORCED — no transition rule consumes it, so its value cannot change any
// Phase-1 decision. This test locks both halves so the seam cannot silently
// become enforcement without failing here. ---

func TestSeam_CellKindDefaultedWhenAbsent(t *testing.T) {
	// Every shipped fixture omits cell_kind; LoadFixture must record the
	// defaulting explicitly rather than leaving it silent/empty.
	snap, err := LoadFixture("testdata/review-with-pr.json")
	if err != nil {
		t.Fatal(err)
	}
	if snap.CellKind.Observed != "" {
		t.Errorf("Observed = %q, want empty (no source parsed in Phase 1)", snap.CellKind.Observed)
	}
	if snap.CellKind.Source != "absent" {
		t.Errorf("Source = %q, want \"absent\"", snap.CellKind.Source)
	}
	if snap.CellKind.DefaultedTo != "implementation" {
		t.Errorf("DefaultedTo = %q, want \"implementation\"", snap.CellKind.DefaultedTo)
	}
}

func TestSeam_CellKindNotEnforced(t *testing.T) {
	tab := loadRealTable(t)
	// Evaluate the same facts under several cell_kind values (including the
	// non-implementation kinds cnos#570 will define). The Decision must be
	// byte-identical every time: Phase 1 observes cell_kind but no rule
	// consumes it, so it cannot change the outcome, target, or action.
	for _, fx := range []string{
		"testdata/review-empty.json",
		"testdata/in-progress-dead-with-commits.json",
		"testdata/changes-with-repair.json",
	} {
		base, err := LoadFixture(fx)
		if err != nil {
			t.Fatalf("%s: %v", fx, err)
		}
		want, err := Evaluate(tab, base)
		if err != nil {
			t.Fatalf("%s: %v", fx, err)
		}
		for _, kind := range []string{"implementation", "issue_authoring", "wave", "cleanup", "recovery"} {
			snap := base
			snap.CellKind = CellKind{Observed: kind, Source: "issue_body"}
			got, err := Evaluate(tab, snap)
			if err != nil {
				t.Fatalf("%s / %s: %v", fx, kind, err)
			}
			if got.Outcome != want.Outcome || got.TargetState != want.TargetState || got.Action != want.Action {
				t.Errorf("%s: cell_kind=%q changed the decision (outcome %q→%q, target %q→%q, action %q→%q) — seam must not be enforced",
					fx, kind, want.Outcome, got.Outcome, want.TargetState, got.TargetState, want.Action, got.Action)
			}
		}
	}
}

// --- cnos#570 observation wiring: parseCellKind reads the `cell_kind:`
// recording line CELL-KINDS.md §"Recording point" names as canonical
// (.cdd/unreleased/{N}/gamma-scaffold.md), in both the bold-markdown form γ
// scaffolds use and the plain form. This is the helper assembleLive calls;
// testing it directly needs no network/GitHub API access. ---

func TestParseCellKind(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{
			name: "bold markdown form with trailing prose",
			body: "**cell_kind:** `doctrine` (this cell produces CDD doctrine + an observation-wiring change)",
			want: "doctrine",
		},
		{
			name: "plain form",
			body: "cell_kind: implementation",
			want: "implementation",
		},
		{
			name: "embedded in a larger scaffold body",
			body: "# γ scaffold\n\n**Branch:** `cycle/570`\n**cell_kind:** `wave`\n\n## Surfaces",
			want: "wave",
		},
		{
			name: "absent",
			body: "# γ scaffold\n\nNo cell kind line here.",
			want: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := parseCellKind(c.body); got != c.want {
				t.Errorf("parseCellKind(%q) = %q, want %q", c.body, got, c.want)
			}
		})
	}
}

// TestAssembleLive_ObservesCellKindFromGammaScaffold drives assembleLive
// (repo="" so it never makes a network call) against a temp directory
// carrying a .cdd/unreleased/{N}/gamma-scaffold.md fixture, and asserts the
// observation populates CellKind.Observed / Source without needing GitHub
// API access. Mirrors this package's existing fixture-based test idiom
// (LoadFixture + testdata/*.json) but exercises the live-path parse
// specifically, since LoadFixture never reads gamma-scaffold.md.
func TestAssembleLive_ObservesCellKindFromGammaScaffold(t *testing.T) {
	dir := t.TempDir()
	issue := 570
	scaffoldDir := filepath.Join(dir, ".cdd", "unreleased", strconv.Itoa(issue))
	if err := os.MkdirAll(scaffoldDir, 0o755); err != nil {
		t.Fatal(err)
	}
	scaffold := "# γ scaffold — cnos#570\n\n**cell_kind:** `doctrine` (worked example)\n"
	if err := os.WriteFile(filepath.Join(scaffoldDir, "gamma-scaffold.md"), []byte(scaffold), 0o644); err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatal(err)
		}
	}()

	snap, err := assembleLive(context.Background(), "", issue, "")
	if err != nil {
		t.Fatalf("assembleLive: %v", err)
	}
	if snap.CellKind.Observed != "doctrine" {
		t.Errorf("CellKind.Observed = %q, want %q", snap.CellKind.Observed, "doctrine")
	}
	if snap.CellKind.Source != "cdd_artifact" {
		t.Errorf("CellKind.Source = %q, want %q", snap.CellKind.Source, "cdd_artifact")
	}
	// normalizeCellKind must not overwrite an observed value with the
	// absent-default (DefaultedTo stays unset because Observed != "").
	if snap.CellKind.DefaultedTo != "" {
		t.Errorf("CellKind.DefaultedTo = %q, want empty (an observed kind must not also carry a default)", snap.CellKind.DefaultedTo)
	}
}
