package cell

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- ParseReturnArgs tests ---

func TestParseReturnArgs_Valid(t *testing.T) {
	args, err := ParseReturnArgs([]string{"--issue", "500", "--verdict", "iterate", "--review", "/tmp/review.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if args.Issue != 500 {
		t.Errorf("Issue: want 500, got %d", args.Issue)
	}
	if args.Verdict != "iterate" {
		t.Errorf("Verdict: want iterate, got %q", args.Verdict)
	}
	if args.ReviewPath != "/tmp/review.md" {
		t.Errorf("ReviewPath: want /tmp/review.md, got %q", args.ReviewPath)
	}
}

func TestParseReturnArgs_AllVerdicts(t *testing.T) {
	for _, v := range []string{"converge", "iterate", "reject", "clarify"} {
		args, err := ParseReturnArgs([]string{"--issue", "1", "--verdict", v, "--review", "x.md"})
		if err != nil {
			t.Errorf("verdict %q: unexpected error: %v", v, err)
		}
		if args.Verdict != v {
			t.Errorf("verdict %q: got %q", v, args.Verdict)
		}
	}
}

func TestParseReturnArgs_InvalidVerdict(t *testing.T) {
	_, err := ParseReturnArgs([]string{"--issue", "1", "--verdict", "bogus", "--review", "x.md"})
	if err == nil {
		t.Error("expected error for invalid verdict")
	}
}

func TestParseReturnArgs_MissingIssue(t *testing.T) {
	_, err := ParseReturnArgs([]string{"--verdict", "iterate", "--review", "x.md"})
	if err == nil {
		t.Error("expected error for missing --issue")
	}
}

func TestParseReturnArgs_MissingVerdict(t *testing.T) {
	_, err := ParseReturnArgs([]string{"--issue", "1", "--review", "x.md"})
	if err == nil {
		t.Error("expected error for missing --verdict")
	}
}

func TestParseReturnArgs_MissingReview(t *testing.T) {
	_, err := ParseReturnArgs([]string{"--issue", "1", "--verdict", "iterate"})
	if err == nil {
		t.Error("expected error for missing --review")
	}
}

func TestParseReturnArgs_NonPositiveIssue(t *testing.T) {
	_, err := ParseReturnArgs([]string{"--issue", "0", "--verdict", "iterate", "--review", "x.md"})
	if err == nil {
		t.Error("expected error for issue=0")
	}
}

func TestParseReturnArgs_UnknownFlag(t *testing.T) {
	_, err := ParseReturnArgs([]string{"--issue", "1", "--verdict", "iterate", "--review", "x.md", "--extra", "foo"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}

// --- ParseResumeArgs tests ---

func TestParseResumeArgs_Valid(t *testing.T) {
	args, err := ParseResumeArgs([]string{"--issue", "500"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if args.Issue != 500 {
		t.Errorf("Issue: want 500, got %d", args.Issue)
	}
}

func TestParseResumeArgs_MissingIssue(t *testing.T) {
	_, err := ParseResumeArgs([]string{})
	if err == nil {
		t.Error("expected error for missing --issue")
	}
}

func TestParseResumeArgs_UnknownFlag(t *testing.T) {
	_, err := ParseResumeArgs([]string{"--issue", "1", "--extra", "foo"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}

// --- readSchemaField tests ---

func TestReadSchemaField_Present(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\n---\n\nbody text\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	schema, err := readSchemaField(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if schema != "cn.operator-review.v1" {
		t.Errorf("schema: want cn.operator-review.v1, got %q", schema)
	}
}

func TestReadSchemaField_Missing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "review.md")
	content := "---\nissue: 500\n---\n\nbody text\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := readSchemaField(path)
	if err == nil {
		t.Error("expected error for missing schema field")
	}
}

func TestReadSchemaField_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "review.md")
	content := "# No frontmatter\n\nsome text\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := readSchemaField(path)
	if err == nil {
		t.Error("expected error for no frontmatter")
	}
}

func TestReadSchemaField_FileNotFound(t *testing.T) {
	_, err := readSchemaField("/nonexistent/path/review.md")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// --- nextRoundNumber tests ---

func TestNextRoundNumber_NoFile(t *testing.T) {
	n, err := nextRoundNumber("/nonexistent/self-coherence.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("want 1 (no file), got %d", n)
	}
}

func TestNextRoundNumber_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "self-coherence.md")
	if err := os.WriteFile(path, []byte("# Self-coherence\n\n## §Gap\n\nsome gap text\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := nextRoundNumber(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("want 1 (no §R sections), got %d", n)
	}
}

func TestNextRoundNumber_OneRound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "self-coherence.md")
	content := "# Self-coherence\n\n## §R0\n\nsome text\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := nextRoundNumber(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("want 1 (after R0), got %d", n)
	}
}

func TestNextRoundNumber_MultipleRounds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "self-coherence.md")
	content := "# Self-coherence\n\n## §R0\n\ntext\n\n## §R1\n\nmore text\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := nextRoundNumber(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("want 2 (after R0+R1), got %d", n)
	}
}

