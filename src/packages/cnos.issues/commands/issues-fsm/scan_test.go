package issuesfsm

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
)

// fakeAssembleFacts serves pre-loaded FactSnapshots keyed by issue number,
// standing in for a live GitHub/git observation (mirrors cell.Finalizer's
// injectable-dependency testing pattern).
func fakeAssembleFacts(byIssue map[int]FactSnapshot) func(ctx context.Context, repo string, issue int, token string) (FactSnapshot, error) {
	return func(_ context.Context, _ string, issue int, _ string) (FactSnapshot, error) {
		snap, ok := byIssue[issue]
		if !ok {
			return FactSnapshot{}, fmt.Errorf("no fixture registered for issue #%d", issue)
		}
		return snap, nil
	}
}

func fakeListIssues(issues ...int) func(ctx context.Context, repo, protocol, token string) ([]int, error) {
	return func(context.Context, string, string, string) ([]int, error) {
		return issues, nil
	}
}

// fakeFinalizer records every invocation and returns the given stdout/error
// for each issue, letting tests assert scan's finalize-invocation count and
// its "new PR" vs "idempotent no-op" stdout-parsing decision without a real
// subprocess.
type fakeFinalizer struct {
	mu      sync.Mutex
	calls   []int
	stdouts map[int]string
	errs    map[int]error
}

func (f *fakeFinalizer) run(_ context.Context, issue int) (string, error) {
	f.mu.Lock()
	f.calls = append(f.calls, issue)
	f.mu.Unlock()
	return f.stdouts[issue], f.errs[issue]
}

func (f *fakeFinalizer) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

// fakeCommenter records every posted comment.
type fakeCommenter struct {
	mu    sync.Mutex
	posts []int
}

func (c *fakeCommenter) post(_ context.Context, _ string, issue int, _ string) error {
	c.mu.Lock()
	c.posts = append(c.posts, issue)
	c.mu.Unlock()
	return nil
}

func (c *fakeCommenter) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.posts)
}

// --- Case 1 / AC2: a live run is a no-op -------------------------------

func TestScan_LiveRunActive_NoOp(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/in-progress-active.json")
	fin := &fakeFinalizer{}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	if len(report.Results) != 1 {
		t.Fatalf("results = %v, want 1", report.Results)
	}
	res := report.Results[0]
	if res.Outcome != "valid" || res.Reconciled || res.Finalized || res.Commented {
		t.Fatalf("expected a clean no-op for a live run, got %+v", res)
	}
	if fin.callCount() != 0 || com.count() != 0 {
		t.Fatalf("expected zero finalize/comment calls for a live run, got finalize=%d comment=%d", fin.callCount(), com.count())
	}
}

func mustLoadFixture(t *testing.T, path string) FactSnapshot {
	t.Helper()
	snap, err := LoadFixture(path)
	if err != nil {
		t.Fatalf("LoadFixture(%s): %v", path, err)
	}
	return snap
}

// --- Case 2 / "died-before-commits": dead run, no matter -> requeue ----

func TestScan_DeadNoMatter_RequeuesToTodo(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/in-progress-dead-no-matter.json")

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
	fin := &fakeFinalizer{}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	res := report.Results[0]
	if !res.Reconciled || res.TargetState != "todo" {
		t.Fatalf("expected reconciled -> todo, got %+v", res)
	}
	if fin.callCount() != 0 {
		t.Fatalf("expected zero finalize calls for a no-matter dead run, got %d", fin.callCount())
	}
	if com.count() != 1 {
		t.Fatalf("expected exactly one recovery comment, got %d", com.count())
	}
	mu.Lock()
	n := len(requests)
	mu.Unlock()
	if n != 2 {
		t.Fatalf("expected a remove+add label request pair, got %d: %v", n, requests)
	}
}

// --- Case 3 / "died-after-commits-before-PR": dead run, matter, no PR --
// -> finalize creates a NEW draft PR -> comment posted.

