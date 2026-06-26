// Package cell implements the cn cell subcommand family: return and resume.
//
// cn cell return --issue N --verdict V --review path
//   Reads a cn.operator-review.v1 artifact; on iterate/reject verdicts,
//   transitions the GitHub issue from status:review to status:changes.
//
// cn cell resume --issue N
//   Re-arms an existing cycle on cycle/{N}: verifies branch and artifact
//   directory exist, appends R[N+1] header to self-coherence.md so the
//   next dispatched round knows it is a continuation.
//
// Design authority: cnos#500; implementation contract pinned in gamma-scaffold.md.
package cell

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ReturnArgs holds the parsed flags for cn cell return.
type ReturnArgs struct {
	Issue      int
	Verdict    string // converge | iterate | reject | clarify
	ReviewPath string // path to cn.operator-review.v1 artifact
}

// ResumeArgs holds the parsed flags for cn cell resume.
type ResumeArgs struct {
	Issue int
}

// ParseReturnArgs parses argv for cn cell return.
// Expected: --issue N --verdict V --review path
func ParseReturnArgs(argv []string) (ReturnArgs, error) {
	var a ReturnArgs
	for i := 0; i < len(argv); i++ {
		switch argv[i] {
		case "--issue":
			if i+1 >= len(argv) {
				return a, fmt.Errorf("--issue requires a value")
			}
			i++
			n, err := strconv.Atoi(argv[i])
			if err != nil || n <= 0 {
				return a, fmt.Errorf("--issue must be a positive integer, got %q", argv[i])
			}
			a.Issue = n
		case "--verdict":
			if i+1 >= len(argv) {
				return a, fmt.Errorf("--verdict requires a value")
			}
			i++
			a.Verdict = argv[i]
		case "--review":
			if i+1 >= len(argv) {
				return a, fmt.Errorf("--review requires a value")
			}
			i++
			a.ReviewPath = argv[i]
		default:
			return a, fmt.Errorf("unknown flag %q", argv[i])
		}
	}
	if a.Issue == 0 {
		return a, fmt.Errorf("--issue is required")
	}
	if a.Verdict == "" {
		return a, fmt.Errorf("--verdict is required")
	}
	if a.ReviewPath == "" {
		return a, fmt.Errorf("--review is required")
	}
	switch a.Verdict {
	case "converge", "iterate", "reject", "clarify":
		// valid
	default:
		return a, fmt.Errorf("--verdict must be one of: converge, iterate, reject, clarify; got %q", a.Verdict)
	}
	return a, nil
}

// ParseResumeArgs parses argv for cn cell resume.
// Expected: --issue N
func ParseResumeArgs(argv []string) (ResumeArgs, error) {
	var a ResumeArgs
	for i := 0; i < len(argv); i++ {
		switch argv[i] {
		case "--issue":
			if i+1 >= len(argv) {
				return a, fmt.Errorf("--issue requires a value")
			}
			i++
			n, err := strconv.Atoi(argv[i])
			if err != nil || n <= 0 {
				return a, fmt.Errorf("--issue must be a positive integer, got %q", argv[i])
			}
			a.Issue = n
		default:
			return a, fmt.Errorf("unknown flag %q", argv[i])
		}
	}
	if a.Issue == 0 {
		return a, fmt.Errorf("--issue is required")
	}
	return a, nil
}

// Returner executes the cn cell return operation.
type Returner struct {
	Repo   string // GitHub repo slug, e.g. "usurobor/cnos"
	Stdout io.Writer
	Stderr io.Writer
	// RunGH is the function used to invoke the gh CLI for label-mutation
	// (write) commands. If nil, the package-level runGH is used. Inject a
	// test double to unit-test the label-mutation path.
	RunGH func(ctx context.Context, args []string, w io.Writer) error
	// RunGHJSON is the function used to invoke the gh CLI for read commands
	// that return JSON on stdout (e.g. `gh issue view --json labels,state`).
	// If nil, the package-level runGHJSON is used. Inject a test double to
	// unit-test the preflight path (F2 / F3, cycle/500 R1).
	RunGHJSON func(ctx context.Context, args []string) ([]byte, error)
}