// --- appendRoundHeader tests ---

func TestAppendRoundHeader_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "self-coherence.md")
	if err := appendRoundHeader(path, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)
	if !strings.Contains(s, "## §R1") {
		t.Errorf("expected §R1 in output, got: %q", s)
	}
}

func TestAppendRoundHeader_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "self-coherence.md")
	existing := "## §R0\n\nsome existing text\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := appendRoundHeader(path, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)
	if !strings.Contains(s, "## §R0") {
		t.Error("existing §R0 must be preserved")
	}
	if !strings.Contains(s, "## §R1") {
		t.Error("new §R1 must be appended")
	}
}

// --- Smoke: Returner.Return with converge verdict (no gh call) ---

func TestReturner_Return_Converge(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: converge\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout strings.Builder
	var stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr}
	args := ReturnArgs{Issue: 500, Verdict: "converge", ReviewPath: reviewPath}
	if err := r.Return(t.Context(), args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout.String(), "converge") {
		t.Errorf("expected 'converge' in output, got: %q", stdout.String())
	}
	// Converge does NOT call gh, so no label mutation can occur.
	// Verified by absence of "status:changes" in stdout.
	if strings.Contains(stdout.String(), "status:changes") {
		t.Error("converge verdict must not apply status:changes transition")
	}
}

// --- Smoke: Returner.Return with iterate/reject verdicts (injected RunGH + RunGHJSON) ---

// mockGHJSON_HealthyPreflight returns a RunGHJSON that simulates a healthy
// preflight: issue OPEN, single status:review label; repo has status:changes
// in its label list. F2/F3 (cycle/500 R1) preflight path.
func mockGHJSON_HealthyPreflight() func(context.Context, []string) ([]byte, error) {
	return func(_ context.Context, args []string) ([]byte, error) {
		joined := strings.Join(args, " ")
		switch {
		case strings.Contains(joined, "issue view"):
			return []byte(`{"state":"OPEN","labels":[{"name":"status:review"},{"name":"dispatch:cell"}]}`), nil
		case strings.Contains(joined, "label list"):
			return []byte(`[{"name":"status:review"},{"name":"status:changes"},{"name":"status:in-progress"}]`), nil
		default:
			return nil, fmt.Errorf("unexpected gh json call: %s", joined)
		}
	}
}