func TestScan_DeadWithMatterNoPR_FinalizesAndComments(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/in-progress-dead-with-commits.json")
	fin := &fakeFinalizer{
		stdouts: map[int]string{snap.Issue: "cn cell finalize: matter detected for cycle/603 (issue #603) — checkpointing.\n  ✓ opened draft PR for cycle/603 (Refs #603)\n"},
	}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	res := report.Results[0]
	if res.Action != "propose_delta_recovery" || !res.Finalized || res.Reconciled || res.TargetState != "" {
		t.Fatalf("expected a finalize-only delta-recovery action, got %+v", res)
	}
	if !res.Commented || com.count() != 1 {
		t.Fatalf("expected exactly one comment for a newly-opened PR, got commented=%v count=%d", res.Commented, com.count())
	}
	if fin.callCount() != 1 {
		t.Fatalf("expected exactly one finalize call, got %d", fin.callCount())
	}
}

// --- cnos#630 defensive fallback: propose_delta_recovery's "PR already
// existed" branch (scan.go's else-branch after RunFinalize) is reached
// only when this tick's observed facts said pr_exists==false (routing to
// propose_delta_recovery, not the new propose_status_todo_with_matter
// rule) but a PR already existed by the time RunFinalize actually ran --
// an observe-vs-apply race, not the cnos#630 steady-state wedge (which
// the new transitions.json rule now intercepts before this branch is
// ever reached). This test keeps that defensive branch under coverage
// after cnos#630's redirect of the old
// TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment (which used to
// exercise it via a pr_exists==true fixture that no longer reaches this
// case at all). ---

func TestScan_DeltaRecoveryObserveApplyRace_FinalizeNoOpNoComment(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/in-progress-dead-with-commits.json") // pr_exists==false at observation time
	fin := &fakeFinalizer{
		stdouts: map[int]string{snap.Issue: "cn cell finalize: PR already open for cycle/603: https://github.com/acme/widgets/pull/1 (idempotent no-op; no duplicate created)\n"},
	}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	res := report.Results[0]
	if res.Action != "propose_delta_recovery" || !res.Finalized {
		t.Fatalf("expected finalize to still be invoked (delta-recovery path -- pr_exists was false at observation time), got %+v", res)
	}
	if res.Commented || com.count() != 0 {
		t.Fatalf("expected NO comment for an already-checkpointed matter state (finalize's own idempotent no-op), got commented=%v count=%d", res.Commented, com.count())
	}
}

// --- Case 4 / "died-after-PR-before-REVIEW-REQUEST": dead run, PR
// already exists, no REVIEW-REQUEST.yml. This is the exact "#614 wedge"
// shape (cnos#630): pre-cnos#630, matter already checkpointed into a PR
// with no REVIEW-REQUEST.yml resolved to propose_delta_recovery FOREVER
// (this test used to assert exactly that permanent no-op, under the name
// TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment) -- the cell
// could neither advance (no REVIEW-REQUEST.yml) nor requeue (matter
// exists) nor get re-claimed (dispatch only claims status:todo), and
// stranded permanently pending a human. cnos#630 REDIRECTS this test: it
// now asserts the mechanical exit. See
// TestAC630_WedgePreFixRuleReproducesStrand (issuesfsm_test.go) for the
// preserved proof that the pre-cnos#630 rule shape reproduces the
// strand against this exact fixture -- the "must fail before the fix"
// half of AC2 that a single-branch cycle cannot demonstrate by literally
// running old vs new code.

func TestScan_DeadWithCheckpointedPR_RequeuesToTodoPreservingMatter(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/scan-died-after-pr-before-review-request.json")

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
	fin := &fakeFinalizer{}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	res := report.Results[0]
	if res.Action != "propose_status_todo_with_matter" || !res.Reconciled || res.TargetState != "todo" {
		t.Fatalf("expected reconciled -> todo via propose_status_todo_with_matter (cnos#630 mechanical exit), got %+v", res)
	}
	if fin.callCount() != 0 {
		t.Fatalf("expected zero finalize calls -- matter is already checkpointed, no further checkpointing needed, got %d", fin.callCount())
	}
	if !res.Commented || com.count() != 1 {
		t.Fatalf("expected exactly one audit-note comment (cnos#630 AC4), got commented=%v count=%d", res.Commented, com.count())
	}
	mu.Lock()
	n := len(requests)
	reqs := append([]string(nil), requests...)
	mu.Unlock()
	if n != 2 {
		t.Fatalf("expected a remove(status:in-progress)+add(status:todo) label request pair -- matter (branch/PR) is never deleted, only the label moves, got %d: %v", n, reqs)
	}
}

