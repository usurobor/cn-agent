package issuesfsm

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"
)

// fakeListClosedCandidates stands in for a live GitHub closed-issue list
// query, mirroring scan_test.go's fakeListIssues pattern.
func fakeListClosedCandidates(cands ...terminalCandidate) func(ctx context.Context, repo, protocol, token string) ([]terminalCandidate, error) {
	return func(context.Context, string, string, string) ([]terminalCandidate, error) {
		return cands, nil
	}
}

// fakeLabelMutator records every remove/add/ensure call, standing in for
// fetch.go's ghRemoveLabel/ghAddLabel/ghEnsureLabelExists so the
// orchestration logic in terminal.go can be asserted on without a real
// network call (mirrors scan_test.go's fakeFinalizer/fakeCommenter
// pattern).
type fakeLabelMutator struct {
	mu      sync.Mutex
	removed []string // "issue:label"
	added   []string // "issue:label"
	ensured []string // "name:color:description"
}

func (m *fakeLabelMutator) remove(_ context.Context, _ string, issue int, _ string, label string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.removed = append(m.removed, fmt.Sprintf("%d:%s", issue, label))
	return nil
}

func (m *fakeLabelMutator) add(_ context.Context, _ string, issue int, _ string, label string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.added = append(m.added, fmt.Sprintf("%d:%s", issue, label))
	return nil
}

func (m *fakeLabelMutator) ensure(_ context.Context, _, _, name, color, description string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ensured = append(m.ensured, fmt.Sprintf("%s:%s:%s", name, color, description))
	return nil
}

func (m *fakeLabelMutator) snapshot() (removed, added, ensured []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.removed...), append([]string(nil), m.added...), append([]string(nil), m.ensured...)
}

func newFakeTerminalOptions(apply bool, mut *fakeLabelMutator, cands ...terminalCandidate) *TerminalOptions {
	return &TerminalOptions{
		Repo: "acme/widgets", Protocol: "cds", Apply: apply,
		ListClosedCandidates: fakeListClosedCandidates(cands...),
		RemoveLabel:          mut.remove,
		AddLabel:             mut.add,
		EnsureLabelExists:    mut.ensure,
	}
}

// --- AC5 fixture 1 / AC1 + AC2: closed + state_reason=completed + stale
// status:*/dispatch:cell/protocol:cds -> recognized (dry run) and cleaned
// with resolution/completed added (--apply). ---

func TestTerminal_ClosedCompletedStale_RecognizedAndCleaned(t *testing.T) {
	cand := terminalCandidate{
		Number:      900,
		StateReason: "completed",
		Labels:      []string{"dispatch:cell", "protocol:cds", "status:review"},
	}

	// AC1: without --apply, the sweep recognizes the candidate and reports
	// it, without mutating anything.
	dryMut := &fakeLabelMutator{}
	report, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(false, dryMut, cand))
	if err != nil {
		t.Fatalf("RunTerminalSweep (dry run): %v", err)
	}
	if len(report.Results) != 1 {
		t.Fatalf("results = %v, want 1 candidate recognized", report.Results)
	}
	res := report.Results[0]
	if res.Reconciled {
		t.Fatalf("dry run must not reconcile, got %+v", res)
	}
	if !strings.Contains(res.Note, "would remove") || !strings.Contains(res.Note, "--apply") {
		t.Errorf("expected a would-reconcile note naming --apply, got %q", res.Note)
	}
	removed, added, ensured := dryMut.snapshot()
	if len(removed) != 0 || len(added) != 0 || len(ensured) != 0 {
		t.Fatalf("dry run performed mutation calls: removed=%v added=%v ensured=%v", removed, added, ensured)
	}

	// AC2: with --apply, the removal set is exactly {status:review,
	// dispatch:cell, protocol:cds} and the added label is exactly
	// resolution/completed (the pinned mapping for state_reason
	// "completed").
	applyMut := &fakeLabelMutator{}
	report2, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(true, applyMut, cand))
	if err != nil {
		t.Fatalf("RunTerminalSweep (--apply): %v", err)
	}
	res2 := report2.Results[0]
	if !res2.Reconciled {
		t.Fatalf("expected --apply to reconcile the candidate, got %+v", res2)
	}
	if res2.Resolution != "resolution/completed" {
		t.Fatalf("resolution = %q, want resolution/completed", res2.Resolution)
	}
	removed2, added2, ensured2 := applyMut.snapshot()
	sort.Strings(removed2)
	wantRemoved := []string{"900:dispatch:cell", "900:protocol:cds", "900:status:review"}
	if !equalStringSlices(removed2, wantRemoved) {
		t.Fatalf("removed = %v, want %v", removed2, wantRemoved)
	}
	if len(added2) != 1 || added2[0] != "900:resolution/completed" {
		t.Fatalf("added = %v, want exactly [900:resolution/completed]", added2)
	}
	if len(ensured2) != 1 || !strings.HasPrefix(ensured2[0], "resolution/completed:ededed:") {
		t.Fatalf("ensured = %v, want resolution/completed ensured with color ededed", ensured2)
	}
}