// Return reads the operator-review artifact, validates it, and on
// iterate/reject applies the status:review → status:changes label transition.
// On converge, it reports the converge path and takes no label action.
// On clarify, it reports that clarification is needed and takes no label action.
func (r *Returner) Return(ctx context.Context, args ReturnArgs) error {
	// Validate the review artifact exists.
	if _, err := os.Stat(args.ReviewPath); err != nil {
		return fmt.Errorf("review artifact not found at %q: %w", args.ReviewPath, err)
	}

	// Parse schema + issue + verdict declarations from the review artifact
	// frontmatter. The artifact is the authority; CLI flags select and confirm
	// it. They must not override silently. F1 (cycle/500 R1): the artifact
	// must match the CLI flags before any label mutation.
	fm, err := readReviewFrontmatter(args.ReviewPath)
	if err != nil {
		return fmt.Errorf("reading review artifact: %w", err)
	}
	if fm.Schema != "cn.operator-review.v1" {
		return fmt.Errorf("review artifact has unexpected schema %q; expected cn.operator-review.v1", fm.Schema)
	}
	if fm.Issue == 0 {
		return fmt.Errorf("review artifact at %q is missing the 'issue:' frontmatter field (review_return_artifact_invalid)", args.ReviewPath)
	}
	if fm.Verdict == "" {
		return fmt.Errorf("review artifact at %q is missing the 'verdict:' frontmatter field (review_return_artifact_invalid)", args.ReviewPath)
	}
	if fm.Issue != args.Issue {
		return fmt.Errorf("review artifact issue mismatch: --issue=%d but artifact frontmatter says issue=%d (review_return_artifact_mismatch)", args.Issue, fm.Issue)
	}
	if fm.Verdict != args.Verdict {
		return fmt.Errorf("review artifact verdict mismatch: --verdict=%s but artifact frontmatter says verdict=%s (review_return_artifact_mismatch)", args.Verdict, fm.Verdict)
	}

	fmt.Fprintf(r.Stdout, "✓ review artifact: %s (schema: %s)\n", args.ReviewPath, fm.Schema)
	fmt.Fprintf(r.Stdout, "  issue: %d  verdict: %s (artifact ↔ flags match)\n", args.Issue, args.Verdict)

	switch args.Verdict {
	case "converge":
		fmt.Fprintf(r.Stdout, "\nVerdict: converge — no label transition applied.\n")
		fmt.Fprintf(r.Stdout, "Operator merges the PR to complete the cell lifecycle.\n")
		return nil
	case "clarify":
		fmt.Fprintf(r.Stdout, "\nVerdict: clarify — no label transition applied.\n")
		fmt.Fprintf(r.Stdout, "HI surfaces the clarification request to the operator. No cycle transition until clarification is received.\n")
		return nil
	case "iterate", "reject":
		// Preflight (F2 / F3, cycle/500 R1): inspect issue state and verify
		// the target label exists in the repo BEFORE any label mutation.
		// This is what makes cn cell return a primitive (rather than a thin
		// wrapper around naive gh issue edit).
		if err := r.preflightIssue(ctx, args); err != nil {
			return err
		}
		if err := r.preflightTargetLabel(ctx, args); err != nil {
			return err
		}
		// Apply status:review → status:changes (drift-aware; F3).
		return r.applyLabelTransition(ctx, args)
	default:
		// Already validated in ParseReturnArgs; belt-and-suspenders.
		return fmt.Errorf("unexpected verdict %q", args.Verdict)
	}
}

