package issuesdispatch

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
)

// fakeIssueStore is an injectable GitHub-API double, mirroring
// issues-fsm's fakeLabelMutator (terminal_test.go) / stateful fixture
// pattern, extended with a body so both halves of the atomic operation
// (body edit + label flip) can be exercised and asserted on without a
// real network call. It is stateful (unlike fakeLabelMutator's pure
// call-recorder) so the AC1 idempotency oracle -- running Dispatch twice
// against the same issue produces the same end state both times -- can be
// proven by actually calling Dispatch twice against the same store.
type fakeIssueStore struct {
	mu      sync.Mutex
	body    string
	labels  []string
	edits   int
	removed []string
	added   []string

	// editErr, when set, makes EditBody fail -- used to prove the
	// body-first ordering safety (AC1's "no half-applied state" oracle).
	editErr error
}

func newFakeIssueStore(body string, labels ...string) *fakeIssueStore {
	return &fakeIssueStore{body: body, labels: append([]string(nil), labels...)}
}

func (s *fakeIssueStore) get(_ context.Context, _ string, _ int, _ string) (issueState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return issueState{Body: s.body, Labels: append([]string(nil), s.labels...)}, nil
}

func (s *fakeIssueStore) editBody(_ context.Context, _ string, _ int, _ string, body string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.edits++
	if s.editErr != nil {
		return s.editErr
	}
	s.body = body
	return nil
}

func (s *fakeIssueStore) removeLabel(_ context.Context, _ string, _ int, _ string, label string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removed = append(s.removed, label)
	out := s.labels[:0:0]
	for _, l := range s.labels {
		if l != label {
			out = append(out, l)
		}
	}
	s.labels = out
	return nil
}

func (s *fakeIssueStore) addLabel(_ context.Context, _ string, _ int, _ string, label string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.added = append(s.added, label)
	s.labels = append(s.labels, label)
	return nil
}

func (s *fakeIssueStore) snapshot() (body string, labels []string, edits int, removed, added []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.body, append([]string(nil), s.labels...), s.edits, append([]string(nil), s.removed...), append([]string(nil), s.added...)
}

func newOptions(store *fakeIssueStore, apply bool) *Options {
	return &Options{
		Repo:        "acme/widgets",
		Issue:       640,
		Apply:       apply,
		GetIssue:    store.get,
		EditBody:    store.editBody,
		RemoveLabel: store.removeLabel,
		AddLabel:    store.addLabel,
	}
}

// --- AC1 positive case: the exact #614/#633 shape -- body carries the
// hold phrase AND the label is already status:todo (the drifted,
// contradictory state: authorization was already given via the label, but
// the body still claims to be held). Dispatch must correct the body
// without re-touching the already-correct label. ---

func TestDispatch_ContradictoryShape_BodyCorrectedLabelUntouched(t *testing.T) {
	store := newFakeIssueStore(contradictoryBody, "dispatch:cell", "protocol:cds", "status:todo")
	res, err := Dispatch(context.Background(), newOptions(store, true))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if !res.BodyChanged {
		t.Fatalf("expected BodyChanged=true for the #614/#633 shape, got %+v", res)
	}
	if res.LabelFlipped {
		t.Fatalf("expected LabelFlipped=false -- label was already status:todo, nothing to flip, got %+v", res)
	}
	if !res.Applied {
		t.Fatalf("expected Applied=true with --apply, got %+v", res)
	}
	body, labels, edits, removed, added := store.snapshot()
	if body != contradictoryBodyWant {
		t.Fatalf("store body after Dispatch = %q, want %q", body, contradictoryBodyWant)
	}
	if edits != 1 {
		t.Fatalf("expected exactly 1 body edit call, got %d", edits)
	}
	if len(removed) != 0 || len(added) != 0 {
		t.Fatalf("expected zero label mutation calls (label was already status:todo), got removed=%v added=%v", removed, added)
	}
	if !contains(labels, "status:todo") {
		t.Fatalf("expected status:todo to remain, got labels=%v", labels)
	}
}

// --- Fresh-dispatch path: status:ready + hold phrase -> both the body
// edit AND the label flip happen, body first. ---