// --- Case 5 / "died-after-REVIEW-REQUEST-before-status:review": PR +
// REVIEW-REQUEST.yml present -> apply status:review via the FSM.

func TestScan_ReviewRequestWithMatter_AppliesReview(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/scan-died-after-review-request-before-review.json")

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
	fin := &fakeFinalizer{}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	res := report.Results[0]
	if !res.Reconciled || res.TargetState != "review" {
		t.Fatalf("expected reconciled -> review, got %+v", res)
	}
	if fin.callCount() != 0 {
		t.Fatalf("expected zero finalize calls when the label itself is the reconciliation, got %d", fin.callCount())
	}
	if com.count() != 1 {
		t.Fatalf("expected exactly one recovery comment, got %d", com.count())
	}
}

// --- Case 6 / "stale-review-no-proof": status:review with missing
// evidence -> blocked; never mutated, never commented.

func TestScan_StaleReviewNoProof_BlockedNoComment(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/review-empty.json")
	fin := &fakeFinalizer{}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	res := report.Results[0]
	if res.Outcome != "blocked" || res.Reconciled || res.Finalized || res.Commented {
		t.Fatalf("expected a reported-only blocked outcome, got %+v", res)
	}
	if fin.callCount() != 0 || com.count() != 0 {
		t.Fatalf("expected zero side effects for a blocked review, got finalize=%d comment=%d", fin.callCount(), com.count())
	}
}

// --- AC9 / cnos#630 AC6: idempotence. Running the scan twice against the
// SAME unchanged fixture set must not grow the comment/finalize/label
// counts on the second pass for states that already settled (blocked, or
// a dead-run requeue -- whether no-matter or, per cnos#630,
// with-checkpointed-matter) -- and must never double-apply a label
// transition since the transitioned issue is no longer in the active
// scan-state set after the first pass.
//
// cnos#630 updates this test's "checkpointed" fixture expectations: pre-
// cnos#630, a dead run with an already-checkpointed PR (pr_exists==true)
// re-matched propose_delta_recovery forever, so idempotence for that
// fixture meant "finalize is called again every tick, but produces no new
// comment." Post-cnos#630, that same fixture now reconciles -- in-progress
// -> todo, matter preserved -- on its FIRST observed tick (the mechanical
// exit from the wedge), exactly like the deadNoMatter fixture already did;
// idempotence for it is now proven the same way deadNoMatter's already
// was: once reconciled, the issue is no longer status:in-progress, so a
// live re-list would not return it to the scanner on the next pass. ---

