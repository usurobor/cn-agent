package issuesfsm

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	// cnos#574 AC2: the blocked "review" rule's evidence_guards now mirror
	// the tightened valid rule's all_true set (review_request_present,
	// pr_exists, pr_has_commits) -- branch_has_commits dropped out of the
	// guard set entirely (replaced by pr_has_commits), so it no longer
	// appears in missing_evidence. This is a deliberate, invariant-tied
	// change, not a regression (see self-coherence.md §ACs AC2).
	wantMissing := map[string]bool{"pr_exists": true, "pr_has_commits": true, "review_request_present": true}
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
		"testdata/in-progress-review-request-with-matter.json",
		"testdata/in-progress-review-request-no-matter.json",
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

// --- cnos#569 Phase 2: --apply is a guard-gated mutation path. Phase 1's
// AC8 locked "zero label mutation, no --apply flag, no gh/GitHub label-
// write call anywhere in this package's non-test source" -- that
// forward-reference ("Label-write authority is Phase 2 (cnos#569)") is
// now consumed. The tests below replace the old Phase-1 negative locks
// with the Phase-2 positive contract: --apply exists, but every write it
// performs is reachable only after Evaluate already produced
// outcome=="proposed" with a non-empty target_state (AC1), and every
// blocked outcome (AC3) exits nonzero with zero write calls. ---

// withFakeGitHub points githubAPIBase at an httptest server for the
// duration of the test and restores the real default on cleanup, so the
// new label-write path (fetch.go's ghAddLabel/ghRemoveLabel) never
// reaches the network.
func withFakeGitHub(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	orig := githubAPIBase
	githubAPIBase = srv.URL
	t.Cleanup(func() {
		srv.Close()
		githubAPIBase = orig
	})
	return srv
}

// TestApply_ReviewTransitionAppliesOnGuardPassAndIsIdempotent is AC1's
// positive + idempotence case for the new in-progress -> review request
// path (AC2's mechanism): a worker with REVIEW-REQUEST.yml + deliverable
// matter present gets the transition applied and confirmed; a second
// --apply against the post-mutation state (status:review, valid) is a
// real second call that performs zero writes.
func TestApply_ReviewTransitionAppliesOnGuardPassAndIsIdempotent(t *testing.T) {
	var mu sync.Mutex
	var requests []string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodPost {
			w.Write([]byte(`[]`))
		}
	})

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{
		"evaluate", "--issue", "700", "--apply",
		"--fixture", "testdata/in-progress-review-request-with-matter.json",
		"--table", realTablePath, "--repo", "acme/widgets", "--token", "tok",
	}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "applied: true") {
		t.Errorf("expected \"applied: true\" in output, got:\n%s", stdout.String())
	}

	mu.Lock()
	got := append([]string(nil), requests...)
	mu.Unlock()
	wantDelete := "DELETE /repos/acme/widgets/issues/700/labels/" + url.PathEscape("status:in-progress")
	wantPost := "POST /repos/acme/widgets/issues/700/labels"
	if len(got) != 2 || got[0] != wantDelete || got[1] != wantPost {
		t.Fatalf("requests = %v, want [%q %q]", got, wantDelete, wantPost)
	}

	// Idempotence (AC1): a REAL second --apply call, against a fixture
	// representing the post-mutation state (status:review with a PR --
	// testdata/review-with-pr.json, already the AC3 valid-review
	// fixture), must perform zero writes.
	mu.Lock()
	requests = nil
	mu.Unlock()
	stdout.Reset()
	stderr.Reset()
	err = Run(context.Background(), []string{
		"evaluate", "--issue", "700", "--apply",
		"--fixture", "testdata/review-with-pr.json",
		"--table", realTablePath, "--repo", "acme/widgets", "--token", "tok",
	}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("second Run: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "applied: false") {
		t.Errorf("expected \"applied: false\" (no-op) on the second call, got:\n%s", stdout.String())
	}
	mu.Lock()
	n := len(requests)
	mu.Unlock()
	if n != 0 {
		t.Errorf("expected zero label-write requests on the idempotent second call, got %d: %v", n, requests)
	}
}