func TestReturner_Return_Iterate_AppliesLabelTransition(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var called [][]string
	mockGH := func(_ context.Context, args []string, _ io.Writer) error {
		called = append(called, append([]string(nil), args...))
		return nil
	}

	var stdout strings.Builder
	var stderr strings.Builder
	r := &Returner{
		Stdout:    &stdout,
		Stderr:    &stderr,
		RunGH:     mockGH,
		RunGHJSON: mockGHJSON_HealthyPreflight(),
	}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	if err := r.Return(t.Context(), args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// F3 (cycle/500 R1) atomicity: one gh call with both --remove-label and
	// --add-label flags (not two separate calls).
	if len(called) != 1 {
		t.Fatalf("expected 1 atomic gh call, got %d: %v", len(called), called)
	}
	editCall := strings.Join(called[0], " ")
	if !strings.Contains(editCall, "--remove-label") || !strings.Contains(editCall, "status:review") {
		t.Errorf("gh call should remove status:review, got: %q", editCall)
	}
	if !strings.Contains(editCall, "--add-label") || !strings.Contains(editCall, "status:changes") {
		t.Errorf("gh call should add status:changes, got: %q", editCall)
	}
	if !strings.Contains(stdout.String(), "status:changes") {
		t.Errorf("expected 'status:changes' in output, got: %q", stdout.String())
	}
}

func TestReturner_Return_Reject_AppliesLabelTransition(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: reject\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var called [][]string
	mockGH := func(_ context.Context, args []string, _ io.Writer) error {
		called = append(called, append([]string(nil), args...))
		return nil
	}

	var stdout strings.Builder
	var stderr strings.Builder
	r := &Returner{
		Stdout:    &stdout,
		Stderr:    &stderr,
		RunGH:     mockGH,
		RunGHJSON: mockGHJSON_HealthyPreflight(),
	}
	args := ReturnArgs{Issue: 500, Verdict: "reject", ReviewPath: reviewPath}
	if err := r.Return(t.Context(), args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// One atomic call (F3).
	if len(called) != 1 {
		t.Fatalf("expected 1 atomic gh call, got %d: %v", len(called), called)
	}
}

// --- F2: preflight issue state (cycle/500 R1) ---

// TestReturner_Return_Preflight_WrongStatusLabel proves an issue at a
// status:* label other than status:review is refused. F2 invariant: exactly
// status:review; nothing else.
func TestReturner_Return_Preflight_WrongStatusLabel(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	var called [][]string
	mockGH := func(_ context.Context, args []string, _ io.Writer) error {
		called = append(called, args)
		return nil
	}
	mockJSON := func(_ context.Context, args []string) ([]byte, error) {
		if strings.Contains(strings.Join(args, " "), "issue view") {
			return []byte(`{"state":"OPEN","labels":[{"name":"status:in-progress"}]}`), nil
		}
		return []byte(`[{"name":"status:changes"}]`), nil
	}
	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr, RunGH: mockGH, RunGHJSON: mockJSON}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error for wrong status label")
	}
	if !strings.Contains(err.Error(), "review_return_state_invalid") {
		t.Errorf("expected review_return_state_invalid, got: %v", err)
	}
	if len(called) != 0 {
		t.Errorf("preflight must reject before any label mutation; got %d gh calls", len(called))
	}
}

// TestReturner_Return_Preflight_MultipleStatusLabels proves an issue with
// more than one status:* label is refused.
func TestReturner_Return_Preflight_MultipleStatusLabels(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	mockGH := func(_ context.Context, _ []string, _ io.Writer) error { return nil }
	mockJSON := func(_ context.Context, args []string) ([]byte, error) {
		if strings.Contains(strings.Join(args, " "), "issue view") {
			return []byte(`{"state":"OPEN","labels":[{"name":"status:review"},{"name":"status:changes"}]}`), nil
		}
		return []byte(`[{"name":"status:changes"}]`), nil
	}
	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr, RunGH: mockGH, RunGHJSON: mockJSON}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error for multiple status labels")
	}
	if !strings.Contains(err.Error(), "review_return_state_invalid") {
		t.Errorf("expected review_return_state_invalid, got: %v", err)
	}
}

