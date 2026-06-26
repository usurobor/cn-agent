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
	"context"
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
	Repo   string    // GitHub repo slug, e.g. "usurobor/cnos"
	Stdout io.Writer
	Stderr io.Writer
	// RunGH is the function used to invoke the gh CLI. If nil, the package-level
	// runGH is used. Inject a test double to unit-test the label-mutation path.
	RunGH func(ctx context.Context, args []string, w io.Writer) error
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

	// Parse schema declaration from the review artifact frontmatter.
	schema, err := readSchemaField(args.ReviewPath)
	if err != nil {
		return fmt.Errorf("reading review artifact: %w", err)
	}
	if schema != "cn.operator-review.v1" {
		return fmt.Errorf("review artifact has unexpected schema %q; expected cn.operator-review.v1", schema)
	}

	fmt.Fprintf(r.Stdout, "✓ review artifact: %s (schema: %s)\n", args.ReviewPath, schema)
	fmt.Fprintf(r.Stdout, "  issue: %d  verdict: %s\n", args.Issue, args.Verdict)

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
		// Apply status:review → status:changes.
		return r.applyLabelTransition(ctx, args)
	default:
		// Already validated in ParseReturnArgs; belt-and-suspenders.
		return fmt.Errorf("unexpected verdict %q", args.Verdict)
	}
}

// applyLabelTransition removes status:review and adds status:changes on the
// GitHub issue. Uses the gh CLI which is an assumed runtime dependency.
func (r *Returner) applyLabelTransition(ctx context.Context, args ReturnArgs) error {
	issueStr := strconv.Itoa(args.Issue)
	repoFlag := r.Repo

	// Use injected RunGH if set, otherwise fall back to the package-level runGH.
	ghFn := r.RunGH
	if ghFn == nil {
		ghFn = runGH
	}

	// Remove status:review label.
	fmt.Fprintf(r.Stdout, "\nApplying label transition: status:review → status:changes ...\n")

	removeArgs := []string{"issue", "edit", issueStr, "--remove-label", "status:review"}
	if repoFlag != "" {
		removeArgs = append([]string{"--repo", repoFlag}, removeArgs...)
	}
	if err := ghFn(ctx, removeArgs, r.Stderr); err != nil {
		return fmt.Errorf("removing status:review label: %w", err)
	}
	fmt.Fprintf(r.Stdout, "  ✓ removed status:review\n")

	// Add status:changes label.
	addArgs := []string{"issue", "edit", issueStr, "--add-label", "status:changes"}
	if repoFlag != "" {
		addArgs = append([]string{"--repo", repoFlag}, addArgs...)
	}
	if err := ghFn(ctx, addArgs, r.Stderr); err != nil {
		return fmt.Errorf("adding status:changes label: %w", err)
	}
	fmt.Fprintf(r.Stdout, "  ✓ added status:changes\n")
	fmt.Fprintf(r.Stdout, "\nLabel transition complete: issue #%d is now status:changes.\n", args.Issue)
	fmt.Fprintf(r.Stdout, "Use 'cn cell resume --issue %d' to re-arm the cycle.\n", args.Issue)
	return nil
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

// verifyBranchExists checks that origin/{branch} exists using git ls-remote.
func verifyBranchExists(ctx context.Context, branch string, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--exit-code", "origin", "refs/heads/"+branch)
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git ls-remote: %w", err)
	}
	return nil
}

// readSchemaField reads the YAML frontmatter of the file at path and returns
// the value of the `schema:` field. Returns an error if the field is absent.
func readSchemaField(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inFrontmatter := false
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
				break // end of frontmatter
			}
			if strings.HasPrefix(line, "schema:") {
				val := strings.TrimSpace(strings.TrimPrefix(line, "schema:"))
				return val, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("schema field not found in frontmatter of %q", path)
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