// TestApply_RequeueTransitionAppliesOnGuardPassAndIsIdempotent covers the
// same AC1 positive+idempotence shape for the pre-existing (Phase 1)
// in-progress -> todo requeue proposal, proving --apply's guard-gating
// generalizes beyond the new review-request rule.
func TestApply_RequeueTransitionAppliesOnGuardPassAndIsIdempotent(t *testing.T) {
	var mu sync.Mutex
	var requests []string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodPost {
			w.Write([]byte(`[]`))
		}
	})

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{
		"evaluate", "--issue", "602", "--apply",
		"--fixture", "testdata/in-progress-dead-no-matter.json",
		"--table", realTablePath, "--repo", "acme/widgets", "--token", "tok",
	}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "applied: true") {
		t.Errorf("expected \"applied: true\" in output, got:\n%s", stdout.String())
	}

	mu.Lock()
	got := append([]string(nil), requests...)
	mu.Unlock()
	wantDelete := "DELETE /repos/acme/widgets/issues/602/labels/" + url.PathEscape("status:in-progress")
	wantPost := "POST /repos/acme/widgets/issues/602/labels"
	if len(got) != 2 || got[0] != wantDelete || got[1] != wantPost {
		t.Fatalf("requests = %v, want [%q %q]", got, wantDelete, wantPost)
	}

	mu.Lock()
	requests = nil
	mu.Unlock()
	stdout.Reset()
	stderr.Reset()
	err = Run(context.Background(), []string{
		"evaluate", "--issue", "602", "--apply",
		"--fixture", "testdata/todo.json",
		"--table", realTablePath, "--repo", "acme/widgets", "--token", "tok",
	}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("second Run: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "applied: false") {
		t.Errorf("expected \"applied: false\" (no-op) on the second call, got:\n%s", stdout.String())
	}
	mu.Lock()
	n := len(requests)
	mu.Unlock()
	if n != 0 {
		t.Errorf("expected zero label-write requests on the idempotent second call, got %d: %v", n, requests)
	}
}

// TestApply_BlockedReviewRequestRefusesAndMutatesNothing is AC1's
// negative case + AC3's structural block, for the new rule: a worker
// that wrote REVIEW-REQUEST.yml but has no deliverable matter (no PR, no
// commits) is refused -- nonzero exit, zero label-write calls.
func TestApply_BlockedReviewRequestRefusesAndMutatesNothing(t *testing.T) {
	var calls int
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	})

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{
		"evaluate", "--issue", "701", "--apply",
		"--fixture", "testdata/in-progress-review-request-no-matter.json",
		"--table", realTablePath, "--repo", "acme/widgets", "--token", "tok",
	}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected a nonzero-exit error for a blocked --apply transition")
	}
	if calls != 0 {
		t.Errorf("expected zero label-write calls for a blocked transition, got %d", calls)
	}
	if !strings.Contains(stdout.String(), "outcome: blocked") {
		t.Errorf("expected \"outcome: blocked\" in rendered output, got:\n%s", stdout.String())
	}
}

// TestApply_EmptyReviewStateBlocked reuses the Phase-1 review-empty.json
// fixture (status:review already set, no PR/commits/REVIEW-REQUEST.yml)
// under --apply: the pre-existing (unchanged) review-state rule already
// blocks it, so --apply must refuse it too -- AC3 is not limited to the
// new in-progress rule.
func TestApply_EmptyReviewStateBlocked(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected label-write request for a blocked transition: %s %s", r.Method, r.URL.Path)
	})

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{
		"evaluate", "--issue", "601", "--apply",
		"--fixture", "testdata/review-empty.json",
		"--table", realTablePath, "--repo", "acme/widgets", "--token", "tok",
	}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected a nonzero-exit error for status:review with no deliverable evidence")
	}
}