// TestReturner_Return_Preflight_ClosedIssue proves a closed issue is
// refused.
func TestReturner_Return_Preflight_ClosedIssue(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	mockGH := func(_ context.Context, _ []string, _ io.Writer) error { return nil }
	mockJSON := func(_ context.Context, args []string) ([]byte, error) {
		if strings.Contains(strings.Join(args, " "), "issue view") {
			return []byte(`{"state":"CLOSED","labels":[{"name":"status:review"}]}`), nil
		}
		return []byte(`[{"name":"status:changes"}]`), nil
	}
	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr, RunGH: mockGH, RunGHJSON: mockJSON}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error for closed issue")
	}
	if !strings.Contains(err.Error(), "review_return_state_invalid") {
		t.Errorf("expected review_return_state_invalid, got: %v", err)
	}
}

// --- F3: atomic label transition + drift handling (cycle/500 R1) ---

// TestReturner_Return_AtomicTransition_OneGHCall proves the transition is
// a single gh issue edit call carrying both --remove-label and --add-label
// (preserving the single-status invariant via one API request, not two).
func TestReturner_Return_AtomicTransition_OneGHCall(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	var called [][]string
	mockGH := func(_ context.Context, args []string, _ io.Writer) error {
		called = append(called, args)
		return nil
	}
	var stdout, stderr strings.Builder
	r := &Returner{
		Stdout: &stdout, Stderr: &stderr,
		RunGH:     mockGH,
		RunGHJSON: mockGHJSON_HealthyPreflight(),
	}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	if err := r.Return(t.Context(), args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(called) != 1 {
		t.Fatalf("expected 1 atomic gh edit call, got %d: %v", len(called), called)
	}
	call := strings.Join(called[0], " ")
	if !strings.Contains(call, "--remove-label status:review") {
		t.Errorf("atomic call must remove status:review, got: %q", call)
	}
	if !strings.Contains(call, "--add-label status:changes") {
		t.Errorf("atomic call must add status:changes, got: %q", call)
	}
}

// TestReturner_Return_TargetLabelMissing proves the transition refuses
// when status:changes is not defined in the repo label set, BEFORE any
// destructive label mutation. F3 bonus consideration (cnos#493 empirical
// witness: the label-doctor gap left target labels missing TODAY during
// the cnos#500 bootstrap recovery).
func TestReturner_Return_TargetLabelMissing(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	var called [][]string
	mockGH := func(_ context.Context, args []string, _ io.Writer) error {
		called = append(called, args)
		return nil
	}
	mockJSON := func(_ context.Context, args []string) ([]byte, error) {
		joined := strings.Join(args, " ")
		if strings.Contains(joined, "issue view") {
			return []byte(`{"state":"OPEN","labels":[{"name":"status:review"}]}`), nil
		}
		// label list: status:changes is MISSING from the repo.
		return []byte(`[{"name":"status:review"},{"name":"status:in-progress"}]`), nil
	}
	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr, RunGH: mockGH, RunGHJSON: mockJSON}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error when status:changes is missing from repo")
	}
	if !strings.Contains(err.Error(), "review_return_target_label_missing") {
		t.Errorf("expected review_return_target_label_missing, got: %v", err)
	}
	if len(called) != 0 {
		t.Errorf("must refuse before destructive label mutation; got %d gh calls", len(called))
	}
}

