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
//
// Cycle/500 R1 operator-review corrections (operator-review.md F1–F5):
//   F1 — readReviewFrontmatter parses issue + verdict; Return verifies the
//        artifact matches the CLI flags before any mutation.
//   F2 — preflightIssue inspects state + status:* labels before mutation.
//   F3 — applyLabelTransition is a single atomic gh edit call carrying
//        both --remove-label and --add-label; preflightTargetLabel verifies
//        status:changes exists in the repo BEFORE mutation; on failure,
//        assessPostFailureDrift reports a precise drift marker. Empirical
//        witness for the F3 sub-concern: cnos#493 (label-doctor gap).
//   F4 — Resumer.Resume requires the caller to already be on cycle/{N}
//        (Option B v0; see verifyOnCycleBranch).
//   F5 — operator-review/SKILL.md §1.4 canonical captured_by examples
//        updated to κ-era identities.
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

// FinalizeArgs holds the parsed flags for cn cell finalize.
//
// Issue is optional (cnos#591 implementation contract): when zero, Finalize
// self-detects the target issue number from the current branch name
// (cycle/{N}) or, failing that, by scanning .cdd/unreleased/{N}/ for the
// issue with matter. An explicit --issue override exists for testability
// and manual invocation; the renderer-emitted mechanical step never passes
// it (it cannot know N in advance).
type FinalizeArgs struct {
	Issue   int
	BaseSHA string // optional; passed through to matter-scan (newly-committed paths since BaseSHA)
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

// ParseFinalizeArgs parses argv for cn cell finalize.
// Expected: [--issue N] [--base-sha SHA] — both optional (cnos#591; no
// required flags at all — an empty argv is valid and lets Finalize
// self-detect the issue from the current branch or the CDD tree).
func ParseFinalizeArgs(argv []string) (FinalizeArgs, error) {
	var a FinalizeArgs
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
		case "--base-sha":
			if i+1 >= len(argv) {
				return a, fmt.Errorf("--base-sha requires a value")
			}
			i++
			a.BaseSHA = argv[i]
		default:
			return a, fmt.Errorf("unknown flag %q", argv[i])
		}
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
	RepoRoot string // absolute path to the repository root
	Stdout   io.Writer
	Stderr   io.Writer
	// CurrentBranch returns the current local git branch name (no leading
	// "refs/heads/"). If nil, the package-level currentLocalBranch is used.
	// Inject a test double to unit-test the F4 (cycle/500 R1) cycle-branch
	// preflight without standing up a real git working tree.
	CurrentBranch func(ctx context.Context) (string, error)
}

// Resume re-arms an existing cycle on cycle/{N}:
//   - preflight: refuse unless the caller is already on cycle/{N}
//     (F4 design choice / cycle/500 R1: Option B v0 — local-only helper)
//   - verifies the cycle branch exists on origin
//   - verifies the .cdd/unreleased/{N}/ artifact directory exists
//   - appends an R[N+1] section header to self-coherence.md
//
// Caller MUST commit and push the §R[N+1] marker after Resume returns.
// Option A (fetch + auto-checkout + commit + push) is deferred to a future
// cycle; documented in self-coherence §R1 + Debt.
func (r *Resumer) Resume(ctx context.Context, args ResumeArgs) error {
	branch := fmt.Sprintf("cycle/%d", args.Issue)
	artifactDir := filepath.Join(r.RepoRoot, ".cdd", "unreleased", strconv.Itoa(args.Issue))

	// F4 (cycle/500 R1) preflight: the caller must already be on cycle/{N}.
	// Without this, the §R[N+1] marker could be appended to self-coherence.md
	// on main (or any other branch), corrupting cycle attribution. This is
	// the Option B v0 contract.
	fmt.Fprintf(r.Stdout, "Resuming cycle/%d ...\n", args.Issue)
	curFn := r.CurrentBranch
	if curFn == nil {
		curFn = currentLocalBranch
	}
	cur, err := curFn(ctx)
	if err != nil {
		return fmt.Errorf("review_resume_branch_unknown: cannot determine current git branch: %w", err)
	}
	if cur != branch {
		return fmt.Errorf("review_resume_wrong_branch: cn cell resume must be invoked from %q; current branch is %q. Run 'git checkout %s' first", branch, cur, branch)
	}
	fmt.Fprintf(r.Stdout, "  ✓ on branch %s (preflight pass)\n", branch)

	// Verify the cycle branch exists on origin.
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
	// F4 / Option B v0 (cycle/500 R1): the caller is responsible for the
	// commit and push of the §R[N+1] marker. Option A (the fetch + checkout
	// + commit + push variant) is a future-cycle item.
	fmt.Fprintf(r.Stdout, "\nNEXT: commit the §R%d header and push to origin/%s:\n", nextRound, branch)
	fmt.Fprintf(r.Stdout, "  git add %s\n", filepath.Join(".cdd", "unreleased", strconv.Itoa(args.Issue), "self-coherence.md"))
	fmt.Fprintf(r.Stdout, "  git commit -m \"alpha-%d: open §R%d\"\n", args.Issue, nextRound)
	fmt.Fprintf(r.Stdout, "  git push origin %s\n", branch)
	return nil
}

