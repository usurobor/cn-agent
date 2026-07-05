// Package issuesfsm implements `cn issues fsm evaluate`: an issue-state
// reconciler. It observes an explicit fact snapshot for one issue (live via
// the GitHub REST API + local git, or from a --fixture for offline/test
// use), evaluates it against a declarative transition table, and prints the
// current state, observed facts, enabled transition, blocked reason, and
// proposed action.
//
// Phase 1 (cnos#568) shipped this read-only: no --apply flag, no code path
// writing a GitHub label, every proposed action (table.go's Rule.Action)
// text printed to stdout and never executed.
//
// Phase 2 (cnos#569) adds the guarded mutation path: an optional --apply
// flag on this same "evaluate" verb (not a new subcommand — see runEvaluate)
// that writes the proposed status label, but ONLY when Evaluate already
// produced outcome=="proposed" with a non-empty TargetState -- i.e. only
// when the declarative transition table's guards already passed. Without
// --apply, this command is exactly as read-only as Phase 1 shipped it
// (byte-identical decision logic and output for any fixture Phase 1's test
// suite already covers) -- see AC1 in .cdd/unreleased/569/gamma-scaffold.md.
//
// Ownership split (AC1): this package owns the generic fact-snapshot model,
// the evaluator engine, and the CLI. The CDS-specific transition table is
// data owned by cnos.cds (src/packages/cnos.cds/skills/cds/fsm/
// transitions.json) — Evaluate (table.go) is unaware of "CDS" as a concept
// beyond reading whatever table file it is given; it never hardcodes a CDS
// state name.
//
// Design authority: cnos#568. The dispatch wrapper lives in
// src/go/internal/cli/cmd_issues_fsm.go and holds no domain logic, per the
// dispatch boundary (INVARIANTS.md T-002, eng/go §2.18) — exactly mirroring
// the `cn issues map` / commands/issues-map precedent (cnos#556, cnos#392).
package issuesfsm

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// defaultTableRelPath is the repo-root-relative location of the CDS
// transition table (AC1: package-owned declarative data, not Go literals).
const defaultTableRelPath = "src/packages/cnos.cds/skills/cds/fsm/transitions.json"

// Run is the entry point for `cn issues fsm`. This package supports two
// sub-verbs: "evaluate" (one issue, named by the caller) and "scan" (every
// active dispatch:cell+protocol:{P} issue, cnos#593 Sub C's mechanical
// recovery scanner) — args[0] after the "issues fsm" noun-verb pair is
// resolved by the CLI dispatcher (see cmd_issues_fsm.go). Mutation authority
// (cnos#569 Phase 2 for evaluate; cnos#593 for scan) is a --apply FLAG on
// each sub-verb, not a separate "apply" sub-verb — the issue's own oracle
// lines (`cn issues fsm evaluate --issue {N} --apply`, `cn issues fsm scan
// --apply`) pin this shape.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "Usage: cn issues fsm evaluate --issue N [--apply] [--fixture path] [--table path]")
		fmt.Fprintln(stderr, "       cn issues fsm scan --protocol P [--apply] [--table path]")
		return fmt.Errorf("cn issues fsm: a subcommand is required")
	}
	switch args[0] {
	case "evaluate":
		return runEvaluate(ctx, args[1:], stdout, stderr)
	case "scan":
		return runScan(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "✗ cn issues fsm: unknown subcommand %q\n\n", args[0])
		fmt.Fprintln(stderr, "cn issues fsm supports only: evaluate, scan")
		fmt.Fprintln(stderr, "(mutation authority is the --apply flag on each sub-verb — cnos#569 Phase 2 for evaluate, cnos#593 for scan — not a separate 'apply' subcommand.)")
		return fmt.Errorf("cn issues fsm: unknown subcommand %q", args[0])
	}
}