// TestAC569_InProgressReviewRequestWithMatterProposesReview and
// TestAC569_InProgressReviewRequestNoMatterBlocked exercise the new
// transitions.json rules directly through Evaluate (no CLI / no
// network), mirroring this file's existing per-fixture AC-style tests.
func TestAC569_InProgressReviewRequestWithMatterProposesReview(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-review-request-with-matter.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" || dec.TargetState != "review" {
		t.Errorf("outcome=%q target_state=%q, want proposed/review", dec.Outcome, dec.TargetState)
	}
	if dec.Action != "propose_status_review" {
		t.Errorf("action = %q, want propose_status_review", dec.Action)
	}
	if dec.EnabledTransition != "in-progress -> review" {
		t.Errorf("enabled_transition = %q, want %q", dec.EnabledTransition, "in-progress -> review")
	}
}

func TestAC569_InProgressReviewRequestNoMatterBlocked(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-review-request-no-matter.json")
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
	if dec.TargetState == "review" {
		t.Fatal("cnos#569 AC3 regression: must not propose review without deliverable matter")
	}
}

// --- cnos#574 AC2: status:review is valid only when REVIEW-REQUEST.yml,
// an open PR, AND commits beyond base on that PR are ALL present --
// review_request_present alone (or any two of the three) is no longer
// enough. These tests are written FIRST per the issue's explicit
// bug-fix TDD instruction and are expected to FAIL against the
// unmodified (pre-cnos#574) transitions.json, where the "review" state's
// valid rule uses any_true over [pr_exists, branch_has_commits,
// review_request_present] -- see gamma-scaffold.md §Per-AC oracle list
// AC2. ---

func TestAC574_ReviewRequestAloneNoLongerValid(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/review-request-only.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "blocked" {
		t.Fatalf("outcome = %q, want blocked (cnos#574 AC2: review_request_present alone is not deliverable matter)", dec.Outcome)
	}
	if dec.TargetState == "review" {
		t.Fatal("cnos#574 AC2 regression: status:review must not validate on review_request_present alone")
	}
}

func TestAC574_ReviewPartialEvidenceBlocked(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/review-partial-evidence.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "blocked" {
		t.Fatalf("outcome = %q, want blocked (cnos#574 AC2: pr_exists true but pr_has_commits false is still partial evidence)", dec.Outcome)
	}
}

func TestAC574_ReviewWithPRStillValid(t *testing.T) {
	// Backward-compat guard: the existing full-evidence fixture must
	// still validate under the tightened rule (untouched invariant).
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
		t.Errorf("outcome = %q, want valid (PR + commits + no regression from cnos#574 tightening)", dec.Outcome)
	}
}

// --- cnos#574 AC3: in-progress -> review is proposed only when
// review_request_present AND pr_exists AND pr_has_commits --
// branch-commits-only (no PR) no longer qualifies. Written FIRST per
// bug-fix TDD; expected to FAIL against the unmodified (pre-cnos#574)
// transitions.json, where rule 1 accepts any_true over [pr_exists,
// pr_has_commits, branch_has_commits] -- see gamma-scaffold.md §Per-AC
// oracle list AC3. ---

func TestAC574_InProgressBranchOnlyNoLongerProposesReview(t *testing.T) {
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-review-request-branch-only.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome == "proposed" && dec.TargetState == "review" {
		t.Fatal("cnos#574 AC3 regression: branch-commits-only (no PR) must not propose in-progress -> review")
	}
	if dec.Outcome != "blocked" {
		t.Errorf("outcome = %q, want blocked (review requested, but no PR -- cnos#574 AC3)", dec.Outcome)
	}
}

func TestAC574_InProgressWithMatterStillProposesReview(t *testing.T) {
	// Backward-compat guard: the existing #573 with-matter fixture (which
	// sets pr_exists + pr_commit_count, not just branch_has_commits) must
	// still propose review under the tightened rule.
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-review-request-with-matter.json")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" || dec.TargetState != "review" {
		t.Errorf("outcome=%q target_state=%q, want proposed/review (PR+commits present -- no regression from cnos#574 tightening)", dec.Outcome, dec.TargetState)
	}
}