// FinalizeFacts is the finalizer's entire matter-detection input: an
// explicit, hidden-inference-free observation of one cell's branch/PR
// state (cnos#591). Mirrors the shape idea of issues-fsm's FactSnapshot
// (src/packages/cnos.issues/commands/issues-fsm/snapshot.go): decision
// logic (HasMatter) consumes a FinalizeFacts value directly and never
// re-derives facts by shelling out itself. This is what makes AC2 (the
// four-signal matter-detection table test) and AC9 (the strand-class
// fixture) constructible in-process, with no real git/gh subprocess.
type FinalizeFacts struct {
	// Issue is the issue number this snapshot describes, or 0 if no
	// candidate issue number could be resolved at all (self-detection
	// failure; see ParseFinalizeArgs / assembleFacts).
	Issue int

	// UncommittedChanges reports whether the working tree has staged or
	// unstaged changes (signal 1 of 4).
	UncommittedChanges bool

	// CommitsBeyondBase is the number of commits on cycle/{Issue} beyond
	// its base (main). Nonzero is signal 2 of 4.
	CommitsBeyondBase int

	// CDDArtifacts lists filenames present under .cdd/unreleased/{Issue}/.
	// Non-empty is signal 3 of 4.
	CDDArtifacts []string

	// BranchExists reports whether cycle/{Issue} exists (locally or on
	// origin). Signal 4 of 4.
	BranchExists bool

	// PRExists, PRNumber, PRURL describe an already-open pull request for
	// cycle/{Issue}, if one exists (gh pr list --head cycle/{Issue}
	// --state open). Used by the idempotent create/update decision (AC4),
	// not by matter-detection itself.
	PRExists bool
	PRNumber int
	PRURL    string
}

// HasMatter reports whether f represents "matter" per the governing rule
// restated in gamma-scaffold.md (cnos#591): the OR of four independently-
// sufficient signals — uncommitted file changes, commits beyond base on
// cycle/{N}, CDD artifacts under .cdd/unreleased/{N}/, or an existing
// cycle/{N} branch. Pure function over an explicit facts struct: no git/gh
// subprocess, no side effects. AC2 / AC9 oracle target.
func HasMatter(f FinalizeFacts) bool {
	return f.UncommittedChanges ||
		f.CommitsBeyondBase > 0 ||
		len(f.CDDArtifacts) > 0 ||
		f.BranchExists
}