func TestDispatch_ReadyWithHoldPhrase_CleansBodyAndFlipsLabel(t *testing.T) {
	store := newFakeIssueStore(contradictoryBody, "dispatch:cell", "protocol:cds", "status:ready")
	res, err := Dispatch(context.Background(), newOptions(store, true))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if !res.BodyChanged || !res.LabelFlipped {
		t.Fatalf("expected both BodyChanged and LabelFlipped, got %+v", res)
	}
	if res.PreviousStatus != "ready" || res.NewStatus != "todo" {
		t.Fatalf("expected ready->todo, got previous=%q new=%q", res.PreviousStatus, res.NewStatus)
	}
	body, labels, edits, removed, added := store.snapshot()
	if body != contradictoryBodyWant {
		t.Fatalf("store body after Dispatch = %q, want %q", body, contradictoryBodyWant)
	}
	if edits != 1 {
		t.Fatalf("expected exactly 1 body edit call, got %d", edits)
	}
	if len(removed) != 1 || removed[0] != "status:ready" {
		t.Fatalf("expected exactly one remove(status:ready), got %v", removed)
	}
	if len(added) != 1 || added[0] != "status:todo" {
		t.Fatalf("expected exactly one add(status:todo), got %v", added)
	}
	if contains(labels, "status:ready") || !contains(labels, "status:todo") {
		t.Fatalf("expected status:ready removed and status:todo present, got %v", labels)
	}
}

// --- AC1 negative case: an already-clean issue (no hold phrase, already
// status:todo) is a safe no-op -- zero mutation calls of any kind. ---

func TestDispatch_AlreadyCleanAndTodo_SafeNoOp(t *testing.T) {
	clean := "## Design-first change\n\nAlready authorized, already clean."
	store := newFakeIssueStore(clean, "dispatch:cell", "protocol:cds", "status:todo")
	res, err := Dispatch(context.Background(), newOptions(store, true))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if !res.NoOp {
		t.Fatalf("expected NoOp=true, got %+v", res)
	}
	if res.BodyChanged || res.LabelFlipped || res.Applied {
		t.Fatalf("expected no changes at all, got %+v", res)
	}
	body, labels, edits, removed, added := store.snapshot()
	if body != clean {
		t.Fatalf("body must be untouched, got %q", body)
	}
	if edits != 0 || len(removed) != 0 || len(added) != 0 {
		t.Fatalf("expected zero mutation calls, got edits=%d removed=%v added=%v", edits, removed, added)
	}
	if !contains(labels, "status:todo") {
		t.Fatalf("labels should be unchanged, got %v", labels)
	}
}

// --- Future-doctrine path: status:ready with an already-clean body (no
// hold phrase to strip at all, the shape every NEW issue should have once
// AC2's doctrine lands) -- only the label flips; no body-edit call is
// ever made. ---

func TestDispatch_ReadyCleanBody_FlipsLabelOnlyNoBodyCall(t *testing.T) {
	clean := "## Design-first change\n\nNothing to clean here."
	store := newFakeIssueStore(clean, "dispatch:cell", "protocol:cds", "status:ready")
	res, err := Dispatch(context.Background(), newOptions(store, true))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if res.BodyChanged {
		t.Fatalf("expected BodyChanged=false, got %+v", res)
	}
	if !res.LabelFlipped {
		t.Fatalf("expected LabelFlipped=true, got %+v", res)
	}
	_, _, edits, removed, added := store.snapshot()
	if edits != 0 {
		t.Fatalf("expected zero body-edit calls when body is already clean, got %d", edits)
	}
	if len(removed) != 1 || len(added) != 1 {
		t.Fatalf("expected exactly one remove + one add, got removed=%v added=%v", removed, added)
	}
}

// --- AC1 idempotency oracle: running Dispatch twice against the same
// (stateful) issue produces the same end state both times -- the second
// run, observing the post-first-run state, takes no action. ---