// preflightIssue verifies the GitHub issue is in a state that admits the
// status:review → status:changes transition. F2 (cycle/500 R1) invariants:
//   - exactly one status:* label present
//   - status:review present
//   - status:changes absent
//   - issue open
//
// On violation, returns review_return_state_invalid with a precise reason.
func (r *Returner) preflightIssue(ctx context.Context, args ReturnArgs) error {
	issueStr := strconv.Itoa(args.Issue)
	ghArgs := []string{"issue", "view", issueStr, "--json", "state,labels"}
	if r.Repo != "" {
		ghArgs = append([]string{"--repo", r.Repo}, ghArgs...)
	}
	ghJSON := r.RunGHJSON
	if ghJSON == nil {
		ghJSON = runGHJSON
	}
	fmt.Fprintf(r.Stdout, "\nPreflighting issue #%d state ...\n", args.Issue)
	out, err := ghJSON(ctx, ghArgs)
	if err != nil {
		return fmt.Errorf("inspecting issue #%d: %w", args.Issue, err)
	}
	state, labels, perr := parseIssueStateLabels(out)
	if perr != nil {
		return fmt.Errorf("parsing issue #%d response: %w", args.Issue, perr)
	}
	if !strings.EqualFold(state, "OPEN") {
		return fmt.Errorf("review_return_state_invalid: issue #%d is %s; transition requires open issue", args.Issue, state)
	}
	statusLabels := filterStatusLabels(labels)
	if len(statusLabels) == 0 {
		return fmt.Errorf("review_return_state_invalid: issue #%d carries no status:* label; transition requires exactly status:review", args.Issue)
	}
	if len(statusLabels) > 1 {
		return fmt.Errorf("review_return_state_invalid: issue #%d carries multiple status:* labels %v; transition requires exactly status:review", args.Issue, statusLabels)
	}
	if statusLabels[0] != "status:review" {
		return fmt.Errorf("review_return_state_invalid: issue #%d is at %s; transition requires status:review", args.Issue, statusLabels[0])
	}
	// statusLabels[0] == "status:review" — by construction status:changes is
	// not present (single-status invariant just verified).
	fmt.Fprintf(r.Stdout, "  ✓ issue #%d is OPEN at status:review (preflight pass)\n", args.Issue)
	return nil
}

// preflightTargetLabel verifies status:changes exists in the repo's label
// set before the transition is attempted. F3 (cycle/500 R1) empirical
// witness: cnos#493 — the label-doctor gap previously left target labels
// missing; the destructive remove-then-add path would have stranded the
// issue statusless. Failing here BEFORE any destructive action.
func (r *Returner) preflightTargetLabel(ctx context.Context, args ReturnArgs) error {
	ghArgs := []string{"label", "list", "--limit", "200", "--json", "name"}
	if r.Repo != "" {
		ghArgs = append([]string{"--repo", r.Repo}, ghArgs...)
	}
	ghJSON := r.RunGHJSON
	if ghJSON == nil {
		ghJSON = runGHJSON
	}
	out, err := ghJSON(ctx, ghArgs)
	if err != nil {
		return fmt.Errorf("inspecting repo labels: %w", err)
	}
	names, perr := parseLabelNames(out)
	if perr != nil {
		return fmt.Errorf("parsing label list response: %w", perr)
	}
	for _, n := range names {
		if n == "status:changes" {
			fmt.Fprintf(r.Stdout, "  ✓ target label status:changes exists in repo (preflight pass)\n")
			return nil
		}
	}
	return fmt.Errorf("review_return_target_label_missing: status:changes label is not defined in the repo; refusing to remove status:review and strand the issue. Run label-doctor before retrying")
}