// Finalizer executes the cn cell finalize operation: mechanical
// checkpoint (commit + push) and idempotent draft-PR create/update,
// gated purely on HasMatter. It never writes a status label and never
// requests an FSM transition — that authority remains with the FSM
// (issues-fsm) and stays entirely out of this code path (cnos#591 AC7/AC8;
// see gamma-scaffold.md Friction note 2).
type Finalizer struct {
	Repo     string // GitHub repo slug, e.g. "usurobor/cnos"; empty uses gh's default
	RepoRoot string // absolute path to the repository root
	Stdout   io.Writer
	Stderr   io.Writer

	// RunGH is the function used to invoke the gh CLI for write commands
	// (gh pr create). If nil, the package-level runGH is used. Mirrors
	// Returner's injection pattern exactly.
	RunGH func(ctx context.Context, args []string, w io.Writer) error
	// RunGHJSON is the function used to invoke the gh CLI for read
	// commands that return JSON on stdout (gh pr list). If nil, the
	// package-level runGHJSON is used. Mirrors Returner's injection
	// pattern exactly.
	RunGHJSON func(ctx context.Context, args []string) ([]byte, error)

	// AssembleFacts builds the FinalizeFacts value for a Finalize call.
	// If nil, the real git-shelling implementation (assembleFacts) is
	// used. Tests inject a fake returning a FinalizeFacts literal
	// directly, so the checkpoint+PR decision logic is exercised with
	// zero real git/gh subprocess calls (AC3/AC4/AC5/AC6/AC9).
	AssembleFacts func(ctx context.Context, args FinalizeArgs) (FinalizeFacts, error)

	// Checkpoint performs the mechanical commit+push step once matter is
	// confirmed: ensures cycle/{issue} exists, commits any uncommitted
	// changes, and pushes. If nil, the real git-shelling implementation
	// (realCheckpoint) is used. Tests inject a no-op fake so the PR
	// create/update decision can be exercised in isolation.
	Checkpoint func(ctx context.Context, issue int, facts FinalizeFacts) error
}

// Finalize is the entry point for cn cell finalize. It assembles (or
// accepts injected) matter-detection facts, no-ops cleanly when no matter
// exists (AC5), otherwise checkpoints the branch and ensures a draft PR
// exists — creating one if absent (AC3/AC9) or leaving an existing open
// PR untouched (AC4; see gamma-scaffold.md Friction note 5 for the
// no-op-vs-update design choice this cycle makes). On checkpoint or
// PR-open failure, Finalize never claims success: it writes a
// FINALIZE-STOP.md evidence file and returns a precise, package-
// convention-following error marker (AC6). It never writes a label and
// never calls `cn issues fsm evaluate` (AC7/AC8).
func (f *Finalizer) Finalize(ctx context.Context, args FinalizeArgs) error {
	assemble := f.AssembleFacts
	if assemble == nil {
		assemble = f.assembleFacts
	}
	facts, err := assemble(ctx, args)
	if err != nil {
		return fmt.Errorf("cell_finalize_facts_unavailable: assembling matter-detection facts: %w", err)
	}

	issue := args.Issue
	if issue == 0 {
		issue = facts.Issue
	}

	if !HasMatter(facts) {
		fmt.Fprintf(f.Stdout, "cn cell finalize: no matter detected — no-op (nothing to checkpoint, no PR touched).\n")
		return nil
	}
	if issue == 0 {
		// Matter signals present but no issue number resolvable at all —
		// defensive no-op rather than guessing a branch/PR target.
		fmt.Fprintf(f.Stdout, "cn cell finalize: matter signals present but no issue number could be resolved — no-op.\n")
		return nil
	}

	branch := fmt.Sprintf("cycle/%d", issue)
	fmt.Fprintf(f.Stdout, "cn cell finalize: matter detected for %s (issue #%d) — checkpointing.\n", branch, issue)

	checkpoint := f.Checkpoint
	if checkpoint == nil {
		checkpoint = f.realCheckpoint
	}
	if err := checkpoint(ctx, issue, facts); err != nil {
		return f.stopWithEvidence(issue, "cell_finalize_checkpoint_failed", err)
	}
	fmt.Fprintf(f.Stdout, "  ✓ checkpoint complete: %s committed and pushed\n", branch)

	return f.ensurePR(ctx, issue, branch)
}