// TestReturner_Return_LabelDrift_OnGHFailure proves the drift assessment
// reports a precise marker when the atomic gh call fails and the issue is
// post-hoc statusless.
func TestReturner_Return_LabelDrift_OnGHFailure(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	mockGH := func(_ context.Context, _ []string, _ io.Writer) error {
		return fmt.Errorf("simulated gh failure")
	}
	calls := 0
	mockJSON := func(_ context.Context, args []string) ([]byte, error) {
		joined := strings.Join(args, " ")
		calls++
		if strings.Contains(joined, "label list") {
			return []byte(`[{"name":"status:review"},{"name":"status:changes"}]`), nil
		}
		// First issue view: preflight (healthy). Second: drift assessment
		// (statusless after the failed transition).
		if strings.Contains(joined, "issue view") {
			if calls <= 2 {
				return []byte(`{"state":"OPEN","labels":[{"name":"status:review"}]}`), nil
			}
			return []byte(`{"state":"OPEN","labels":[]}`), nil
		}
		return nil, fmt.Errorf("unexpected: %s", joined)
	}
	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr, RunGH: mockGH, RunGHJSON: mockJSON}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error when gh fails")
	}
	if !strings.Contains(err.Error(), "review_return_label_drift") {
		t.Errorf("expected review_return_label_drift in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "statusless") {
		t.Errorf("expected statusless marker in drift report, got: %v", err)
	}
}

// --- F1: artifact-as-authority validation (issue + verdict must match flags) ---

// TestReturner_Return_IssueMismatch_RejectsBeforeLabelMutation proves an
// artifact whose issue field does NOT match --issue is rejected before any
// gh call is made. F1 (cycle/500 R1): the artifact is the authority; CLI
// flags select and confirm it.
func TestReturner_Return_IssueMismatch_RejectsBeforeLabelMutation(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	// Artifact says issue: 497 — caller passes --issue 500.
	content := "---\nschema: cn.operator-review.v1\nissue: 497\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var called [][]string
	mockGH := func(_ context.Context, args []string, _ io.Writer) error {
		called = append(called, append([]string(nil), args...))
		return nil
	}

	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr, RunGH: mockGH}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error for issue mismatch")
	}
	if !strings.Contains(err.Error(), "review_return_artifact_mismatch") {
		t.Errorf("expected review_return_artifact_mismatch error, got: %v", err)
	}
	if len(called) != 0 {
		t.Errorf("no gh call must occur before mismatch is detected; got %d calls: %v", len(called), called)
	}
}

// TestReturner_Return_VerdictMismatch_RejectsBeforeLabelMutation proves an
// artifact whose verdict field does NOT match --verdict is rejected before
// any gh call is made.
func TestReturner_Return_VerdictMismatch_RejectsBeforeLabelMutation(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	// Artifact says verdict: converge — caller passes --verdict iterate.
	content := "---\nschema: cn.operator-review.v1\nissue: 500\nverdict: converge\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var called [][]string
	mockGH := func(_ context.Context, args []string, _ io.Writer) error {
		called = append(called, append([]string(nil), args...))
		return nil
	}

	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr, RunGH: mockGH}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error for verdict mismatch")
	}
	if !strings.Contains(err.Error(), "review_return_artifact_mismatch") {
		t.Errorf("expected review_return_artifact_mismatch error, got: %v", err)
	}
	if len(called) != 0 {
		t.Errorf("no gh call must occur before mismatch is detected; got %d calls: %v", len(called), called)
	}
}

// TestReturner_Return_MissingIssueField_Rejects proves an artifact missing
// the `issue:` frontmatter field is rejected.
func TestReturner_Return_MissingIssueField_Rejects(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "operator-review.md")
	content := "---\nschema: cn.operator-review.v1\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Fatal("expected error for missing issue field")
	}
	if !strings.Contains(err.Error(), "review_return_artifact_invalid") {
		t.Errorf("expected review_return_artifact_invalid error, got: %v", err)
	}
}

// --- Smoke: Returner.Return with wrong schema ---