func TestDispatch_Idempotent_SecondRunIsNoOp(t *testing.T) {
	store := newFakeIssueStore(contradictoryBody, "dispatch:cell", "protocol:cds", "status:ready")

	res1, err := Dispatch(context.Background(), newOptions(store, true))
	if err != nil {
		t.Fatalf("first Dispatch: %v", err)
	}
	if !res1.Applied {
		t.Fatalf("expected first run to apply, got %+v", res1)
	}
	bodyAfter1, labelsAfter1, editsAfter1, removedAfter1, addedAfter1 := store.snapshot()

	res2, err := Dispatch(context.Background(), newOptions(store, true))
	if err != nil {
		t.Fatalf("second Dispatch: %v", err)
	}
	if !res2.NoOp {
		t.Fatalf("expected second run to be a no-op against the post-first-run state, got %+v", res2)
	}
	bodyAfter2, labelsAfter2, editsAfter2, removedAfter2, addedAfter2 := store.snapshot()

	if bodyAfter1 != bodyAfter2 {
		t.Fatalf("body changed on second run: %q -> %q", bodyAfter1, bodyAfter2)
	}
	if !equalStrSlices(labelsAfter1, labelsAfter2) {
		t.Fatalf("labels changed on second run: %v -> %v", labelsAfter1, labelsAfter2)
	}
	if editsAfter2 != editsAfter1 {
		t.Fatalf("expected no new body-edit calls on the second run, got %d -> %d", editsAfter1, editsAfter2)
	}
	if len(removedAfter2) != len(removedAfter1) || len(addedAfter2) != len(addedAfter1) {
		t.Fatalf("expected no new label mutation calls on the second run, removed %v->%v added %v->%v", removedAfter1, removedAfter2, addedAfter1, addedAfter2)
	}
}

// --- Ordering safety (AC1 "no half-applied state" oracle): when the body
// edit fails, the label is NEVER flipped -- Dispatch must not produce the
// exact contradiction (status:todo + un-cleaned hold phrase) it exists to
// fix. ---

func TestDispatch_BodyEditFails_LabelNeverFlipped(t *testing.T) {
	store := newFakeIssueStore(contradictoryBody, "dispatch:cell", "protocol:cds", "status:ready")
	store.editErr = errors.New("transient github api error")

	_, err := Dispatch(context.Background(), newOptions(store, true))
	if err == nil {
		t.Fatalf("expected an error when the body edit fails")
	}
	if !strings.Contains(err.Error(), "label not touched") {
		t.Errorf("expected the error to name the safety property, got: %v", err)
	}
	_, labels, _, removed, added := store.snapshot()
	if len(removed) != 0 || len(added) != 0 {
		t.Fatalf("expected zero label mutation calls after a body-edit failure, got removed=%v added=%v", removed, added)
	}
	if !contains(labels, "status:ready") || contains(labels, "status:todo") {
		t.Fatalf("expected status:ready to remain untouched (never flipped to status:todo), got %v", labels)
	}
}

// --- Dry run: without --apply, Dispatch reports what it would do and
// makes zero mutation calls. ---

func TestDispatch_DryRun_NoMutationCalls(t *testing.T) {
	store := newFakeIssueStore(contradictoryBody, "dispatch:cell", "protocol:cds", "status:ready")
	res, err := Dispatch(context.Background(), newOptions(store, false))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if res.Applied {
		t.Fatalf("expected Applied=false without --apply, got %+v", res)
	}
	if !res.BodyChanged || !res.LabelFlipped {
		t.Fatalf("expected the dry-run report to still name both pending actions, got %+v", res)
	}
	body, labels, edits, removed, added := store.snapshot()
	if body != contradictoryBody {
		t.Fatalf("dry run must not touch the body, got %q", body)
	}
	if edits != 0 || len(removed) != 0 || len(added) != 0 {
		t.Fatalf("expected zero mutation calls in dry run, got edits=%d removed=%v added=%v", edits, removed, added)
	}
	if !contains(labels, "status:ready") {
		t.Fatalf("dry run must not touch labels, got %v", labels)
	}
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

func equalStrSlices(a, b []string) bool {
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

// --- CLI wiring smoke test: `cn issues dispatch` requires --issue. ---

func TestRun_RequiresIssueFlag(t *testing.T) {
	var stdout, stderr strings.Builder
	err := Run(context.Background(), []string{"--repo", "acme/widgets"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when --issue is omitted")
	}
	if !strings.Contains(err.Error(), "--issue is required") {
		t.Errorf("expected a precise --issue-required error, got: %v", err)
	}
}

func TestRun_RequiresRepoFlag(t *testing.T) {
	t.Setenv("GITHUB_REPOSITORY", "")
	var stdout, stderr strings.Builder
	err := Run(context.Background(), []string{"--issue", "640"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when --repo is omitted and $GITHUB_REPOSITORY is unset")
	}
	if !strings.Contains(err.Error(), "--repo") {
		t.Errorf("expected a precise --repo-required error, got: %v", err)
	}
}