// ensurePR performs the idempotent draft-PR create/update decision (AC4):
// an existing open PR for cycle/{issue} is left as a clean no-op (this
// cycle's chosen shape — see gamma-scaffold.md Friction note 5); absent
// one, a new draft PR is opened targeting main, referencing the issue
// (AC3/AC7/AC9). It never touches REVIEW-REQUEST.yml, never calls
// `cn issues fsm evaluate`, and never mutates a label (AC7/AC8).
func (f *Finalizer) ensurePR(ctx context.Context, issue int, branch string) error {
	ghJSON := f.RunGHJSON
	if ghJSON == nil {
		ghJSON = runGHJSON
	}
	listArgs := []string{"pr", "list", "--head", branch, "--state", "open", "--json", "number,url"}
	if f.Repo != "" {
		listArgs = append([]string{"--repo", f.Repo}, listArgs...)
	}
	out, err := ghJSON(ctx, listArgs)
	if err != nil {
		return f.stopWithEvidence(issue, "cell_finalize_pr_list_failed", err)
	}
	existing, perr := parsePRList(out)
	if perr != nil {
		return f.stopWithEvidence(issue, "cell_finalize_pr_list_unparseable", perr)
	}
	if len(existing) > 0 {
		fmt.Fprintf(f.Stdout, "cn cell finalize: PR already open for %s: %s (idempotent no-op; no duplicate created)\n", branch, existing[0].URL)
		return nil
	}

	ghFn := f.RunGH
	if ghFn == nil {
		ghFn = runGH
	}
	title := fmt.Sprintf("cycle/%d", issue)
	body := fmt.Sprintf("Refs #%d\n\nAutomated checkpoint opened by the mechanical cell finalizer (cnos#591). Draft — not yet review-ready.\n", issue)
	createArgs := []string{"pr", "create", "--draft", "--base", "main", "--head", branch, "--title", title, "--body", body}
	if f.Repo != "" {
		createArgs = append([]string{"--repo", f.Repo}, createArgs...)
	}
	if err := ghFn(ctx, createArgs, f.Stderr); err != nil {
		return f.stopWithEvidence(issue, "cell_finalize_pr_open_failed", err)
	}
	fmt.Fprintf(f.Stdout, "  ✓ opened draft PR for %s (Refs #%d)\n", branch, issue)
	return nil
}

// stopWithEvidence records a STOP-evidence artifact at
// .cdd/unreleased/{issue}/FINALIZE-STOP.md naming the failure class and
// the raw error (AC6), then returns a precise error carrying the same
// marker, following this package's review_return_*/review_resume_*
// naming convention. It never touches a label or REVIEW-REQUEST.yml —
// evidence is for a human or a later δ hard-block judgment, not an
// auto-escalation (gamma-scaffold.md Friction note 2).
func (f *Finalizer) stopWithEvidence(issue int, marker string, cause error) error {
	dir := filepath.Join(f.RepoRoot, ".cdd", "unreleased", strconv.Itoa(issue))
	if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
		return fmt.Errorf("%s: %w (additionally failed to create evidence directory %q: %v)", marker, cause, dir, mkErr)
	}
	path := filepath.Join(dir, "FINALIZE-STOP.md")
	content := fmt.Sprintf(
		"# FINALIZE-STOP\n\nmarker: %s\nissue: %d\nerror: %s\n\nThis is mechanical checkpoint/PR-open evidence written by `cn cell finalize`\n(cnos#591). It does not itself change any status label; a human or a\nlater hard-block judgment consumes this evidence.\n",
		marker, issue, cause,
	)
	if wErr := os.WriteFile(path, []byte(content), 0o644); wErr != nil {
		return fmt.Errorf("%s: %w (additionally failed to write evidence file %q: %v)", marker, cause, path, wErr)
	}
	fmt.Fprintf(f.Stderr, "✗ cn cell finalize: %s: %v (evidence: %s)\n", marker, cause, path)
	return fmt.Errorf("%s: %w", marker, cause)
}