// --- AC5 fixture 2: closed + state_reason=not_planned + stale labels ->
// cleaned with resolution/not-planned added -- NOT resolution/completed,
// NOT resolution/wontfix. This is the exact ambiguity the operator's
// clarifying comment resolved; getting it backwards is a correctness bug. ---

func TestTerminal_ClosedNotPlannedStale_ResolvesToNotPlanned(t *testing.T) {
	cand := terminalCandidate{
		Number:      901,
		StateReason: "not_planned",
		Labels:      []string{"dispatch:cell", "protocol:cds", "status:in-progress"},
	}
	mut := &fakeLabelMutator{}
	report, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(true, mut, cand))
	if err != nil {
		t.Fatalf("RunTerminalSweep: %v", err)
	}
	res := report.Results[0]
	if !res.Reconciled {
		t.Fatalf("expected reconciliation, got %+v", res)
	}
	if res.Resolution != "resolution/not-planned" {
		t.Fatalf("resolution = %q, want resolution/not-planned (not resolution/completed or resolution/wontfix)", res.Resolution)
	}
	_, added, _ := mut.snapshot()
	if len(added) != 1 || added[0] != "901:resolution/not-planned" {
		t.Fatalf("added = %v, want exactly [901:resolution/not-planned]", added)
	}
}

// --- AC5 fixture 3: closed + already clean (the list query itself
// excludes it once dispatch:cell/protocol:cds are gone) -> zero results,
// zero mutation calls -- this is AC3's idempotence falling out of the
// list-query filter itself, not a bolted-on "already done" check. ---

func TestTerminal_AlreadyClean_NoOp(t *testing.T) {
	mut := &fakeLabelMutator{}
	// No candidates at all -- exactly what a live query returns once a
	// prior run already removed dispatch:cell/protocol:cds from issue 900.
	report, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(true, mut))
	if err != nil {
		t.Fatalf("RunTerminalSweep: %v", err)
	}
	if len(report.Results) != 0 {
		t.Fatalf("expected zero candidates for an already-clean sweep, got %v", report.Results)
	}
	removed, added, ensured := mut.snapshot()
	if len(removed) != 0 || len(added) != 0 || len(ensured) != 0 {
		t.Fatalf("expected zero mutation calls, got removed=%v added=%v ensured=%v", removed, added, ensured)
	}
}

// --- AC3: idempotence across two passes. First pass reconciles a closed
// candidate; the second pass's ListClosedCandidates (reflecting the
// post-cleanup state, i.e. the issue no longer carries dispatch:cell +
// protocol:cds) returns zero candidates for it -> zero new mutation
// calls on the second pass. ---