func runEvaluate(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("issues fsm evaluate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	issue := fs.Int("issue", 0, "issue number to evaluate")
	fixture := fs.String("fixture", "", "read the fact snapshot from this JSON file instead of live GitHub/git observation (offline)")
	table := fs.String("table", "", "path to the transition table JSON (default: repo-root-relative "+defaultTableRelPath+")")
	repo := fs.String("repo", "", "target repository as owner/name (default: $GITHUB_REPOSITORY)")
	token := fs.String("token", "", "GitHub token (default: $GITHUB_TOKEN, then $GH_TOKEN)")
	apply := fs.Bool("apply", false, "apply the proposed transition (cnos#569 Phase 2): writes the status label ONLY when the decision's outcome is \"proposed\" with a non-empty target_state (i.e. only when the transition table's guards already passed). Off by default -- without --apply this command is exactly as read-only as Phase 1 (cnos#568) shipped it. Requires --repo (or $GITHUB_REPOSITORY) to know where to write; works with --fixture for hermetic testing (facts come from the fixture, the write still targets the real --repo).")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: cn issues fsm evaluate --issue N [--apply] [--fixture path] [--table path] [--repo owner/name]")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Without --apply (default, cnos#568 Phase 1 behavior): read-only. Prints")
		fmt.Fprintln(stderr, "current state, observed facts, enabled transition, blocked reason, and")
		fmt.Fprintln(stderr, "proposed action. Never mutates a label.")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "With --apply (cnos#569 Phase 2): additionally writes the proposed status")
		fmt.Fprintln(stderr, "label, but ONLY when the evaluated outcome is \"proposed\" with a non-empty")
		fmt.Fprintln(stderr, "target_state -- i.e. only when the transition table's guards passed. A")
		fmt.Fprintln(stderr, "\"blocked\" outcome exits nonzero and writes nothing (e.g. an empty")
		fmt.Fprintln(stderr, "status:review request with no PR/commits/REVIEW-REQUEST.yml evidence).")
		fmt.Fprintln(stderr)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	r := *repo
	if r == "" {
		r = os.Getenv("GITHUB_REPOSITORY")
	}
	tk := *token
	if tk == "" {
		tk = os.Getenv("GITHUB_TOKEN")
	}
	if tk == "" {
		tk = os.Getenv("GH_TOKEN")
	}

	var snap FactSnapshot
	var err error
	switch {
	case *fixture != "":
		snap, err = LoadFixture(*fixture)
		if err != nil {
			return err
		}
	default:
		if *issue == 0 {
			return fmt.Errorf("cn issues fsm evaluate: --issue is required (or use --fixture for offline mode)")
		}
		snap, err = assembleLive(ctx, r, *issue, tk)
		if err != nil {
			return err
		}
	}
	if *issue != 0 {
		snap.Issue = *issue
	}

	tablePath := *table
	if tablePath == "" {
		tablePath, err = resolveDefaultTablePath()
		if err != nil {
			return err
		}
	}
	t, err := LoadTable(tablePath)
	if err != nil {
		return err
	}

	dec, err := Evaluate(t, snap)
	if err != nil {
		return err
	}

	var applyErr error
	if *apply {
		dec.ApplyAttempted = true
		switch {
		case dec.Outcome == "proposed" && dec.TargetState != "":
			// Guard-gated mutation: reachable only because Evaluate already
			// matched a rule whose outcome is "proposed" -- the table's
			// guards passed (AC1). This is the package's only write path.
			if r == "" {
				applyErr = fmt.Errorf("cn issues fsm evaluate --apply: --repo (or $GITHUB_REPOSITORY) is required to apply a transition")
				break
			}
			if err := applyStatusLabel(ctx, r, dec.Issue, tk, dec.CurrentState, dec.TargetState); err != nil {
				applyErr = fmt.Errorf("cn issues fsm evaluate --apply: %w", err)
				break
			}
			dec.Applied = true
		case dec.Outcome == "blocked":
			// AC3: empty-review (or any other blocked transition) is
			// structurally refused -- no write attempted, nonzero exit.
			applyErr = fmt.Errorf("cn issues fsm evaluate --apply: transition blocked (%s); refusing to mutate any label", orNoneStr(dec.BlockedReason, "no reason given"))
		default:
			// outcome == "valid" (nothing proposed), or a "proposed"
			// action with no direct status target (e.g.
			// propose_delta_recovery): nothing this command can apply as
			// a label mutation. Not an error -- this is the AC1
			// idempotence case: re-running --apply against the
			// post-mutation state finds outcome=="valid" and is a no-op.
		}
	}

	dec.Render(stdout)
	if applyErr != nil {
		return applyErr
	}
	return nil
}