// applyLabelTransition transitions the issue from status:review to
// status:changes. Uses the gh CLI which is an assumed runtime dependency.
//
// F3 (cycle/500 R1) atomicity: gh issue edit accepts --remove-label and
// --add-label in a single invocation, which the GitHub API serves as one
// labels-PATCH call. This preserves the single-status invariant and removes
// the prior "remove succeeded; add failed; issue statusless" failure mode
// (empirical witness: cnos#493 label-doctor gap; runtime exercise of the
// previous code path stranded cnos#500 at no status during the bootstrap
// recovery TODAY).
//
// If the gh call fails, we re-inspect the issue state to report whether
// the transition partially applied (review_return_label_drift) or remained
// at status:review (no drift; safe to retry).
func (r *Returner) applyLabelTransition(ctx context.Context, args ReturnArgs) error {
	issueStr := strconv.Itoa(args.Issue)
	repoFlag := r.Repo

	// Use injected RunGH if set, otherwise fall back to the package-level runGH.
	ghFn := r.RunGH
	if ghFn == nil {
		ghFn = runGH
	}

	fmt.Fprintf(r.Stdout, "\nApplying label transition: status:review → status:changes (atomic) ...\n")

	// Single gh call: both --remove-label and --add-label. gh + the labels
	// API resolve this as a single PATCH against the issue's labels[],
	// preserving the single-status invariant.
	editArgs := []string{"issue", "edit", issueStr,
		"--remove-label", "status:review",
		"--add-label", "status:changes",
	}
	if repoFlag != "" {
		editArgs = append([]string{"--repo", repoFlag}, editArgs...)
	}
	if err := ghFn(ctx, editArgs, r.Stderr); err != nil {
		// Atomic call failed. Re-inspect to report drift status precisely.
		drift := r.assessPostFailureDrift(ctx, args)
		return fmt.Errorf("applying label transition for issue #%d: %w (%s)", args.Issue, err, drift)
	}
	fmt.Fprintf(r.Stdout, "  ✓ removed status:review; added status:changes\n")
	fmt.Fprintf(r.Stdout, "\nLabel transition complete: issue #%d is now status:changes.\n", args.Issue)
	fmt.Fprintf(r.Stdout, "Use 'cn cell resume --issue %d' to re-arm the cycle.\n", args.Issue)
	return nil
}

// assessPostFailureDrift inspects the issue after a transition failure and
// returns a structured marker describing whether the issue is statusless,
// at status:changes, or still at status:review. Best-effort: if the
// inspection itself fails, returns a marker noting that the post-failure
// state is unknown so the operator inspects manually.
func (r *Returner) assessPostFailureDrift(ctx context.Context, args ReturnArgs) string {
	issueStr := strconv.Itoa(args.Issue)
	ghArgs := []string{"issue", "view", issueStr, "--json", "state,labels"}
	if r.Repo != "" {
		ghArgs = append([]string{"--repo", r.Repo}, ghArgs...)
	}
	ghJSON := r.RunGHJSON
	if ghJSON == nil {
		ghJSON = runGHJSON
	}
	out, err := ghJSON(ctx, ghArgs)
	if err != nil {
		return "review_return_label_drift: post-failure state unknown — inspect issue manually"
	}
	_, labels, perr := parseIssueStateLabels(out)
	if perr != nil {
		return "review_return_label_drift: post-failure state unparseable — inspect issue manually"
	}
	statusLabels := filterStatusLabels(labels)
	switch {
	case len(statusLabels) == 0:
		return "review_return_label_drift: issue is now statusless — manual recovery required"
	case len(statusLabels) == 1 && statusLabels[0] == "status:changes":
		return "review_return_label_drift: transition partially applied; issue is at status:changes but the runtime reported failure; safe to treat as transitioned"
	case len(statusLabels) == 1 && statusLabels[0] == "status:review":
		return "no drift: issue remained at status:review; safe to retry"
	default:
		return fmt.Sprintf("review_return_label_drift: issue carries unexpected status labels %v — manual recovery required", statusLabels)
	}
}