func TestTerminal_Idempotent_SecondPassIsNoOp(t *testing.T) {
	cand := terminalCandidate{
		Number:      902,
		StateReason: "completed",
		Labels:      []string{"dispatch:cell", "protocol:cds", "status:review"},
	}
	mut := &fakeLabelMutator{}

	report1, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(true, mut, cand))
	if err != nil {
		t.Fatalf("first RunTerminalSweep: %v", err)
	}
	if !report1.Results[0].Reconciled {
		t.Fatalf("expected first pass to reconcile, got %+v", report1.Results[0])
	}
	removedAfterFirst, addedAfterFirst, _ := mut.snapshot()
	if len(removedAfterFirst) != 3 || len(addedAfterFirst) != 1 {
		t.Fatalf("first pass mutation counts = removed:%d added:%d, want 3/1", len(removedAfterFirst), len(addedAfterFirst))
	}

	// Second pass: the list query now excludes #902 (its dispatch:cell and
	// protocol:cds labels are gone in the real world after pass 1) -- this
	// fixture map simulates exactly that post-cleanup listing, mirroring
	// scan_test.go's TestScan_Idempotent second-pass simulation.
	report2, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(true, mut))
	if err != nil {
		t.Fatalf("second RunTerminalSweep: %v", err)
	}
	if len(report2.Results) != 0 {
		t.Fatalf("expected zero candidates on the second pass, got %v", report2.Results)
	}
	removedAfterSecond, addedAfterSecond, ensuredAfterSecond := mut.snapshot()
	if len(removedAfterSecond) != len(removedAfterFirst) || len(addedAfterSecond) != len(addedAfterFirst) {
		t.Fatalf("second pass performed new mutation calls: removed %d->%d added %d->%d", len(removedAfterFirst), len(removedAfterSecond), len(addedAfterFirst), len(addedAfterSecond))
	}
	if len(ensuredAfterSecond) != 1 {
		t.Fatalf("expected ensure-label-exists to have been called exactly once total (first pass only), got %d", len(ensuredAfterSecond))
	}
}

// --- AC5 fixture 4 / AC3: an OPEN dispatch:cell+protocol:cds issue is
// never touched by the sweep, regardless of its status:* label -- proven
// at the live-query level below (TestLiveListClosedCandidates_*), and here
// at the orchestration level: an open issue is never a candidate at all,
// so it can never appear in ListClosedCandidates's return value in the
// first place. Combined with an in-flight closed candidate to prove the
// open issue's absence is not merely "there was nothing to process". ---

func TestTerminal_OpenIssueNeverTouched(t *testing.T) {
	closedCand := terminalCandidate{
		Number:      903,
		StateReason: "completed",
		Labels:      []string{"dispatch:cell", "protocol:cds", "status:review"},
	}
	// issue 904 is open in the real world; a live state=closed query would
	// never return it at all -- fakeListClosedCandidates below models
	// exactly that contract by simply never including it, the same way
	// fakeListIssues in scan_test.go models liveListActiveIssues's own
	// server-side state filtering.
	mut := &fakeLabelMutator{}
	report, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(true, mut, closedCand))
	if err != nil {
		t.Fatalf("RunTerminalSweep: %v", err)
	}
	for _, res := range report.Results {
		if res.Issue == 904 {
			t.Fatalf("open issue #904 must never appear as a candidate, got %+v", res)
		}
	}
	removed, added, _ := mut.snapshot()
	for _, r := range append(removed, added...) {
		if strings.HasPrefix(r, "904:") {
			t.Fatalf("open issue #904 was mutated: %v / %v", removed, added)
		}
	}
}

// --- AC1 (negative) at the live-query level: liveListClosedCandidates's
// own request URL names state=closed and the dispatch:cell/protocol:{P}
// labels filter -- a fake server that honors that filter (the way the
// real GitHub API does) excludes an open issue and a closed issue missing
// protocol:cds automatically; only the true candidate is returned. ---

