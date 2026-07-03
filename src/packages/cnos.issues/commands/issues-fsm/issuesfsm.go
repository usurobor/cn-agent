// Package issuesfsm implements `cn issues fsm evaluate`: a read-only issue-
// state reconciler (cnos#568 Phase 1). It observes an explicit fact snapshot
// for one issue (live via the GitHub REST API + local git, or from a
// --fixture for offline/test use), evaluates it against a declarative
// transition table, and prints the current state, observed facts, enabled
// transition, blocked reason, and proposed action.
//
// Phase 1 is read-only by design: there is no --apply flag and no code path
// anywhere in this package writes a GitHub label. Every proposed action
// (table.go's Rule.Action) is text printed to stdout, never executed. Label-
// write authority is Phase 2 (cnos#569) — see AC8.
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

// Run is the entry point for `cn issues fsm`. Phase 1 supports exactly one
// sub-verb, "evaluate" — args[0] after the "issues fsm" noun-verb pair is
// resolved by the CLI dispatcher (see cmd_issues_fsm.go). Do not add
// "apply" here: it is explicitly out of scope for Phase 1 (cnos#568,
// cnos#569 is Phase 2).
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "Usage: cn issues fsm evaluate --issue N [--fixture path] [--table path]")
		return fmt.Errorf("cn issues fsm: a subcommand is required")
	}
	switch args[0] {
	case "evaluate":
		return runEvaluate(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "✗ cn issues fsm: unknown subcommand %q\n\n", args[0])
		fmt.Fprintln(stderr, "Phase 1 (cnos#568) supports only: evaluate")
		fmt.Fprintln(stderr, "(mutation subcommands like 'apply' are Phase 2 — cnos#569 — and do not exist here.)")
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
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: cn issues fsm evaluate --issue N [--fixture path] [--table path] [--repo owner/name]")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Read-only (cnos#568 Phase 1): prints current state, observed facts,")
		fmt.Fprintln(stderr, "enabled transition, blocked reason, and proposed action. Never mutates")
		fmt.Fprintln(stderr, "a label. There is no --apply flag.")
		fmt.Fprintln(stderr)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
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
	dec.Render(stdout)
	return nil
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