// parseIssueStateLabels extracts state ("OPEN"/"CLOSED") and label names
// from the JSON returned by `gh issue view --json state,labels`.
func parseIssueStateLabels(jsonBytes []byte) (state string, labels []string, err error) {
	type label struct {
		Name string `json:"name"`
	}
	var resp struct {
		State  string  `json:"state"`
		Labels []label `json:"labels"`
	}
	if jerr := decodeJSON(jsonBytes, &resp); jerr != nil {
		return "", nil, jerr
	}
	names := make([]string, 0, len(resp.Labels))
	for _, l := range resp.Labels {
		names = append(names, l.Name)
	}
	return resp.State, names, nil
}

// parseLabelNames extracts label names from the JSON returned by
// `gh label list --json name`.
func parseLabelNames(jsonBytes []byte) ([]string, error) {
	type label struct {
		Name string `json:"name"`
	}
	var resp []label
	if err := decodeJSON(jsonBytes, &resp); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(resp))
	for _, l := range resp {
		names = append(names, l.Name)
	}
	return names, nil
}

// filterStatusLabels returns labels whose name starts with "status:".
func filterStatusLabels(labels []string) []string {
	out := make([]string, 0, len(labels))
	for _, l := range labels {
		if strings.HasPrefix(l, "status:") {
			out = append(out, l)
		}
	}
	return out
}

// Resumer executes the cn cell resume operation.
type Resumer struct {
	RepoRoot string    // absolute path to the repository root
	Stdout   io.Writer
	Stderr   io.Writer
}

// Resume re-arms an existing cycle on cycle/{N}:
// - verifies the cycle branch exists on origin
// - verifies the .cdd/unreleased/{N}/ artifact directory exists
// - appends an R[N+1] section header to self-coherence.md
func (r *Resumer) Resume(ctx context.Context, args ResumeArgs) error {
	branch := fmt.Sprintf("cycle/%d", args.Issue)
	artifactDir := filepath.Join(r.RepoRoot, ".cdd", "unreleased", strconv.Itoa(args.Issue))

	// Verify the cycle branch exists on origin.
	fmt.Fprintf(r.Stdout, "Resuming cycle/%d ...\n", args.Issue)
	if err := verifyBranchExists(ctx, branch, r.Stderr); err != nil {
		return fmt.Errorf("cycle branch %q not found on origin: %w", branch, err)
	}
	fmt.Fprintf(r.Stdout, "  ✓ branch %s exists on origin\n", branch)

	// Verify the artifact directory exists.
	if _, err := os.Stat(artifactDir); err != nil {
		return fmt.Errorf("artifact directory %q not found: %w", artifactDir, err)
	}
	fmt.Fprintf(r.Stdout, "  ✓ artifact directory .cdd/unreleased/%d/ exists\n", args.Issue)

	// Determine current round by counting existing §R[N] sections.
	selfCoherencePath := filepath.Join(artifactDir, "self-coherence.md")
	nextRound, err := nextRoundNumber(selfCoherencePath)
	if err != nil {
		return fmt.Errorf("reading self-coherence.md: %w", err)
	}

	// Append R[N+1] section header to self-coherence.md.
	if err := appendRoundHeader(selfCoherencePath, nextRound); err != nil {
		return fmt.Errorf("appending round header to self-coherence.md: %w", err)
	}
	fmt.Fprintf(r.Stdout, "  ✓ appended §R%d header to self-coherence.md\n", nextRound)

	fmt.Fprintf(r.Stdout, "\nCycle re-armed: branch %s, artifact directory .cdd/unreleased/%d/, round R%d.\n",
		branch, args.Issue, nextRound)
	fmt.Fprintf(r.Stdout, "Existing artifacts preserved. δ routes α R%d → β R%d → γ closeout amendment.\n",
		nextRound, nextRound)
	return nil
}

// runGH executes the gh CLI with the given args. stderr is forwarded to w.
func runGH(ctx context.Context, args []string, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "gh", args...)
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

// runGHJSON executes the gh CLI with the given args, captures stdout and
// returns it as bytes. Used for read-only inspection commands like
// `gh issue view --json state,labels` (F2 / F3 preflight, cycle/500 R1).
func runGHJSON(ctx context.Context, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gh %s: %w (stderr: %s)", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.Bytes(), nil
}