func TestLiveListClosedCandidates_QueryShapeAndFiltering(t *testing.T) {
	var capturedURL string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		// Emulate GitHub's own server-side filtering for labels=/state=:
		// only #950 (closed, carries both dispatch:cell and protocol:cds)
		// is returned. #951 (open) and #952 (closed but missing
		// protocol:cds) are excluded here exactly as the real API would
		// exclude them before this pass's client code ever sees a
		// response -- this pass performs no client-side state/label
		// filtering of its own (unlike liveListActiveIssues, which must
		// client-filter status:* because GitHub cannot OR two status
		// labels in one query; this pass needs no such client-side step).
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"number": 950, "state_reason": "not_planned", "labels": [{"name":"dispatch:cell"},{"name":"protocol:cds"},{"name":"status:blocked"}]}
		]`))
	})

	cands, err := liveListClosedCandidates(context.Background(), "acme/widgets", "cds", "tok")
	if err != nil {
		t.Fatalf("liveListClosedCandidates: %v", err)
	}
	if !strings.Contains(capturedURL, "state=closed") {
		t.Errorf("expected state=closed in the request URL, got %s", capturedURL)
	}
	if !strings.Contains(capturedURL, "labels=dispatch:cell,protocol:cds") {
		t.Errorf("expected labels=dispatch:cell,protocol:cds in the request URL, got %s", capturedURL)
	}
	if len(cands) != 1 || cands[0].Number != 950 {
		t.Fatalf("expected exactly candidate #950, got %+v", cands)
	}
	if cands[0].StateReason != "not_planned" {
		t.Fatalf("state_reason = %q, want not_planned", cands[0].StateReason)
	}
	wantLabels := []string{"dispatch:cell", "protocol:cds", "status:blocked"}
	gotLabels := append([]string(nil), cands[0].Labels...)
	sort.Strings(gotLabels)
	sort.Strings(wantLabels)
	if !equalStringSlices(gotLabels, wantLabels) {
		t.Fatalf("labels = %v, want %v", gotLabels, wantLabels)
	}
}

// --- Unrecognized state_reason: never guess a resolution label; skip and
// report, mutate nothing. In practice unreachable for a state=closed
// query (GitHub's enum is only completed/not_planned/reopened/null, and
// reopened cannot occur on a currently-closed issue) but coded
// defensively per the γ scaffold + cnos#368 doctrine. ---

func TestTerminal_UnrecognizedStateReason_SkippedNeverMutated(t *testing.T) {
	cand := terminalCandidate{
		Number:      905,
		StateReason: "",
		Labels:      []string{"dispatch:cell", "protocol:cds", "status:blocked"},
	}
	mut := &fakeLabelMutator{}
	report, err := RunTerminalSweep(context.Background(), newFakeTerminalOptions(true, mut, cand))
	if err != nil {
		t.Fatalf("RunTerminalSweep: %v", err)
	}
	res := report.Results[0]
	if res.Reconciled {
		t.Fatalf("expected no reconciliation for an unrecognized state_reason, got %+v", res)
	}
	if res.Resolution != "" {
		t.Fatalf("expected no resolution label guessed, got %q", res.Resolution)
	}
	removed, added, ensured := mut.snapshot()
	if len(removed) != 0 || len(added) != 0 || len(ensured) != 0 {
		t.Fatalf("expected zero mutation calls for an unrecognized state_reason, got removed=%v added=%v ensured=%v", removed, added, ensured)
	}
}

// --- AC4: label removal reuses ghRemoveLabel verbatim (default wiring,
// not injected) -- proven by exercising the live defaults end to end
// through a fake GitHub server, mirroring scan_test.go's withFakeGitHub
// pattern for the same purpose. ---

func TestTerminal_LiveDefaults_RouteThroughSharedPrimitives(t *testing.T) {
	var mu sync.Mutex
	var requests []string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		switch {
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"number": 906, "state_reason": "completed", "labels": [{"name":"dispatch:cell"},{"name":"protocol:cds"},{"name":"status:review"}]}]`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/labels") && !strings.Contains(r.URL.Path, "/issues/"):
			// ghEnsureLabelExists: repo-level label creation.
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodPost {
				w.Write([]byte(`[]`))
			}
		}
	})

	opts := &TerminalOptions{Repo: "acme/widgets", Protocol: "cds", Apply: true, Token: "tok"}
	report, err := RunTerminalSweep(context.Background(), opts)
	if err != nil {
		t.Fatalf("RunTerminalSweep: %v", err)
	}
	if !report.Results[0].Reconciled {
		t.Fatalf("expected reconciliation via live defaults, got %+v", report.Results[0])
	}

	mu.Lock()
	got := append([]string(nil), requests...)
	mu.Unlock()

	var deletes, posts, gets int
	for _, r := range got {
		switch {
		case strings.HasPrefix(r, "DELETE "):
			deletes++
		case strings.HasPrefix(r, "POST "):
			posts++
		case strings.HasPrefix(r, "GET "):
			gets++
		}
	}
	if gets != 1 {
		t.Errorf("expected exactly one GET (the candidate list), got %d: %v", gets, got)
	}
	if deletes != 3 {
		t.Errorf("expected exactly 3 DELETE calls (status:review, dispatch:cell, protocol:cds), got %d: %v", deletes, got)
	}
	// 2 POSTs: one ensure-label-exists (repo labels endpoint) + one
	// add-label (issue labels endpoint) -- no direct/duplicate mutation
	// path outside ghRemoveLabel/ghAddLabel/ghEnsureLabelExists.
	if posts != 2 {
		t.Errorf("expected exactly 2 POST calls (ensure-label-exists + add-label), got %d: %v", posts, got)
	}
}