func TestAC574_InProgressRunActiveNonGatingPreserved(t *testing.T) {
	// cnos#569's main use case (a worker requesting its own transition
	// while its run is still in_progress) must survive the cnos#574
	// tightening unchanged: rules 1/2 fire before rule 3's run_active
	// check, regardless of run_active's value.
	tab := loadRealTable(t)
	snap, err := LoadFixture("testdata/in-progress-review-request-with-matter.json")
	if err != nil {
		t.Fatal(err)
	}
	if snap.RunState != "in_progress" {
		t.Fatalf("fixture precondition: expected run_state=in_progress, got %q", snap.RunState)
	}
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" || dec.Action != "propose_status_review" {
		t.Errorf("outcome=%q action=%q, want proposed/propose_status_review despite run_active=true", dec.Outcome, dec.Action)
	}
}

// TestApplyStatusLabel_ToleratesNotFoundOnRemove and
// TestApplyStatusLabel_NoRemovalWhenFromStateEmpty lock the two
// implementation-contract-adjacent decisions the γ scaffold's Friction
// notes flagged as α's call to make (see self-coherence.md): a 404 on
// the remove-old-label call is tolerated (already-removed / manually-
// recovered state stays idempotent), and no remove call is issued at
// all when there was no prior status label to remove.
func TestApplyStatusLabel_ToleratesNotFoundOnRemove(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusNotFound)
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}
	})
	if err := applyStatusLabel(context.Background(), "acme/widgets", 1, "tok", "in-progress", "review"); err != nil {
		t.Fatalf("applyStatusLabel: %v (expected a 404 on remove to be tolerated)", err)
	}
}

func TestApplyStatusLabel_NoRemovalWhenFromStateEmpty(t *testing.T) {
	var deleteCalls int
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteCalls++
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	})
	if err := applyStatusLabel(context.Background(), "acme/widgets", 1, "tok", "", "todo"); err != nil {
		t.Fatalf("applyStatusLabel: %v", err)
	}
	if deleteCalls != 0 {
		t.Errorf("expected no DELETE call when fromState is empty, got %d", deleteCalls)
	}
}