// prListItem is one entry of `gh pr list --json number,url` output.
type prListItem struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// parsePRList decodes the JSON returned by `gh pr list --json number,url`.
func parsePRList(jsonBytes []byte) ([]prListItem, error) {
	var items []prListItem
	if err := decodeJSON(jsonBytes, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// assembleFacts is the production, real-git fact-assembly path (not
// exercised by unit tests — they inject Finalizer.AssembleFacts directly
// per AC2/AC9). It resolves the target issue number per the
// implementation contract: (1) current branch matches cycle/(\d+); else
// (2) scan .cdd/unreleased/{N}/ for the N with uncommitted or (given
// --base-sha) newly-committed matter. If neither resolves an issue
// number, Issue is left 0 and Finalize takes the AC5 no-matter path.
func (f *Finalizer) assembleFacts(ctx context.Context, args FinalizeArgs) (FinalizeFacts, error) {
	var facts FinalizeFacts

	issue := args.Issue
	if issue == 0 {
		if cur, err := currentLocalBranch(ctx); err == nil {
			if n, ok := parseCycleBranchName(cur); ok {
				issue = n
			}
		}
	}
	if issue == 0 {
		issue = f.scanForIssueWithMatter(ctx, args)
	}
	facts.Issue = issue
	if issue == 0 {
		return facts, nil
	}

	branch := fmt.Sprintf("cycle/%d", issue)
	facts.UncommittedChanges = gitHasUncommittedChanges(ctx, f.RepoRoot)
	facts.BranchExists = gitBranchExists(ctx, f.RepoRoot, branch)
	if facts.BranchExists {
		facts.CommitsBeyondBase = gitCommitsBeyondBase(ctx, f.RepoRoot, branch, "main")
	}
	facts.CDDArtifacts = listCDDArtifacts(f.RepoRoot, issue)
	return facts, nil
}

// scanForIssueWithMatter implements the second self-detection path: when
// the current branch is not cycle/{N}, scan .cdd/unreleased/*/ for
// uncommitted paths (always) or newly-committed paths since --base-sha
// (when given), and return the first N with a matching path. Returns 0
// if nothing is found (the caller's AC5 no-matter path).
func (f *Finalizer) scanForIssueWithMatter(ctx context.Context, args FinalizeArgs) int {
	candidates := gitPorcelainPaths(ctx, f.RepoRoot, ".cdd/unreleased/")
	if args.BaseSHA != "" {
		candidates = append(candidates, gitDiffNamesSinceBase(ctx, f.RepoRoot, args.BaseSHA, ".cdd/unreleased/")...)
	}
	for _, p := range candidates {
		if n, ok := extractIssueFromCDDPath(p); ok {
			return n
		}
	}
	return 0
}

// realCheckpoint is the production, real-git checkpoint implementation
// (not exercised by unit tests — they inject Finalizer.Checkpoint as a
// no-op fake). Ensures cycle/{issue} exists, commits any uncommitted
// changes, and pushes to origin.
func (f *Finalizer) realCheckpoint(ctx context.Context, issue int, facts FinalizeFacts) error {
	branch := fmt.Sprintf("cycle/%d", issue)

	if !facts.BranchExists {
		if err := runGitCmd(ctx, f.RepoRoot, f.Stderr, "checkout", "-B", branch); err != nil {
			return fmt.Errorf("creating branch %s: %w", branch, err)
		}
	}
	if facts.UncommittedChanges {
		if err := runGitCmd(ctx, f.RepoRoot, f.Stderr, "add", "-A"); err != nil {
			return fmt.Errorf("staging changes on %s: %w", branch, err)
		}
		msg := fmt.Sprintf("cell-finalize: checkpoint matter for cycle/%d", issue)
		if err := runGitCmd(ctx, f.RepoRoot, f.Stderr, "commit", "-m", msg); err != nil {
			return fmt.Errorf("committing checkpoint on %s: %w", branch, err)
		}
	}
	if err := runGitCmd(ctx, f.RepoRoot, f.Stderr, "push", "-u", "origin", branch); err != nil {
		return fmt.Errorf("pushing %s: %w", branch, err)
	}
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

// currentLocalBranch returns the current local git branch (no "refs/heads/"
// prefix). Returns an error if HEAD is detached or git is not available.
// F4 (cycle/500 R1) preflight surface; tests inject a mock instead.
func currentLocalBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git rev-parse: %w (stderr: %s)", err, stderr.String())
	}
	name := strings.TrimSpace(stdout.String())
	if name == "" || name == "HEAD" {
		return "", fmt.Errorf("detached HEAD; cn cell resume requires being on cycle/{N}")
	}
	return name, nil
}

// --- Finalizer git plumbing (cnos#591) ---
//
// These are real exec.Command calls used only by the production
// assembleFacts/realCheckpoint paths above. Unit tests never reach them:
// they inject Finalizer.AssembleFacts / Finalizer.Checkpoint fakes
// instead, so the matter-detection facts struct (and the checkpoint+PR
// decision it drives) is constructible in tests with no git subprocess
// at all (AC2/AC3/AC4/AC5/AC6/AC9).

// runGitCmd runs git with the given args in dir, forwarding stderr to w.
func runGitCmd(ctx context.Context, dir string, w io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, "git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

// runGitCmdOutput runs git with the given args in dir and returns stdout.
func runGitCmdOutput(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w (stderr: %s)", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String(), nil
}

// gitHasUncommittedChanges reports whether the working tree at dir has
// any staged or unstaged changes. Errors are treated as "no" (the caller
// falls back to the other three matter signals; a broken git invocation
// here should not itself manufacture a false matter signal).
func gitHasUncommittedChanges(ctx context.Context, dir string) bool {
	out, err := runGitCmdOutput(ctx, dir, "status", "--porcelain")
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) != ""
}

// gitBranchExists reports whether branch exists locally or on origin.
func gitBranchExists(ctx context.Context, dir, branch string) bool {
	if err := runGitCmd(ctx, dir, io.Discard, "rev-parse", "--verify", "--quiet", branch); err == nil {
		return true
	}
	if err := runGitCmd(ctx, dir, io.Discard, "rev-parse", "--verify", "--quiet", "origin/"+branch); err == nil {
		return true
	}
	return false
}

// gitCommitsBeyondBase returns the count of commits on branch beyond base
// (e.g. "main..cycle/591"). Returns 0 on any error (branch/base missing).
func gitCommitsBeyondBase(ctx context.Context, dir, branch, base string) int {
	out, err := runGitCmdOutput(ctx, dir, "rev-list", "--count", base+".."+branch)
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0
	}
	return n
}