// decodeJSON unmarshals jsonBytes into v.
func decodeJSON(jsonBytes []byte, v any) error {
	return json.Unmarshal(jsonBytes, v)
}

// verifyBranchExists checks that origin/{branch} exists using git ls-remote.
func verifyBranchExists(ctx context.Context, branch string, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--exit-code", "origin", "refs/heads/"+branch)
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git ls-remote: %w", err)
	}
	return nil
}

// reviewFrontmatter captures the operator-review fields the runtime cares
// about: schema (which gate selects the right shape), and issue + verdict
// (which the artifact-as-authority doctrine pins as the source of truth for
// the `cn cell return` invocation; F1 / cycle/500 R1).
type reviewFrontmatter struct {
	Schema  string
	Issue   int
	Verdict string
}

// readSchemaField reads the YAML frontmatter of the file at path and returns
// the value of the `schema:` field. Returns an error if the field is absent.
// Retained for back-compat / focused unit tests; production callers use
// readReviewFrontmatter which captures issue + verdict in addition to schema.
func readSchemaField(path string) (string, error) {
	fm, err := readReviewFrontmatter(path)
	if err != nil {
		return "", err
	}
	if fm.Schema == "" {
		return "", fmt.Errorf("schema field not found in frontmatter of %q", path)
	}
	return fm.Schema, nil
}

// readReviewFrontmatter parses the YAML frontmatter of the file at path and
// extracts the schema / issue / verdict fields the runtime needs to verify
// the artifact-as-authority contract. Returns an error if the file has no
// frontmatter block at all; absent individual fields are returned as zero
// values for the caller to interpret.
func readReviewFrontmatter(path string) (reviewFrontmatter, error) {
	var fm reviewFrontmatter
	f, err := os.Open(path)
	if err != nil {
		return fm, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	frontmatterClosed := false
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if lineNum == 1 {
			if line == "---" {
				inFrontmatter = true
			}
			continue
		}
		if inFrontmatter {
			if line == "---" {
				frontmatterClosed = true
				break // end of frontmatter
			}
			switch {
			case strings.HasPrefix(line, "schema:"):
				fm.Schema = strings.TrimSpace(strings.TrimPrefix(line, "schema:"))
			case strings.HasPrefix(line, "issue:"):
				val := strings.TrimSpace(strings.TrimPrefix(line, "issue:"))
				if n, err := strconv.Atoi(val); err == nil {
					fm.Issue = n
				}
			case strings.HasPrefix(line, "verdict:"):
				fm.Verdict = strings.TrimSpace(strings.TrimPrefix(line, "verdict:"))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fm, err
	}
	if !inFrontmatter || !frontmatterClosed {
		return fm, fmt.Errorf("YAML frontmatter not found in %q (expected leading and trailing '---')", path)
	}
	return fm, nil
}

// nextRoundNumber returns the next round number by counting §R[N] section
// markers in self-coherence.md. If the file does not exist, returns 1.
// Pattern: lines starting with "## §R" followed by a digit.
func nextRoundNumber(path string) (int, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	defer f.Close()

	max := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Match "## §R<digits>" section headers.
		if strings.HasPrefix(line, "## §R") {
			rest := strings.TrimPrefix(line, "## §R")
			// Take the leading digits.
			end := 0
			for end < len(rest) && rest[end] >= '0' && rest[end] <= '9' {
				end++
			}
			if end > 0 {
				n, err := strconv.Atoi(rest[:end])
				if err == nil && n > max {
					max = n
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return max + 1, nil
}

// appendRoundHeader appends a §R[N] section marker to self-coherence.md.
func appendRoundHeader(path string, round int) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	header := fmt.Sprintf("\n---\n\n## §R%d\n\n<!-- α: append R%d artifacts here (AC verification, review-ready signal) -->\n", round, round)
	_, err = fmt.Fprint(f, header)
	return err
}