// runScan is the CLI wrapper for `cn issues fsm scan` (cnos#593 Sub C): the
// mechanical recovery scanner that reconciles every open
// dispatch:cell+protocol:{P} issue currently at status:in-progress or
// status:review, per RunScan (scan.go). Mirrors runEvaluate's flag-parsing
// shape (repo/token/table/apply) plus a required --protocol (scan has no
// single issue to derive it from; evaluate needs no protocol at all since
// its --issue already names the one issue to observe).
func runScan(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("issues fsm scan", flag.ContinueOnError)
	fs.SetOutput(stderr)
	protocol := fs.String("protocol", "", "the dispatch wake's owning protocol qualifier (e.g. \"cds\"); required -- selects which protocol:{P} labeled issues this scan reconciles")
	table := fs.String("table", "", "path to the transition table JSON (default: repo-root-relative "+defaultTableRelPath)
	repo := fs.String("repo", "", "target repository as owner/name (default: $GITHUB_REPOSITORY)")
	token := fs.String("token", "", "GitHub token (default: $GITHUB_TOKEN, then $GH_TOKEN)")
	apply := fs.Bool("apply", false, "apply proposed reconciliations: writes status labels via the same guarded primitive `evaluate --apply` uses, and invokes `cn cell finalize` + posts a recovery comment for dead runs with matter. Off by default -- without --apply this command is read-only (reports what it would do).")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: cn issues fsm scan --protocol P [--apply] [--table path] [--repo owner/name]")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Reconciles every open dispatch:cell+protocol:{P} issue currently at")
		fmt.Fprintln(stderr, "status:in-progress or status:review, per cnos#593: no active run + no")
		fmt.Fprintln(stderr, "matter -> requeue to status:todo; no active run + matter + no PR ->")
		fmt.Fprintln(stderr, "checkpoint via `cn cell finalize`; PR + REVIEW-REQUEST.yml present ->")
		fmt.Fprintln(stderr, "status:review; a live run, or a blocked/under-evidenced state, is left alone.")
		fmt.Fprintln(stderr)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *protocol == "" {
		return fmt.Errorf("cn issues fsm scan: --protocol is required")
	}

	r := *repo
	if r == "" {
		r = os.Getenv("GITHUB_REPOSITORY")
	}
	tk := *token
	if tk == "" {
		tk = os.Getenv("GITHUB_TOKEN")
	}
	if tk == "" {
		tk = os.Getenv("GH_TOKEN")
	}

	tablePath := *table
	var err error
	if tablePath == "" {
		tablePath, err = resolveDefaultTablePath()
		if err != nil {
			return err
		}
	}
	t, err := LoadTable(tablePath)
	if err != nil {
		return err
	}

	opts := &ScanOptions{
		Repo:     r,
		Token:    tk,
		Protocol: *protocol,
		Apply:    *apply,
		Table:    t,
	}
	report, runErr := RunScan(ctx, opts)
	renderScan(stdout, report)
	return runErr
}

// resolveDefaultTablePath walks up from the current working directory
// looking for defaultTableRelPath, so `cn issues fsm evaluate` works from
// any directory inside a cnos checkout without requiring --table.
func resolveDefaultTablePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, defaultTableRelPath)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cn issues fsm evaluate: could not find %s above %s; pass --table explicitly", defaultTableRelPath, mustGetwd())
		}
		dir = parent
	}
}

func mustGetwd() string {
	d, _ := os.Getwd()
	return d
}