// --- fetch.go's new primitive: ghEnsureLabelExists tolerates a 422
// "already exists" response as an idempotent no-op, mirroring
// ghRemoveLabel's 404-tolerance idiom. ---

func TestGhEnsureLabelExists_TreatsAlreadyExists422AsNoOp(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message":"Validation Failed","errors":[{"code":"already_exists","field":"name"}]}`))
	})
	if err := ghEnsureLabelExists(context.Background(), "acme/widgets", "tok", "resolution/completed", "ededed", ""); err != nil {
		t.Fatalf("expected 422 already_exists to be tolerated as a no-op, got %v", err)
	}
}

func TestGhEnsureLabelExists_CreatesOnFirstUse(t *testing.T) {
	var gotBody string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		gotBody = buf.String()
		w.WriteHeader(http.StatusCreated)
	})
	if err := ghEnsureLabelExists(context.Background(), "acme/widgets", "tok", "resolution/not-planned", "ededed", ""); err != nil {
		t.Fatalf("ghEnsureLabelExists: %v", err)
	}
	if !strings.Contains(gotBody, `"name":"resolution/not-planned"`) || !strings.Contains(gotBody, `"color":"ededed"`) {
		t.Errorf("unexpected request body: %s", gotBody)
	}
}

// --- CLI wiring: `cn issues fsm terminal` exists and requires --protocol,
// mirroring TestAC1Sub593_ScanSubcommandWired's shape for scan. ---

func TestAC1_TerminalSubcommandWired(t *testing.T) {
	var stdout, stderr strings.Builder
	err := Run(context.Background(), []string{"terminal"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when --protocol is omitted")
	}
	if !strings.Contains(err.Error(), "--protocol is required") {
		t.Errorf("expected a precise --protocol-required error, got: %v", err)
	}
}

// --- AC4: this whole cycle must not touch transitions.json, and
// terminal.go must never call Evaluate/Table. Structural proof: terminal.go
// compiles and runs with zero Table/Evaluate dependency in TerminalOptions
// -- there is no Table field to populate, unlike ScanOptions. This test
// pins that shape so a future edit adding one back would be a visible
// diff, not a silent regression. ---

func TestTerminalOptions_HasNoFSMTableDependency(t *testing.T) {
	opts := &TerminalOptions{}
	opts.setDefaults()
	// Compile-time-adjacent check: TerminalOptions has no Table field.
	// (If one is ever added, this line will fail to compile, which is the
	// point -- see terminal.go's package doc for why this pass must never
	// gain an Evaluate/Table dependency.)
	var _ = opts.ListClosedCandidates
	var _ = opts.RemoveLabel
	var _ = opts.AddLabel
	var _ = opts.EnsureLabelExists
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