func TestScan_Idempotent(t *testing.T) {
	tab := loadRealTable(t)
	deadNoMatter := mustLoadFixture(t, "testdata/in-progress-dead-no-matter.json")
	checkpointed := mustLoadFixture(t, "testdata/scan-died-after-pr-before-review-request.json")
	blocked := mustLoadFixture(t, "testdata/review-empty.json")

	var mu sync.Mutex
	var labelRequests []string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		labelRequests = append(labelRequests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodPost {
			w.Write([]byte(`[]`))
		}
	})

	fin := &fakeFinalizer{}
	com := &fakeCommenter{}

	byIssue := map[int]FactSnapshot{
		deadNoMatter.Issue: deadNoMatter,
		checkpointed.Issue: checkpointed,
		blocked.Issue:      blocked,
	}
	opts := func() *ScanOptions {
		return &ScanOptions{
			Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
			ListActiveIssues: fakeListIssues(deadNoMatter.Issue, checkpointed.Issue, blocked.Issue),
			AssembleFacts:    fakeAssembleFacts(byIssue),
			RunFinalize:      fin.run,
			PostComment:      com.post,
		}
	}

	if _, err := RunScan(context.Background(), opts()); err != nil {
		t.Fatalf("first RunScan: %v", err)
	}
	firstFinalizeCount := fin.callCount()
	firstCommentCount := com.count()
	mu.Lock()
	firstLabelCount := len(labelRequests)
	mu.Unlock()

	if firstCommentCount != 2 { // deadNoMatter's requeue + checkpointed's requeue-with-matter, both cnos#630-audited comments
		t.Fatalf("first pass: expected exactly two comments (both requeues), got %d", firstCommentCount)
	}
	if firstFinalizeCount != 0 { // neither requeue path invokes the finalizer (cnos#630: matter already checkpointed, or no matter to checkpoint)
		t.Fatalf("first pass: expected zero finalize calls, got %d", firstFinalizeCount)
	}

	// Second pass: BOTH deadNoMatter and checkpointed have ALREADY moved
	// to status:todo in the real world after the first pass's
	// reconciliation (this fixture map doesn't mutate, so we simulate the
	// post-reconciliation world by dropping them from the active-issue
	// listing -- exactly what a live re-list would do, since todo is not
	// an active-scan state).
	byIssue2 := map[int]FactSnapshot{
		blocked.Issue: blocked,
	}
	opts2 := &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: true, Table: tab,
		ListActiveIssues: fakeListIssues(blocked.Issue),
		AssembleFacts:    fakeAssembleFacts(byIssue2),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	}
	if _, err := RunScan(context.Background(), opts2); err != nil {
		t.Fatalf("second RunScan: %v", err)
	}

	if com.count() != firstCommentCount {
		t.Fatalf("second pass posted a NEW comment (idempotence violation): before=%d after=%d", firstCommentCount, com.count())
	}
	if fin.callCount() != firstFinalizeCount {
		t.Fatalf("second pass invoked a NEW finalize call (idempotence violation): before=%d after=%d", firstFinalizeCount, fin.callCount())
	}
	mu.Lock()
	secondLabelCount := len(labelRequests) - firstLabelCount
	mu.Unlock()
	if secondLabelCount != 0 {
		t.Fatalf("second pass performed %d NEW label requests (idempotence violation)", secondLabelCount)
	}
}

// --- Dry-run: without --apply, scan reports proposed reconciliations but
// performs zero mutations. ---

func TestScan_NoApply_ReportsWithoutMutating(t *testing.T) {
	tab := loadRealTable(t)
	snap := mustLoadFixture(t, "testdata/in-progress-dead-no-matter.json")

	var mu sync.Mutex
	var requests []string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})
	fin := &fakeFinalizer{}
	com := &fakeCommenter{}

	report, err := RunScan(context.Background(), &ScanOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: false, Table: tab,
		ListActiveIssues: fakeListIssues(snap.Issue),
		AssembleFacts:    fakeAssembleFacts(map[int]FactSnapshot{snap.Issue: snap}),
		RunFinalize:      fin.run,
		PostComment:      com.post,
	})
	if err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	res := report.Results[0]
	if res.Reconciled || res.Commented {
		t.Fatalf("expected a dry-run report with zero mutation, got %+v", res)
	}
	if !strings.Contains(res.Note, "--apply") {
		t.Errorf("expected the dry-run note to mention --apply, got %q", res.Note)
	}
	mu.Lock()
	n := len(requests)
	mu.Unlock()
	if n != 0 {
		t.Fatalf("expected zero label requests without --apply, got %d", n)
	}
	if fin.callCount() != 0 || com.count() != 0 {
		t.Fatalf("expected zero finalize/comment calls without --apply, got finalize=%d comment=%d", fin.callCount(), com.count())
	}
}

// --- AC1: `cn issues fsm scan` exists and is wired into the CLI. ---

func TestAC1Sub593_ScanSubcommandWired(t *testing.T) {
	var stdout, stderr strings.Builder
	err := Run(context.Background(), []string{"scan"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when --protocol is omitted")
	}
	if !strings.Contains(err.Error(), "--protocol is required") {
		t.Errorf("expected a precise --protocol-required error, got: %v", err)
	}
}