// gitPorcelainPaths returns the paths reported by `git status --porcelain`
// scoped to pathspec (uncommitted changes only).
func gitPorcelainPaths(ctx context.Context, dir, pathspec string) []string {
	out, err := runGitCmdOutput(ctx, dir, "status", "--porcelain", "--", pathspec)
	if err != nil {
		return nil
	}
	var paths []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		paths = append(paths, fields[len(fields)-1])
	}
	return paths
}

// gitDiffNamesSinceBase returns paths under pathspec changed between
// baseSHA and HEAD (newly-committed matter since the wake's baseline).
func gitDiffNamesSinceBase(ctx context.Context, dir, baseSHA, pathspec string) []string {
	out, err := runGitCmdOutput(ctx, dir, "diff", "--name-only", baseSHA+"..HEAD", "--", pathspec)
	if err != nil {
		return nil
	}
	var paths []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

// extractIssueFromCDDPath extracts N from a path containing
// ".../unreleased/{N}/...". Returns ok=false if no such segment is found.
func extractIssueFromCDDPath(p string) (int, bool) {
	parts := strings.Split(filepath.ToSlash(p), "/")
	for i, part := range parts {
		if part == "unreleased" && i+1 < len(parts) {
			if n, err := strconv.Atoi(parts[i+1]); err == nil && n > 0 {
				return n, true
			}
		}
	}
	return 0, false
}

// listCDDArtifacts lists filenames present under
// repoRoot/.cdd/unreleased/{issue}/ (non-recursive; files only). Returns
// nil if the directory does not exist.
func listCDDArtifacts(repoRoot string, issue int) []string {
	dir := filepath.Join(repoRoot, ".cdd", "unreleased", strconv.Itoa(issue))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names
}

// parseCycleBranchName extracts N from a branch name of the form
// "cycle/{N}". Returns ok=false for any other shape.
func parseCycleBranchName(name string) (int, bool) {
	const prefix = "cycle/"
	if !strings.HasPrefix(name, prefix) {
		return 0, false
	}
	n, err := strconv.Atoi(name[len(prefix):])
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
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