func TestReturner_Return_WrongSchema(t *testing.T) {
	dir := t.TempDir()
	reviewPath := filepath.Join(dir, "bad-review.md")
	content := "---\nschema: cn.something-else.v1\nissue: 500\nverdict: iterate\n---\n\nbody\n"
	if err := os.WriteFile(reviewPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout strings.Builder
	var stderr strings.Builder
	r := &Returner{Stdout: &stdout, Stderr: &stderr}
	args := ReturnArgs{Issue: 500, Verdict: "iterate", ReviewPath: reviewPath}
	err := r.Return(t.Context(), args)
	if err == nil {
		t.Error("expected error for wrong schema")
	}
	if !strings.Contains(err.Error(), "unexpected schema") {
		t.Errorf("expected 'unexpected schema' in error, got: %v", err)
	}
}

// --- F4: cn cell resume cycle-branch preflight (cycle/500 R1) ---

// TestResumer_Resume_RefusesWhenOnMain proves the F4 preflight refuses
// to operate when the caller is not on cycle/{N}. Previously the code
// would have appended §R[N+1] to whatever working tree happened to be
// checked out — corrupting the wrong branch.
func TestResumer_Resume_RefusesWhenOnMain(t *testing.T) {
	dir := t.TempDir()
	artifactDir := filepath.Join(dir, ".cdd", "unreleased", "500")
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	scPath := filepath.Join(artifactDir, "self-coherence.md")
	if err := os.WriteFile(scPath, []byte("## §R0\n\ntext\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr strings.Builder
	r := &Resumer{
		RepoRoot: dir,
		Stdout:   &stdout,
		Stderr:   &stderr,
		CurrentBranch: func(_ context.Context) (string, error) {
			return "main", nil
		},
	}
	err := r.Resume(t.Context(), ResumeArgs{Issue: 500})
	if err == nil {
		t.Fatal("expected refusal when on main")
	}
	if !strings.Contains(err.Error(), "review_resume_wrong_branch") {
		t.Errorf("expected review_resume_wrong_branch, got: %v", err)
	}
	// self-coherence.md must NOT have been modified.
	content, _ := os.ReadFile(scPath)
	if strings.Contains(string(content), "§R1") {
		t.Error("preflight failure must not append §R1 to self-coherence.md")
	}
}

// TestResumer_Resume_RefusesOnDetachedHead proves the preflight refuses
// on a detached HEAD.
func TestResumer_Resume_RefusesOnDetachedHead(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".cdd", "unreleased", "500"), 0o755); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr strings.Builder
	r := &Resumer{
		RepoRoot: dir,
		Stdout:   &stdout,
		Stderr:   &stderr,
		CurrentBranch: func(_ context.Context) (string, error) {
			return "", fmt.Errorf("detached HEAD; cn cell resume requires being on cycle/{N}")
		},
	}
	err := r.Resume(t.Context(), ResumeArgs{Issue: 500})
	if err == nil {
		t.Fatal("expected refusal on detached HEAD")
	}
	if !strings.Contains(err.Error(), "review_resume_branch_unknown") {
		t.Errorf("expected review_resume_branch_unknown, got: %v", err)
	}
}

// --- Smoke: Resumer.Resume with existing self-coherence.md ---

func TestResumer_Resume_AppendRound(t *testing.T) {
	// Build a fake repo structure: .cdd/unreleased/9999/self-coherence.md
	// and a fake "origin" git branch (skipping actual git call for unit test).
	// We test only the artifact-directory and round-append logic here;
	// verifyBranchExists is integration-tested against a real git repo.
	dir := t.TempDir()
	artifactDir := filepath.Join(dir, ".cdd", "unreleased", "9999")
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	scPath := filepath.Join(artifactDir, "self-coherence.md")
	if err := os.WriteFile(scPath, []byte("## §R0\n\ntext\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// verifyBranchExists will fail (not a real git repo), so we test
	// appendRoundHeader directly as the sub-function under test.
	nextRound, err := nextRoundNumber(scPath)
	if err != nil {
		t.Fatal(err)
	}
	if nextRound != 1 {
		t.Errorf("nextRound: want 1, got %d", nextRound)
	}
	if err := appendRoundHeader(scPath, nextRound); err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(scPath)
	s := string(content)
	if !strings.Contains(s, "## §R0") {
		t.Error("§R0 must be preserved")
	}
	if !strings.Contains(s, "## §R1") {
		t.Error("§R1 must be appended")
	}
}