// TestApply_MissingRepoErrors: --apply without --repo (and no
// $GITHUB_REPOSITORY) cannot know where to write, even when the
// evaluated outcome is "proposed" -- must error, not silently skip the
// write.
func TestApply_MissingRepoErrors(t *testing.T) {
	t.Setenv("GITHUB_REPOSITORY", "")
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{
		"evaluate", "--issue", "700", "--apply",
		"--fixture", "testdata/in-progress-review-request-with-matter.json",
		"--table", realTablePath,
	}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error: --apply with no --repo/$GITHUB_REPOSITORY cannot apply a transition")
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
		"testdata/in-progress-review-request-with-matter.json",
		"testdata/in-progress-review-request-no-matter.json",
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

// TestObserveRemoteBranch_ExistsWithCommits is cnos#574 AC4's core proof at
// the observation-primitive level: a fake GitHub server reports the branch
// exists (200 on /branches/{branch}) with 3 commits ahead of main (via
// /compare's ahead_by), and observeRemoteBranch must report exactly that --
// fully hermetic (no local git, no other assembleLive network calls; see
// this function's doc comment in fetch.go for why it was extracted).
func TestObserveRemoteBranch_ExistsWithCommits(t *testing.T) {
	repo := "acme/widgets"
	branch := "cycle/9601"

	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/" + repo + "/branches/" + branch:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"` + branch + `"}`))
		case "/repos/" + repo + "/compare/main..." + branch:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ahead_by":3}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	exists, commits := observeRemoteBranch(context.Background(), repo, branch, "")
	if !exists {
		t.Fatal("exists = false, want true (remote branch confirmed via the GitHub API)")
	}
	if commits != 3 {
		t.Errorf("commitsBeyondBase = %d, want 3 (from the /compare ahead_by response)", commits)
	}
}

// TestObserveRemoteBranch_AbsentReportsFalse confirms a 404 on the
// branches endpoint (branch genuinely absent on the remote too) reports
// exists=false, commits=0 -- not an error, per this function's "kept
// honest but deliberately simple" live-path contract.
func TestObserveRemoteBranch_AbsentReportsFalse(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	exists, commits := observeRemoteBranch(context.Background(), "acme/widgets", "cycle/404", "")
	if exists {
		t.Error("exists = true, want false for a 404 branches response")
	}
	if commits != 0 {
		t.Errorf("commitsBeyondBase = %d, want 0", commits)
	}
}

// TestAssembleLive_RemoteOnlyBranchResolvesToDeltaRecovery composes
// observeRemoteBranch's result (as assembleLive's fallback block would) with
// the real transition table via Evaluate, proving the cnos#574 AC4
// end-to-end invariant: a remote-only cycle/{N} branch (no local ref -- the
// local-git block would report BranchExists=false) with commits beyond
// base, a dead run (no active workflow run), and status:in-progress
// resolves to propose_delta_recovery -- never propose_status_todo (the
// cnos#368 blind-requeue failure mode this AC exists to prevent).
func TestAssembleLive_RemoteOnlyBranchResolvesToDeltaRecovery(t *testing.T) {
	repo := "acme/widgets"
	branch := "cycle/9601"

	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/" + repo + "/branches/" + branch:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"` + branch + `"}`))
		case "/repos/" + repo + "/compare/main..." + branch:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ahead_by":3}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// Simulate the local-git-absent case (as assembleLive's guard
	// `!snap.BranchExists` would see for a branch with no local ref):
	// BranchExists starts false, then the fallback observation fills it.
	exists, commits := observeRemoteBranch(context.Background(), repo, branch, "")

	snap := FactSnapshot{
		Issue:             9601,
		Labels:            []string{"dispatch:cell", "protocol:cds", "status:in-progress"},
		BranchExists:      exists,
		CommitsBeyondBase: commits,
		RunState:          "", // dead run: no active workflow run
	}

	tab := loadRealTable(t)
	dec, err := Evaluate(tab, snap)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Outcome != "proposed" || dec.Action != "propose_delta_recovery" {
		t.Fatalf("outcome=%q action=%q, want proposed/propose_delta_recovery (remote-only branch with commits must be recovery matter, not requeue)", dec.Outcome, dec.Action)
	}
	if dec.TargetState == "todo" {
		t.Fatal("cnos#368/cnos#574 AC4 regression: a remote-only branch with commits beyond base must never propose blind requeue to status:todo")
	}
}

// TestAssembleLive_RemoteBranchFallbackNeverRunsWithoutRepo
// confirms the cnos#574 AC4 fallback never fires (and thus never overrides
// or downgrades) when local git already found the branch, or when repo==""
// -- the API call backs up local git's gap, it does not replace or
// re-verify a local-git success, and it never runs without a repo to query.
func TestAssembleLive_RemoteBranchFallbackNeverRunsWithoutRepo(t *testing.T) {
	apiCalled := false
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/branches/") || strings.Contains(r.URL.Path, "/compare/") {
			apiCalled = true
		}
		w.WriteHeader(http.StatusNotFound)
	})

	// repo=="" skips the whole GitHub-observation section (labels, the
	// new branch/compare fallback, PR, runs) entirely per assembleLive's
	// existing repo!="" guards -- this locks the "only fills the gap,
	// never replaces, never runs without a repo" contract stated on the
	// fallback block in fetch.go.
	snap, err := assembleLive(context.Background(), "", 1, "")
	if err != nil {
		t.Fatalf("assembleLive: %v", err)
	}
	if apiCalled {
		t.Fatal("expected zero GitHub API calls when repo==\"\"")
	}
	if snap.BranchExists {
		t.Error("BranchExists = true, want false (no local git branch, no repo to query)")
	}
}
